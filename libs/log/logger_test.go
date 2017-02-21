package logger_test

import (
	"fmt"
	//"os"
	//"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"smartgo/libs/log"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_GetLogFileAndNextTime(t *testing.T) {
	fn, next, err := logger.GetFileNameAndNextTime("perquarter", "smartgo")
	t.Log(fn, next, err)

	fn, next, err = logger.GetFileNameAndNextTime("halfhour", "smartgo")
	t.Log(fn, next, err)

	fn, next, err = logger.GetFileNameAndNextTime("perhour", "smartgo")
	t.Log(fn, next, err)

	fn, next, err = logger.GetFileNameAndNextTime("perday", "smartgo")
	t.Log(fn, next, err)

	fn, next, err = logger.GetFileNameAndNextTime("laoliu", "smartgo")
	t.Log(fn, next, err)

	fn, next, err = logger.GetFileNameAndNextTime("", "smartgo")
	t.Log(fn, next, err)

	/*
	   logPath := filepath.Join("/Users/liujia/go/src/ygo/src/libs/log", "ddd.log")
	   t.Log(logPath)

	   logPath = filepath.Join("/Users/liujia/go/src/ygo/src/libs/log/", "ddd.log")
	   t.Log(logPath)

	   logPath, err = logger.PathJoin("./log", "ddd.log")
	   t.Log(logPath)

	   logPath = filepath.Join("./log", "ddd.log")
	   t.Log(logPath)
	*/
}

func Test_Base(t *testing.T) {
	fmt.Println("1111!")
	conf := logger.NewLogConfig("log.conf")
	err := conf.LoadConfig()
	if err != nil {
		t.Errorf("load log conf file failed: %v", err)
	}

	fmt.Println(conf)

	fmt.Println("Starting!")

	logger := logger.NewLogger(conf)
	var wg sync.WaitGroup
	wg.Add(1)

	time.AfterFunc(5*time.Second, func() {
		fmt.Println("5 second in select ")
		logger.Fatal("liujia: ", time.Now())
		wg.Done()
		return
	})

	go func() {
		for {
			select {
			case <-time.After(200 * time.Millisecond):
				fmt.Println("200 msecond in select ")
				logger.Debugf("liujia: %v", time.Now())
			}
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				fmt.Println("100 msecond in select ")
				logger.Infof("liujia: %v", time.Now())
			}
		}
	}()

	fmt.Println("waiting!")
	wg.Wait()

	fmt.Println("finished!")
}
