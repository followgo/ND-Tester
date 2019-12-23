package sshclient

import (
	"time"

	"github.com/followgo/ND-Tester/public/errors"
)

// Cmd sends command and returns output
func (sc *sshClient) Cmd(cmd string) (s string, err error) {
	if err := sc.Write([]byte(cmd)); err != nil {
		return "",  errors.Wrapf(err, "write the %q command", cmd)
	}
	return sc.readUntilRe(sc.promptRe)
}

// Write is the same as WriteRaw, but adds CRLF to given string
func (sc *sshClient) Write(bytes []byte) error {
	bytes = append(bytes, sc.LineBreaks...)

	if sc.sessionWriter!=nil{
		sc.sessionWriter.Write( bytes)
	}

	return sc.WriteRaw(bytes)
}

// WriteRaw writes raw bytes to tcp connection
func (sc *sshClient) WriteRaw(bytes []byte) (err error) {
	_, err = sc.stdin.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "write bytes to TCP connection")
	}
	time.Sleep(time.Second)
	return nil
}
