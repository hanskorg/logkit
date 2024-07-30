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
