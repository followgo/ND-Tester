package driver

import (
	"io"
	"time"
)

// Driver 驱动接口
type Driver interface {
	// Login 连接主机，并且输入用户名和密码，并且等待返回
	// 连接主机后，依次写入 words，因为有些主机的登陆要求是输入用户名和密码，有些的主机只要输入密码。
	Login(words ...string) error

	// Close 关闭连接
	Close() error

	// Write 向 TCP 连接写入数据
	// 写入前不清空回显的数据流
	Write(data []byte) error

	// Cmd 写入命令，并且等待返回（直到数据流匹配 prompt 正则表达式为止）
	// 每次写入前清空回显的数据流
	// 如果超时，则先执行 ResetWriteBuffer()，然后返回。
	Cmd(cmd string, timeout time.Duration) (string, error)

	// ResetWriteBuffer 复位写入数据流的 Buffer
	// 对于某些主机，如果写入的命令是无效的，是不会清空写入 buffer 的。
	ResetWriteBuffer() error

	// ReadUntil 读取数据流，直到 waitfor 或 prompt 正则表达式匹配为止
	// 如果 waitfor 为空，则只匹配 prompt 正则表达式
	ReadUntil(waitfor string, timeout time.Duration) (string, error)

	// SetLogSession 设置会话记录器，所有的回显数据流都会写入到这里
	SetLogSessionWriter(w io.WriteCloser)

	// SetPromptPattern 设置命令提示符的 regex 匹配字符串
	SetPromptPattern(pattern string)

	// SetPromptPattern 设置写入换行符
	SetLineBreak(pattern string)

	// SetTimeout 设置读写操作的超时时间
	SetTimeout(dur time.Duration)

	// RegisterMatchingCallback 注册一个回调函数，使用 regex 匹配进行触发。当回显的输出匹配到 pattern 则执行
	// 一般应用于屏显翻页
	RegisterMatchingCallback(pattern string, callback func()) error
}
