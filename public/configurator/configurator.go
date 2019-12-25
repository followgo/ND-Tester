// configurator 配置装载器，配置的加载、保存、文件监视。
// 支持的文件类型: yaml, toml, json
package configurator

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"github.com/followgo/ND-Tester/public/errors"
)

// watchConfigs 已加载的配置文件，当接收到 HUP 进程信号会重载
var watchConfigs = make(map[string]interface{})

// init 启动一个协程监视进程信号(kill -HUP)，如果收到则执行 Load()
// func init() {
// 	go func() {
// 		HUP := make(chan os.Signal)
// 		for {
// 			signal.Notify(HUP, syscall.SIGHUP)
// 			<-HUP
//
// 			logrus.Info("获取一个 HUP 进程信号，开始重新加载配置文件...")
// 			for k, v := range watchConfigs {
// 				c := NewConfigurator(k, v)
// 				if err := c.Load(); err != nil {
// 					logrus.WithField("filename", c.Filename).WithError(err).Error("重载配置文件")
// 				}
// 			}
// 		}
// 	}()
// }

type configurator struct {
	// Filename 配置文件名
	Filename string

	// TargetObj 存储配置的对象
	TargetObj interface{}

	// fileType 配置文件类型，从文件扩展名中获取。可选 .json, .yaml, .toml
	fileType string
}

// NewConfigurator 创建一个配置器
func NewConfigurator(filename string, targetObj interface{}) *configurator {
	return &configurator{
		Filename:  filename,
		fileType:  filepath.Ext(filename),
		TargetObj: targetObj,
	}
}

// Load 从指定文件中加载配置
func (c *configurator) Load() error {
	data, err := ioutil.ReadFile(c.Filename)
	if err != nil {
		return errors.Wrapf(err, "read the %q file", c.Filename)
	}

	// 添加到监视 hashMap 中
	watchConfigs[c.Filename] = c.TargetObj

	switch typ := strings.ToLower(c.fileType); typ {
	case ".yaml":
		return yaml.Unmarshal(data, c.TargetObj)
	case ".toml":
		return toml.Unmarshal(data, c.TargetObj)
	case ".json":
		return json.Unmarshal(data, c.TargetObj)
	default:
		return errors.Errorf("the file type of %q is not supported", c.Filename)
	}
}

// Save 保存配置到指定的文件
func (c *configurator) Save(fileComment string) error {
	f, err := os.OpenFile(c.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "open the %q file", c.Filename)
	}
	defer f.Close()

	var buf = new(bytes.Buffer)
	switch typ := strings.ToLower(c.fileType); typ {
	case ".yaml":
		data, err := yaml.Marshal(c.TargetObj)
		if err != nil {
			return errors.Wrapf(err, "save config to %q file", c.Filename)
		}
		buf.WriteString("# " + fileComment + "\n\n")
		buf.Write(data)

	case ".toml":
		buf.WriteString("# " + fileComment + "\n\n")
		if err := toml.NewEncoder(buf).Encode(c.TargetObj); err != nil {
			return errors.Wrapf(err, "save config to %q file", c.Filename)
		}

	case ".json": // json 不支持注释
		data, err := json.Marshal(c.TargetObj)
		if err != nil {
			return errors.Wrapf(err, "save config to %q file", c.Filename)
		}
		buf.Write(data)

	default:
		return errors.Errorf("the file type of %q is not supported", c.Filename)
	}

	if _, err := buf.WriteTo(f); err != nil {
		return errors.Wrapf(err, "save config to %q file", c.Filename)
	}
	return nil
}

// Watching 监视配置文件，如果配置文件发生变化，则执行 Load()
// 忽略错误，即如果
// func (c *configurator) Watching() {
// 	go func() {
// 		var lastModTime time.Time
//
// 		for range time.Tick(10 * time.Second) {
// 			fInfo, err := os.Stat(c.Filename)
// 			if err != nil {
// 				logrus.WithField("filename", c.Filename).WithError(err).Error("获取文件信息")
// 			}
// 			newModTime := fInfo.ModTime()
//
// 			if !lastModTime.IsZero() && lastModTime != newModTime {
// 				if err := c.Load(); err != nil {
// 					logrus.WithField("filename", c.Filename).WithError(err).Error("重载配置文件")
// 				}
// 			}
//
// 			lastModTime = newModTime
// 		}
// 	}()
// }
