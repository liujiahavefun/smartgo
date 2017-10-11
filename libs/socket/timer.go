package socket

import (
	"time"
)

type Timer struct {
	tick 		*time.Ticker
	donechan 	chan struct{}
	done        *AtomicBoolean
}

func (self *Timer) Stop() {
	if self.done.CompareAndSet(false, true) == false {
		self.donechan <- struct{}{}
	}
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
		done: NewAtomicBoolean(false),
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