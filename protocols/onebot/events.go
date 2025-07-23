package onebot

import (
	"strconv"
	"time"
	"yora/internal/adapter"
	"yora/internal/event"
)

// OneBot消息结构体
type Event struct {
	UserIDInt       int          `json:"user_id"`
	Anonymous       any          `json:"anonymous"`
	RawMessageValue string       `json:"raw_message"`
	Font            int          `json:"font"`
	SelfIDInt       int          `json:"self_id"`
	PostType        PostType     `json:"post_type"`
	MessageType     MessageType  `json:"message_type"`
	SubTypeValue    SubType      `json:"sub_type"`
	MessageIDInt    int          `json:"message_id"`
	GroupIDInt      int          `json:"group_id"`
	MessageValue    any          `json:"message"`
	SenderValue     *Sender      `json:"sender"`
	MessageStyle    MessageStyle `json:"message_style"`
	TimeStamp       int          `json:"time"`
}

type MessageStyle struct {
	BubbleID              int  `json:"bubble_id"`
	PendantID             int  `json:"pendant_id"`
	FontID                int  `json:"font_id"`
	FontEffectID          int  `json:"font_effect_id"`
	IsCSFontEffectEnabled bool `json:"is_cs_font_effect_enabled"`
	BubbleDIYTextID       int  `json:"bubble_diy_text_id"`
}

// Platform implements event.Event.
func (e *Event) Platform() adapter.PlatformInfo {
	return &PlatformInfo{
		Platform:        "onebot",
		PlatformVersion: "",
		AppVersion:      "",
	}
}

var _ event.Event = (*Event)(nil)

// Raw implements event.Event.
func (e *Event) Raw() any {
	return nil
}

// SelfID implements event.Event.
func (e *Event) SelfID() string {
	return strconv.Itoa(e.SelfIDInt)
}

// SubType implements event.Event.
func (e *Event) SubType() string {
	return string(e.SubTypeValue)
}

// Time implements event.Event.
func (e *Event) Time() time.Time {
	return time.Unix(int64(e.TimeStamp), 0)
}

// Type implements event.Event.
func (e *Event) Type() string {
	return string(e.PostType)
}

// 实现event.MessageEvent
func (e *Event) UserID() string {
	return strconv.Itoa(e.UserIDInt)
}
func (e *Event) ChatID() string {
	return strconv.Itoa(e.GroupIDInt)
}

func (e *Event) Message() event.Message {
	return e.Message()
}
func (e *Event) RawMessage() string {
	return e.RawMessageValue
}

func (e *Event) Sender() Sender {
	return *e.SenderValue
}

func (e *Event) ChatType() string {
	panic("unimplemented")
}

func (e *Event) IsPrivate() bool {
	return e.PostType == "message" && e.MessageType == "private"
}
func (e *Event) IsGroup() bool {
	panic("unimplemented")
}
func (e *Event) MessageID() string {
	return strconv.Itoa(e.MessageIDInt)
}
func (e *Event) ReplyTo() string {
	return ""
}

// 实现event.MetaEvent

func (e *Event) MetaEventType() string {
	panic("unimplemented")
}
func (e *Event) Status() map[string]any {
	panic("unimplemented")
}

// 实现event.NoticeEvent
func (e *Event) NoticeType() string {
	panic("unimplemented")
}

func (e *Event) Extra() map[string]any {
	panic("unimplemented")
}
func (e *Event) OperatorID() string {
	return strconv.Itoa(e.UserIDInt)
}

// 实现event.NoticeEvent
func (e *Event) GroupID() string {
	return strconv.Itoa(e.GroupIDInt)
}
func (e *Event) Comment() string {
	return strconv.Itoa(e.GroupIDInt)
}

func (e *Event) Flag() string {
	panic("unimplemented")
}

func (e *Event) RequestType() string {
	panic("unimplemented")
}
