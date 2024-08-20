package logkit

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"time"
)

type Writer interface {
	//Write 写日志
	Write(msg []byte) (int, error)
	//Close 日志退出
	Close() error
}

type logrusWriter struct {
	*logrus.Logger
	level logrus.Level
}

func (w *logrusWriter) Write(msg []byte) (int, error) {
	return w.WriterLevel(w.level).Write(msg)
}
func (w *logrusWriter) WriteWithLevel(level logrus.Level, msg []byte) (wd int, e error) {
	wd, e = w.WriterLevel(level).Write(msg)
	return
}
func (w *logrusWriter) Close() error {
	return w.Writer().Close()
}

type stdWriter struct {
	level  Level
	caller Caller
}

func (sw *stdWriter) Write(msg []byte) (int, error) {
	return fmt.Println(string(msg))
}
func (*stdWriter) Close() error {
	return nil
}

func format(level Level, caller Caller, msg string) string {
	if caller != NONE {
		var (
			context string
			pc      uintptr
			file    string
			line    int
		)
		pc, file, line, _ = runtime.Caller(3)
		switch caller {
		case FullPATHFunc:
			context = fmt.Sprintf("%s:%03d::%-30s", file, line, path.Base(runtime.FuncForPC(pc).Name()))
		case BasePathFunc:
			context = fmt.Sprintf("%s:%03d::%-15s", path.Base(file), line, path.Base(runtime.FuncForPC(pc).Name()))
		case BasePath:
			context = fmt.Sprintf("%s:%03d", path.Base(file), line)
		default:
			context = fmt.Sprintf("%s:%03d", path.Base(file), line)
		}
		return fmt.Sprintf("%s [%4s] %s %-44s\r", time.Now().Format("2006-01-02 15:04:05.999"), getLevelName(level), context, msg)
	} else {
		return fmt.Sprintf("%s [%4s] %-44s\r", time.Now().Format("2006-01-02 15:04:05.999"), getLevelName(level), msg)
	}
}

func write(l *Logger, level Level, msg string) {
	for c, w := range l.writers {
		var err error
		if c == LOGRUS {
			_, err = w.Write([]byte(msg))
		} else {
			_, err = w.Write([]byte(format(level, l.caller, msg)))
		}
		if err != nil {
			println(fmt.Sprintf("logkit write fail, %s", err.Error()))
		}
	}

}
