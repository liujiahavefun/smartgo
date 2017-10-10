package test

import (
    "testing"
    "time"
)

type SignalTester struct {
    *testing.T
    signal chan int
    timeout time.Duration
}

func (self *SignalTester) SetTimeout(duration time.Duration) {
    self.timeout = duration
}

func (self *SignalTester) WaitAndExpect(value int, msg string) bool {
    select {
    case v := <-self.signal:
        if v != value {
            self.Fail()
            self.Logf("%s\n", msg)
            return false
        }

    case <-time.After(self.timeout):
        self.Logf("signal timeout: %d %s", value, msg)
        self.Fail()
        return false
    }

    return true
}

func (self *SignalTester) Done(value int) {
    self.signal <- value
}

func NewSignalTester(t *testing.T) *SignalTester {
    return &SignalTester{
        T:       t,
        timeout: 5 * time.Second,
        signal:  make(chan int),
    }
}

func NewSignalTesterTimeout(t *testing.T, secs time.Duration) *SignalTester {
    return &SignalTester{
        T:       t,
        timeout: secs * time.Second,
        signal:  make(chan int),
    }
}