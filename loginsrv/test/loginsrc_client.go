package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"os"
	"smartgo/loginsrv/proto/thrift/gen-go/login/rpc"
	"time"
)

const (
	LOGINSRV_RPC_ADDR = "123.56.88.196:9100"
)

func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

func Usage() {
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
}

func main() {
	flag.Usage = Usage

	protocol := flag.String("protocol", "binary", "Specify the protocol (binary, compact, json, simplejson)")
	framed := flag.Bool("framed", true, "Use framed transport")
	buffered := flag.Bool("buffered", true, "Use buffered transport")
	addr := flag.String("addr", LOGINSRV_RPC_ADDR, "Address to listen to")
	secure := flag.Bool("secure", false, "Use tls secure transport")

	flag.Parse()

	fmt.Printf("protocol : %v \n", *protocol)
	fmt.Printf("framed : %v \n", *framed)
	fmt.Printf("buffered : %v \n", *buffered)
	fmt.Printf("addr : %v \n", *addr)
	fmt.Printf("secure : %v \n", *secure)

	/*
	   这里指定传输时序列化的协议。
	   TBinaryProtocol 二进制格式
	   TCompactProtocol    压缩格式
	   TJSONProtocol   JSON格式
	   TSimpleJSONProtocol 提供JSON只写协议, 生成的文件很容易通过脚本语言解析。
	   TDebugProtocol  使用易懂的可读的文本格式，以便于debug
	*/
	var protocolFactory thrift.TProtocolFactory
	switch *protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		fmt.Fprint(os.Stderr, "Invalid protocol specified", protocol, "\n")
		Usage()
		os.Exit(1)
	}

	/*
	   数据传输协议 感觉File和Memory适合本机RPC？
	   直觉Framed靠谱，因为服务端使用NonBlocking模式的话，客户端必须用Framed，这样才能到更高性能
	   TSocket             阻塞式socker
	   TFileTransport      以文件形式进行传输。---通常用于日志上传？
	   THttpTransport      采用Http传输协议进行数据传输
	   TZlibTransport      使用zlib进行压缩， 与其他传输方式联合使用。当前无java实现。

	   TMemoryTransport    将内存用于I/O. java实现时内部实际使用了简单的ByteArrayOutputStream。
	   TFramedTransport    以frame为单位进行传输，非阻塞式服务中使用。以frame为单位进行传输，非阻塞式服务中使用。
	                       同TBufferedTransport类似，也会对相关数据进行buffer，同时，它支持定长数据发送和接收。
	   TBufferedTransport：对某个Transport对象操作的数据进行buffer，即从buffer中读取数据进行传输，或者将数据直接写入buffer
	*/
	var transportFactory thrift.TTransportFactory
	if *buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if *framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	if err := runClient(transportFactory, protocolFactory, *addr, *secure); err != nil {
		fmt.Println("error running client:", err)
	}
}

func handleClient(client *rpc.LoginServiceClient) (err error) {
	token, err := client.LoginByPasswd("liujia", "123456", "pc", nil)
	fmt.Println("token: ", token, " err: ", err)
	return err
}

func runClient(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string, secure bool) error {
	var transport thrift.TTransport
	var err error

	if secure {
		cfg := new(tls.Config)
		cfg.InsecureSkipVerify = true
		transport, err = thrift.NewTSSLSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTSocket(addr) //NewTSocketTimeout?
	}

	if err != nil {
		fmt.Println("Error opening socket:", err)
		return err
	}

	transport = transportFactory.GetTransport(transport)
	defer transport.Close()
	if err := transport.Open(); err != nil {
		return err
	}

	return handleClient(rpc.NewLoginServiceClientFactory(transport, protocolFactory))
}
