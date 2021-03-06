package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	//"smartgo/libs/proto/thrift/gen-go/test/demo"
	"smartgo/libs/proto/thrift/gen-go/test2/rpc"
	//"github.com/apache/thrift/lib/go/thrift"
	"crypto/tls"
	"flag"
	"os"
	"time"
)

const (
	HOST = "127.0.0.1"
	PORT = "9090"
)

/*
func main() {
	startTime := currentTimeMillis()

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(net.JoinHostPort(HOST, PORT))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}

	useTransport := transportFactory.GetTransport(transport)
	client := demo.NewTest2ThriftClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to "+HOST+":"+PORT, " ", err)
		os.Exit(1)
	}
	defer transport.Close()

	for i := 0; i < 10; i++ {
		paramMap := make(map[string]string)
		paramMap["a"] = "batu.demo"
		paramMap["b"] = "test" + strconv.Itoa(i+1)
		r1, _ := client.CallBack(time.Now().Unix(), "go client", paramMap)
		fmt.Println("GOClient Call->", r1)
	}

	model := demo.Article{1, "Go第一篇文章", "我在这里", "liuxinming"}
	client.Put(&model)
	endTime := currentTimeMillis()
	fmt.Printf("本次调用用时:%d-%d=%d毫秒\n", endTime, startTime, (endTime - startTime))
}
*/

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
	//server := flag.Bool("server", false, "Run server")
	protocol := flag.String("P", "binary", "Specify the protocol (binary, compact, json, simplejson)")
	framed := flag.Bool("framed", false, "Use framed transport")
	buffered := flag.Bool("buffered", false, "Use buffered transport")
	addr := flag.String("addr", "localhost:9090", "Address to listen to")
	secure := flag.Bool("secure", false, "Use tls secure transport")

	flag.Parse()

	/*
		这里指定传输时序列化的协议。
		TBinaryProtocol	二进制格式
		TCompactProtocol	压缩格式
		TJSONProtocol	JSON格式
		TSimpleJSONProtocol	提供JSON只写协议, 生成的文件很容易通过脚本语言解析。
		TDebugProtocol	使用易懂的可读的文本格式，以便于debug
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
		TSocket				阻塞式socker
		TFileTransport		以文件形式进行传输。---通常用于日志上传？
		THttpTransport		采用Http传输协议进行数据传输
		TZlibTransport		使用zlib进行压缩， 与其他传输方式联合使用。当前无java实现。

		TMemoryTransport	将内存用于I/O. java实现时内部实际使用了简单的ByteArrayOutputStream。
		TFramedTransport	以frame为单位进行传输，非阻塞式服务中使用。以frame为单位进行传输，非阻塞式服务中使用。
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

func handleClient(client *rpc.CalculatorClient) (err error) {
	err = client.Ping()
	fmt.Println("ping(): ", err)

	sum, _ := client.Add(1, 1)
	//fmt.Print("1+1=", sum, "\n")
	if err != nil {
		fmt.Println("Unable to Add:", err)
	} else {
		fmt.Println("Add return: ", sum)
	}

	work := rpc.NewWork()
	work.Op = rpc.Operation_DIVIDE
	work.Num1 = 1
	work.Num2 = 0
	quotient, err := client.Calculate(1, work)
	if err != nil {
		switch v := err.(type) {
		case *rpc.InvalidOperation:
			fmt.Println("Invalid operation:", v)
		default:
			fmt.Println("Error during operation:", err)
		}
		//return err
	} else {
		fmt.Println("Whoa we can divide by 0 with new value:", quotient)
	}

	work.Op = rpc.Operation_SUBTRACT
	work.Num1 = 15
	work.Num2 = 10
	diff, err := client.Calculate(1, work)
	if err != nil {
		switch v := err.(type) {
		case *rpc.InvalidOperation:
			fmt.Println("Invalid operation:", v)
		default:
			fmt.Println("Error during operation:", err)
		}
		//return err
	} else {
		fmt.Print("15-10=", diff, "\n")
	}

	err = client.Zip()
	if err != nil {
		fmt.Println("Unable to Zip():", err)
	} else {
		fmt.Println("Zip done")
	}

	// derived interface from shared.thrift
	log, err := client.GetStruct(1)
	if err != nil {
		fmt.Println("Unable to get struct:", err)
	} else {
		fmt.Println("Check log return: ", log)
	}

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
	return handleClient(rpc.NewCalculatorClientFactory(transport, protocolFactory))
}
