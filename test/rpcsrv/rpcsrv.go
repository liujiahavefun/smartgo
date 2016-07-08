package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	//"github.com/apache/thrift/lib/go/thrift"
	"crypto/tls"
	"flag"
	"os"
	"smartgo/libs/proto/thrift/gen-go/test2/rpc"
)

/*
const (
	NetworkAddr = "127.0.0.1:9101"
)

type RpcServiceImpl struct {
}

func (this *RpcServiceImpl) FunCall(callTime int64, funCode string, paramMap map[string]string) (r []string, err error) {
	fmt.Println("-->FunCall:", callTime, funCode, paramMap)

	for k, v := range paramMap {
		r = append(r, k+v)
	}
	return
}

func main() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	//protocolFactory := thrift.NewTCompactProtocolFactory()

	serverTransport, err := thrift.NewTServerSocket(NetworkAddr)
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}

	handler := &RpcServiceImpl{}
	processor := rpc.NewRpcServiceProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("thrift server in", NetworkAddr)
	server.Serve()
}
*/

func Usage() {
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
}

func main() {
	flag.Usage = Usage
	//server := flag.Bool("server", false, "Run server")
	protocol := flag.String("P", "binary", "Specify the protocol (binary, compact, json, simplejson)")
	framed := flag.Bool("framed", false, "Use framed transport")
	buffered := flag.Bool("buffered", false, "Use buffered transport")
	addr := flag.String("addr", "localhost:9090", "Address to listen to")
	secure := flag.Bool("secure", false, "Use tls secure transport")

	flag.Parse()

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

	if err := runServer(transportFactory, protocolFactory, *addr, *secure); err != nil {
		fmt.Println("error running server:", err)
	}
}

func runServer(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string, secure bool) error {
	var transport thrift.TServerTransport
	var err error
	if secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		transport, err = thrift.NewTSSLServerSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(addr)
	}

	if err != nil {
		return err
	}
	fmt.Printf("%T\n", transport)
	handler := NewCalculatorHandler()
	processor := rpc.NewCalculatorProcessor(handler)

	/*
		go里只有SimpleServer一种模式，已经是异步非阻塞的，直接用这个就好了
		TSimpleServer – 简单的单线程服务模型，常用于测试
		TThreadPoolServer – 多线程服务模型，使用标准的阻塞式IO。
		TNonblockingServer – 多线程服务模型，使用非阻塞式IO（需使用TFramedTransport数据传输方式）
	*/
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting the simple server... on ", addr)
	return server.Serve()
}
