package event

import (
	"strconv"
	"yora/internal/adapter"
	"yora/internal/event"
)

var _ event.Sender = (*Sender)(nil)
var _ event.GroupSender = (*GroupSender)(nil)

type Sender struct {
	UserID   int    `json:"user_id"`
	Nickname string `json:"nickname"`
	Card     string `json:"card"`
	Sex      string `json:"sex"`
	Age      int    `json:"age"`
	Area     string `json:"area"`
	Level    string `json:"level"`
	RoleStr  string `json:"role"`
	Title    string `json:"title"`
}

// Extra implements event.Sender.
func (s Sender) Extra() map[string]any {
	var m = make(map[string]any)
	m["sex"] = s.Sex
	m["age"] = s.Age
	m["area"] = s.Area
	m["level"] = s.Level
	m["role"] = s.Role
	m["title"] = s.Title
	return m
}

// ID implements event.Sender.
func (s Sender) ID() string {
	return strconv.Itoa(s.UserID)
}

// Username implements event.Sender.
func (s Sender) Username() string {
	return s.Nickname
}

// DisplayName implements event.Sender.
func (s Sender) DisplayName() string {
	return s.Card
}

// TODO
func (s Sender) AvatarURL() string {
	return ""
}

func (s Sender) Role() string {
	return s.RoleStr
}

// TODO
func (s Sender) IsAnonymous() bool {
	return false
}

// Protocol implements event.Sender.
func (s Sender) Protocol() adapter.Protocol {
	return adapter.ProtocolOneBot
}

// Raw implements event.Sender.
func (s Sender) Raw() any {
	return s
}

type GroupSender struct {
	Sender
	GroupID int    `json:"group_id"`
	RoleStr string `json:"role"`
}

func (s GroupSender) ChatID() string {
	return strconv.Itoa(s.GroupID)
}

func (s GroupSender) Role() string {
	return s.RoleStr
}
