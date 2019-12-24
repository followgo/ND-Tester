package driver

// DriName 名称
type DriName string

const (
	// 驱动名称
	TelnetDriver DriName = "telnet"
	SSHDriver    DriName = "ssh"
	SerialDriver DriName = "serial"
)
