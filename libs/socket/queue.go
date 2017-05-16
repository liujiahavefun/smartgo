package socket

import (
	"runtime/debug"
	"time"
)

type EventQueue interface {
	StartLoop()

	StopLoop(result int)

	//等待退出
	Wait() int

	//投递事件, 通过队列到达消费者端
	Post(evd EventDispatcher, data interface{})

	//延时投递
	PostDelayed(evd EventDispatcher, dur time.Duration, data interface{})
}

type queueData struct {
	dispatcher  EventDispatcher
	data interface{}
}

type eventQueue struct {
	queue chan queueData

	exitSignal chan int

	capturePanic bool
}

func NewEventQueue() EventQueue {
	self := &eventQueue{
		queue:      make(chan queueData, 10000),
		exitSignal: make(chan int),
	}

	return self
}

func (self *eventQueue) StartLoop() {
	go func() {
		for v := range self.queue {
			self.protectedCall(v.dispatcher, v.data)
		}
	}()
}

func (self *eventQueue) StopLoop(result int) {
	self.exitSignal <- result
}

func (self *eventQueue) Wait() int {
	return <-self.exitSignal
}

func (self *eventQueue) Post(evd EventDispatcher, data interface{}) {
	self.queue <- queueData{dispatcher: evd, data: data}
}

func (self *eventQueue) PostDelayed(evd EventDispatcher, delayed time.Duration, data interface{}) {
	go func() {
		time.AfterFunc(delayed, func() {
			self.Post(evd, data)
		})
	}()
}

func (self *eventQueue) protectedCall(evd EventDispatcher, data interface{}) {
	if self.capturePanic {
		defer func() {
			if err := recover(); err != nil {
				logFatalln(err)
				debug.PrintStack()
			}
		}()
	}

	if evd != nil {
		evd.CallData(data)
	} else if f, ok := data.(func()); ok {
		f()
	}
}
