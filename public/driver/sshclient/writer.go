package sshclient

import (
	"fmt"
)

// Cmd sends command and returns output
func (sc *sshClient) Cmd(cmd string) (s string, err error) {
	if err := sc.Write([]byte(cmd)); err != nil {
		return "", err
	}
	return sc.readUntilRe(sc.promptRe)
}

// Write is the same as WriteRaw, but adds CRLF to given string
func (sc *sshClient) Write(bytes []byte) error {
	bytes = append(bytes, '\n')
	return sc.WriteRaw(bytes)
}

// WriteRaw writes raw bytes to tcp connection
func (sc *sshClient) WriteRaw(bytes []byte) (err error) {
	_, err = sc.stdin.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to writeRaw(): [%w]", err)
	}
	return nil
}
