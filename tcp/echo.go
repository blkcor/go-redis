package tcp

import (
	"bufio"
	"context"
	"github.com/blkcor/go-redis/lib/logger"
	"github.com/blkcor/go-redis/lib/sync/atomic"
	"github.com/blkcor/go-redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (e *EchoClient) Close() error {
	e.Waiting.WaitWithTimeout(10 * time.Second)
	_ = e.Conn.Close()
	return nil
}

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	//if we are closing, we close the connection
	if handler.closing.Get() {
		_ = conn.Close()
		return
	}
	client := &EchoClient{
		Conn: conn,
	}

	//we add the client to the active connections
	handler.activeConn.Store(client, struct{}{})
	//now, we could read and write to the connection
	reader := bufio.NewReader(conn)
	for {
		//read the line
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("Connection closed")
				//we close the connection and remove the client from the active connections
				handler.activeConn.Delete(client)
			}
			logger.Warn("Error reading from connection: %s", err.Error())
			return
		}
		client.Waiting.Add(1)
		//write the message back
		_, err = conn.Write([]byte(msg))
		client.Waiting.Done()
	}
}

func (handler *EchoHandler) Close() error {
	logger.Info("Closing echo handler")
	handler.closing.Set(true)
	handler.activeConn.Range(func(key, value interface{}) bool {
		//go through all the active connections and close them
		client := key.(*EchoClient)
		_ = client.Conn.Close()
		return true
	})
	return nil
}
