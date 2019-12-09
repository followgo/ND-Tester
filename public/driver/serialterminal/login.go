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
		return errors.Wrap(err, "打开端口")
	}

	return st.Login()
}

// Login is a simple wrapper for login/password auth
func (st *serialTerminal) Login() (err error) {
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
							return errors.Wrap(err, "按回车")
						}
						continue LOOP
					}
					return errors.Wrap(err, "等待用户名提示符")
				}

				if err := st.Write([]byte(st.Username)); err != nil {
					return errors.Wrap(err, "写入登陆用户名")
				}
			}

			if st.Password != "" {
				if _, err := st.readUntilRe(st.passwordPromptRe); err != nil {
					if IsTimeout(err) {
						if err := st.WriteRaw([]byte{'\n'}); err != nil {
							return errors.Wrap(err, "按回车")
						}
						continue LOOP
					}
					return errors.Wrap(err, "等待密码提示符")
				}

				if err := st.Write([]byte(st.Password)); err != nil {
					return errors.Wrap(err, "写入登陆密码")
				}
			}

			// and wait for prompt
			_, err = st.readUntilRe(st.promptRe)
			return errors.Wrap(err, "等待提示符")
		}
	}
}
