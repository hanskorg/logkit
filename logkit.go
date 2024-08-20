package logkit

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	defaultFlushInterval = time.Second
	defaultFileSplitSize = uint64(1204 * 1024 * 1800)
	defaultBuffSize      = 1024 * 4
)

type Logger struct {
	level         Level
	writers       map[Channel]Writer
	logPath       string
	channels      []Channel
	caller        Caller
	adapter       string
	flushInterval time.Duration
	fileSplitSize uint64
	fileBuffSize  uint64
	defaultWriter Writer
}

var defaultLogger = &Logger{
	level: LevelInfo,
	writers: map[Channel]Writer{
		STDOUT: &stdWriter{
			level:  LevelInfo,
			caller: BasePath,
		}},
	caller:        BasePath,
	flushInterval: defaultFlushInterval,
	fileSplitSize: defaultFileSplitSize,
	fileBuffSize:  defaultBuffSize,
	logPath:       "./logkit.log",
}

func (l *Logger) Log(level Level, args ...interface{}) {
	switch level {
	case LevelDebug:
		l.Debugs(args...)
	case LevelInfo:
		l.Infos(args...)
	case LevelWarn:
		l.Warns(args...)
	case LevelError:
		l.Errors(args...)
	case LevelFatal:
		l.Fatal(args...)
	default:
		l.Debugs(args...)
	}
}

func NewLogger(opts ...Option) *Logger {
	return setupLogger(&Logger{
		level: LevelInfo,
		writers: map[Channel]Writer{
			STDOUT: &stdWriter{
				level:  LevelInfo,
				caller: BasePath,
			}},
		caller:        BasePath,
		flushInterval: defaultFlushInterval,
		fileSplitSize: defaultFileSplitSize,
		fileBuffSize:  defaultBuffSize,
		logPath:       "./logkit.log",
	}, opts...)
}

func (l *Logger) Close() {
	closes(l.writers)
}
func SetLogger(opts ...Option) {
	setupLogger(defaultLogger, opts...)
}
func setupLogger(_logger *Logger, options ...Option) *Logger {
	var (
		e             error
		withStdout    bool
		fileWriter    Writer
		logrusAdapter *logrus.Logger
	)
	for _, opt := range options {
		opt(_logger)
	}

	for _, c := range _logger.channels {
		if c == STDOUT {
			withStdout = true
		}
		if (c == FIlE || c == LOGRUS) && fileWriter == nil {
			if fileWriter, e = NewFileLogger(_logger.logPath, _logger.flushInterval, _logger.fileSplitSize, _logger.fileBuffSize); e != nil {
				println(fmt.Sprintf("logkit new filelogger fail, %s", e.Error()))
			}
			if logrusAdapter != nil {
				logrusAdapter.Out = fileWriter
			}
			_logger.defaultWriter = fileWriter
		}
		if c == FIlE {
			_logger.writers[FIlE] = fileWriter
		}
		if c == LOGRUS {
			logrusAdapter = &logrus.Logger{
				Out:       _logger.defaultWriter,
				Formatter: new(logrus.TextFormatter),
			}
			logrusAdapter.Level, _ = logrus.ParseLevel(_logger.level.String())
			_logger.writers[LOGRUS] = &logrusWriter{
				logrusAdapter,
				logrusAdapter.Level,
			}
		}

	}
	if !withStdout {
		delete(_logger.writers, STDOUT)
	}
	return _logger
}

func Close() {
	closes(defaultLogger.writers)
}
func closes(ws map[Channel]Writer) {
	for c, w := range ws {
		if e := w.Close(); e != nil {
			println("logkit writer [%s] close fail, %s", c.String(), e.Error())
		}
	}
}
func GetWriter(channel Channel) Writer {
	if w, has := defaultLogger.writers[channel]; has {
		return w
	}
	return defaultLogger.defaultWriter
}

func getLevelName(level Level) string {
	levelName, _ := levelToNames[level]
	return levelName
}

func level() Level {
	return defaultLogger.level
}

func (l *Logger) Debugs(args ...interface{}) {
	if l.level <= LevelDebug {
		write(l, LevelDebug, fmt.Sprintln(args...))
	}
}

func (l *Logger) Infos(args ...interface{}) {
	if level() <= LevelInfo {
		write(l, LevelInfo, fmt.Sprintln(args...))
	}
}

func (l *Logger) Warns(args ...interface{}) {
	if level() <= LevelWarn {
		write(l, LevelWarn, fmt.Sprintln(args...))
	}
}

func (l *Logger) Errors(args ...interface{}) {
	if level() <= LevelError {
		write(l, LevelError, fmt.Sprintln(args...))
	}
}

func (l *Logger) Fatal(args ...interface{}) {
	write(l, LevelFatal, fmt.Sprintln(args...))
	l.Close()
	os.Exit(1)
}

func Debug(str string) {
	if level() <= LevelDebug {
		write(defaultLogger, LevelDebug, str)
	}
}

func Debugf(format string, args ...interface{}) {
	if level() <= LevelDebug {
		write(defaultLogger, LevelDebug, fmt.Sprintf(format, args...))
	}
}

func Info(str string) {
	if level() <= LevelInfo {
		write(defaultLogger, LevelInfo, str)
	}
}

func Infof(format string, args ...interface{}) {
	if level() <= LevelInfo {
		write(defaultLogger, LevelInfo, fmt.Sprintf(format, args...))
	}
}

func Warn(str string) {
	if level() <= LevelWarn {
		write(defaultLogger, LevelWarn, str)
	}
}

func Warnf(format string, args ...interface{}) {
	if level() <= LevelWarn {
		write(defaultLogger, LevelWarn, fmt.Sprintf(format, args...))
	}
}

func Error(str string) {
	if level() <= LevelError {
		write(defaultLogger, LevelError, str)
	}
}

func Errorf(format string, args ...interface{}) {
	if level() <= LevelError {
		write(defaultLogger, LevelError, fmt.Sprintf(format, args...))
	}
}
