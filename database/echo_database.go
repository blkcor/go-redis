package database

import (
	"github.com/blkcor/go-redis/interface/database"
	"github.com/blkcor/go-redis/interface/resp"
	"github.com/blkcor/go-redis/resp/reply"
)

type EchoDatabase struct {
}

func (e EchoDatabase) Exec(conn resp.Connection, args database.CmdLine) resp.Reply {
	return reply.MakeMultiBulkReply(args)
}

func (e EchoDatabase) Close() {

}

func (e EchoDatabase) AfterClientClose(conn resp.Connection) {

}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}
