package driver

import (
	"io"
)

// Driver 驱动接口
type Driver interface {
	// OpenAndLogin 连接主机，并且输入用户名和密码，并且等待返回
	// 连接主机后，依次写入 words，因为有些主机的登陆要求是输入用户名和密码，有些的主机只要输入密码。
	OpenAndLogin() error

	// Close 关闭连接
	Close()

	// Write 向 TCP 连接写入数据
	WriteRaw(data []byte) error

	// Write 向 TCP 连接写入数据，自动追加换行符
	Write(data []byte) error

	// Cmd 写入命令，并读取直到匹配 prompt 正则表达式
	Cmd(cmd string) (string, error)

	// ReadAll 读取所有
	ReadAll() (string, error)

	// ReadUntil 读取直到匹配 waitfor 正则表达式
	ReadUntil(waitfor string) (string, error)

	// 设置会话记录器
	SetSessionFile(file string) error
	SetSessionWriter(io.WriteCloser)
}
