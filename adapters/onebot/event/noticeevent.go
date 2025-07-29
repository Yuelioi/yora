package event

import (
	"strconv"
	"yora/pkg/event"
)

var _ event.NoticeEvent = (*NoticeEvent)(nil)

type NoticeEvent struct {
	Event
}

func (n *NoticeEvent) ChatID() string {
	return strconv.Itoa(n.GroupIDInt)
}

// todo
func (n *NoticeEvent) Extra() map[string]any {
	return nil
}

// todo
func (n *NoticeEvent) OperatorID() string {
	return ""
}

// UserID implements event.NoticeEvent.
// Subtle: this method shadows the method (Event).UserID of NoticeEvent.Event.
func (n *NoticeEvent) UserID() string {
	return strconv.Itoa(n.UserIDInt)
}
