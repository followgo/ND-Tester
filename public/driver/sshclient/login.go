package sshclient

import (
	"fmt"

	"golang.org/x/crypto/ssh"

	"github.com/followgo/ND-Tester/public/errors"
)

// OpenAndLogin 连接主机并登录
func (sc *sshClient) OpenAndLogin() (err error) {
	config := &ssh.ClientConfig{
		Timeout: sc.Timeout,
		User:    sc.Username,
		Auth:    []ssh.AuthMethod{},

		// allow any host key to be used (non-prod)
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// HostKeyCallback: ssh.FixedHostKey(sc.Key) // verify host public key

		// optional host key algo list
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
		},
	}

	if sc.Password != "" {
		config.Auth = append(config.Auth, ssh.Password(sc.Password))
	}

	if len(sc.Key) > 0 {
		signer, err := ssh.ParsePrivateKey(sc.Key)
		if err != nil {
			return errors.Wrapf(err, "parse the %q private key", sc.Key)
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	addr := fmt.Sprintf("%s:%d", sc.Host, sc.Port)
	sc.client, err = ssh.Dial("tcp", addr, config)
	if err != nil {
		return errors.Wrap(err, "dial to")
	}

	if err := sc.createSession(); err != nil {
		return errors.Wrap(err, "create session")
	}

	// and wait for prompt
	_, err = sc.readUntilRe(sc.promptRe)
	return errors.Wrap(err, "wait for prompt")
}
