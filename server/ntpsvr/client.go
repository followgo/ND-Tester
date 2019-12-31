package ntpsvr

import (
	"fmt"
	"net"
	"time"

	"github.com/followgo/ND-Tester/public/errors"
)

// GetTime 获取从NTP服务器获取时间
func GetTime(ip string, port uint16, timeout time.Duration) (t time.Time, err error) {
	rAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return time.Time{}, errors.Wrap(err, "resolve addr")
	}

	conn, err := net.DialUDP("udp", nil, rAddr)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "dial to server")
	}
	defer conn.Close()

	t1 := time.Now()
	t1NtpTime := toNtpTime(t1)
	req := ntpPacket{
		Flags:             0x23,
		TransmitTimestamp: t1NtpTime,
	}
	data, err := marshalPacket(req)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "marshal packet")
	}

	if err := conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return time.Time{}, errors.Wrap(err, "set write deadline to connection")
	}
	if _, err := conn.Write(data); err != nil {
		return time.Time{}, errors.Wrap(err, "write request packet to connection")
	}

	// Receive the response.
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return time.Time{}, errors.Wrap(err, "set read deadline to connection")
	}
	var receData = make([]byte, 256)
	n, rAddr, err := conn.ReadFromUDP(receData)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "read from connection")
	}

	// 校验报文大小是否符合预期
	if n < minimumPktSize {
		return time.Time{}, errors.New("the packet size is too small")
	}

	// 反序列化
	req, err = decodePacket(receData[:n])
	if err != nil {
		return time.Time{}, err
	}

	if req.TransmitTimestamp == 0x0 {
		return time.Time{}, errors.New("invalid transmit time in response")
	}
	if req.OriginTimestamp != t1NtpTime {
		return time.Time{}, errors.New("server response mismatch")
	}

	// rtt 时间
	rtt := time.Now().Sub(t1)
	serverTime := fromNtpTime(req.TransmitTimestamp)
	return serverTime.Add(rtt), err
}
