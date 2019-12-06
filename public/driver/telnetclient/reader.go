package telnetclient

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ReadUntil reads tcp stream from server until 'waitfor' regex matches.
// Returns gathered output and error, if any.
// Any escape sequences are cutted out during reading for providing clean output for parsing/reading.
func (c *telnetClient) ReadUntil(waitfor string) (s string, err error) {
	waitForRe, err := regexp.Compile(waitfor)
	if err != nil {
		return "", fmt.Errorf("cannot compile 'waitfor' regexp: [%w]", err)
	}
	return c.readUntilRe(waitForRe)
}

// ReadAll read all from tcp stream
func (c *telnetClient) ReadAll() (s string, err error) {
	return c.readUntilRe(nil)
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
		if waitForRe == nil && err == ErrReadTimeout { // no timeout if read all
			err = nil
		}

		buf.Write(lastLine.Bytes())
		s := buf.String()

		if c.sessionWriter != nil {
			_, _ = c.sessionWriter.Write([]byte(s))
		}

		return s, err
	}

	timeout := c.Timeout
	if waitForRe == nil {
		timeout = time.Second
	}
	after := time.After(timeout) // 如果遇到翻页，则重现赋值

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
					return returnFn(fmt.Errorf("error while reading escape sequence: [%w]", err))
				}
				seq = append(seq, b1)

				if b1 == cmdSB { // subNegotiation
					// read all until subNeg. end.
					for {
						bn, err := c.readByte()
						if err != nil {
							return returnFn(fmt.Errorf("error while reading escape subnegotiation sequence: [%w]", err))
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
						return returnFn(fmt.Errorf("error while reading IAC sequence: [%w]", err))
					}
					seq = append(seq, bn)
				}

				// Sequence finished, do something with it:
				if err := c.negotiate(seq); err != nil {
					return returnFn(fmt.Errorf("failed to negotiate connection: [%w]", err))
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
		}
	}
}

// read one byte from tcp stream
func (c *telnetClient) readByte() (b byte, err error) {
	if err := c.conn.SetReadDeadline(time.Now().Add(c.Timeout)); err != nil {
		return 0, err
	}
	data := make([]byte, 1, 1)
	_, err = c.conn.Read(data)
	return data[0], err
}
