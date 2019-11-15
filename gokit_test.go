package logkit

import (
	"flag"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func init()  {
	var l Level
	flag.Int("level", int(l) , "l")
	flag.Parse()
}

func BenchmarkGoKit(b *testing.B) {
	defer Exit()
	Init(FIlE, "test", LevelDebug, true, BasePath)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("test " + strconv.FormatInt(int64(i), 10))
	}
}

func TestInfo(t *testing.T) {
	defer Exit()
	fmt.Println("start")
	for i := 0; i < 10; i++ {
		go func(i int) {
			Info("test " + strconv.FormatInt(int64(i), 10))
		}(i)
	}
	fmt.Println("end")
	time.Sleep(time.Second * 2)
	//
	//for i := 0; i< 1000 ; i++ {
	//	Info("test "+ strconv.FormatInt(int64(i), 10))
	//}

	time.Sleep(time.Second * 1)
	Info("test 2")
	time.Sleep(time.Second * 1)
	Info("test 3")
}

func TestBuffer(t *testing.T) {
	defer Exit()

	//for i := 0 ;  i < 1024; i++  {
	//	str += strconv.FormatInt(int64(i),10)
	//}
	Infof("test %s --- %s", "1", "23")
	//Exit()
}

func TestFlag(t *testing.T)  {
}
