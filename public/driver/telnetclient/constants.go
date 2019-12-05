package telnetclient

/*
	telnet的命令格式
	-----------------------
	  IAC | 命令码 | 选项码 |
	-----------------------
*/
// telnet的命令
const (
	cmdEOF   = 236 // 文件结束符
	cmdSUSP  = 237 // 挂起当前进程（作业控制）
	cmdABORT = 238 // 异常中止进程
	cmdEOR   = 239 // 记录结束符
	cmdSE    = 240 // 子选项结束
	cmdNOP   = 241 // 空操作
	cmdDM    = 242 // 数据标记
	cmdBRK   = 243 // 终止符（break）
	cmdIP    = 244 // 终止进程
	cmdAO    = 245 // 异常中止输出
	cmdAYT   = 246 // 请求应答
	cmdEC    = 247 // 终止符
	cmdEL    = 248 // 删除行
	cmdGA    = 249 // 继续进行
	cmdSB    = 250 // 子选项开始
	cmdWILL  = 251 // 选项协商: 同意启动（enable）选项
	cmdWONT  = 252 // 选项协商: 拒绝启动选项
	cmdDO    = 253 // 选项协商: 认可选项请求
	cmdDONT  = 254 // 选项协商: 拒绝选项请求
	cmdIAC   = 255 // 字符0XFF
)

/*
	选项协商：4种请求
	1）WILL：发送方本身将激活选项
	2）DO：发送方想叫接受端激活选项
	3）WONT：发送方本身想禁止选项
	4）DONT：发送方想让接受端去禁止选项
	发送者	接收者	说明
	WILL    DO      发送者想激活某选项，接受者接收该选项请求
	WILL    DONT    发送者想激活某选项，接受者拒绝该选项请求
	DO      WILL    发送者希望接收者激活某选项，接受者接受该请求
	DO      DONT    发送者希望接收6者激活某选项，接受者拒绝该请求
	WONT    DONT    发送者希望使某选项无效，接受者必须接受该请求
	DONT    WONT    发送者希望对方使某选项无效，接受者必须接受该请求
*/
// telnet 选项码
const (
	optBIN      = 0  // 二进制传输
	optECHO     = 1  // 回显
	optREC      = 2  // 重连
	optSGA      = 3  // 抑制继续进行
	optSTATUS   = 5  // 状态
	optTIMER    = 6  // 时钟标识
	optLOG      = 18 // Logout
	optTERMTYPE = 24 // 终端类型
	optWINSIZE  = 31 // 窗口大小
	optTSP      = 32 // 终端速度 Terminal Speed
	optRFC      = 33 // 远程流量控制 Remote Flow Control
	optLINEMODE = 34 // 行方式
	optENVVAR   = 36 // 环境变量

	// opt_SB_SEND SEND subneg
	opt_SB_SEND = 1
	// TELOPT_SB_IS IS subneg
	opt_SB_IS = 0
	// TELOPT_SB_NEV_ENVIRON
	opt_SB_NEV_ENVIRON = 39
)
