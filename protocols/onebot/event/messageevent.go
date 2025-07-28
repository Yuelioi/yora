package event

import (
	"strconv"
	"sync"
	"yora/internal/event"
	"yora/protocols/onebot/message"

	basemsg "yora/internal/message"
)

var _ event.MessageEvent = (*MessageEvent)(nil)

type MessageEvent struct {
	*Event
	messageCache basemsg.Message
	once         sync.Once
}

func (e *Event) UserID() string {
	return strconv.Itoa(e.UserIDInt)
}
func (e *Event) ChatID() string {
	return strconv.Itoa(e.GroupIDInt)
}

// Extra implements event.MessageEvent.
func (m *MessageEvent) Extra() map[string]any {
	panic("unimplemented")
}

func (m *MessageEvent) IsGroup() bool {
	return m.MessageType == "group"
}

func (m *MessageEvent) IsPrivate() bool {
	return m.MessageType == "private"
}
func (m *MessageEvent) Message() basemsg.Message {
	m.once.Do(func() {
		m.messageCache = message.New(m.MessageValue)
	})
	return m.messageCache
}

func (m *MessageEvent) MessageID() string {
	return strconv.Itoa(m.MessageIDInt)
}

func (m *MessageEvent) RawMessage() string {
	return m.RawMessageValue
}

// TODO 获取reply里的at
func (m *MessageEvent) ReplyTo() string {
	for _, seg := range m.Message().Segments() {
		if seg.Type() == "at" {
			if qq, ok := seg.GetData("qq"); ok {
				return qq.(string)
			}
		}
	}
	return ""
}

// Sender implements event.MessageEvent.
func (m *MessageEvent) Sender() basemsg.Sender {
	return m.SenderValue
}

// UserID implements event.MessageEvent.
func (m *MessageEvent) UserID() string {
	return strconv.Itoa(m.UserIDInt)
}
