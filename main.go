package main

import (
	"github.com/sirupsen/logrus"

	"github.com/followgo/ND-Tester/config"
)

func main() {
	// 加载配置文件
	if err := config.LoadConfigs(); err != nil {
		logrus.WithError(err).Fatalln("加载配置文件失败")
	} else {
		logrus.Infoln("已经加载配置文件")
	}

	// 初始化日志记录器
}
