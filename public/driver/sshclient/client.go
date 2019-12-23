package sshclient

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/followgo/ND-Tester/public/errors"
)

// callbackPattern 钩子函数
type callbackPattern struct {
	Re *regexp.Regexp
	Cb func()
}

// sshClient ssh 客户端
type sshClient struct {
	Host        string
	Port        uint16
	Username    string
	Password    string
	Key         []byte
	Timeout     time.Duration
	LineBreaks  []byte
	ByeCommands []string // 关闭连接前执行的命令

	// 命令行提示符
	promptRe         *regexp.Regexp
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
		Host:       host,
		Port:       22,
		Username:   username,
		Password:   password,
		Timeout:    5 * time.Second,
		LineBreaks: []byte{'\n'},

		promptRe:         regexp.MustCompile(`(?msi:[\$%#>]$)`),
		callbacks:        make([]callbackPattern, 0, 5),

		stdout: bytes.NewBuffer(nil),
	}
	return sc
}

// SetSessionFile 设置会话记录文件
func (sc *sshClient) SetSessionFile(filename string) (err error) {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return errors.Wrapf(err, "cannot access the directory where the %s file placed", filename)
	}
	sc.sessionWriter, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	return errors.Wrapf(err, "open the %s file", filename)
}

// SetSessionFile 设置会话记录器
func (sc *sshClient) SetSessionWriter(w io.WriteCloser) {
	sc.sessionWriter = w
}

// SetPrompt allows you to change prompt without re-creating ssh client. Default is `(?msi:[\$%#>]$)`
func (sc *sshClient) SetPromptExpr(pattern string) (err error) {
	sc.promptRe, err = regexp.Compile(pattern)
	return errors.Wrapf(err, "compile the %q regular expression", pattern)
}

// Close 关闭连接
func (sc *sshClient) Close()  {
	if sc.session != nil {
		for _, byeCmd := range sc.ByeCommands {
			_ = sc.Write([]byte(byeCmd))
		}

		_ = sc.session.Close()
	}

	if sc.stdin != nil {
		_ = sc.stdin.Close()
	}

	if sc.client != nil {
		_ = sc.client.Close()
	}

	if sc.stdout != nil && sc.sessionWriter != nil {
		_, _ = sc.stdout.WriteTo(sc.sessionWriter)
	}

	if sc.sessionWriter != nil {
		_ = sc.sessionWriter.Close()
	}

	return
}

// RegisterTurnPageCallback registers new callback based on regex string. When current output string matches given
// regex, callback is called. Returns error if regex cannot be compiled.
func (sc *sshClient) RegisterTurnPageCallback(pattern string, callback func()) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return errors.Wrapf(err, "compile the %q regular expression", pattern)
	}

	sc.callbacks = append(sc.callbacks, callbackPattern{
		Cb: callback,
		Re: re,
	})

	return nil
}
