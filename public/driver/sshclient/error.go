package sshclient

import (
	"errors"
)

var (
	ErrReadTimeout  = errors.New("read timeout")
)
