package sshclient

import (
	"golang.org/x/crypto/ssh"

	"github.com/followgo/ND-Tester/public/errors"
)

// createSession 创建会话，并且开始 shell
func (sc *sshClient) createSession() (err error) {
	sc.session, err = sc.client.NewSession()
	if err != nil {
		return errors.Wrap(err, "new session")
	}

	// stdin, stdout, stderr
	if sc.stdin, err = sc.session.StdinPipe(); err != nil {
		return errors.Wrap(err, "get the stdin pipe")
	}
	sc.session.Stdout = sc.stdout
	sc.session.Stderr = sc.stdout

	// Request pseudo terminal
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echo
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := sc.session.RequestPty("xterm", 500, 200, modes); err != nil {
		return errors.Wrap(err, "set terminal mode")
	}

	// Start remote shell
	err = sc.session.Shell()
	return errors.Wrap(err, "start shell")
}
