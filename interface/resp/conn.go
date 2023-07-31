package resp

type Connection interface {
	//写入数据
	Write([]byte) error
	//得到数据库的索引 即当前使用的数据库
	GetDBIndex() int
	//选择数据库
	SelectDB(int)
}
