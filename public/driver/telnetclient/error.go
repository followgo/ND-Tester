package telnetclient

import (
	"errors"
)

var (
	ErrReadTimeout  = errors.New("read timeout")
	ErrWriteTimeout = errors.New("write timeout")
)