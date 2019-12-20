package telnetclient

import (
	"fmt"
	"net"
)

// OpenAndLogin Open tcp connection and login
func (c *telnetClient) OpenAndLogin() (err error) {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	c.conn, err = net.DialTimeout("tcp", addr, c.Timeout)
	if err != nil {
		return err
	}

	// 设置窗口为最大，减少翻页次数
	if err := c.WriteRaw([]byte{cmdIAC, cmdSB, optWINSIZE, 0x00, 0xfe, 0x00, 0xfe, cmdIAC, cmdSE}); err != nil {
		return err
	}

	return c.login()
}

// Login is a simple wrapper for login/password auth
func (c *telnetClient) login() (err error) {
	if c.Username != "" {
		if _, err := c.readUntilRe(c.loginPromptRe); err != nil {
			return err
		}

		if err := c.Write([]byte(c.Username)); err != nil {
			return err
		}
	}

	if c.Password != "" {
		if _, err := c.readUntilRe(c.passwordPromptRe); err != nil {
			return err
		}

		if err := c.Write([]byte(c.Password)); err != nil {
			return err
		}
	}

	// and wait for prompt
	_, err = c.readUntilRe(c.promptRe)
	return err
}