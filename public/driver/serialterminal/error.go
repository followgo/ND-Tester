package serialterminal

import (
	"github.com/goburrow/serial"

	"github.com/followgo/ND-Tester/public/errors"
)

var (
	errLoginTimeout = errors.New("login timeout")
)

// IsTimeout 返回一个布尔值，当错误包含 serial.ErrTimeout 则返回 true, 或者返回假
func IsTimeout(err error) bool {
	return errors.Is(err, serial.ErrTimeout) || errors.Is(err, errLoginTimeout)
}
