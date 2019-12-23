package telnetclient

import (
	"io"
	"time"

	"github.com/followgo/ND-Tester/public/errors"
)

// Cmd sends command and returns output
func (c *telnetClient) Cmd(cmd string) (s string, err error) {
	if err := c.Write([]byte(cmd)); err != nil {
		return "", errors.Wrapf(err, "write the %q command", cmd)
	}
	return c.readUntilRe(c.promptRe)
}

// Write is the same as WriteRaw, but adds CRLF to given string
func (c *telnetClient) Write(bytes []byte) error {
	bytes = append(bytes, c.LineBreaks...)
	return c.WriteRaw(bytes)
}

// WriteRaw writes raw bytes to tcp connection
func (c *telnetClient) WriteRaw(bytes []byte) (err error) {
	err = c.conn.SetWriteDeadline(time.Now().Add(c.Timeout))
	if err != nil {
		return errors.Wrap(err, "set write deadline")
	}

	_, err = c.conn.Write(bytes)
	if err != nil && err != io.ErrShortWrite {
		return errors.Wrap(err, "write bytes to TCP connection")
	}
	return nil
}
