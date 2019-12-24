package config

// Host 运行测试框架的的默认主机配置
var Host = HostConfig{
	SerialPort: "com1",
}

// HostConfig 运行测试框架的主机配置
type HostConfig struct {
	SerialPort string `json:"serial_port" yaml:"serial_port" tome:"serial_port"`
}
