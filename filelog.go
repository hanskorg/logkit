package logkit

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type bufferNode struct {
	bytes.Buffer
	next *bufferNode
}

type mFileLogger struct {
	writer *bufferWriter
	mu     sync.Mutex

	filepath      string
	freeList      *bufferNode
	freeListMu    sync.Mutex
	flushInterval time.Duration
	fileSplitSize uint64
	bufferSize    uint64
}

type bufferWriter struct {
	*bufio.Writer
	logPath     string
	logName     string
	file        *os.File
	slot        int
	startTime   time.Time
	byteSize    uint64 // The number of bytes written to this file
	maxFileSize uint64
	bufferSize  uint64
}

func (w *mFileLogger) getBuffer() *bufferNode {
	w.freeListMu.Lock()
	b := w.freeList
	if b != nil {
		w.freeList = b.next
	}
	w.freeListMu.Unlock()
	if b == nil {
		b = new(bufferNode)
	} else {
		b.next = nil
		b.Reset()
	}
	return b
}

func (w *mFileLogger) putBuffer(b *bufferNode) {
	if b.Len() >= 256 {
		// Let big buffers die with gc.
		return
	}
	w.freeListMu.Lock()
	b.next = w.freeList
	w.freeList = b
	w.freeListMu.Unlock()
}

func (w *mFileLogger) Close() error {
	return w.flush()
}

func NewFileLogger(path string, flushInterval time.Duration, fileSplitSize uint64, bufferSize uint64) (writer Writer, err error) {
	writer = &mFileLogger{
		filepath:      path,
		flushInterval: flushInterval,
		fileSplitSize: fileSplitSize,
		bufferSize:    bufferSize,
		writer: &bufferWriter{
			logPath:     path,
			maxFileSize: fileSplitSize,
			bufferSize:  bufferSize,
		},
	}
	err = writer.(*mFileLogger).writer.rotate(time.Now(), 0)
	if err != nil {
		return
	}
	go writer.(*mFileLogger).flushDaemon()
	return
}

func (w *mFileLogger) flushDaemon() {
	for _ = range time.NewTicker(w.flushInterval).C {
		w.flush()
	}
}

func (w *mFileLogger) flush() (err error) {
	if w.writer == nil || w.writer.Writer == nil {
		return
	}
	err = w.writer.Flush()
	if err != nil {
		return
	}
	err = w.writer.Sync()
	return
}

func (w *mFileLogger) Write(msg []byte) (n int, err error) {

	buf := w.getBuffer()
	w.mu.Lock()
	defer w.mu.Unlock()
	buf.Write(msg)
	writer := w.writer

	if writer == nil {
		w.writer = &bufferWriter{
			logPath:     w.filepath,
			maxFileSize: w.fileSplitSize,
			bufferSize:  w.bufferSize,
		}
		writer = w.writer
	}

	if err = writer.checkRotate(time.Now()); err != nil {
		fmt.Println("[logkit] check rotate err: " + err.Error())
		return
	}

	w.putBuffer(buf)
	return writer.Write(buf.Bytes())
}

func (bufferW *bufferWriter) Write(p []byte) (n int, err error) {
	bufferW.Writer.Available()
	n, err = bufferW.Writer.Write(p)
	bufferW.byteSize += uint64(n)
	if bufferW.Buffered() >= int(bufferW.bufferSize) {
		err = bufferW.Writer.Flush()
	}
	return n, err
}

func (bufferW *bufferWriter) Sync() error {
	return bufferW.file.Sync()
}

func (bufferW *bufferWriter) checkRotate(now time.Time) error {
	if bufferW.file == nil {
		return bufferW.rotate(now, 0)
	}
	sYear, sMonth, sDay := bufferW.startTime.Date()
	year, month, day := now.Date()
	if year != sYear || month != sMonth || day != sDay {
		return bufferW.rotate(now, 0)
	}
	if bufferW.byteSize >= bufferW.maxFileSize {
		return bufferW.rotate(now, bufferW.slot+1)
	}
	return nil
}

func (bufferW *bufferWriter) rotate(oldTime time.Time, slot int) (err error) {
	var (
		newFileName      string
		year, month, day = oldTime.Date()
		fileInfo         os.FileInfo
	)

	if bufferW.file != nil {
		err = bufferW.Flush()
		if err != nil {
			return
		}
		err = bufferW.file.Close()
		if err != nil {
			return
		}
		if slot > 0 {
			newFileName = fmt.Sprintf("%s-%02d%02d%02d.%02d", bufferW.logPath, year, month, day, slot-1)
		} else {
			newFileName = fmt.Sprintf("%s-%02d%02d%02d", bufferW.logPath, year, month, day)
		}
		err = os.Rename(bufferW.logPath, newFileName)
	}

	if err := bufferW.openFile(bufferW.logPath, bufferW.logName); err != nil {
		return fmt.Errorf("rotate file error: %s", err.Error())
	}
	fileInfo, err = bufferW.file.Stat()
	if err != nil {
		return
	}
	bufferW.byteSize = uint64(fileInfo.Size())
	bufferW.Writer = bufio.NewWriterSize(bufferW.file, int(bufferW.bufferSize))
	bufferW.slot = slot
	bufferW.startTime = time.Now()
	bufferW.byteSize = 0

	return nil
}

func (bufferW *bufferWriter) openFile(fileName, logName string) (err error) {
	var (
		file *os.File
	)
	for {
		file, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err == nil {
			break
		}
		// try to create all the parent directories for specified log file
		// if it doesn't exist
		if os.IsNotExist(err) {
			err2 := os.MkdirAll(filepath.Dir(fileName), 0755)
			if err2 != nil {
				return err
			}
			continue
		}
		return err
	}
	bufferW.file = file
	return nil
}
