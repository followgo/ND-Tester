package driver

// 换行符和回车符
const (
	CR = byte('\r')
	LF = byte('\n')
)

// 默认端口
const (
	DefaultPortForTelnet = 23
	DefaultPortForSSH    = 22
)

// 驱动名称
const (
	DriveNameOfTelnet = "telnet"
	DriveNameOfSSH    = "ssh"
	DriveNameOfSerial = "serial"
)