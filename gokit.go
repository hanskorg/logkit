package logkit

import (
	"flag"
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	inited        bool
	auto          bool
	logWriter     io.Writer
	flushInterval time.Duration
	fileSplitSize uint64
	logLevel      = LevelInfo
	logLevelName  string
	logName       string
	logPath       string
	channel       Channel
	alsoStdout    bool
	withCaller    Caller

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

func (l *Level) String() string {
	return levelToNames[*l]
}

// Get is part of the flag.Value interface.
func (l *Level) Get() interface{} {
	return *l
}

func (l *Level) Set(value string) error {
	for i, name := range levelToNames {
		if strings.ToUpper(value) == name {
			*l = i
		}
	}
	if *l == Default {
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		*l = Level(v)
	}
	if *l == Default {
		*l = LevelDebug
	}
	return nil
}

type Channel byte

const (
	FIlE Channel = iota
	SYSLOG
	KAFKA
)

func (c *Channel) String() string {
	switch *c {
	case FIlE:
		return "file"
	case SYSLOG:
		return "syslog"
	}
	return "file"
}
func (c *Channel) Set(value string) error {
	switch value {
	case "file":
		*c = FIlE
	case "syslog":
		*c = SYSLOG
	default:
		*c = FIlE
	}
	return nil
}

type Caller byte

const (
	_ Caller = iota
	NONE
	FullPATHFunc
	BasePathFunc
	BasePath
)

func (c *Caller) String() string {
	switch *c {
	case NONE:
		return "none"
	case FullPATHFunc:
		return "full"
	case BasePathFunc:
		return "file_func"
	case BasePath:
		return "file"
	}
	return "file"
}
func (c *Caller) Set(value string) error {
	switch value {
	case "file":
		*c = BasePath
	case "file_func":
		*c = BasePathFunc
	case "full":
		*c = FullPATHFunc
	default:
		*c = BasePathFunc
	}
	return nil
}

type Writer interface {
	//Write 写日志
	Write(msg []byte) (int, error)
	//Exit 日志退出
	Close() error
}

func GetWriter() io.Closer {
	return logWriter.(io.Closer)
}

func Exit() {
	logWriter.(io.Closer).Close()
}

func init() {
	flag.Var(&logLevel, "log.level", "log level, default `INFO`, it can be `DEBUG, INFO, WARN, ERROR, FATAL`")
	flag.Var(&channel, "log.channel", "write to , it can be `file syslog`")
	flag.Var(&withCaller, "log.withCaller", "call context, by default filename and func name, it can be `file, file_func, full`")

	flag.BoolVar(&alsoStdout, "log.alsoStdout", false, "log out to stand error as well, default `false`")
	flag.StringVar(&logName, "log.name", "log", "log name, by default log will out to `/data/logs/{name}.log`")
	flag.BoolVar(&auto, "log.autoInit", true, "log will be init automatic")
	flag.DurationVar(&flushInterval, "log.interval", time.Second*5, "duration time on flush to disk")
	flag.Uint64Var(&fileSplitSize, "log.split", uint64(1204*1024*1800), "log fail split on bytes")
}

func Init() (writer io.Writer, err error) {
	if inited {
		return nil, fmt.Errorf("logkit has been inited")
	}

	if logName == "" {
		return nil, fmt.Errorf("log name must not be empty")
	}
	if logWriter == nil && channel == FIlE {
		if logPath == "" {
			logPath = "/data/logs/" + logName + ".log"
		}
		logWriter, err = NewFileLogger(logPath, logName, flushInterval, fileSplitSize, 4*1024)
		if err != nil {
			return
		}
	}
	if logWriter == nil && channel == SYSLOG {
		logWriter, _ = NewSyslogWriter("", "", logLevel, logName)
	}
	inited = true
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
			context = fmt.Sprintf("%s:%03d::%30s", path.Base(file), line, path.Base(runtime.FuncForPC(pc).Name()))
		case BasePath:
			context = fmt.Sprintf("%s:%03d", path.Base(file), line)
		default:
			context = fmt.Sprintf("%s:%03d", path.Base(file), line)
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
