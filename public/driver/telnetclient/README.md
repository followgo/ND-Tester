# telnetclient Telnet Client Package

Simple telnet client lib, written in golang.

## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/followgo/ND-Tester/public/driver/telnetclient"
)

func main() {
	// SC3700 只有密码，没有用户名
	c := telnetclient.New("192.168.118.1", "", "firstmile")
	// 设置提示符: <SC3700>, [SC3700]
	if err := c.SetPromptExpr(`(?msi:[\]>]$)`); err != nil {
		log.Fatalln(err)
	}
	// 设置会话记录文件
	if err := c.SetSessionFile("session.txt"); err != nil {
		log.Fatalln(err)
	}
	// 华为 SC3700 翻页提示符 ---- More ----，按空格键下一页
	if err := c.RegisterTurnPageCallback(`----\sMore\s----`, func() {
		_ = c.WriteRaw([]byte(" "))
	}); err != nil {
		log.Fatalln(err)
	}
	// 连接并登陆
	if err := c.DialAndLogin(); err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	// 输入命令，并打印回显
	s, err := c.Cmd("display version")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("系统信息: \n%s\n", s)

	s, err = c.Cmd("display current-configuration")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("所有配置: \n%s\n", s)
}
```