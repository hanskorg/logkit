package logkit

import (
	"strconv"
	"strings"
)

// Level 日志等级
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

var (
	levelToNames = map[Level]string{
		LevelFatal: "FATAL",
		LevelError: "ERR",
		LevelWarn:  "WAN",
		LevelInfo:  "INF",
		LevelDebug: "DEG",
		LevelTrace: "TRE",
	}
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

type Caller byte

const (
	_ Caller = iota
	NONE
	FullPATHFunc
	BasePathFunc
	BasePath
)

type Channel byte

const (
	STDOUT Channel = iota
	SYSLOG         = 0b0001
	KAFKA          = 0b0010
	FIlE           = 0b0100
	LOGRUS         = 0b1000
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
func (c *Channel) String() string {
	switch *c {
	case FIlE:
		return "file"
	case SYSLOG:
		return "syslog"
	case STDOUT:
		return "stdout"
	case LOGRUS:
		return "logrus"
	}
	return "file"
}
func (c Channel) Set(value string) Channel {
	switch strings.ToLower(value) {
	case "file":
		c = FIlE
	case "syslog":
		c = SYSLOG
	case "logurs":
		c = LOGRUS
	case "none":
		c = STDOUT
	default:
		c = FIlE
	}
	return c
}

type Option func(*Logger)

func SetLevel(_level Level) Option {
	return func(l *Logger) {
		l.level = _level
	}
}

// WithChannel add Output channel
func WithChannel(channel Channel) Option {
	return func(l *Logger) {
		if channel == KAFKA {
			println("logkit ignore kafka channel, not implement")
			return
		}
		l.channels = append(l.channels, channel)
	}
}

// SetPath set log filename
func SetPath(path string) Option {
	return func(l *Logger) {
		l.logPath = path
	}
}

// SetWithCaller set caller flag
// set before inited
func SetWithCaller(withWho Caller) Option {
	return func(l *Logger) {
		l.caller = withWho
	}
}

func SetFileSplitSize(size uint64) Option {
	return func(l *Logger) {
		l.fileSplitSize = size
	}
}

func SetFileBuffSize(size uint64) Option {
	return func(l *Logger) {
		l.fileBuffSize = size
	}
}

func SetSysLogAddr(addr string) Option {
	return func(l *Logger) {
		l.syslogAddr = addr
	}
}

func SetLogName(name string) Option {
	return func(l *Logger) {
		l.name = name
	}
}
