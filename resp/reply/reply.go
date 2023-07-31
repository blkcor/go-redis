package reply

import (
	"bytes"
	"github.com/blkcor/go-redis/interface/resp"
	"strconv"
)

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

/*
BulkReply 块回复（这里将字符串转换成redis协议支持的格式）
@Arg 回复的字符串
*/
type BulkReply struct {
	Arg []byte
}

func (b BulkReply) ToBytes() []byte {
	if b.Arg == nil {
		return nullBulkReplyBytes
	}
	reply := "$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF
	return []byte(reply)
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{arg}
}

/*
MultiBulkReply 多块回复
@Args 回复的字符串数组
*/

type MultiBulkReply struct {
	Args [][]byte
}

func (m MultiBulkReply) ToBytes() []byte {
	argLen := len(m.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	for _, arg := range m.Args {
		if arg == nil {
			buf.WriteString(string(nullBulkReplyBytes) + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{args}
}

/*
StatusReply 状态回复
@Status 回复的状态
*/

type StatusReply struct {
	Status string
}

func (s StatusReply) ToBytes() []byte {
	return []byte("+" + s.Status + CRLF)
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{Status: status}
}

/*
IntReply 整数回复
@Code 回复的状态码
*/

type IntReply struct {
	Code int64
}

func (i IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)
}

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

func MakeIntReply(code int) *IntReply {
	return &IntReply{Code: int64(code)}
}

/*
StandardErrorReply 标准错误回复
@Status 回复的状态
*/

type StandardErrorReply struct {
	Status string
}

func (s StandardErrorReply) Error() string {
	return s.Status
}

func (s StandardErrorReply) ToBytes() []byte {
	return []byte("-" + s.Status + CRLF)
}

func MakeStandardErrorReply(status string) *StandardErrorReply {
	return &StandardErrorReply{Status: status}
}

/*
IsErrorReply 判断回复是否是错误回复
*/
func IsErrorReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
