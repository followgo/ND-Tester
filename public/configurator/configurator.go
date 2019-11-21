// configurator 配置装载器，配置的加载、保存、文件监视。
// 支持的文件类型: yaml, toml, json
package configurator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

// Logger 日志记录器
var Logger = log.New(os.Stdout, "[config] ", log.LstdFlags)

type configurator struct {
	// Filename 配置文件名
	Filename string

	// TargetObj 存储配置的对象
	TargetObj interface{}

	// fileType 配置文件类型，从文件扩展名中获取。可选 .json, .yaml, .toml
	fileType string

	// sigWatcher 进程信号监视器。防止一个文件启动多个监视器
	once sync.Once
}

// NewConfigurator 创建一个配置器
func NewConfigurator(filename string, targetObj interface{}) *configurator {
	return &configurator{
		Filename:  filename,
		fileType:  filepath.Ext(filename),
		TargetObj: targetObj,
	}
}

// Load 从指定文件中加载配置，会开启一个协程监视文件，如果发生修改则重新加载
func (c *configurator) Load() error {
	data, err := ioutil.ReadFile(c.Filename)
	if err != nil {
		return fmt.Errorf("读 %q 配置文件: [%w]", c.Filename, err)
	}

	// 监视进程信号
	c.once.Do(func() { go c.watchSig() })

	switch typ := strings.ToLower(c.fileType); typ {
	case ".yaml":
		return yaml.Unmarshal(data, c.TargetObj)
	case ".toml":
		return toml.Unmarshal(data, c.TargetObj)
	case ".json":
		return json.Unmarshal(data, c.TargetObj)
	default:
		return fmt.Errorf("不能加载 %q 文件，不支持此文件类型", c.Filename)
	}
}

// Save 保存配置到指定的文件
func (c *configurator) Save(fileComment string) error {

	f, err := os.OpenFile(c.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("写 %q 配置文件: [%w]", c.Filename, err)
	}
	defer f.Close()

	var buf bytes.Buffer
	switch typ := strings.ToLower(c.fileType); typ {
	case ".yaml":
		data, err := yaml.Marshal(c.TargetObj)
		if err != nil {
			return fmt.Errorf("保存 %q 配置文件: [%w]", c.Filename, err)
		}
		buf.WriteString("# " + fileComment + "\n\n")
		buf.Write(data)

	case ".toml":
		data, err := toml.Marshal(c.TargetObj)
		if err != nil {
			return fmt.Errorf("保存 %q 配置文件: [%w]", c.Filename, err)
		}
		buf.WriteString("# " + fileComment + "\n\n")
		buf.Write(data)

	case ".json": // json 不支持注释
		data, err := json.Marshal(c.TargetObj)
		if err != nil {
			return fmt.Errorf("保存 %q 配置文件: [%w]", c.Filename, err)
		}
		buf.Write(data)

	default:
		return fmt.Errorf("不能写入 %q 文件，不支持此文件类型", c.Filename)
	}

	if _, err := buf.WriteTo(f); err != nil {
		return fmt.Errorf("写 %q 配置文件: [%w]", c.Filename, err)
	}
	return nil
}

// Watching 监视配置文件，如果配置文件发生变化，则执行 Load()
// 忽略错误，即如果
func (c *configurator) Watching() {
	go func() {
		var lastModTime time.Time

		for range time.Tick(10 * time.Second) {
			fInfo, err := os.Stat(c.Filename)
			if err != nil {
				Logger.Println(fmt.Errorf("获取 %q 文件信息: [%w]", c.Filename, err))
			}
			newModTime := fInfo.ModTime()

			if !lastModTime.IsZero() && lastModTime != newModTime {
				if err := c.Load(); err != nil {
					Logger.Println(fmt.Errorf("不能重新加载 %q 配置文件: [%w]", c.Filename, err))
				}
			}

			lastModTime = newModTime
		}
	}()
}

// watchSig 监视进程信号(kill -HUP)，如果收到则执行 Load()
func (c *configurator) watchSig() {
	HUP := make(chan os.Signal)
	for {
		signal.Notify(HUP, syscall.SIGHUP)
		<-HUP
		Logger.Printf("获取一个 HUP 进程信号，开始重新加载 %q 文件", c.Filename)

		if err := c.Load(); err != nil {
			Logger.Println(fmt.Errorf("不能重新加载 %q 配置文件: [%w]", c.Filename, err))
		}
	}
}
