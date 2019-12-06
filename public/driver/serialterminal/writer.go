package serialterminal

import (
	"time"
)

// Cmd sends command and returns output
func (st *serialTerminal) Cmd(cmd string) (s string, err error) {
	if err := st.Write([]byte(cmd)); err != nil {
		return "", err
	}
	return st.readUntilRe(st.promptRe)
}

// Write is the same as WriteRaw, but adds CRLF to given string
func (st *serialTerminal) Write(bytes []byte) error {
	bytes = append(bytes, '\n')
	time.Sleep(100 * time.Millisecond)
	return st.WriteRaw(bytes)
}

// WriteRaw writes raw bytes to port
func (st *serialTerminal) WriteRaw(bytes []byte) (err error) {
	for i := range bytes {
		if _, err := st.p.Write([]byte{bytes[i]}); err != nil {
			return err
		}
		time.Sleep(30 * time.Millisecond)
	}
	return nil
}
