package fileout

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志输出接口
type WriteSyncer interface {
	io.Writer
	Sync() error
}

// 创建分割日志的write
func NewRollingFile(path, srvname string, maxSize, maxAge, maxBackups int) WriteSyncer {
	if err := os.MkdirAll(path, 0766); err != nil {
		panic(err)
	}

	return newLumberjackWriteSyncer(&lumberjack.Logger{
		Filename:   filepath.Join(path, srvname+".log"), // 日志文件路径
		MaxSize:    maxSize,                             // 每个日志文件保存的最大尺寸 单位：M
		MaxAge:     maxAge,                              // 文件最多保存多少天
		MaxBackups: maxBackups,                          // 日志文件最多保存多少个备份
		LocalTime:  true,                                // 是否使用本地时间
		Compress:   false,                               // 是否压缩
	})
}

type lumberjackWriteSyncer struct {
	*lumberjack.Logger
	buf       *bytes.Buffer
	logChan   chan []byte
	closeChan chan interface{}
	maxSize   int
}

func (l *lumberjackWriteSyncer) Write(bs []byte) (int, error) {
	b := make([]byte, len(bs))
	for i, c := range bs {
		b[i] = c
	}
	l.logChan <- b
	return 0, nil
}

func (l *lumberjackWriteSyncer) Sync() error {
	return nil
}

func (l *lumberjackWriteSyncer) Stop() {
	close(l.closeChan)
}

func (l *lumberjackWriteSyncer) sync() error {
	defer l.buf.Reset()
	_, err := l.Logger.Write(l.buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (l *lumberjackWriteSyncer) run() {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			if l.buf.Len() > 0 {
				l.sync()
			}
		case bs := <-l.logChan:
			_, err := l.buf.Write(bs)
			if err != nil {
				continue
			}
			if l.buf.Len() > l.maxSize {
				l.sync()
			}
		case <-l.closeChan:
			l.sync()
			return
		}
	}
}

func newLumberjackWriteSyncer(l *lumberjack.Logger) *lumberjackWriteSyncer {
	ws := &lumberjackWriteSyncer{
		Logger:    l,
		buf:       bytes.NewBuffer([]byte{}),
		logChan:   make(chan []byte, 5000),
		closeChan: make(chan interface{}),
		maxSize:   1024,
	}
	go ws.run()
	return ws
}
