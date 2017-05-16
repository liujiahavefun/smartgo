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

func logDebugln(args ...interface{}) {
    gLog.Debug(args...)
}

func logDebugf(format string, args ...interface{}) {
    gLog.Debugf(format, args...)
}

func logInfoln(args ...interface{}) {
    gLog.Info(args...)
}

func logInfof(format string, args ...interface{}) {
    gLog.Infof(format, args...)
}

func logWarningln(args ...interface{}) {
    gLog.Warning(args...)
}

func logWarningf(format string, args ...interface{}) {
    gLog.Warningf(format, args...)
}

func logErrorln(args ...interface{}) {
    gLog.Error(args...)
}

func logErrorf(format string, args ...interface{}) {
    gLog.Errorf(format, args...)
}

func logFatalln(args ...interface{}) {
    gLog.Fatal(args...)
}

func logFatalf(format string, args ...interface{}) {
    gLog.Fatalf(format, args...)
}
