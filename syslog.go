package logkit

import "log/syslog"

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

func (self *SyslogWriter) Write(level Level, msg string) {

	switch level {
	case LevelFatal:
		self.writer.Crit(msg)
	case LevelError:
		self.writer.Err(msg)
	case LevelWarn:
		self.writer.Warning(msg)
	case LevelInfo:
		self.writer.Info(msg)
	case LevelDebug:
		self.writer.Debug(msg)
	default:
		self.writer.Write([]byte(msg))
	}
}

func (self *SyslogWriter) Exit() {
	// ignore the error return code
	self.writer.Close()
}
