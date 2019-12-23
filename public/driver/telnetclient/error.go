package telnetclient

import (
	"github.com/followgo/ND-Tester/public/errors"
)

var (
	ErrReadTimeout  = errors.New("read timeout")
)