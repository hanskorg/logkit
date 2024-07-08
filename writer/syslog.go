package writer

import (
	"log/syslog"
)

type SyslogWriter struct {
	network  string
	raddr    string
	priority syslog.Priority
	tag      string
	writer   *syslog.Writer
}

func NewSyslogWriter(network, raddr string, level Level, tag string) (Writer, error) {
	var priority syslog.Priority
	switch level {
	case LevelDebug:
		priority = syslog.LOG_DEBUG
		break
	case LevelInfo:
		priority = syslog.LOG_INFO
		break
	case LevelWarn:
		priority = syslog.LOG_WARNING
		break
	case LevelError:
		priority = syslog.LOG_ERR
		break
	case LevelFatal:
		priority = syslog.LOG_ALERT
		break
	default:
		priority = syslog.LOG_INFO
	}
	writer, err := syslog.Dial(network, raddr, priority, tag)
	if err != nil {
		return nil, err
	}
	object := &SyslogWriter{
		network:  network,
		raddr:    raddr,
		priority: priority,
		tag:      tag,
		writer:   writer,
	}
	return object, nil
}

func (self *SyslogWriter) Write(msg []byte) (int, error) {
	return self.writer.Write([]byte(msg))
}

func (self *SyslogWriter) Close() error {
	// ignore the error return code
	return self.writer.Close()
}
