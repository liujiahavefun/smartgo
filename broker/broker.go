package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
)

//liujia: 编译命令 go build -ldflags "-X main.BUILD_VERSION=0.0.1 -X main.BUILD_TIME=`date +%Y年%m月%d日-%H:%M:%S`"
//ldflags -X 用于linker指定package内为初始化的变量的值
var (
	BUILD_VERSION string
	BUILD_TIME    string
)

var configFile = flag.String("config_file", "broker.conf", "input broker config file name")

func version() {
	if len(BUILD_VERSION) == 0 {
		BUILD_VERSION = "unknown"
	}
	if len(BUILD_TIME) == 0 {
		BUILD_TIME = "unknown"
	}
	fmt.Printf("broker version %s built on %s,  Copyright(c) 2016 liujia@smartgo.com\n", BUILD_VERSION, BUILD_TIME)
}

func main() {
	version()

	flag.Parse()

	//init logger
	fmt.Println("init logger")
	initLogger()

	//load config
	cfg := NewBrokerConfig(*configFile)
	err := cfg.LoadConfig()
	if err != nil {
		logErrorf(err.Error())
		panic("failed to load config")
	}

	//init SessionMgr
	gSessionMgr = NewSessionMgr()

	//start server
	go func() {
		logInfo("to start broker server", cfg.ListenOn)
		start(cfg.ListenOn)
	}()

	//进程收到的退出信号
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	select {
	case <-signals:
		logInfo("Received OS signal, stop broker server")
	}

	logInfo("exit broker process")
}
