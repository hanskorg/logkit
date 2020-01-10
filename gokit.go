package logkit

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"time"
)

var (
	inited       bool
	logWriter    io.Writer
	logLevel     = LevelInfo
	logLevelName string
	logName      string
	logPath      string
	channel      Channel
	alsoStdout   bool
	withCaller   Caller

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
type Caller byte

const (
	FIlE Channel = iota
	SYSLOG
	KAFKA
)

const (
	_ Caller = iota
	NONE
	FullPATHFunc
	BasePathFunc
	BasePath
)

type Writer interface {
	//Write 写日志
	Write(msg []byte) (int, error)
	//Exit 日志退出
	Close() error
}

func Exit() {
	logWriter.(io.Closer).Close()
}

func Init(_channel Channel, name string, level Level, _alsoStdout bool, _withCaller Caller) (writer io.Writer, err error ){
	if inited {
		return nil, fmt.Errorf("logkit has been inited")
	}

	if name != "" {
		logName = name
	} else {
		return nil, fmt.Errorf("log name must not be empty")
	}

	if logWriter == nil && channel == FIlE {
		if logPath == "" {
			logPath = "/data/logs/" + logName + ".log"
		}
		logWriter, err  = NewFileLogger(logPath, logName, time.Second*5, 1204*1024*1800, 4*1024)
		if err != nil {
			return
		}
	}
	if logWriter == nil && channel == SYSLOG {
		logWriter, _ = NewSyslogWriter("", "", level, logName)
	}
	inited = true

	logLevel = level
	channel = _channel
	alsoStdout = _alsoStdout
	withCaller = _withCaller
	return logWriter, nil
}

func SetPath(path string) {
	logPath = path
}
func getLevelName(level Level) string {
	levelName, _ := levelToNames[level]
	return levelName
}

func format(level Level, msg string) string {
	if withCaller != NONE {
		var (
			context string
			pc      uintptr
			file    string
			line    int
		)
		pc, file, line, _ = runtime.Caller(3)
		switch withCaller {
		case FullPATHFunc:
			context = fmt.Sprintf("%s:%03d::%30s", file, line, path.Base(runtime.FuncForPC(pc).Name()))
		case BasePathFunc:
			context = fmt.Sprintf("%s:%03d::%30s",   path.Base(file), line,  path.Base(runtime.FuncForPC(pc).Name()))
		case BasePath:
			context = fmt.Sprintf("%s:%03d",   path.Base(file), line)
		default:
			context = fmt.Sprintf("%s:%03d",   path.Base(file), line)
		}

		return fmt.Sprintf("%s\t[%4s]\t%s\t%s\n", time.Now().Format("2006-01-02 15:04:05.999"), getLevelName(level), context, msg)
	} else {
		return fmt.Sprintf("%s\t[%4s]\t%s\n", time.Now().Format("2006-01-02 15:04:05.999"), getLevelName(level), msg)
	}
}

func write(level Level, msg string) (err error) {
	if !inited {
		return fmt.Errorf("logkit has been inited")
	}
	messageStr := format(level, msg)
	_, err = logWriter.Write([]byte(messageStr))
	if alsoStdout {
		fmt.Print(messageStr)
	}
	return
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
