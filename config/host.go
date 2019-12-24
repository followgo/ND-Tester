package config

// Host 运行测试框架的的默认主机配置
var Host = HostConfig{
	SerialPort: "com1",
	Logger: hostLoggerConfig{
		BaseFile:      "./logs/main.log",
		Level:         "INFO",
		MaxMegaSize:   100,
		MaxBackups:    7,
		MaxAgeDays:    7,
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
	MaxMegaSize   int    `json:"max_mega_size" yaml:"max_mega_size" toml:"max_mega_size"`
	MaxBackups    int    `json:"max_backups" yaml:"max_backups" toml:"max_backups"`
	MaxAgeDays    int    `json:"max_age_days" yaml:"max_age_days" toml:"max_age_days"`
	OutputConsole bool   `json:"output_console" yaml:"output_console" toml:"output_console"`
}
