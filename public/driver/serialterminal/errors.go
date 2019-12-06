package serialterminal

import (
	"errors"
)

var (
	ErrLoginTimeout = errors.New("login timeout")
	ErrReadTimeout  = errors.New("read timeout")
	ErrWriteTimeout = errors.New("write timeout")
)
