package telnetclient

import (
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// callbackPattern 钩子函数
type callbackPattern struct {
	Re *regexp.Regexp
	Cb func()
}

// telnetClient Telnet 客户端数据结构
type telnetClient struct {
	Host     string
	Port     uint16
	Username string
	Password string
	Timeout  time.Duration

	// 命令行提示符
	promptRe         *regexp.Regexp
	loginPromptRe    *regexp.Regexp
	passwordPromptRe *regexp.Regexp

	// conn is TCP socket object
	conn net.Conn

	sessionWriter io.WriteCloser
	callbacks     []callbackPattern
}

// New 创建一个 telnetClient 实例
func New(host, username, password string) *telnetClient {
	c := &telnetClient{
		Host:     host,
		Port:     23,
		Username: username,
		Password: password,
		Timeout:  5 * time.Second,

		promptRe:         regexp.MustCompile(`(?msi:[\$%#>]$)`),
		loginPromptRe:    regexp.MustCompile(`[Uu]ser(\s)?[Nn]ame\:(\s+)?$`),
		passwordPromptRe: regexp.MustCompile(`[Pp]ass[Ww]ord\:$`),
		callbacks:        nil,
	}
	return c
}

// SetSessionFile 设置会话记录文件
func (c *telnetClient) SetSessionFile(filename string) (err error) {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	c.sessionWriter, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	return err
}

// SetSessionFile 设置会话记录器
func (c *telnetClient) SetSessionWriter(w io.WriteCloser) {
	c.sessionWriter = w
}

// SetPrompt allows you to change prompt without re-creating ssh client. Default is `(?msi:[\$%#>]$)`
func (c *telnetClient) SetPromptRe(pattern string) (err error) {
	c.promptRe, err = regexp.Compile(pattern)
	return err
}

// SetLoginPrompt sets custom login prompt. Default is `[Uu]ser(\s)?[Nn]ame\:(\s+)?$`
func (c *telnetClient) SetLoginPromptRe(pattern string) (err error) {
	c.loginPromptRe, err = regexp.Compile(pattern)
	return err
}

// SetPasswordPrompt sets custom password prompt. Default is `[Pp]ass[Ww]ord\:$`
func (c *telnetClient) SetPasswordPromptRe(pattern string) (err error) {
	c.passwordPromptRe, err = regexp.Compile(pattern)
	return err
}

// Close closes telnet connection.
func (c *telnetClient) Close() (err error) {
	if c.conn != nil {
		if _, err := c.readUntilRe(nil); err != nil {
			return err
		}

		if c.sessionWriter != nil {
			_ = c.sessionWriter.Close()
		}

		return c.conn.Close()
	}

	return nil
}

// RegisterTurnPageCallback registers new callback based on regex string. When current output string matches given
// regex, callback is called. Returns error if regex cannot be compiled.
func (c *telnetClient) RegisterTurnPageCallback(pattern string, callback func()) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	c.callbacks = append(c.callbacks, callbackPattern{
		Cb: callback,
		Re: re,
	})

	return nil
}
