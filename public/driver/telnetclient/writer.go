package telnetclient

import (
	"fmt"
	"io"
	"time"
)

// Cmd sends command and returns output
func (c *telnetClient) Cmd(cmd string) (s string, err error) {
	if err := c.Write([]byte(cmd)); err != nil {
		return "", err
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
		return ErrWriteTimeout
	}

	_, err = c.conn.Write(bytes)
	if err != nil && err != io.ErrShortWrite {
		return fmt.Errorf("failed to writeRaw(): [%w]", err)
	}
	return nil
}
