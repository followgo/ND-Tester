package telnetclient

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/followgo/ND-Tester/public/errors"
)

// ReadAll read all from tcp stream
func (c *telnetClient) ReadAll() (s string, err error) {
	s, err = c.readUntilRe(nil)
	return s, errors.Wrap(err, "read all from tcp stream")
}

// ReadUntil reads tcp stream from server until 'waitfor' regex matches.
// Returns gathered output and error, if any.
// Any escape sequences are cutted out during reading for providing clean output for parsing/reading.
func (c *telnetClient) ReadUntil(waitfor string) (s string, err error) {
	waitForRe, err := regexp.Compile(waitfor)
	if err != nil {
		return "", errors.Wrapf(err, "compile the %q regular expression", waitfor)
	}
	return c.readUntilRe(waitForRe)
}

// ReadUntilRe reads tcp stream from server until 'waitForRe' regex.Regexp
// if waitForRe is nil, read all from tcp stream
// Returns gathered output and error, if any.
// Any escape sequences are cutted out during reading for providing clean output for parsing/reading.
func (c *telnetClient) readUntilRe(waitForRe *regexp.Regexp) (s string, err error) {
	var (
		buf        bytes.Buffer
		lastLine   bytes.Buffer
		inSequence bool // 转义符
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

		if c.sessionWriter != nil {
			_, _ = c.sessionWriter.Write([]byte(s))
		}

		return s, errors.Wrap(err, "read until a regular expression")
	}

	timeout := c.Timeout
	if waitForRe == nil {
		timeout = time.Second
	}
	after := time.After(timeout) // 翻页操作会重新赋值

	LOOP:
	for {
		select {
		case <-after:
			return returnFn(ErrReadTimeout)
		default:
			b, err := c.readByte()
			if err != nil {
				return returnFn(err)
			}

			// catch escape sequences
			if b == cmdIAC {
				seq := []byte{b}

				b1, err := c.readByte()
				if err != nil {
					return returnFn(err)
				}
				seq = append(seq, b1)

				if b1 == cmdSB { // subNegotiation
					// read all until subNeg. end.
					for {
						bn, err := c.readByte()
						if err != nil {
							return returnFn(err)
						}
						seq = append(seq, bn)
						if bn == cmdSE {
							break
						}
					}
				} else {
					// not subsequence.
					bn, err := c.readByte()
					if err != nil {
						return returnFn(err)
					}
					seq = append(seq, bn)
				}

				// Sequence finished, do something with it:
				if err := c.negotiate(seq); err != nil {
					return returnFn(errors.Wrap(err, "negotiation failed"))
				}
			}

			// cut out escape sequences
			if b == 27 {
				inSequence = true
				continue
			}
			if inSequence {
				// 2) 0-?, @-~, ' ' - / === 48-63, 32-47, finish with 64-126
				if b == 91 {
					continue
				}
				if b == 20 {
					continue
				}
				if b >= 32 && b <= 63 {
					// just skip it
					continue
				}
				if b >= 64 && b <= 126 {
					// finish sequence
					inSequence = false
					lastLine.Reset()
					continue
				}
			}

			// not IAC sequence, but IAC char =\
			if b == cmdIAC {
				continue
			}

			// remove \r ; remove backspaces
			if b == 8 {
				if lastLine.Len() > 0 {
					lastLine.Truncate(lastLine.Len() - 1)
				}
				continue
			}
			if b == '\r' {
				continue
			}

			if b != '\n' {
				lastLine.Write([]byte{b})
			}

			// check for regex matching. Execute callback if matched.
			if len(c.callbacks) > 0 {
				for i := range c.callbacks {
					if c.callbacks[i].Re.Match(bytes.TrimSpace(lastLine.Bytes())) {
						after = time.After(timeout)
						c.callbacks[i].Cb()
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

			continue LOOP
		}
	}
}

// read one byte from tcp stream
func (c *telnetClient) readByte() (b byte, err error) {
	if err := c.conn.SetReadDeadline(time.Now().Add(c.Timeout)); err != nil {
		return 0, errors.Wrap(err, "set read deadline")
	}
	data := make([]byte, 1, 1)
	_, err = c.conn.Read(data)

	return data[0], errors.Wrap(err, "read a byte from TCP stream")
}
