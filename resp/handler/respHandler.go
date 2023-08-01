package handler

import (
	"context"
	database2 "github.com/blkcor/go-redis/database"
	"github.com/blkcor/go-redis/interface/database"
	"github.com/blkcor/go-redis/lib/logger"
	"github.com/blkcor/go-redis/lib/sync/atomic"
	"github.com/blkcor/go-redis/resp/connection"
	"github.com/blkcor/go-redis/resp/parser"
	"github.com/blkcor/go-redis/resp/reply"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	UnknownErrReplyBytes = reply.MakeStandardErrorReply("-Err unknown\r\n").ToBytes()
)

type RespHandler struct {
	activeConn sync.Map
	db         database.Database
	closing    atomic.Boolean
}

func MakeHandler() *RespHandler {
	var db database.Database
	db = database2.NewEchoDatabase()
	return &RespHandler{
		db: db,
	}
}

/*
CloseClient 关闭单个客户端连接
*/
func (r *RespHandler) CloseClient(conn *connection.Connection) {
	_ = conn.Close()
	r.db.AfterClientClose(conn)
	r.activeConn.Delete(conn)
}

/*
Close 关闭所有客户端连接
*/
func (r *RespHandler) Close() error {
	logger.Info("shutting down......")
	r.closing.Set(true)
	r.activeConn.Range(func(key, value interface{}) bool {
		conn := value.(*connection.Connection)
		_ = conn.Close()
		return true
	})
	r.db.Close()
	return nil
}

/*
Handle 处理客户端连接
*/
func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		_ = conn.Close()
	}
	client := connection.NewConnection(conn)
	r.activeConn.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
	for payload := range ch {
		//error
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				r.CloseClient(client)
				logger.Errorf("closed connection from %s", client.RemoteAddr().String())
				return
			} else {
				//failed when writing to client
				errorReply := reply.MakeStandardErrorReply(payload.Err.Error())
				err := client.Write(errorReply.ToBytes())
				if err != nil {
					r.CloseClient(client)
					logger.Errorf("closed connection from %s", client.RemoteAddr().String())
					return
				}
			}
		}
		//exec
		if payload.Data == nil {
			continue
		}
		data, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Info("require a multi bulk reply")
			continue
		}
		result := r.db.Exec(client, data.Args)
		if result == nil {
			_ = client.Write(result.ToBytes())
		} else {
			_ = client.Write(UnknownErrReplyBytes)
		}
	}
}
