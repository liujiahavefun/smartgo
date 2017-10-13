package socket

import (
	"time"
	"sync"
)

type Timer struct {
	tick 		*time.Ticker
	donechan 	chan struct{}
	done        sync.Once
}

func (self *Timer) Stop() {
	self.done.Do(func() {
		self.donechan <- struct{}{}
	})
}

func NewTimer(eq EventQueue, dur time.Duration, callback func(*Timer)) *Timer {
	return NewTimer2(nil, eq, dur, callback)
}

/*
 * 其实有没有evd这个参数区别不大，都是在eq的派遣线程上下文被执行
 */
func NewTimer2(evd EventDispatcher, eq EventQueue, dur time.Duration, callback func(*Timer)) *Timer {
	self := &Timer{
		tick: time.NewTicker(dur),
		donechan: make(chan struct{}),
	}

	go func() {
		defer self.tick.Stop()
		for {
			select {
			case <-self.tick.C:
				eq.Post(evd, func() {
					callback(self)
				})
			case <-self.donechan:
				return
			}
		}
	}()

	return self
}