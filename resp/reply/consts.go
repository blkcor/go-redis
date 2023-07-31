package reply

/*
*
Pong Reply
*/
type PongReply struct{}

var pongBytes = []byte("+PONG\r\n")

func (p PongReply) ToBytes() []byte {
	return pongBytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}

/**
Ok Reply
*/

type OkReply struct{}

var okBytes = []byte("+ok\r\n")

func (o OkReply) ToBytes() []byte {
	return okBytes
}

func MakeOkReply() *OkReply {
	return &OkReply{}
}

/**
NullBulkReply
*/

type NullBulkReply struct{}

var nullBulkBytes = []byte("$-1\r\n")

func (n NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

/**
EmptyMultiBulkReply 空结构体
*/

type EmptyMultiBulkReply struct{}

var emptyMultiBulkBytes = []byte("*0\r\n")

func (e EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}
func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

/**
NoReply
*/

type NoReply struct{}

var noBytes = []byte("")

func (n NoReply) ToBytes() []byte {
	return noBytes
}
func MakeNoReply() *NoReply {
	return &NoReply{}
}
