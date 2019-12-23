# serialterminal Serial Terminal Package

Simple serial lib, written in golang.

## Example

访问博达工业交换机例子

```go
func main() {
	s := serialterminal.New("com1", "admin", "admin", 9600)
	if err := s.SetSessionFile("session.txt"); err != nil {
		logrus.Fatalln(err)
	}
	s.ByeCommands = []string{"exit"}

	// 设置翻页钩子
	err := s.RegisterTurnPageCallback(`--More--`, func() {
		_ = s.WriteRaw([]byte(" "))
	})
	if err != nil {
		logrus.Println("设置翻页钩子", err)
	}

	// 打开端口和登陆
	if err := s.OpenAndLogin(); err != nil {
		logrus.Fatalln(err)
	}
	defer s.Close()

	if _, err := s.Cmd("enable"); err != nil {
		logrus.Fatalf("进入特权模式: %s", err.Error())
	}

	if s, err := s.Cmd("show running-config"); err != nil {
		logrus.Fatalf("打印运行配置: %s", err.Error())
	} else {
		fmt.Println(s)
	}

	if _, err := s.Cmd("exit"); err != nil {
		logrus.Fatalf("退出特权模式: %s", err.Error())
	}
}
```