package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	//"smartgo/libs/net"
	//"smartgo/libs/pool"
)

//liujia: 编译命令 go build -ldflags "-X main.BUILD_VERSION=0.0.1 -X main.BUILD_TIME=`date +%Y年%m月%d日-%H:%M:%S`"
//ldflags -X 用于linker指定package内为初始化的变量的值
var (
	BUILD_VERSION string
	BUILD_TIME    string
)

var configFile = flag.String("config_file", "broker_config.json", "input broker config file name")

func version() {
	if len(BUILD_VERSION) == 0 {
		BUILD_VERSION = "unknown"
	}
	if len(BUILD_TIME) == 0 {
		BUILD_TIME = "unknown"
	}
	fmt.Printf("broker version %s built on %s,  Copyright(c) 2016 liujia@smartgo.com \n", BUILD_VERSION, BUILD_TIME)
}

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "true")
}

func main() {
	version()

	flag.Parse()

	cfg := NewBrokerConfig(*configFile)
	err := cfg.LoadConfig()
	if err != nil {
		glog.Error(err.Error())
		return
	}
}
