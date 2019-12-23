# sshclient SSH Client Package

Simple ssh client lib, written in golang.

## Example

访问某一台 Linux 主机

```go
func main() {
	s := sshclient.New("192.168.0.15", "root", "firstmile")
	_ = s.SetPromptExpr(`#$`)
	s.ByeCommands = []string{"logout"}
	if err := s.SetSessionFile("session.txt"); err != nil {
		logrus.Println(err)
	}
	s.Timeout = 10 * time.Second

	if err := s.OpenAndLogin(); err != nil {
		logrus.Fatalf("连接并登陆: %s", err)
	}
	defer s.Close()

	if s, err := s.Cmd("cat /proc/cpuinfo"); err != nil {
		logrus.Fatalf("查看CPU信息: %s", err)
	} else {
		fmt.Println(s)
	}

	if s, err := s.Cmd("ls /"); err != nil {
		logrus.Fatalf("查看根目录文件: %s", err)
	} else {
		fmt.Println(s)
	}
}
```