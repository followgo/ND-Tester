package main

import (
	"github.com/sirupsen/logrus"

	"github.com/followgo/ND-Tester/config"
	"github.com/followgo/ND-Tester/public/errors"
	"github.com/followgo/ND-Tester/public/mylogrus"
)

func main() {
	// 加载Host配置文件
	if err := config.LoadHostConfig(); err != nil {
		logrus.WithError(err).Fatalln("加载 Host 配置文件失败")
	} else {
		logrus.Infoln("已经加载 Host 配置文件")
	}

	// 初始化日志记录器
	initLogger()
	logrus.Infoln("已经初始化日志记录器")

	// 加载 Dut 配置文件
	if err := config.LoadDutConfig(); err != nil {
		logrus.WithError(err).Fatalln("加载 Dut 配置文件失败")
	} else {
		logrus.Infoln("已经加载 Dut 配置文件")
	}

	logrus.WithError(errors.New("test")).Errorln("测试")
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
	opt.OverWrite = config.Host.Logger.OverWrite
	opt.OutputConsole = config.Host.Logger.OutputConsole
	opt.UseRotate = false
	opt.UseJSONFormat = false
	mylogrus.SetStdLogrus(opt)
}
