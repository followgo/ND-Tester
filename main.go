package main

import (
	"github.com/sirupsen/logrus"

	"github.com/followgo/ND-Tester/config"
	"github.com/followgo/ND-Tester/public/mylogrus"
)

func main() {
	// 加载配置文件
	if err := config.LoadConfigs(); err != nil {
		logrus.WithError(err).Fatalln("加载配置文件失败")
	} else {
		logrus.Infoln("已经加载配置文件")
	}

	// 初始化日志记录器
	initLogger()
	logrus.Infoln("已经初始化日志记录器")
}

// initLogger 初始化日志记录器
func initLogger() {
	opt := mylogrus.DefaultOption
	lvl, err := logrus.ParseLevel(config.Host.Logger.Level)
	if err != nil {
		logrus.WithError(err).Errorln("不认识的日志级别")
	} else {
		opt.Level = lvl
	}

	opt.BaseFile = config.Host.Logger.BaseFile
	opt.MaxMegaSize = config.Host.Logger.MaxMegaSize
	opt.MaxBackups = config.Host.Logger.MaxBackups
	opt.MaxAgeDays = config.Host.Logger.MaxAgeDays
	opt.Console = config.Host.Logger.OutputConsole
	mylogrus.SetStdLogrus(opt)
}
