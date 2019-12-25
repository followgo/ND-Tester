package config

// Host 运行测试框架的的默认主机配置
var Host = HostConfig{
	SerialPort: "com1",
	Logger: hostLoggerConfig{
		BaseFile:      "./logs/main.log",
		Level:         "INFO",
		OverWrite:     false,
		OutputConsole: true,
	},
}

// HostConfig 运行测试框架的主机配置
type HostConfig struct {
	SerialPort string           `json:"serial_port" yaml:"serial_port" toml:"serial_port"`
	Logger     hostLoggerConfig `json:"logger" yaml:"logger" toml:"logger"`
}

// hostLoggerConfig 日志记录器
type hostLoggerConfig struct {
	BaseFile      string `json:"base_file" yaml:"base_file" toml:"base_file"`
	Level         string `json:"level" yaml:"level" toml:"level"`
	OverWrite     bool   `json:"over_write" yaml:"over_write" toml:"over_write"`
	OutputConsole bool   `json:"output_console" yaml:"output_console" toml:"output_console"`
}
