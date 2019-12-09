package serialterminal

import (
	"bytes"
	"regexp"
	"strings"
	"time"

	"github.com/followgo/ND-Tester/public/errors"
)

// ReadAll read all from tcp stream
func (st *serialTerminal) ReadAll() (s string, err error) {
	s, err = st.readUntilRe(nil)
	return s, errors.Wrap(err, "读取所有数据")
}

// ReadUntil reads tcp stream from server until 'waitfor' regex matches.
// Returns gathered output and error, if any.
// Any escape sequences are cutted out during reading for providing clean output for parsing/reading.
func (st *serialTerminal) ReadUntil(waitfor string) (s string, err error) {
	waitForRe, err := regexp.Compile(waitfor)
	if err != nil {
		return "", errors.Wrapf(err, "编译 %q", waitfor)
	}
	s, err = st.readUntilRe(waitForRe)
	return s, errors.Wrapf(err, "读取直到匹配 %q", waitfor)
}

// ReadUntilRe reads tcp stream from server until 'waitForRe' regex.Regexp
// if waitForRe is nil, read all from tcp stream
// Returns gathered output and error, if any.
// Any escape sequences are cutted out during reading for providing clean output for parsing/reading.
func (st *serialTerminal) readUntilRe(waitForRe *regexp.Regexp) (s string, err error) {
	var (
		buf        bytes.Buffer
		lastLine   bytes.Buffer
		inSequence bool
	)

	// 读取完成后的操作
	var returnFn = func(err error) (string, error) {
		buf.Write(lastLine.Bytes())
		s := buf.String()

		if st.sessionWriter != nil {
			_, _ = st.sessionWriter.Write([]byte(s))
		}

		return s, errors.Wrapf(err, "读取直到匹配 %q 表达式", waitForRe.String())
	}

	timeout := st.Timeout
	if waitForRe == nil {
		timeout = time.Second
	}
	after := time.After(timeout) // 如果遇到翻页，则重现赋值

	for {
		select {
		case <-after:
			return returnFn(nil)
		default:
			b, err := st.readByte()
			if err != nil {
				return returnFn(err)
			}

			// 忽略掉 FirstMile 交换机翻页后出现的长空白
			if b == 0x0 {
				inSequence = true
				continue
			}
			if inSequence == true {
				switch b {
				case 0xd, 0x20:
					continue
				default:
					inSequence = false
					lastLine.Reset()
				}
			}

			if b == '\r' || b == 0x1b {
				continue
			}

			if b != '\n' {
				lastLine.Write([]byte{b})
			}

			// check for regex matching. Execute callback if matched.
			if len(st.callbacks) > 0 {
				for i := range st.callbacks {
					if st.callbacks[i].Re.Match(bytes.TrimSpace(lastLine.Bytes())) {
						after = time.After(timeout)
						st.callbacks[i].Cb()
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

// read one byte from port
func (st *serialTerminal) readByte() (b byte, err error) {
	data := make([]byte, 1, 1)
	_, err = st.p.Read(data)
	return data[0], errors.Wrap(err, "读取一个字节")
}
