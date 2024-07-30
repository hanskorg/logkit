package logkit

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"

	w "github.com/hanskorg/logkit/writer"
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
	channel       w.Channel
	alsoStdout    bool
	withCaller    Caller
	stdOut        io.Writer
	levelToNames  = map[Level]string{
		LevelFatal: "FATAL",
		LevelError: "ERROR",
		LevelWarn:  "WARN",
		LevelInfo:  "INFO",
		LevelDebug: "DEBUG",
		LevelTrace: "TRACE",
	}
)

func GetWriter() (io.Writer, error) {
	if logWriter == nil {
		return nil, fmt.Errorf("logkit not inited")
	}
	return logWriter, nil
}

func Exit() {
	if logWriter == nil {
		return
	}
	logWriter.(io.Closer).Close()
}

// func init() {
// 	flag.Var(&logLevel, "log.level", "log level, default `INFO`, it can be `DEBUG, INFO, WARN, ERROR, FATAL`")
// 	flag.Var(&channel, "log.channel", "write to , it can be `file syslog`")
// 	flag.Var(&withCaller, "log.withCaller", "call context, by default filename and func name, it can be `file, file_func, full`")

// 	flag.BoolVar(&alsoStdout, "log.alsoStdout", false, "log out to stand error as well, default `false`")
// 	flag.StringVar(&logName, "log.name", "log", "log name, by default log will out to `/data/logs/{name}.log`")
// 	flag.BoolVar(&auto, "log.autoInit", true, "log will be init automatic")
// 	flag.DurationVar(&flushInterval, "log.interval", time.Second*5, "duration time on flush to disk")
// 	flag.Uint64Var(&fileSplitSize, "log.split", uint64(1204*1024*1800), "log fail split on bytes")
// }

// SetDebug set logger debug output
func SetDebug(debug bool) {
	if debug {
		alsoStdout = true
		withCaller = BasePathFunc
		logLevel = LevelDebug
	}
}

// SetPath set log filename
// set before inited
func SetPath(path string) {
	logPath = path
}

// SetName set logname
// set before inited
func SetName(name string) {
	logName = name
}

// SetWithCaller set caller flag
// set before inited
func SetWithCaller(withWho string) {
	withCaller.Set(withWho)
}

// SetAlsoStdout set stdout or not
// set before inited
func SetAlsoStdout(stdout bool) {
	alsoStdout = stdout
}

// SetChannel set channel
// set before inited
func SetChannel(channelName string) {
	channel.Set(channelName)
}

// SetLevel set level
func SetLevel(level Level) {
	logLevel = level
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
			logPath = "/var/log/" + logName + ".log"
		}

		logWriter, err = NewFileLogger(logPath, logName, flushInterval, fileSplitSize, 4*1024)
		if err != nil {
			return
		}
	}
	if logWriter == nil && channel == SYSLOG {
		logWriter, _ = NewSyslogWriter("", "", logLevel, logName)
	}
	if logWriter == nil && channel == STDOUT {
		logWriter = os.Stdout
	}
	if alsoStdout {
		if channel == STDOUT {
			stdOut = logWriter
		} else {
			stdOut = os.Stdout
		}
	}
	inited = true
	return logWriter, nil
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
			context = fmt.Sprintf("%s:%03d::%-30s", file, line, path.Base(runtime.FuncForPC(pc).Name()))
		case BasePathFunc:
			context = fmt.Sprintf("%s:%03d::%-15s", path.Base(file), line, path.Base(runtime.FuncForPC(pc).Name()))
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
	if !flag.Parsed() {
		return fmt.Errorf("logkit write must been flag parsed")
	}
	if auto && !inited {
		Init()
	}
	if !inited {
		return fmt.Errorf("logkit has been inited")
	}
	messageStr := format(level, msg)
	_, err = logWriter.Write([]byte(messageStr))
	if alsoStdout {
		stdOut.Write([]byte(messageStr))
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
