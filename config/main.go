package config

import (
	"os"
	"path/filepath"

	"github.com/followgo/ND-Tester/public/configurator"
	"github.com/followgo/ND-Tester/public/errors"
	"github.com/followgo/ND-Tester/public/helper"
)

// LoadConfigs 加载所有配置
// 如果配置不存在，则生成默认配置文件
func LoadConfigs() error {
	var (
		cfgDir                  = "./config"
		hostCfgFile, dutCfgFile = filepath.Join(cfgDir, "host.yaml"), filepath.Join(cfgDir, "dut.yaml")
	)

	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		return errors.Wrapf(err, "不能访问 %q 目录", cfgDir)
	}

	// load host config file
	if found, err := helper.HasFile(hostCfgFile); err != nil {
		return errors.Wrap(err, "不能读取配置文件")
	} else {
		c := configurator.NewConfigurator(hostCfgFile, &Host)
		if found {
			if err := c.Load(); err != nil {
				return errors.Wrap(err, "加载配置文件")
			}
		} else {
			if err := c.Save("运行测试框架的主机配置文件"); err != nil {
				return errors.Wrap(err, "保存配置文件")
			}
		}
	}

	// load dut config file
	if found, err := helper.HasFile(dutCfgFile); err != nil {
		return errors.Wrap(err, "不能读取配置文件")
	} else {
		c := configurator.NewConfigurator(dutCfgFile, &Dut)
		if found {
			if err := c.Load(); err != nil {
				return errors.Wrap(err, "加载配置文件")
			}
		} else {
			if err := c.Save("待测设备的配置文件"); err != nil {
				return errors.Wrap(err, "保存配置文件")
			}
		}
	}

	return nil
}
