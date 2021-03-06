package main

import (
    "testing"
    "time"
    "sync/atomic"
    "syscall"
    "fmt"
    "os"

    sessproto "smartgo/proto/session_event"
    testproto "smartgo/proto/test_event"
    "smartgo/libs/socket"
    "smartgo/libs/socket/example"
)

var signal *test.SignalTester

// 测试地址
const benchmarkAddress = "127.0.0.1:6789"

// 客户端并发数量
const clientCount = 200

// 测试时间(秒)
const benchmarkSeconds = 20

var(
    acceptedServerSession int64
    connectedServerSession int64
    connectedClientCount int64
)

func server() {
    queue := socket.NewEventQueue()
    qpstester := test.NewQPSTester(queue, func(qps int) {
        fmt.Printf("QPS: %d, Accepted Client: %d, Connected Client: %d \n", qps, acceptedServerSession, connectedServerSession)
    })

    server := socket.NewTcpServer(queue).Start(benchmarkAddress)
    socket.RegisterMessage(server, "test_event.TestEchoACK", func(content interface{}, ses socket.Session) {
        if qpstester.Acc() > benchmarkSeconds {
            signal.Done(1)
            fmt.Printf("Average QPS: %d, Accepted Client: %d, Connected Client: %d \n", qpstester.Average(), acceptedServerSession, connectedServerSession)
        }

        ses.Send(&testproto.TestEchoACK{})
    })

    socket.RegisterMessage(server, "session_event.SessionAccepted", func(content interface{}, ses socket.Session) {
        atomic.AddInt64(&acceptedServerSession, 1)
    })

    socket.RegisterMessage(server, "session_event.SessionAcceptFailed", func(content interface{}, ses socket.Session) {
        msg := content.(*sessproto.SessionAcceptFailed)
        fmt.Println("SessionAcceptFailed, err: ", msg.Reason)
    })

    socket.RegisterMessage(server, "session_event.SessionConnected", func(content interface{}, ses socket.Session) {
        atomic.AddInt64(&connectedServerSession, 1)
    })

    socket.RegisterMessage(server, "session_event.SessionConnectFailed", func(content interface{}, ses socket.Session) {
        fmt.Println("SessionConnectFailed")
    })

    socket.RegisterMessage(server, "sessevent.SessionError", func(content interface{}, ses socket.Session) {
        msg := content.(sessproto.SessionError)
        fmt.Println("SessionError: ", msg.Reason)
    })

    queue.StartLoop()
}

func client() {
    queue := socket.NewEventQueue()
    connector := socket.NewConnector(queue)

    socket.RegisterMessage(connector, "session_event.SessionConnected", func(content interface{}, ses socket.Session) {
        atomic.AddInt64(&connectedClientCount, 1)
        ses.Send(&testproto.TestEchoACK{})
    })

    socket.RegisterMessage(connector, "session_event.SessionError", func(content interface{}, ses socket.Session) {
        msg := content.(*sessproto.SessionError)
        fmt.Println("session error:", msg.Reason)
    })

    socket.RegisterMessage(connector, "test_event.TestEchoACK", func(content interface{}, ses socket.Session) {
        ses.Send(&testproto.TestEchoACK{})
    })

    connector.Start(benchmarkAddress)

    queue.StartLoop()
}

func sessprotoRegisterMessage() {
    fmt.Println("sessevent RegisterMessage")
    // session.proto
    socket.RegisterMessageMeta("session_event.SessionAccepted", (*sessproto.SessionAccepted)(nil), 348117910)
    socket.RegisterMessageMeta("session_event.SessionAcceptFailed", (*sessproto.SessionAcceptFailed)(nil), 1978788392)
    socket.RegisterMessageMeta("session_event.SessionConnected", (*sessproto.SessionConnected)(nil), 3543838007)
    socket.RegisterMessageMeta("session_event.SessionConnectFailed", (*sessproto.SessionConnectFailed)(nil), 1720533237)
    socket.RegisterMessageMeta("session_event.SessionClosed", (*sessproto.SessionClosed)(nil), 90181607)
    socket.RegisterMessageMeta("session_event.SessionError", (*sessproto.SessionError)(nil), 1937281175)
}

func testprotoRegisterMessage() {
    fmt.Println("test_event RegisterMessage")
    // test.proto
    socket.RegisterMessageMeta("test_event.TestEchoACK", (*testproto.TestEchoACK)(nil), 509149489)
}

func TestIO(t *testing.T) {
    EnableManyFiles()

    sessprotoRegisterMessage()
    testprotoRegisterMessage()
    socket.Init()

    signal = test.NewSignalTester(t)

    // 超时时间为测试时间延迟一会
    signal.SetTimeout((benchmarkSeconds + 5) * time.Second)

    fmt.Println("start server")

    server()

    fmt.Println("start all clients")
    for i := 0; i < clientCount; i++ {
        time.Sleep(50*time.Millisecond)
        go client()
    }

    fmt.Println("all clients started")

    signal.WaitAndExpect(1, "recv time out")
}

func EnableManyFiles() {
    var rlim syscall.Rlimit
    rlim.Cur = 50000
    rlim.Max = 50000

    err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlim)
    if err != nil {
        fmt.Println("set rlimit error: " + err.Error())
        os.Exit(1)
    }

    err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
    if err != nil {
        fmt.Println("get rlimit error: " + err.Error())
        os.Exit(1)
    }

    fmt.Println("rlim.Curr", rlim.Cur)
    fmt.Println("rlim.Max", rlim.Max)
}