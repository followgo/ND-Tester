package config

import (
	"os"
	"path/filepath"

	"github.com/followgo/ND-Tester/public/configurator"
	"github.com/followgo/ND-Tester/public/errors"
	"github.com/followgo/ND-Tester/public/helper"
)

const (
	// cfgFileDir 配置文件目录
	cfgFileDir  = "./config"
	hostCfgFile = "host.toml"
	dutCfgFile  = "dut.toml"
)

// LoadHostConfig 加载主机的配置
func LoadHostConfig() error { return tryLoadCfg(hostCfgFile, &Host) }

// LoadDutConfig 加载待测设备的配置
func LoadDutConfig() error { return tryLoadCfg(dutCfgFile, &Dut) }

// tryLoadCfg 加载配置文件到 target 对象
func tryLoadCfg(file string, target interface{}) error {
	if err := os.MkdirAll(cfgFileDir, 0755); err != nil {
		return errors.Wrapf(err, "不能访问 %q 目录", cfgFileDir)
	}
	filePth := filepath.Join(cfgFileDir, file)

	if found, err := helper.HasFile(filePth); err != nil {
		return errors.Wrap(err, "不能读取配置文件")
	} else {
		c := configurator.NewConfigurator(filePth, target)
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

	return nil
}
