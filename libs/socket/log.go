package socket

import (
    "fmt"
    "smartgo/libs/log"
)

const (
    LOG_CONFIG_FILE = "socket.conf"
)

var (
    gLog *logger.Logger
)

func initLogger() {
    conf := logger.NewLogConfig(LOG_CONFIG_FILE)
    err := conf.LoadConfig()
    if err != nil {
        fmt.Errorf("load log conf file failed: %v", err)
    }
    gLog = logger.NewLogger(conf)
    gLog.Infof("init logger %s done, %v", "ddd", conf)
}

func logInitialized() bool {
    return gLog != nil
}

func logDebugln(args ...interface{}) {
    if logInitialized() == false {
        return
    }
    gLog.Debug(args...)
}

func logDebugf(format string, args ...interface{}) {
    if logInitialized() == false {
        return
    }
    gLog.Debugf(format, args...)
}

func logInfoln(args ...interface{}) {
    if logInitialized() == false {
        return
    }
    gLog.Info(args...)
}

func logInfof(format string, args ...interface{}) {
    if logInitialized() == false {
        return
    }
    gLog.Infof(format, args...)
}

func logWarningln(args ...interface{}) {
    if logInitialized() == false {
        return
    }
    gLog.Warning(args...)
}

func logWarningf(format string, args ...interface{}) {
    if logInitialized() == false {
        return
    }
    gLog.Warningf(format, args...)
}

func logErrorln(args ...interface{}) {
    if logInitialized() == false {
        return
    }
    gLog.Error(args...)
}

func logErrorf(format string, args ...interface{}) {
    if logInitialized() == false {
        return
    }
    gLog.Errorf(format, args...)
}

func logFatalln(args ...interface{}) {
    if logInitialized() == false {
        return
    }
    gLog.Fatal(args...)
}

func logFatalf(format string, args ...interface{}) {
    if logInitialized() == false {
        return
    }
    gLog.Fatalf(format, args...)
}
