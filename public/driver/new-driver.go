package driver

import (
	"io"
	"io/ioutil"
	"time"

	"github.com/followgo/ND-Tester/config"
	"github.com/followgo/ND-Tester/public/driver/serialterminal"
	"github.com/followgo/ND-Tester/public/driver/sshclient"
	"github.com/followgo/ND-Tester/public/driver/telnetclient"
	"github.com/followgo/ND-Tester/public/errors"
)

// NewDriver 返回一个 Driver 接口
func NewDriver(name DriName) (Driver, error) {
	username, password := config.Dut.Username, config.Dut.Password

	switch name {
	case TelnetDriver:
		if config.Dut.Telnet.Username != "" {
			username = config.Dut.Telnet.Username
			password = config.Dut.Telnet.Password
		}
		dri := telnetclient.New(config.Dut.IP, username, password)
		dri.Port = config.Dut.Telnet.Port
		dri.Timeout = time.Duration(config.Dut.Telnet.TimeoutMs) * time.Millisecond
		dri.LineBreaks = []byte(config.Dut.Telnet.Linebreak)
		dri.ByeCommands = config.Dut.Telnet.ByeCommands

		if err := dri.SetPromptExpr(config.Dut.Telnet.PromptPattern); err != nil {
			return nil, errors.Wrap(err, "设置 telnet 的命令行提示符")
		}
		if err := dri.SetLoginPromptExpr(config.Dut.Telnet.LoginPromptPattern); err != nil {
			return nil, errors.Wrap(err, "设置 telnet 的登陆提示符")
		}
		if err := dri.SetPasswordPromptExpr(config.Dut.Telnet.PasswordPromptPattern); err != nil {
			return nil, errors.Wrap(err, "设置 telnet 的密码提示符")
		}

		for _, cb := range config.Dut.Telnet.CallbackPatterns {
			if err := dri.RegisterTurnPageCallback(cb.MatchingPattern, func() {
				_ = dri.WriteRaw(cb.EnterChars)
			}); err != nil {
				return nil, errors.Wrap(err, "注册 telnet 回调函数")
			}
		}

		return dri, nil

	case SerialDriver:
		dri := serialterminal.New(config.Host.SerialPort, username, password, config.Dut.Serial.BaudRate)
		dri.DataBits = config.Dut.Serial.DataBits
		dri.Parity = config.Dut.Serial.Parity
		dri.StopBits = config.Dut.Serial.StopBits
		dri.Timeout = time.Duration(config.Dut.Serial.TimeoutMs) * time.Millisecond
		dri.LineBreaks = []byte(config.Dut.Serial.Linebreak)
		dri.ByeCommands = config.Dut.Serial.ByeCommands

		if err := dri.SetPromptExpr(config.Dut.Serial.PromptPattern); err != nil {
			return nil, errors.Wrap(err, "设置 serial 的命令行提示符")
		}
		if err := dri.SetLoginPromptExpr(config.Dut.Serial.LoginPromptPattern); err != nil {
			return nil, errors.Wrap(err, "设置 serial 的命令行提示符")
		}
		if err := dri.SetPasswordPromptExpr(config.Dut.Serial.PasswordPromptPattern); err != nil {
			return nil, errors.Wrap(err, "设置 serial 的命令行提示符")
		}

		for _, cb := range config.Dut.Serial.CallbackPatterns {
			if err := dri.RegisterTurnPageCallback(cb.MatchingPattern, func() {
				_ = dri.WriteRaw(cb.EnterChars)
			}); err != nil {
				return nil, errors.Wrap(err, "注册 serial 回调函数")
			}
		}

		return dri, nil

	case SSHDriver:
		dri := sshclient.New(config.Dut.IP, username, password)
		if config.Dut.SSH.PrivateKeyFile != "" {
			Key, err := ioutil.ReadFile(config.Dut.SSH.PrivateKeyFile)
			if err != nil {
				return nil, errors.Wrap(err, "读 SSH 的私有密钥文件")
			}
			dri.Key = Key
		}
		dri.Port = config.Dut.SSH.Port
		dri.Timeout = time.Duration(config.Dut.SSH.TimeoutMs) * time.Millisecond
		dri.LineBreaks = []byte(config.Dut.SSH.Linebreak)
		dri.ByeCommands = config.Dut.SSH.ByeCommands

		if err := dri.SetPromptExpr(config.Dut.SSH.PromptPattern); err != nil {
			return nil, errors.Wrap(err, "设置 serial 的命令行提示符")
		}

		for _, cb := range config.Dut.SSH.CallbackPatterns {
			if err := dri.RegisterTurnPageCallback(cb.MatchingPattern, func() {
				_ = dri.WriteRaw(cb.EnterChars)
			}); err != nil {
				return nil, errors.Wrap(err, "注册 SSH 回调函数")
			}
		}

		return dri, nil
	}

	return nil, errors.New("undefined")
}

// NewDriverWithSessionFile 带会话记录
func NewDriverWithSessionFile(name DriName, file string) (Driver, error) {
	dri, err := NewDriver(name)
	if err != nil {
		return nil, err
	}
	err = dri.SetSessionFile(file)
	return dri, errors.Wrap(err, "设置 driver 的会话记录文件")
}

// NewDriverWithSessionWriter 带会话记录
func NewDriverWithSessionWriter(name DriName, w io.WriteCloser) (Driver, error) {
	dri, err := NewDriver(name)
	if err != nil {
		return nil, err
	}
	dri.SetSessionWriter(w)
	return dri, nil
}
