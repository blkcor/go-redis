package reply

/*
*
UnKnownErrorReply
*/
type UnKnownErrorReply struct{}

var unKnowErrorBytes = []byte("-ERR unknown\r\n")

func (e *UnKnownErrorReply) Error() string {
	return "-ERR unknown"
}

func (e *UnKnownErrorReply) ToBytes() []byte {
	return unKnowErrorBytes
}

/*
*
ArgNumErrorReply
*/
type ArgNumErrorReply struct {
	Cmd string
}

func (e *ArgNumErrorReply) Error() string {
	return "-ERR wrong number of arguments of " + e.Cmd + "command\r\n"
}

func (e *ArgNumErrorReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments of " + e.Cmd + "command\r\n")
}

func MakeArgNumErrorReply(cmd string) *ArgNumErrorReply {
	return &ArgNumErrorReply{cmd}
}

/**
SyntaxErrorReply
*/

type SyntaxErrorReply struct{}

func (e *SyntaxErrorReply) Error() string {
	return "-ERR syntax error\r\n"
}

func (e *SyntaxErrorReply) ToBytes() []byte {
	return []byte("-ERR syntax error\r\n")
}

/**
WrongTypeErrorReply
*/

type WrongTypeErrorReply struct{}

func (e *WrongTypeErrorReply) Error() string {
	return "-ERR wrong type\r\n"
}

func (e *WrongTypeErrorReply) ToBytes() []byte {
	return []byte("-ERR wrong type\r\n")
}

/**
ProtocolErrorReply
*/

type ProtocolErrorReply struct{}

func (e *ProtocolErrorReply) Error() string {
	return "-ERR protocol error\r\n"
}

func (e *ProtocolErrorReply) ToBytes() []byte {
	return []byte("-ERR protocol error\r\n")
}
