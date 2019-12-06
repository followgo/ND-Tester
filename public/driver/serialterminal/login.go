package serialterminal

import (
	"errors"
	"time"

	"github.com/goburrow/serial"
)

// OpenAndLogin Open tcp connection and login
func (st *serialTerminal) Open() (err error) {
	st.p, err = serial.Open(&serial.Config{
		Address:  st.PortName,
		BaudRate: st.BaudRate,
		DataBits: st.DataBits,
		StopBits: st.StopBits,
		Parity:   st.Parity,
		Timeout:  st.Timeout,
	})
	return err
}

// Login is a simple wrapper for login/password auth
func (st *serialTerminal) Login() (err error) {
LOOP:
	for {
		select {
		case <-time.After(st.Timeout * 5):
			return ErrLoginTimeout
		default:
			if st.Username != "" {
				if _, err := st.readUntilRe(st.loginPromptRe); err != nil {
					if errors.As(err, &ErrReadTimeout) {
						if err := st.Write([]byte(" ")); err != nil {
							return err
						}
						continue LOOP
					}
					return err
				}

				if err := st.Write([]byte(st.Username)); err != nil {
					return err
				}
			}

			if st.Password != "" {
				if _, err := st.readUntilRe(st.passwordPromptRe); err != nil {
					if errors.As(err, &ErrReadTimeout) {
						if err := st.Write([]byte(" ")); err != nil {
							return err
						}
						continue LOOP
					}
					return err
				}

				if err := st.Write([]byte(st.Password)); err != nil {
					return err
				}
			}

			// and wait for prompt
			_, err = st.readUntilRe(st.promptRe)
			return err
		}
	}
}
