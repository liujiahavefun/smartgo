/*
*
* This should be the package description/comment, fill this in future, please remember me!
*
 */

package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

const (
	ERROR_FILEPATH_SUFFIX = ".wef"
)

const defaultCallDepth int = 2

type Logger struct {
	rootPath       string   // desc: absolute path
	file           *os.File // desc: log file object
	fileError      *os.File // desc: log file object for warning/error/fatal
	level          int      // option: log level, DEBUG/INFO/WARNING/ERROR/FATAL
	depth          int      // default: 2
	splitpolicy    string   // default: perhour
	fileName       string   // current log file name, if splitError is set to true, error log file name is xxxx.wef
	fileNamePrefix string   // file name prefix add to log file name
	splitError     bool     // whether write warning/error/fatal to another log file
	detailInfo     bool     // print file,line,callstack in log
	nexttime       time.Time
	opChan         chan *cmd
	finishChan     chan *struct{}
}

func NewLogger(config *LogConfig) *Logger {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	l := &Logger{}
	l.depth = defaultCallDepth
	l.rootPath = config.LogRootDir
	l.fileNamePrefix = config.NamePrefix
	l.splitError = config.SplitError
	l.detailInfo = config.DetailInfo

	switch config.LogLevel {
	case "debug":
		l.level = DEBUG
	case "info":
		l.level = INFO
	case "warning":
		l.level = WARNING
	case "error":
		l.level = ERROR
	case "fatal":
		l.level = FATAL
	default:
		l.level = INFO
	}

	switch config.LogSplitPolicy {
	case "perminute":
		l.splitpolicy = "perminute"
	case "perquarter":
		l.splitpolicy = "perquarter"
	case "halfhour":
		l.splitpolicy = "halfhour"
	case "perhour":
		l.splitpolicy = "perhour"
	case "perday":
		l.splitpolicy = "perday"
	default:
		l.splitpolicy = "perhour"
	}

	err := l.getLogFile()
	if err != nil {
		panic(err)
	}

	l.opChan = make(chan *cmd, 1024)
	go l.handler()

	return l
}

func (l *Logger) SetCallDepth(depth int) {
	if depth > 0 {
		l.depth = depth
	}
}

func (l *Logger) Debug(args ...interface{}) {
	if DEBUG < l.level {
		return
	}

	l.writeLogFormat(DEBUG, fmt.Sprintf("%s", args))
}
func (l *Logger) Debugf(format string, args ...interface{}) {
	if DEBUG < l.level {
		return
	}

	l.writeLogFormat(DEBUG, fmt.Sprintf(format, args...))
}
func (l *Logger) Info(args ...interface{}) {
	if INFO < l.level {
		return
	}

	l.writeLogFormat(INFO, fmt.Sprintf("%s", args))
}
func (l *Logger) Infof(format string, args ...interface{}) {
	if INFO < l.level {
		return
	}

	l.writeLogFormat(INFO, fmt.Sprintf(format, args...))
}
func (l *Logger) Warning(args ...interface{}) {
	if WARNING < l.level {
		return
	}

	l.writeLogFormat(WARNING, fmt.Sprintf("%s", args))
}
func (l *Logger) Warningf(format string, args ...interface{}) {
	if WARNING < l.level {
		return
	}

	l.writeLogFormat(WARNING, fmt.Sprintf(format, args...))
}
func (l *Logger) Error(args ...interface{}) {
	if ERROR < l.level {
		return
	}

	l.writeLogFormat(ERROR, fmt.Sprintf("%s", args))
}
func (l *Logger) Errorf(format string, args ...interface{}) {
	if ERROR < l.level {
		return
	}

	l.writeLogFormat(ERROR, fmt.Sprintf(format, args...))
}
func (l *Logger) Fatal(args ...interface{}) {
	if FATAL < l.level {
		return
	}

	l.writeLogFormat(FATAL, fmt.Sprintf("%s", args))
	os.Exit(1)
}
func (l *Logger) Fatalf(format string, args ...interface{}) {
	if FATAL < l.level {
		return
	}

	l.writeLogFormat(FATAL, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *Logger) getLogFile() error {
	rootPath := l.rootPath
	exist, err := IsDirExist(rootPath)
	if err != nil {
		panic(err)
	}

	if exist == false {
		//liujia: 我看文档说，如果目录存在MkdirAll啥都不做，是不是可以不用判断目录存在，直接创建试试？
		err = os.MkdirAll(rootPath, os.ModeDir)
		if err != nil {
			panic(err)
		}

		//这里要修改一下权限，否则当前用户只能读和执行，不能写
		os.Chmod(rootPath, 0777)
	}

	fn, nexttime, err := GetFileNameAndNextTime(l.splitpolicy, l.fileNamePrefix)
	logPath := filepath.Join(l.rootPath, fn)
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to open log file \"%s\"", logPath))
	}

	l.fileName = logPath
	l.file = f
	l.nexttime = nexttime

	if l.splitError == true {
		l.fileError, err = os.OpenFile(logPath+ERROR_FILEPATH_SUFFIX, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to open error log file \"%s\"", logPath))
		}
	}

	return err
}

func (l *Logger) handler() {
	for {
		select {
		case <-l.finishChan:
			return
		case c := <-l.opChan:
			l.handleCmd(c)
		case <-time.After(1 * time.Second):
			l.handleCmd(&cmd{op: "flush"})
		}
	}
}

func (l *Logger) handleCmd(c *cmd) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("handle cmd failed [", c, "]", " err: ", err)
		}
	}()

	switch c.op {
	case "flush":
		l.sync()
	case "write":
		if time.Now().UnixNano() > l.nexttime.UnixNano() {
			err := l.reset()
			if err != nil {
				panic(err)
			}
		}

		err := l.write(c)
		if err != nil {
			panic(err)
		}
	}
}

func (l *Logger) sync() {
	l.file.Sync()
	if l.splitError == true {
		l.fileError.Sync()
	}
}

func (l *Logger) reset() error {
	l.file.Close()
	if l.splitError {
		l.fileError.Close()
	}
	return l.getLogFile()
}

func (l *Logger) write(c *cmd) (err error) {
	if l.splitError && c.level > INFO {
		_, err = l.fileError.WriteString(c.log)
	} else {
		_, err = l.file.WriteString(c.log)
	}

	return err
}

func (l *Logger) postCmd(c *cmd) {
	l.opChan <- c
}

func (l *Logger) postCmdDirectly(c *cmd) {
	l.handleCmd(c)
}

func (l *Logger) writeLogFormat(level int, log string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	time := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05.99999")

	var flag string

	switch level {
	case DEBUG:
		flag = "Debug"
	case INFO:
		flag = "Info"
	case WARNING:
		flag = "Warning"
	case ERROR:
		flag = "Error"
	case FATAL:
		flag = "Fatal"
	}

	pc, file, line, ok := runtime.Caller(l.depth)
	if ok == false {
		panic(errors.New("failed to get call stack"))
	}

	f := runtime.FuncForPC(pc)

	//l.opChan <- &cmd{op: "write", log: fmt.Sprintf("%s [%s] [%s:%d:%s] %s\n", time, flag, file, line, f.Name(), log)}

	var msg string
	if l.detailInfo {
		msg = fmt.Sprintf("%s [%s] [%s:%d] [%s] [%s] %s\n", time, flag, file, line, f.Name(), callstack(), log)
	} else {
		msg = fmt.Sprintf("%s [%s] %s\n", time, flag, log)
	}

	c := &cmd{level: level, op: "write", log: msg}
	if level != FATAL {
		l.postCmd(c)
	} else {
		l.postCmdDirectly(c)
	}
}
