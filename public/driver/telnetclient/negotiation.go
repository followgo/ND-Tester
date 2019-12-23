package telnetclient

import (
	"github.com/followgo/ND-Tester/public/errors"
)

// negotiate 协商
// ignore: TELNET_AYT, TELNET_AO
func (c *telnetClient) negotiate(sequence []byte) (err error) {
	if len(sequence) < 3 {
		return nil // do nothing
	}

	if len(sequence) == 3 {
		switch sequence[1] {
		case cmdDO: // DO sequence
			err = c.WriteRaw([]byte{cmdIAC, cmdWILL, sequence[2]})
		case cmdWONT: // WONT -> DONT
			err = c.WriteRaw([]byte{cmdIAC, cmdDONT, sequence[2]})
		case cmdWILL:
			err = c.WriteRaw([]byte{cmdIAC, cmdDO, sequence[2]})
		}
		if err != nil {
			return errors.Wrapf(err, "answer to the %x command", sequence[1])
		}
	}

	// subSeq SEND request
	if len(sequence) == 6 && sequence[1] == cmdSB && sequence[3] == opt_SB_SEND {
		// what to send?
		switch sequence[2] {
		case optTERMTYPE:
			// set terminal to xterm
			err = c.WriteRaw([]byte{cmdIAC, cmdSB, optTERMTYPE, opt_SB_IS, 'X', 'T', 'E', 'R', 'M', cmdIAC, cmdSE})
			break
		case optWINSIZE:
			// set terminal's window size
			err = c.WriteRaw([]byte{cmdIAC, cmdSB, optWINSIZE, opt_SB_IS, 0xfe, opt_SB_IS, 0xfe, cmdIAC, cmdSE})
			break
		case opt_SB_NEV_ENVIRON:
			// send new-env -> is new env
			err = c.WriteRaw([]byte{cmdIAC, cmdSB, opt_SB_NEV_ENVIRON, opt_SB_IS, cmdIAC, cmdSE})
			break
		default:
			// accept all
			err = c.WriteRaw([]byte{cmdIAC, cmdSB, sequence[2], 0, cmdIAC, cmdSE})
			break
		}

		if err != nil {
			return errors.Wrapf(err, "answer to the %x sub command", sequence[2])
		}
	}

	return nil
}
