package sshclient

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/followgo/ND-Tester/public/errors"
)

// ReadAll read all from tcp stream
func (sc *sshClient) ReadAll() (s string, err error) {
	return sc.readUntilRe(nil)
}

// ReadUntil reads tcp stream from server until 'waitfor' regex matches.
// Returns gathered output and error, if any.
// Any escape sequences are cutted out during reading for providing clean output for parsing/reading.
func (sc *sshClient) ReadUntil(waitfor string) (s string, err error) {
	waitForRe, err := regexp.Compile(waitfor)
	if err != nil {
		return "", errors.Wrapf(err, "compile the %q regular expression", waitfor)
	}
	return sc.readUntilRe(waitForRe)
}

// ReadUntilRe reads tcp stream from server until 'waitForRe' regex.Regexp
// if waitForRe is nil, read all from tcp stream
// Returns gathered output and error, if any.
// Any escape sequences are cutted out during reading for providing clean output for parsing/reading.
// TODO: 首迈交换机会持续打印 0x07，但是使用putty和SecureCRT不会这样。
func (sc *sshClient) readUntilRe(waitForRe *regexp.Regexp) (s string, err error) {
	var (
		buf      bytes.Buffer
		lastLine bytes.Buffer
	)

	// 函数返回的操作
	var returnFn = func(err error) (string, error) {
		if errors.Is(err, io.EOF) {
			err = nil
		}
		if waitForRe == nil && errors.Is(err, ErrReadTimeout) { // no timeout if read all
			err = nil
		}

		buf.Write(lastLine.Bytes())
		s := buf.String()

		if sc.sessionWriter != nil {
			_, _ = sc.sessionWriter.Write([]byte(s))
		}

		return s, errors.Wrapf(err, "read until the %q regular expression", waitForRe.String())
	}

	timeout := sc.Timeout
	if waitForRe == nil {
		timeout = time.Second
	}
	after := time.After(timeout) // 如果遇到翻页，则重现赋值

	for {
		select {
		case <-after:
			return returnFn(ErrReadTimeout)
		default:
			b, err := sc.readByte()
			if err != nil {
				return returnFn(err)
			}

			if b == '\r' {
				continue
			}

			if b != '\n' {
				lastLine.Write([]byte{b})
			}

			// check for regex matching. Execute callback if matched.
			if len(sc.callbacks) > 0 {
				for i := range sc.callbacks {
					if sc.callbacks[i].Re.Match(bytes.TrimSpace(lastLine.Bytes())) {
						after = time.After(timeout)
						sc.callbacks[i].Cb()
						lastLine.Reset()
					}
				}
			}

			// check for CRLF.
			// We need last line to compare with prompt.
			// if b == '\n'
			if b == '\n' {
				if line := lastLine.String(); strings.TrimSpace(line) != "" {
					buf.WriteString(strings.TrimRight(line, " ") + "\n")
				}
				lastLine.Reset()
			}

			// After reading, we should check for regexp every time.
			// Unfortunately, we cant wait only CRLF, because prompt usually comes without CRLF.
			if waitForRe != nil {
				if waitForRe.Match(lastLine.Bytes()) {
					return returnFn(nil)
				}
			}
		}
	}
}

// read one byte from tcp stream
func (sc *sshClient) readByte() (b byte, err error) {
	b, err = sc.stdout.ReadByte()
	return b, errors.Wrap(err, "read a byte from TCP stream")
}
