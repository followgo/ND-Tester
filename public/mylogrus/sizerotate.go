package mylogrus

import (
	"io"

	"github.com/natefinch/lumberjack"
)

// NewWriterWithSizeRotate 新建一个按大小滚动的文件 Writer
func NewWriterWithSizeRotate(baseFile string, maxMegaSize, maxBackups, maxAgeDays int) io.WriteCloser {
	return &lumberjack.Logger{
		Filename:   baseFile,
		MaxSize:    maxMegaSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAgeDays,
		Compress:   true,
		LocalTime:  true,
	}
}
