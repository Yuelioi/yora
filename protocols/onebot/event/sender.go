package event

import (
	"strconv"
	"yora/internal/adapter"
	basemsg "yora/internal/message"
)

var _ basemsg.Sender = (*Sender)(nil)

type Sender struct {
	UserID   int    `json:"user_id"`
	Nickname string `json:"nickname"`
	CardStr  string `json:"card"`
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
	return s.CardStr
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

func (s Sender) Card() string {
	return s.CardStr
}
