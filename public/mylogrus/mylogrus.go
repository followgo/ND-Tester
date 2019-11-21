// mylogrus 日志记录器。
// 包装 logrus, lumberjack，支持日志文件能按大小滚动
package mylogrus

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// NewMyLogrus 创建一个 Logrus 实例，支持日志文件能按大小滚动
func NewMyLogrus(opt Option) *logrus.Logger {
	if opt.IsEmpty() {
		opt = DefaultOption
	}
	logger := logrus.New()

	// rotate and compress writer
	rollingWriter := NewWriterWithSizeRotate(opt.BaseFile, opt.MaxMegaSize, opt.MaxBackups, opt.MaxAgeDays)
	logger.SetLevel(opt.Level)

	// 设置日志结构 text or json
	if opt.UseJSONFormat {
		logger.Formatter = &logrus.JSONFormatter{TimestampFormat: opt.DataFormatter}
	} else {
		logger.Formatter = &logrus.TextFormatter{TimestampFormat: opt.DataFormatter}
	}

	if opt.Console {
		logger.Out = io.MultiWriter(rollingWriter, os.Stdout)
	} else {
		logger.Out = rollingWriter
	}

	return logger
}

func SetStdLogrus(opt Option) {
	if opt.IsEmpty() {
		opt = DefaultOption
	}

	// rotate and compress writer
	rollingWriter := NewWriterWithSizeRotate(opt.BaseFile, opt.MaxMegaSize, opt.MaxBackups, opt.MaxAgeDays)
	logrus.SetLevel(opt.Level)

	// 设置日志结构 text or json
	var formatter logrus.Formatter
	if opt.UseJSONFormat {
		formatter = &logrus.JSONFormatter{TimestampFormat: opt.DataFormatter}
	} else {
		formatter = &logrus.TextFormatter{TimestampFormat: opt.DataFormatter}
	}
	logrus.SetFormatter(formatter)

	if opt.Console {
		logrus.SetOutput(io.MultiWriter(rollingWriter, os.Stdout))
	} else {
		logrus.SetOutput(rollingWriter)
	}
}
