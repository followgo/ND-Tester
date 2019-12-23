package telnetclient

import (
	"fmt"
	"net"

	"github.com/followgo/ND-Tester/public/errors"
)

// OpenAndLogin Open tcp connection and login
func (c *telnetClient) OpenAndLogin() (err error) {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	c.conn, err = net.DialTimeout("tcp", addr, c.Timeout)
	if err != nil {
		return errors.Wrap(err, "dial to")
	}

	// 设置窗口为最大，减少翻页次数
	if err := c.WriteRaw([]byte{cmdIAC, cmdSB, optWINSIZE, 0x00, 0xfe, 0x00, 0xfe, cmdIAC, cmdSE}); err != nil {
		return errors.Wrap(err, "set window size")
	}

	return c.login()
}

// Login is a simple wrapper for login/password auth
func (c *telnetClient) login() (err error) {
	if c.Username != "" {
		if _, err := c.readUntilRe(c.loginPromptRe); err != nil {
			return errors.Wrap(err, "wait for login prompt")
		}

		if err := c.Write([]byte(c.Username)); err != nil {
			return errors.Wrap(err, "write username")
		}
	}

	if c.Password != "" {
		if _, err := c.readUntilRe(c.passwordPromptRe); err != nil {
			return errors.Wrap(err, "wait for password prompt")
		}

		if err := c.Write([]byte(c.Password)); err != nil {
			return errors.Wrap(err, "write password")
		}
	}

	// and wait for prompt
	_, err = c.readUntilRe(c.promptRe)
	return errors.Wrap(err, "wait for prompt")
}
