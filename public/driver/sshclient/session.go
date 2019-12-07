package sshclient

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

// createSession 创建会话，并且开始 shell
func (sc *sshClient) createSession() (err error) {
	sc.session, err = sc.client.NewSession()
	if err != nil {
		return err
	}

	// stdin, stdout, stderr
	if sc.stdin, err = sc.session.StdinPipe(); err != nil {
		return err
	}
	sc.session.Stdout = sc.stdout
	sc.session.Stderr = sc.stdout

	// Request pseudo terminal
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echo
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := sc.session.RequestPty("xterm", 500, 200, modes); err != nil {
		return fmt.Errorf("request for pseudo terminal failed: [%w]", err)
	}

	// Start remote shell
	if err := sc.session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: [%w]", err)
	}

	return nil
}
