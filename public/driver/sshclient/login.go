package sshclient

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

// DialAndLogin 连接主机并登录
func (sc *sshClient) DialAndLogin() (err error) {
	config := &ssh.ClientConfig{
		Timeout: sc.Timeout,
		User:    sc.Username,
		Auth:    []ssh.AuthMethod{},
		// HostKeyCallback: ssh.FixedHostKey(sc.getHostKey(sc.Host)),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
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
			return fmt.Errorf("不能解析私有密钥: [%w]", err)
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	addr := fmt.Sprintf("%s:%d", sc.Host, sc.Port)
	sc.client, err = ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}

	return sc.createSession()
}

// getHostKey parse OpenSSH known_hosts file
// ssh or use ssh-keyscan to get initial key
func (sc *sshClient) getHostKey(host string) ssh.PublicKey {
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey found for %s", host)
	}

	return hostKey
}
