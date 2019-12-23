package serialterminal

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/goburrow/serial"

	"github.com/followgo/ND-Tester/public/errors"
)

// callbackPattern 钩子函数
type callbackPattern struct {
	Re *regexp.Regexp
	Cb func()
}

type serialTerminal struct {
	PortName string
	BaudRate int    // Baud rate (default 115200)
	DataBits int    // Data bits: 5, 6, 7 or 8 (default 8)
	StopBits int    // Stop bits: 1 or 2 (default 1)
	Parity   string // Parity: N - None, E - Even, O - Odd (default N)

	Timeout time.Duration

	Username    string
	Password    string
	LineBreaks  []byte
	ByeCommands []string // 关闭连接前执行的命令

	// 命令行提示符
	promptRe         *regexp.Regexp
	loginPromptRe    *regexp.Regexp
	passwordPromptRe *regexp.Regexp

	p             serial.Port
	sessionWriter io.WriteCloser
	callbacks     []callbackPattern
}

// New 创建一个 serialTerminal 实例
func New(portName, username, password string, baudRate int) *serialTerminal {
	st := &serialTerminal{
		PortName:         portName,
		BaudRate:         baudRate,
		DataBits:         8,
		Parity:           "N",
		StopBits:         1,
		Timeout:          5 * time.Second,
		Username:         username,
		Password:         password,
		LineBreaks:       []byte{'\n'},
		promptRe:         regexp.MustCompile(`(?msi:[\$%#>]$)`),
		loginPromptRe:    regexp.MustCompile(`[Uu]ser(\s)?[Nn]ame\:(\s+)?$`),
		passwordPromptRe: regexp.MustCompile(`[Pp]ass[Ww]ord\:$`),
		callbacks:        make([]callbackPattern, 0, 5),
	}
	return st
}

// SetSessionFile 设置会话记录文件
func (st *serialTerminal) SetSessionFile(filename string) (err error) {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return errors.Wrapf(err, "cannot access the directory where the %s file placed", filename)
	}
	st.sessionWriter, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	return errors.Wrapf(err, "open the %s file", filename)
}

// SetSessionFile 设置会话记录器
func (st *serialTerminal) SetSessionWriter(w io.WriteCloser) {
	st.sessionWriter = w
}

// SetPrompt allows you to change prompt without re-creating ssh client. Default is `(?msi:[\$%#>]$)`
func (st *serialTerminal) SetPromptExpr(pattern string) (err error) {
	st.promptRe, err = regexp.Compile(pattern)
	return errors.Wrapf(err, "compile the %q regular expression", pattern)
}

// SetLoginPrompt sets custom login prompt. Default is `[Uu]ser(\s)?[Nn]ame\:(\s+)?$`
func (st *serialTerminal) SetLoginPromptExpr(pattern string) (err error) {
	st.loginPromptRe, err = regexp.Compile(pattern)
	return errors.Wrapf(err, "compile the %q regular expression", pattern)
}

// SetPasswordPrompt sets custom password prompt. Default is `[Pp]ass[Ww]ord\:$`
func (st *serialTerminal) SetPasswordPromptExpr(pattern string) (err error) {
	st.passwordPromptRe, err = regexp.Compile(pattern)
	return errors.Wrapf(err, "compile the %q regular expression", pattern)
}

// Close closes telnet connection.
func (st *serialTerminal) Close() {
	if st.p != nil {
		for _, byeCmd := range st.ByeCommands {
			_ = st.Write([]byte(byeCmd))
		}

		_, _ = st.ReadAll()
		_ = st.p.Close()
	}

	if st.sessionWriter != nil {
		_ = st.sessionWriter.Close()
	}

	return
}

// RegisterTurnPageCallback registers new callback based on regex string. When current output string matches given
// regex, callback is called. Returns error if regex cannot be compiled.
func (st *serialTerminal) RegisterTurnPageCallback(pattern string, callback func()) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return errors.Wrapf(err, "compile the %q regular expression", pattern)
	}

	st.callbacks = append(st.callbacks, callbackPattern{
		Cb: callback,
		Re: re,
	})

	return nil
}
