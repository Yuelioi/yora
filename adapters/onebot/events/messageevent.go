package events

import (
	"strconv"
	"sync"
	"yora/adapters/onebot/messages"
	"yora/pkg/event"

	"yora/pkg/message"
)

var _ event.MessageEvent = (*MessageEvent)(nil)

type MessageEvent struct {
	*Event
	messageCache message.Message
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
func (m *MessageEvent) Message() message.Message {
	m.once.Do(func() {
		m.messageCache = messages.New(m.MessageValue)
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
func (m *MessageEvent) Sender() message.Sender {
	return m.SenderValue
}

// UserID implements event.MessageEvent.
func (m *MessageEvent) UserID() string {
	return strconv.Itoa(m.UserIDInt)
}
