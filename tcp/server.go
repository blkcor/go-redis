package tcp

import (
	"context"
	"github.com/blkcor/go-redis/interface/tcp"
	"github.com/blkcor/go-redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(config *Config, handler tcp.Handler) error {
	listener, err := net.Listen("tcp", config.Address)
	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	if err != nil {
		return err
	}
	logger.Infof("tcp server listening on %s", config.Address)
	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	//we need to close the listener and the handler when the server is closed
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	//if we receive a close signal, we close the listener
	go func() {
		<-closeChan
		logger.Info("shutting down")
		_ = listener.Close()
		_ = handler.Close()
	}()

	//get the context
	ctx := context.Background()
	waitDone := sync.WaitGroup{}

	//accept loop
	for true {
		//block accept
		conn, err := listener.Accept()
		if err != nil {
			logger.Errorf("accept error %s", err.Error())
			break
		}
		logger.Infof("accept %s", conn.RemoteAddr().String())
		waitDone.Add(1)
		//start a new go routine to handle the connection
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	//wait for all the connection to be closed
	waitDone.Wait()
}
