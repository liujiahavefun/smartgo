package test

import (
    "sync"
    "time"

    "smartgo/libs/socket"
)

type QPSTester struct {
    qpsGuard sync.Mutex
    qps      int
    total    int
    count int
}

func (self *QPSTester) Acc() int {
    self.qpsGuard.Lock()
    defer self.qpsGuard.Unlock()

    self.qps++
    return self.count
}

//一轮计算
func (self *QPSTester) Turn() (ret int) {
    self.qpsGuard.Lock()
    defer self.qpsGuard.Unlock()

    if self.qps > 0 {
        ret = self.qps
    }

    self.total += self.qps
    self.qps = 0
    self.count++

    return
}

// 均值
func (self *QPSTester) Average() int {
    self.qpsGuard.Lock()
    defer self.qpsGuard.Unlock()

    if self.count == 0 {
        return 0
    }

    return self.total / self.count
}

func NewQPSTester(evq socket.EventQueue, callback func(int)) *QPSTester {
    self := &QPSTester{}
    socket.NewTimer(evq, time.Second, func(t *socket.Timer) {
        qps := self.Turn()
        callback(qps)
    })

    return self
}
