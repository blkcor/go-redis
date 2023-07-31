package reply

/*
UnKnownErrorReply 未知错误回复
*/
type UnKnownErrorReply struct{}

var unKnowErrorBytes = []byte("-Err unknown\r\n")

func (e *UnKnownErrorReply) Error() string {
	return "Err unknown"
}

func (e *UnKnownErrorReply) ToBytes() []byte {
	return unKnowErrorBytes
}

/*
ArgNumErrorReply 参数数量错误回复
@Cmd 用户输入的命令
*/
type ArgNumErrorReply struct {
	Cmd string
}

func (e *ArgNumErrorReply) Error() string {
	return "-Err wrong number of arguments of " + e.Cmd + "command\r\n"
}

func (e *ArgNumErrorReply) ToBytes() []byte {
	return []byte("-Err wrong number of arguments of " + e.Cmd + "command\r\n")
}

func MakeArgNumErrorReply(cmd string) *ArgNumErrorReply {
	return &ArgNumErrorReply{cmd}
}

/*
SyntaxErrorReply 语法错误回复
*/

type SyntaxErrorReply struct{}

func (e *SyntaxErrorReply) Error() string {
	return "-Err syntax error\r\n"
}

func (e *SyntaxErrorReply) ToBytes() []byte {
	return []byte("-Err syntax error\r\n")
}

func MakeSyntaxErrorReply() *SyntaxErrorReply {
	return &SyntaxErrorReply{}
}

/**
WrongTypeErrorReply 类型错误回复
*/

type WrongTypeErrorReply struct{}

func (e *WrongTypeErrorReply) Error() string {
	return "-Err wrong type\r\n"
}

func (e *WrongTypeErrorReply) ToBytes() []byte {
	return []byte("-Err wrong type\r\n")
}

func MakeWrongTypeErrorReply() *WrongTypeErrorReply {
	return &WrongTypeErrorReply{}
}

/*
ProtocolErrorReply 协议错误回复
*/

type ProtocolErrorReply struct{}

func (e *ProtocolErrorReply) Error() string {
	return "-Err protocol error\r\n"
}

func (e *ProtocolErrorReply) ToBytes() []byte {
	return []byte("-Err protocol error\r\n")
}
