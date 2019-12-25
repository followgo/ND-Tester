package mylogrus

import (
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
)

// Option 用于初始化日志记录器的选项
type Option struct {
	// BaseFile 日志文件名
	BaseFile string

	// Level 日志级别
	Level logrus.Level

	// UseJSONFormat 为 true 输出 json 结构化日志
	UseJSONFormat bool

	// OutputConsole 日志同时打印到控制台
	OutputConsole bool

	// OverWrite 覆盖之前的日志文件，当 UseRotate=true 此参数无效
	OverWrite bool

	// UseRotate 使用滚动写入
	UseRotate bool

	// MaxMegaSize 日志文件大小超过后将轮转，单位MB
	MaxMegaSize int

	// MaxBackups 日志保留份数
	MaxBackups int

	// MaxAgeDays 日志保留天数
	MaxAgeDays int

	// DataFormatter 日志的日期格式化
	DataFormatter string
}

// IsEmpty 判断没有赋值
func (o Option) IsEmpty() bool {
	return reflect.DeepEqual(o, Option{})
}

// DefaultOption 默认的日志选项
var DefaultOption = Option{
	BaseFile:      "./logs/main.log",
	Level:         logrus.InfoLevel,
	UseJSONFormat: true,
	OutputConsole: true,
	OverWrite:     false,
	UseRotate:     false,
	MaxMegaSize:   100,
	MaxBackups:    7,
	MaxAgeDays:    7,
	DataFormatter: time.RFC3339,
}
