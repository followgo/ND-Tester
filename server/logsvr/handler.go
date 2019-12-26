package logsvr

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/mcuadros/go-syslog.v2/format"

	"github.com/followgo/ND-Tester/config"
)

var (
	// lastMessages 存储最新的消息
	lastMessages     = make([]string, 0, 50)
	lastMessagesLock = new(sync.Mutex)
)

type myHandler struct{}

func (myHandler) Handle(logPart format.LogParts, n int64, e error) {
	if e != nil {
		logrus.WithError(e).Warning("syslog server discarded a bad message")
		return
	}

	var client string
	if v, ok := logPart["client"]; ok {
		client, ok = v.(string)
		if ok {
			client = client[:strings.LastIndexByte(client, ':')]
		}

	} else {
		logrus.Warning("syslog server discarded a bad message without client IP")
		return
	}

	// receive only DUT's IP
	if client != config.Dut.IP {
		logrus.Warning("syslog server discarded a message what it is not DUT's IP")
		return
	}

	jsonData, err := json.Marshal(logPart)
	if err != nil {
		logrus.WithError(err).Errorln("parse log parts to json")
	}

	lastMessagesLock.Lock()
	lastMessages = append(lastMessages, string(jsonData))
	// if len(lastMessages) < lastMessagesMax {
	// 	lastMessages = append(lastMessages, string(jsonData))
	// } else {
	// 	lastMessages = append(lastMessages[1:], string(jsonData))
	// }
	lastMessagesLock.Unlock()

	if sessionWriter != nil {
		_, _ = sessionWriter.Write(jsonData)
		_, _ = sessionWriter.Write([]byte{'\n'})
	}
}

// PourOutLastMessagesString 返回所有的最新的消息，并且清空缓存
func PourOutLastMessagesString() string {
	lastMessagesLock.Lock()
	defer lastMessagesLock.Unlock()

	var bt strings.Builder
	for _, m := range lastMessages {
		bt.WriteString(m)
		bt.WriteByte('\n')
	}

	lastMessages = lastMessages[0:0]
	return bt.String()
}

// PourOutLastMessages 返回所有的最新的消息，并且清空缓存
func PourOutLastMessages() []string {
	lastMessagesLock.Lock()
	defer lastMessagesLock.Unlock()

	var messages = make([]string, len(lastMessages))
	copy(messages, lastMessages)

	lastMessages = lastMessages[0:0]
	return messages
}
