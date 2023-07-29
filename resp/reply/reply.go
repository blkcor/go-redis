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

/**
* BulkReply
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

/**
* MultiBulkReply
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

/**
* StatusReply
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

/**
*IntReply
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

/**
*StandardErrorReply
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

/**
*判断回复是否是错误回复
 */
func IsErrorReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
