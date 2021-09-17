package notice

import (
	"github.com/Luna-CY/v2ray-helper/common/util"
	"sync"
	"time"
)

const (
	MessageTypeInfo = iota + 1
	MessageTypeWarn
	MessageTypeError
)

var m *manager

func GetManager() *manager {
	if nil == m {
		m = new(manager)
		m.notice = []*Notice{}
	}

	return m
}

type manager struct {
	nm     sync.RWMutex
	notice []*Notice
}

func (m *manager) Push(messageType int, title, message string) {
	m.nm.Lock()
	defer m.nm.Unlock()

	if 5 <= len(m.notice) {
		m.notice = m.notice[1:5]
	}

	m.notice = append(m.notice, &Notice{Id: util.GenerateRandomString(8), Type: messageType, Title: title, Message: message, Time: time.Now().Unix()})
}

func (m *manager) GetAll() []*Notice {
	m.nm.RLock()
	defer m.nm.RUnlock()

	return m.notice
}

func (m *manager) Clean() {
	m.nm.Lock()
	defer m.nm.Unlock()

	m.notice = []*Notice{}
}

type Notice struct {
	Id      string `json:"id"`
	Type    int    `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
	IsRead  bool   `json:"is_read"`
}
