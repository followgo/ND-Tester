package serialterminal

import (
	"time"

	"github.com/goburrow/serial"

	"github.com/followgo/ND-Tester/public/errors"
)

// OpenAndLogin Open tcp connection and login
func (st *serialTerminal) OpenAndLogin() (err error) {
	st.p, err = serial.Open(&serial.Config{
		Address:  st.PortName,
		BaudRate: st.BaudRate,
		DataBits: st.DataBits,
		StopBits: st.StopBits,
		Parity:   st.Parity,
		Timeout:  st.Timeout,
	})
	if err != nil {
		return errors.Wrap(err, "open serial port")
	}

	return st.login()
}

// Login is a simple wrapper for login/password auth
func (st *serialTerminal) login() (err error) {
LOOP:
	for {
		select {
		case <-time.After(st.Timeout * 5):
			return errLoginTimeout
		default:
			if st.Username != "" {
				if _, err := st.readUntilRe(st.loginPromptRe); err != nil {
					if IsTimeout(err) {
						if err := st.WriteRaw([]byte{'\n'}); err != nil {
							return errors.Wrap(err, "write enter key")
						}
						continue LOOP
					}
					return errors.Wrap(err, "wait for login prompt")
				}

				if err := st.Write([]byte(st.Username)); err != nil {
					return errors.Wrap(err, "write username")
				}
			}

			if st.Password != "" {
				if _, err := st.readUntilRe(st.passwordPromptRe); err != nil {
					if IsTimeout(err) {
						if err := st.WriteRaw([]byte{'\n'}); err != nil {
							return errors.Wrap(err, "write enter key")
						}
						continue LOOP
					}
					return errors.Wrap(err, "wait for password prompt")
				}

				if err := st.Write([]byte(st.Password)); err != nil {
					return errors.Wrap(err, "write password")
				}
			}

			// and wait for prompt
			_, err = st.readUntilRe(st.promptRe)
			return errors.Wrap(err, "wait for prompt")
		}
	}
}
