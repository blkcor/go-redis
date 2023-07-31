package resp

type Reply interface {
	//将Reply转换为字节数组
	ToBytes() []byte
}
