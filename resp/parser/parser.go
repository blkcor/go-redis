package parser

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/blkcor/go-redis/interface/resp"
	"github.com/blkcor/go-redis/lib/logger"
	"github.com/blkcor/go-redis/resp/reply"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
)

/*
Payload 客户端发送的载荷
@Data	客户端发送的数据
@Err	错误
*/
type Payload struct {
	Data resp.Reply
	Err  error
}

/*
readState 读取状态
@readingMultiLine	是否正在读取多行数据
@expectedArgsCount	期望的参数数量
@msgType	消息类型
@args	参数
@bulkLen	块长度
*/
type readState struct {
	readingMultiLine  bool
	expectedArgsCount int
	msgType           byte
	args              [][]byte
	bulkLen           int64
}

/*
finished 判断解析是否完成
*/
func (r *readState) finished() bool {
	return r.expectedArgsCount > 0 && r.expectedArgsCount == len(r.args)
}

func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

/*
parse0 异步解析数据流
*/
func parse0(reader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(debug.Stack())
		}
	}()
	bufReader := bufio.NewReader(reader)
	var state readState
	var err error
	var msg []byte

	for true {
		var ioErr bool
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			if ioErr {
				ch <- &Payload{Err: err}
				return
			}
			ch <- &Payload{Err: errors.New("protocol error:" + err.Error())}
			state = readState{}
			continue
		}
		//解析多行数据
		if !state.readingMultiLine {
			if msg[0] == '*' {
				err = parseMultiBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{Err: errors.New("protocol error:" + string(msg))}
					state = readState{}
					continue
				}
				//user input an empty array
				if state.expectedArgsCount == 0 {
					ch <- &Payload{Data: reply.EmptyMultiBulkReply{}}
				}
				state = readState{}
				continue
			} else if msg[0] == '$' {
				err = parseBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{Err: errors.New("protocol error:" + string(msg))}
					state = readState{}
					continue
				}
				if state.bulkLen == -1 {
					ch <- &Payload{Data: reply.NullBulkReply{}}
					state = readState{}
					continue
				}
			} else {
				// + - :
				result, err := parseSingleLine(msg)
				ch <- &Payload{Data: result, Err: err}
				state = readState{}
				continue
			}
		} else {
			err = ReadBody(msg, &state)
			if err != nil {
				ch <- &Payload{Err: errors.New("protocol error:" + string(msg))}
				state = readState{}
				continue
			}
			if state.finished() {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply.MakeMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply.MakeBulkReply(state.args[0])
				}
				ch <- &Payload{Data: result}
				state = readState{}
			}
		}
	}
}

/*
readLine 读取单行数据并返回（包括\r\n）
return msg, ifIOError, err
example: *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$4\r\nvalue\r\n
*/
func readLine(buf *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error
	//case1: without preset bulkLen such as *3\r\n
	if state.bulkLen == 0 {
		msg, err = buf.ReadBytes('\n')
		if err != nil {
			return nil, false, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, true, errors.New("protocol error:" + string(msg))
		}
		//case2: with preset bulkLen
	} else {
		//+2 for \r\n
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(buf, msg)
		if err != nil {
			return nil, true, err
		}
		//check \r\n
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error:" + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}

/*
parseMultiBulkHeader 解析多行消息头部
*/
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var err error
	var expectedLine uint64
	re := string(msg[1 : len(msg)-2])
	fmt.Println(re)
	expectedLine, err = strconv.ParseUint(re, 10, 32)
	if err != nil {
		return errors.New("protocol error:" + string(msg))
	}
	//set readState
	if expectedLine == 0 {
		state.expectedArgsCount = 0
	} else if expectedLine > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = int(expectedLine)
		state.args = make([][]byte, 0, expectedLine)
	} else {
		return errors.New("protocol error:" + string(msg))
	}
	return nil
}

/*
parseBulkHeader 解析块消息头部
$3\r\n
*/
func parseBulkHeader(msg []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error:" + string(msg))
	}
	//empty bulk
	if state.bulkLen == -1 {
		return nil
	} else if state.bulkLen > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = 1
		state.args = make([][]byte, 0, 1)
		return nil
	} else {
		return errors.New("protocol error:" + string(msg))
	}
}

/*
parseSingleLine 解析单行，生成reply
*/
func parseSingleLine(msg []byte) (resp.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch str[0] {
	case '+':
		return reply.MakeStatusReply(str[1:]), nil
	case '-':
		return reply.MakeStandardErrorReply(str[1:]), nil
	case ':':
		res, err := strconv.Atoi(str[1:])
		if err != nil {
			return nil, errors.New("protocol error:" + string(msg))
		}
		result = reply.MakeIntReply(res)
	}
	return result, nil
}

/*
ReadBody 实现消息体的读取
*/
func ReadBody(msg []byte, state *readState) error {
	var err error
	//trim the \r\n
	line := msg[:len(msg)-2]
	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol error:" + string(msg))
		}
		if state.bulkLen <= 0 {
			//append an empty []byte
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}
	return nil
}
