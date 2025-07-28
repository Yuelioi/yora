package event

import (
	"strconv"
	"time"
	"yora/internal/event"
	"yora/protocols/onebot/message"
)

var _ event.Event = (*Event)(nil)

// OneBot消息结构体
type Event struct {
	UserIDInt       int             `json:"user_id"`
	Anonymous       any             `json:"anonymous"`
	RawMessageValue string          `json:"raw_message"`
	Font            int             `json:"font"`
	SelfIDInt       int             `json:"self_id"`
	PostType        PostType        `json:"post_type"`
	MessageType     MessageType     `json:"message_type"`
	SubTypeValue    SubType         `json:"sub_type"`
	MessageIDInt    int             `json:"message_id"`
	GroupIDInt      int             `json:"group_id"`
	MessageValue    any             `json:"message"`
	SenderValue     *message.Sender `json:"sender"`
	MessageStyle    MessageStyle    `json:"message_style"`
	TimeStamp       int             `json:"time"`
}

func (e *Event) Raw() any {
	return e
}

func (e *Event) SelfID() string {
	return strconv.Itoa(e.SelfIDInt)
}

func (e *Event) SubType() string {
	return string(e.SubTypeValue)
}

func (e *Event) Time() time.Time {
	return time.Unix(int64(e.TimeStamp), 0)
}

func (e *Event) Type() string {
	return string(e.PostType)
}

type MessageStyle struct {
	BubbleID              int  `json:"bubble_id"`
	PendantID             int  `json:"pendant_id"`
	FontID                int  `json:"font_id"`
	FontEffectID          int  `json:"font_effect_id"`
	IsCSFontEffectEnabled bool `json:"is_cs_font_effect_enabled"`
	BubbleDIYTextID       int  `json:"bubble_diy_text_id"`
}
