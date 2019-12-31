package ntpsvr

import (
	"bytes"
	"encoding/binary"
	"time"
)

const (
	// DefaultPort 默认的 NTP 端口
	DefaultPort = 123

	// ntpVersion NTP版本号
	ntpVersion = 4

	// minimumPktSize 最小的报文长度
	minimumPktSize = 48

	// ntpServerMode 服务端模式 （主动对等体=1，被动对等体=2，服务端模式=4，客户端模式=3，广播/组播模式=5）
	ntpServerMode uint8 = 4

	// ntpClientMode 客户端模式
	ntpClientMode uint8 = 3
)

// ntpPacket 报文结构
// NTP packet format (v3 with optional v4 fields removed)
//
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |LI | VN  |Mode |    Stratum     |     Poll      |  Precision   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                         Root Delay                            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                         Root Dispersion                       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                          Reference ID                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// +                     Reference Timestamp (64)                  +
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// +                      Origin Timestamp (64)                    +
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// +                      Receive Timestamp (64)                   +
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// +                      Transmit Timestamp (64)                  +
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
type ntpPacket struct {
	Flags               uint8   // leap indicator + version number + mode
	PeerClockStratum    uint8   // 系统时钟的层数
	PeerPollingInterval uint8   // 轮询间隔
	PeerClockPrecision  uint8   // 系统时钟的精度
	RootDelay           uint32  // 本地到主参考时钟源的往返时间
	RootDispersion      uint32  // 系统时钟相对于主参考时钟的最大误差
	ReferenceId         [4]byte // 参考时钟源的标识
	ReferenceTimestamp  uint64  // 系统时钟最后一次被设定或更新的时间
	OriginTimestamp     uint64  // NTP请求报文离开发送端时发送端的本地时间
	ReceiveTimestamp    uint64  // NTP请求报文到达接收端时接收端的本地时间
	TransmitTimestamp   uint64  // 应答报文离开应答者时应答者的本地时间
}

var (
	// offset 时间偏移值
	offset = int64(time.Unix(0, 0).Sub(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) / time.Second)
)

// toNtpTime converts the time.Time value t into its 64-bit fixed-point
func toNtpTime(t time.Time) uint64 {
	ret := uint64(t.Unix()+offset) << 32
	ret |= uint64(float64(t.Nanosecond()+1) / 1e9 * (1 << 32))
	return ret
}

// fromNtpTime interprets the fixed-point ntpTime as an absolute time and returns the corresponding time.Time value.
func fromNtpTime(ntpTime uint64) time.Time {
	return time.Unix(int64(ntpTime>>32)-offset, int64(float64(ntpTime&0xFFFFFFFF)/(1<<32)*1e9))
}

// marshalPacket 序列化
func marshalPacket(req ntpPacket) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, req); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// decodePacket 反序列化
func decodePacket(buf []byte) (rsp ntpPacket, err error) {
	err = binary.Read(bytes.NewReader(buf), binary.BigEndian, &rsp)
	return
}
