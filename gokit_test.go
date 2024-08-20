package logkit

import (
	"strconv"
	"testing"
)

func BenchmarkKit(b *testing.B) {
}

func TestStdInfo(t *testing.T) {
	defer Close()
	//NewLogger(WithChannel("file"), WithChannel("logrus"))
	Info("fuck")
	Info("this this a test info")
}

func TestFileInfo(t *testing.T) {
	logger := NewLogger(WithChannel(FIlE), WithChannel(STDOUT), SetWithCaller(FullPATHFunc))
	defer logger.Close()
	logger.Log(LevelInfo, "this this a test info", "just for test")
	logger.Log(LevelDebug, "this this a test info", "just for test")
	logger.Log(LevelWarn, "this this a test info", "just for test")
	logger.Log(LevelFatal, "this this a test info", "just for test")
	logger.Log(LevelFatal, "this this a test info", "just for test")
}

func TestSyslogInfo(t *testing.T) {
	logger := NewLogger(WithChannel(FIlE), WithChannel(SYSLOG), SetLogName("test"), SetSysLogAddr("udp://localhost:514"), SetWithCaller(FullPATHFunc))
	defer logger.Close()
	logger.Log(LevelInfo, "this this a test info", "just for test")
	logger.Log(LevelDebug, "this this a test info", "just for test")
	logger.Log(LevelWarn, "this this a test info", "just for test")
	logger.Log(LevelFatal, "this this a test info", "just for test")
	logger.Log(LevelFatal, "this this a test info", "just for test")
}

func TestDefaultLogger(t *testing.T) {
	defer Close()
	SetLogger(WithChannel(FIlE), WithChannel(STDOUT), SetLogName("test"), SetWithCaller(BasePath))
	Infof("this this a test info")
}

func TestBuffer(t *testing.T) {
	defer Close()
	var str string
	for i := 0; i < 1024; i++ {
		str += strconv.FormatInt(int64(i), 10)
	}
	Infof("test %s --- %s", "1", str)
}
