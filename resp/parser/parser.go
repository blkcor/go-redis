package parser

import (
	"bufio"
	"errors"
	"github.com/blkcor/go-redis/interface/resp"
	"github.com/blkcor/go-redis/resp/reply"
	"io"
	"strconv"
	"strings"
)

type Payload struct {
	Data resp.Reply
	Err  error
}

type readState struct {
	readingMultiLine  bool
	expectedArgsCount int
	msgType           byte
	args              [][]byte
	bulkLen           int64
}

func (r *readState) finished() bool {
	return r.expectedArgsCount > 0 && r.expectedArgsCount == len(r.args)
}

func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {
	defer func() {

		if err := recover(); err != nil {

		}
	}()

	for true {

	}
}

/*
*
读取单行数据并返回（包括\r\n）
*/
func readLine(buf *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error
	//case1: without preset bulkLen
	if state.bulkLen == 0 {
		msg, err = buf.ReadBytes('\n')
		if err != nil {
			return nil, false, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error:" + string(msg))
		}
		//case2: with preset bulkLen
	} else {
		//+2 for \r\n
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(buf, msg)
		if err != nil {
			return nil, false, err
		}
		//check \r\n
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error:" + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, true, nil
}

/*
*解析多行消息头部
 */
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var err error
	var expectedLine uint64
	expectedLine, err = strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 32)
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
		state.args = make([][]byte, 0, 1)
	} else {
		return errors.New(" protocol error:" + string(msg))
	}
	return nil
}

/*
*
解析单行，生成reply
*/
func parseSingleLine(msg []byte) (resp.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch str[0] {
	case '+':
		return reply.MakeStatusReply(string(msg[1:])), nil
	case '-':
		return reply.MakeStandardErrorReply(string(msg[1:])), nil
	case ':':
		res, err := strconv.Atoi(string(msg[1:]))
		if err != nil {
			return nil, errors.New("protocol error:" + string(msg))
		}
		result = reply.MakeIntReply(res)
	}
	return result, nil
}

/*
*
实现消息体的读取
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
