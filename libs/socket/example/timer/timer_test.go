package timer

import (
    "testing"
    "time"
    "fmt"

    //sessproto "smartgo/proto/sessevent"
    //testproto "smartgo/proto/test"
    "smartgo/libs/socket"
    "smartgo/libs/socket/example"
)

func TestTimer(t *testing.T) {
    signal := test.NewSignalTester(t)

    queue := socket.NewEventQueue()
    queue.StartLoop()

    const testTimes = 3
    var count int = testTimes

    socket.NewTimer(queue, time.Second, func(t *socket.Timer) {
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
    signal := test.NewSignalTester(t)

    queue := socket.NewEventQueue()
    queue.StartLoop()

    fmt.Println("delay 1 sec begin")

    queue.PostDelayed(nil, time.Second, func() {
        fmt.Println("delay done")
        signal.Done(1)
    })

    signal.WaitAndExpect(1, "delay not work")
}
