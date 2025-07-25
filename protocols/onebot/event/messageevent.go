package event

import (
	"strconv"
	"yora/internal/event"
	"yora/protocols/onebot/message"
)

var _ event.MessageEvent = (*MessageEvent)(nil)

type MessageEvent struct {
	*Event
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

// IsGroup implements event.MessageEvent.
func (m *MessageEvent) IsGroup() bool {
	return m.SubTypeValue == "group"
}

// IsPrivate implements event.MessageEvent.
func (m *MessageEvent) IsPrivate() bool {
	return m.SubTypeValue == "private"
}

func (e *Event) Message() event.Message {
	return message.New(e.MessageValue)
}

func (e *Event) MessageID() string {
	return strconv.Itoa(e.MessageIDInt)
}

func (e *Event) RawMessage() string {
	return e.RawMessageValue
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
func (m *MessageEvent) Sender() event.Sender {
	return m.SenderValue
}

// UserID implements event.MessageEvent.
func (m *MessageEvent) UserID() string {
	return strconv.Itoa(m.UserIDInt)
}
