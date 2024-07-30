package logkit

import "strings"

type Writer interface {
	//Write 写日志
	Write(msg []byte) (int, error)
	//Close 日志退出
	Close() error
}

type Channel byte

const (
	FIlE   Channel = iota
	SYSLOG         = 0b0001
	KAFKA          = 0b0010
	STDOUT         = 0b0100
	LOGRUS         = 0b1000
)

func (c *Channel) String() string {
	switch *c {
	case FIlE:
		return "file"
	case SYSLOG:
		return "syslog"
	case STDOUT:
		return "none"
	}
	return "file"
}
func (c *Channel) Set(value string) error {
	switch strings.ToLower(value) {
	case "file":
		*c = FIlE
	case "syslog":
		*c = SYSLOG
	case "none":
		*c = STDOUT
	default:
		*c = FIlE
	}
	return nil
}
