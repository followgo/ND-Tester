package config

// Dut 待测设备的默认配置，配置加载程序会重置值
var Dut = DutConfig{
	IP:       "",
	Username: "superuser",
	Password: "123",
	Telnet: dutTelnetConfig{
		Port:                  23,
		TimeoutMs:             5000,
		Linebreak:             "\n",
		ByeCommands:           []string{"exit"},
		PromptPattern:         `(?msi:[\$%#>]$)`,
		LoginPromptPattern:    `(?msi:user(\s)?name\:(\s+)?$)`,
		PasswordPromptPattern: `(?msi:password\:$)`,
	},
	Serial: dutSerialConfig{
		BaudRate:              115200,
		DataBits:              8,
		StopBits:              1,
		Parity:                "N",
		TimeoutMs:             5000,
		Linebreak:             "\n",
		ByeCommands:           []string{"exit"},
		PromptPattern:         `(?msi:[\$%#>]$)`,
		LoginPromptPattern:    `(?msi:user(\s)?name\:(\s+)?$)`,
		PasswordPromptPattern: `(?msi:password\:$)`,
	},
	SSH: dutSSHConfig{
		Port:          23,
		TimeoutMs:     5000,
		Linebreak:     "\n",
		ByeCommands:   []string{"exit"},
		PromptPattern: `(?msi:[\$%#>]$)`,
	},
}

// DutConfig 待测设备的配置
type DutConfig struct {
	IP       string          `json:"ip" yaml:"ip" toml:"ip"`
	Username string          `json:"username" yaml:"username" toml:"username"`
	Password string          `json:"password" yaml:"password" toml:"password"`
	Telnet   dutTelnetConfig `json:"telnet" yaml:"telnet" toml:"telnet"`
	Serial   dutSerialConfig `json:"serial" yaml:"serial" toml:"serial"`
	SSH      dutSSHConfig    `json:"ssh" yaml:"ssh" toml:"ssh"`
}

// dutTelnetConfig 待测设备的 Telnet 接口配置
type dutTelnetConfig struct {
	Port uint16 `json:"port" yaml:"port" toml:"port"`

	Username    string   `json:"username" yaml:"username" toml:"username"` // 如果为空则使用主用户名和密码
	Password    string   `json:"password" yaml:"password" toml:"password"`
	TimeoutMs   int64    `json:"timeout_ms" yaml:"timeout_ms" toml:"timeout_ms"`
	Linebreak   string   `json:"linebreak" yaml:"linebreak" toml:"linebreak"`          // 写入数据的换行
	ByeCommands []string `json:"bye_commands" yaml:"bye_commands" toml:"bye_commands"` // 断开连接前要执行的命令

	PromptPattern         string            `json:"prompt_pattern" yaml:"prompt_pattern" toml:"prompt_pattern"`
	LoginPromptPattern    string            `json:"login_prompt_pattern" yaml:"login_prompt_pattern" toml:"login_prompt_pattern"`
	PasswordPromptPattern string            `json:"password_prompt_pattern" yaml:"password_prompt_pattern" toml:"password_prompt_pattern"`
	CallbackPatterns      []callbackPattern `json:"callback_patterns" yaml:"callback_patterns" toml:"callback_patterns"`
}

// dutSerialConfig 待测设备的串口配置
type dutSerialConfig struct {
	BaudRate int    `json:"baud_rate" yaml:"baud_rate" toml:"baud_rate"` // Baud rate (default 115200)
	DataBits int    `json:"data_bits" yaml:"data_bits" toml:"data_bits"` // Data bits: 5, 6, 7 or 8 (default 8)
	StopBits int    `json:"stop_bits" yaml:"stop_bits" toml:"stop_bits"` // Stop bits: 1 or 2 (default 1)
	Parity   string `json:"parity" yaml:"parity" toml:"parity"`          // Parity: N - None, E - Even, O - Odd (default N)

	Username    string   `json:"username" yaml:"username" toml:"username"` // 如果为空则使用主用户名和密码
	Password    string   `json:"password" yaml:"password" toml:"password"`
	TimeoutMs   int64    `json:"timeout_ms" yaml:"timeout_ms" toml:"timeout_ms"`
	Linebreak   string   `json:"linebreak" yaml:"linebreak" toml:"linebreak"`          // 写入数据的换行
	ByeCommands []string `json:"bye_commands" yaml:"bye_commands" toml:"bye_commands"` // 断开连接前要执行的命令

	PromptPattern         string            `json:"prompt_pattern" yaml:"prompt_pattern" toml:"prompt_pattern"`
	LoginPromptPattern    string            `json:"login_prompt_pattern" yaml:"login_prompt_pattern" toml:"login_prompt_pattern"`
	PasswordPromptPattern string            `json:"password_prompt_pattern" yaml:"password_prompt_pattern" toml:"password_prompt_pattern"`
	CallbackPatterns      []callbackPattern `json:"callback_patterns" yaml:"callback_patterns" toml:"callback_patterns"`
}

// dutSSHConfig 待测设备的 SSH 接口配置
type dutSSHConfig struct {
	Port uint16 `json:"port" yaml:"port" toml:"port"`

	Username       string   `json:"username" yaml:"username" toml:"username"` // 如果为空则使用主用户名和密码
	Password       string   `json:"password" yaml:"password" toml:"password"`
	PrivateKeyFile string   `json:"private_key_file" yaml:"private_key_file" toml:"private_key_file"`
	TimeoutMs      int64    `json:"timeout_ms" yaml:"timeout_ms" toml:"timeout_ms"`
	Linebreak      string   `json:"linebreak" yaml:"linebreak" toml:"linebreak"`          // 写入数据的换行
	ByeCommands    []string `json:"bye_commands" yaml:"bye_commands" toml:"bye_commands"` // 断开连接前要执行的命令

	PromptPattern    string            `json:"prompt_pattern" yaml:"prompt_pattern" toml:"prompt_pattern"`
	CallbackPatterns []callbackPattern `json:"callback_patterns" yaml:"callback_patterns" toml:"callback_patterns"`
}

// callbackPattern 命令行接口的回调图式
// 当输出字符串匹配 MatchingPattern 正则表达式时，回调函数会向接口写入 EnterChars 字符
// application scenario: 遇到翻页提示字符串，让回调函数输入翻页键
type callbackPattern struct {
	MatchingPattern string `json:"matching_pattern" yaml:"matching_pattern" toml:"matching_pattern"`
	EnterChars      []byte `json:"enter_chars" yaml:"enter_chars" toml:"enter_chars"`
}
