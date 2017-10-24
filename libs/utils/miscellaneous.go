package utils

import (
    "time"
    "syscall"
    "fmt"
)

func EnableFileDescriptor(limit uint64) (err error) {
    var rlim syscall.Rlimit

    /*
    err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
    if err != nil {
        fmt.Println("get rlimit error: " + err.Error())
        os.Exit(1)
    }
    */

    rlim.Cur = limit
    rlim.Max = limit

    return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlim)
}

func CurrentTimeMillSecond() int64 {
    return time.Now().UnixNano() / 1000000
}

func CurrentTimeSecond() int64 {
    return time.Now().Unix()
}

func Time2String(ts int64) string  {
    var secs int64 = ts/1000
    var left = ts%1000
    var s string = time.Unix(secs, 0).Format("2006-01-02 15:04:05")
    return s+fmt.Sprintf(":%v", left)
}