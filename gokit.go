package logkit

import (
	"fmt"
	"io"
	"runtime"
	"time"
)

var (
	inited       bool
	logWriter    Writer
	logLevel     = LevelInfo
	logLevelName string
	logName      string
	logPath      string
	channel      Channel
	alsoStdout   bool
	withCaller   bool

	levelToNames = map[Level]string{
		LevelFatal: "FATAL",
		LevelError: "ERROR",
		LevelWarn:  "WARN",
		LevelInfo:  "INFO",
		LevelDebug: "DEBUG",
		LevelTrace: "TRACE",
	}
)

//Level 日志等级
type Level int

const (
	Default Level = iota
	LevelTrace
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type Channel byte

const (
	FIlE Channel = iota
	SYSLOG
	KAFKA
)

type Writer interface {
	//Write 写日志
	Write(Level, string)
	//Exit 日志退出
	Exit()
}

func Exit() {
	logWriter.Exit()

}

func Init(_channel Channel, name string, level Level, _alsoStdout bool, _withCaller bool) error {
	if inited {
		return fmt.Errorf("logkit has been inited")
	}

	if name != "" {
		logName = name
	} else {
		return fmt.Errorf("log name must not be empty")
	}

	logLevel = level
	channel = _channel
	alsoStdout = _alsoStdout
	withCaller = _withCaller
	return nil
}

func SetPath(path string) {
	logPath = path
}
func getLevelName(level Level) string {
	levelName, _ := levelToNames[level]
	return levelName
}

func format(level Level, msg string) string {
	if withCaller {
		_, file, line, _ := runtime.Caller(3)
		return fmt.Sprintf("%s [%s] %s:%d %s \n", time.Now().Format("2006-01-02 15:04:05.999"), getLevelName(level), file, line, msg)
	} else {
		return fmt.Sprintf("%s [%s] %s \n", time.Now().Format("2006-01-02 15:04:05.999"), getLevelName(level), msg)
	}
}

func write(level Level, msg string) {
	if !inited {
		if logWriter == nil && channel == FIlE {
			if logPath == "" {
				logPath = "/data/logs/" + logName + ".log"
			}
			logWriter = NewFileLogger(logPath, logName, time.Second*5, 1204*1024*1800, 256*1024)
		}
		if logWriter == nil && channel == SYSLOG {
			logWriter, _ = NewSyslogWriter("", "", level, logName)
		}
		inited = true
	}
	messageStr := format(level, msg)
	logWriter.Write(level, messageStr)
	if alsoStdout {
		fmt.Print(messageStr)
	}
}

func level() Level {
	return logLevel
}

func Debug(str string) {
	if level() <= LevelDebug {
		write(LevelDebug, str)
	}
}

func Debugs(args ...interface{}) {
	if level() <= LevelDebug {
		write(LevelDebug, fmt.Sprintln(args...))
	}
}

func Debugf(format string, args ...interface{}) {
	if level() <= LevelDebug {
		write(LevelDebug, fmt.Sprintf(format, args...))
	}
}

func Info(str string) {
	if level() <= LevelInfo {
		write(LevelInfo, str)
	}
}

func Infos(args ...interface{}) {
	if level() <= LevelInfo {
		write(LevelInfo, fmt.Sprintln(args...))
	}
}

func Infof(format string, args ...interface{}) {
	if level() <= LevelInfo {
		write(LevelInfo, fmt.Sprintf(format, args...))
	}
}

func Warn(str string) {
	if level() <= LevelWarn {
		write(LevelWarn, str)
	}
}

func Warns(args ...interface{}) {
	if level() <= LevelWarn {
		write(LevelWarn, fmt.Sprintln(args...))
	}
}

func Warnf(format string, args ...interface{}) {
	if level() <= LevelWarn {
		write(LevelWarn, fmt.Sprintf(format, args...))
	}
}

func Error(str string) {
	if level() <= LevelError {
		write(LevelError, str)
	}
}

func Errors(args ...interface{}) {
	if level() <= LevelError {
		write(LevelError, fmt.Sprintln(args...))
	}
}

func Errorf(format string, args ...interface{}) {
	if level() <= LevelError {
		write(LevelError, fmt.Sprintf(format, args...))
	}
}

func NewLogWriter(level Level) io.Writer {
	return &stdWriter{level}
}

type stdWriter struct {
	level Level
}

func (this *stdWriter) Write(data []byte) (int, error) {
	write(this.level, string(data))
	return len(data), nil
}
