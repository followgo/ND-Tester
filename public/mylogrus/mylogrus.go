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

	var w io.WriteCloser
	if opt.UseRotate {
		// rotate and compress writer
		w = NewWriterWithSizeRotate(opt.BaseFile, opt.MaxMegaSize, opt.MaxBackups, opt.MaxAgeDays)
	} else {
		var openFlag int
		if opt.OverWrite {
			openFlag = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
		} else {
			openFlag = os.O_CREATE | os.O_WRONLY | os.O_APPEND
		}

		var err error
		if w, err = os.OpenFile(opt.BaseFile, openFlag, 0644); err != nil {
			logrus.WithError(err).Fatalln("create log file")
		}
	}

	logger.SetLevel(opt.Level)

	// 设置日志结构 text or json
	if opt.UseJSONFormat {
		logger.Formatter = &logrus.JSONFormatter{TimestampFormat: opt.DataFormatter}
	} else {
		logger.Formatter = &logrus.TextFormatter{TimestampFormat: opt.DataFormatter}
	}

	if opt.OutputConsole {
		logger.Out = io.MultiWriter(w, os.Stdout)
	} else {
		logger.Out = w
	}

	return logger
}

func SetStdLogrus(opt Option) {
	if opt.IsEmpty() {
		opt = DefaultOption
	}

	var w io.WriteCloser
	if opt.UseRotate {
		// rotate and compress writer
		w = NewWriterWithSizeRotate(opt.BaseFile, opt.MaxMegaSize, opt.MaxBackups, opt.MaxAgeDays)
	} else {
		var openFlag int
		if opt.OverWrite {
			openFlag = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
		} else {
			openFlag = os.O_CREATE | os.O_WRONLY | os.O_APPEND
		}

		var err error
		if w, err = os.OpenFile(opt.BaseFile, openFlag, 0644); err != nil {
			logrus.WithError(err).Fatalln("create log file")
		}
	}

	logrus.SetLevel(opt.Level)

	// 设置日志结构 text or json
	var formatter logrus.Formatter
	if opt.UseJSONFormat {
		formatter = &logrus.JSONFormatter{TimestampFormat: opt.DataFormatter}
	} else {
		formatter = &logrus.TextFormatter{TimestampFormat: opt.DataFormatter}
	}
	logrus.SetFormatter(formatter)

	if opt.OutputConsole {
		logrus.SetOutput(io.MultiWriter(w, os.Stdout))
	} else {
		logrus.SetOutput(w)
	}
}
