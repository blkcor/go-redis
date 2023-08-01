package database

import "github.com/blkcor/go-redis/interface/resp"

type CmdLine [][]byte

type Database interface {
	Exec(conn resp.Connection, args CmdLine) resp.Reply
	Close()
	AfterClientClose(conn resp.Connection)
}

type DataEntity struct {
	Data interface{}
}
