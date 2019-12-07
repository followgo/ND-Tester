package sshclient

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"golang.org/x/crypto/ssh"
)

// callbackPattern 钩子函数
type callbackPattern struct {
	Re *regexp.Regexp
	Cb func()
}

// sshClient ssh 客户端
type sshClient struct {
	Host     string
	Port     uint16
	Username string
	Password string
	Key      []byte
	Timeout  time.Duration

	// 命令行提示符
	promptRe         *regexp.Regexp
	loginPromptRe    *regexp.Regexp
	passwordPromptRe *regexp.Regexp
	callbacks        []callbackPattern

	client  *ssh.Client
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  *bytes.Buffer // stderr write to stdout
	// stderr        *bytes.Buffer
	sessionWriter io.WriteCloser
}

// New 创建一个 ssh client 实例
func New(host, username, password string) *sshClient {
	sc := &sshClient{
		Host:     host,
		Port:     23,
		Username: username,
		Password: password,
		Timeout:  5 * time.Second,

		promptRe:         regexp.MustCompile(`(?msi:[\$%#>]$)`),
		loginPromptRe:    regexp.MustCompile(`[Uu]ser(\s)?[Nn]ame\:(\s+)?$`),
		passwordPromptRe: regexp.MustCompile(`[Pp]ass[Ww]ord\:$`),
		callbacks:        make([]callbackPattern, 0, 5),

		stdout: bytes.NewBuffer(nil),
	}
	return sc
}

// SetSessionFile 设置会话记录文件
func (sc *sshClient) SetSessionFile(filename string) (err error) {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	sc.sessionWriter, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	return err
}

// SetSessionFile 设置会话记录器
func (sc *sshClient) SetSessionWriter(w io.WriteCloser) {
	sc.sessionWriter = w
}

// SetPrompt allows you to change prompt without re-creating ssh client. Default is `(?msi:[\$%#>]$)`
func (sc *sshClient) SetPromptExpr(pattern string) (err error) {
	sc.promptRe, err = regexp.Compile(pattern)
	return err
}

// SetLoginPrompt sets custom login prompt. Default is `[Uu]ser(\s)?[Nn]ame\:(\s+)?$`
func (sc *sshClient) SetLoginPromptExpr(pattern string) (err error) {
	sc.loginPromptRe, err = regexp.Compile(pattern)
	return err
}

// SetPasswordPrompt sets custom password prompt. Default is `[Pp]ass[Ww]ord\:$`
func (sc *sshClient) SetPasswordPromptExpr(pattern string) (err error) {
	sc.passwordPromptRe, err = regexp.Compile(pattern)
	return err
}

// Close 关闭连接
func (sc *sshClient) Close() error {
	if sc.session != nil {
		if err := sc.session.Close(); err != nil {
			return err
		}
	}

	if sc.stdin != nil {
		_ = sc.stdin.Close()
	}

	if sc.client != nil {
		if err := sc.client.Close(); err != nil {
			return err
		}
	}

	if sc.stdout != nil && sc.sessionWriter != nil {
		if _, err := sc.stdout.WriteTo(sc.sessionWriter); err != nil {
			return err
		}
		if sc.sessionWriter != nil {
			_ = sc.sessionWriter.Close()
		}
	}

	return nil
}

// RegisterTurnPageCallback registers new callback based on regex string. When current output string matches given
// regex, callback is called. Returns error if regex cannot be compiled.
func (sc *sshClient) RegisterTurnPageCallback(pattern string, callback func()) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	sc.callbacks = append(sc.callbacks, callbackPattern{
		Cb: callback,
		Re: re,
	})

	return nil
}
