package socket

import (
	"time"
)

type Timer struct {
	tick *time.Ticker
	done chan struct{}
}

func (self *Timer) Stop() {
	self.done <- struct{}{}
}

func NewTimer(eq EventQueue, dur time.Duration, callback func(*Timer)) *Timer {
	self := &Timer{
		tick: time.NewTicker(dur),
		done: make(chan struct{}),
	}

	go func() {
		defer self.tick.Stop()
		for {
			select {
			case <-self.tick.C:
				eq.Post(nil, func() {
					callback(self)
				})
			case <-self.done:
				return
			}
		}
	}()

	return self
}

/*
 * 其实有没有evd这个参数区别不大，都是在eq的派遣线程上下文被执行
 */
func NewTimer2(evd EventDispatcher, eq EventQueue, dur time.Duration, callback func(*Timer)) *Timer {
	self := &Timer{
		tick: time.NewTicker(dur),
		done: make(chan struct{}),
	}

	go func() {
		defer self.tick.Stop()
		for {
			select {
			case <-self.tick.C:
				eq.Post(evd, func() {
					callback(self)
				})
			case <-self.done:
				return
			}
		}
	}()

	return self
}