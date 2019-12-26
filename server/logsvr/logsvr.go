package logsvr

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/mcuadros/go-syslog.v2"

	"github.com/followgo/ND-Tester/config"
	"github.com/followgo/ND-Tester/public/errors"
)

var (
	server        *syslog.Server
	sessionWriter io.WriteCloser
)

func RunWithSessionFile(file string) error {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	sessionWriter = f
	return Run()
}

func RunWithSessionWriter(w io.WriteCloser) error {
	sessionWriter = w
	return Run()
}

// Run 启动 Syslog 服务，只接收 DUT 的日志消息
func Run() (err error) {
	server = syslog.NewServer()

	switch config.Dut.Syslog.Format {
	case "RFC3164":
		server.SetFormat(syslog.RFC3164)
	case "RFC5424":
		server.SetFormat(syslog.RFC5424)
	case "RFC6587":
		server.SetFormat(syslog.RFC6587)
	default:
		server.SetFormat(syslog.Automatic)
	}

	server.SetHandler(new(myHandler))

	if config.Dut.Syslog.ProtocolType == "tcp" {
		err = server.ListenTCP(fmt.Sprintf(":%d", config.Dut.Syslog.Port))
	} else {
		err = server.ListenUDP(fmt.Sprintf(":%d", config.Dut.Syslog.Port))
	}
	if err != nil {
		return errors.Wrap(err, "listening...")
	}

	if err := server.Boot(); err != nil {
		return errors.Wrap(err, "start syslog service")
	}

	return nil
}

// Stop 停止 Syslog 服务
func Stop() {
	if server != nil {
		_ = server.Kill()
	}

	if sessionWriter != nil {
		_ = sessionWriter.Close()
	}
}

// Wait 等待直到服务器停止
func Wait() { server.Wait() }
