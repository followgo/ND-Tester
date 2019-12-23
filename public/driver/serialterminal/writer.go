package serialterminal

import (
	"time"

	"github.com/followgo/ND-Tester/public/errors"
)

// Cmd sends command and returns output
func (st *serialTerminal) Cmd(cmd string) (s string, err error) {
	if err := st.Write([]byte(cmd)); err != nil {
		return "", errors.Wrapf(err, "write the %q command", cmd)
	}
	return st.readUntilRe(st.promptRe)
}

// Write is the same as WriteRaw, but adds CRLF to given string
func (st *serialTerminal) Write(bytes []byte) error {
	time.Sleep(100 * time.Millisecond)
	bytes = append(bytes, st.LineBreaks...)
	return st.WriteRaw(bytes)
}

// WriteRaw writes raw bytes to port
func (st *serialTerminal) WriteRaw(bytes []byte) (err error) {
	for i := range bytes {
		if _, err := st.p.Write([]byte{bytes[i]}); err != nil {
			return errors.Wrap(err, "write a byte to port")
		}
		time.Sleep(30 * time.Millisecond)
	}
	return nil
}
