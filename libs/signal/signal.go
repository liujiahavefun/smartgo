package signal

import (
	"github.com/golang/glog"
	"os"
	"os/signal"
	"syscall"
)

type SignalHanlder func()

//register signals handler.
func RegsiterSignalHandler(handler SignalHanlder) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		glog.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			//process to quit, do not care
			return
		case syscall.SIGHUP:
			//most times to reload config
			handler()
		default:
			return
		}
	}
}
