package driver

// DriName 名称
type DriName string

const (
	// 换行符和回车符
	CR = '\r'
	LF = '\n'

	// 默认端口
	DefaultPortForTelnet = 23
	DefaultPortForSSH    = 22

	// 驱动名称
	TelnetDriver DriName = "telnet"
	SSHDriver    DriName = "ssh"
	SerialDriver DriName = "serial"
)
