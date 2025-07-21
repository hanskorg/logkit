package logkit

import "os"

type SyslogWriter struct {
	network string
	raddr   string
	//priority syslog.Priority
	tag string
	//writer   *syslog.Writer
}

func NewSyslogWriter(network, raddr string, level Level, tag string) (Writer, error) {
	return os.Stdout, nil
}

func (self *SyslogWriter) Write(msg []byte) (int, error) {
	return 0, nil
}

func (self *SyslogWriter) Close() error {
	// ignore the error return code
	return nil
}
