package main

import (
	"fmt"
	"github.com/blkcor/go-redis/config"
	"github.com/blkcor/go-redis/lib/logger"
	"github.com/blkcor/go-redis/tcp"
	"os"
)

const configFile string = "redis.conf"

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

func fileExists(fileName string) bool {
	file, err := os.Stat(fileName)
	return err == nil && !file.IsDir()
}

func main() {
	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}

	tcpConfig := tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}

	err := tcp.ListenAndServeWithSignal(&tcpConfig, tcp.MakeHandler())
	if err != nil {
		logger.Error(err)
	}
}
