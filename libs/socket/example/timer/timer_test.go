package timer

import (
    "testing"
    "time"
    "fmt"

    //sessproto "smartgo/proto/sessevent"
    //testproto "smartgo/proto/test"
    "smartgo/libs/utils"
    "smartgo/libs/socket"
    "smartgo/libs/socket/example"
)

func TestTimer(t *testing.T) {
    fmt.Printf("TestTimer, go id: %v\n", utils.GoID())
    signal := test.NewSignalTester(t)

    queue := socket.NewEventQueue()
    server := socket.NewTcpServer(queue).Start("127.0.0.1:7201")
    queue.StartLoop()

    const testTimes = 5
    var count int = testTimes

    socket.NewTimer2(server, queue, time.Second, func(t *socket.Timer) {
        fmt.Printf("TestTimer callback, go id: %v\n", utils.GoID())
        fmt.Println("timer 1 sec tick")
        signal.Done(1)

        count--
        if count == 0 {
            t.Stop()
            signal.Done(2)
        }
    })

    for i := 0; i < testTimes; i++ {
        signal.WaitAndExpect(1, "timer not tick")
    }

    signal.WaitAndExpect(2, "timer not stop")
}

func TestDelay(t *testing.T) {
    fmt.Printf("TestDelay, go id: %v\n", utils.GoID())
    signal := test.NewSignalTester(t)

    queue := socket.NewEventQueue()
    server := socket.NewTcpServer(queue).Start("127.0.0.1:7201")
    queue.StartLoop()

    fmt.Println("delay 1 sec begin")

    queue.PostDelayed(server, time.Second, func() {
        fmt.Printf("TestDelay callback, go id: %v\n", utils.GoID())
        fmt.Println("delay done")
        signal.Done(1)
    })

    signal.WaitAndExpect(1, "delay not work")
}
