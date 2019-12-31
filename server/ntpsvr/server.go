// 服务端
package ntpsvr

import (
	"fmt"
	"net"
	"time"

	"github.com/followgo/ND-Tester/public/errors"
)

// serverQuit 服务器退出标志
var serverQuit = make(chan bool)

// StartNTPServer 启动 NTP 服务
func StartServer(listenIP string, port uint16) error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", listenIP, port))
	if err != nil {
		return errors.Wrap(err, "resolve udp addr")
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return errors.Wrap(err, "listen udp addr")
	}
	defer conn.Close()

	var data = make([]byte, 512)
	for {
		select {
		case <-serverQuit:
			return nil
		default:
			n, remote, err := conn.ReadFromUDP(data)
			if err == nil {
				go handle(conn, data[:n], remote)
			} else {
				return errors.Wrap(err, "read data from connection")
			}
		}
	}
}

// StopNTPServer 停止时间服务
func StopServer() {
	close(serverQuit)
}

// 消息收发的处理协成
func handle(conn *net.UDPConn, data []byte, remote *net.UDPAddr) {
	// 校验报文大小是否符合预期
	if len(data) < minimumPktSize {
		return
	}

	// 反序列化
	req, err := decodePacket(data)
	if err != nil {
		return
	}
	clientTimestamp := req.TransmitTimestamp

	// 获取版本号
	versionNumber := req.Flags & 0x38

	// 设置参数
	req.Flags = versionNumber | ntpServerMode
	req.PeerClockStratum = 1 // 系统时钟的层数，取值范围为1～16，它定义了时钟的准确度。层数为1的时钟准确度最高，准确度从1到16依次递减，层数为16的时钟处于未同步状态，不能作为参考时钟。
	req.PeerPollingInterval = 5
	req.PeerClockPrecision = 0x100 - 16 // (2^-16 = 0.000015sec)
	req.RootDelay = 0
	req.RootDispersion = 0
	req.ReferenceId = [4]byte{'L', 'O', 'C', 'L'}

	t := toNtpTime(time.Now())
	req.ReferenceTimestamp = t
	req.OriginTimestamp = clientTimestamp
	req.ReceiveTimestamp = t
	req.TransmitTimestamp = t

	// 回显给客户端
	data, err = marshalPacket(req)
	if err != nil {
		return
	}

	_, _ = conn.WriteToUDP(data, remote)
}
