package message

import (
	"encoding/json"
	"fmt"
	"strings"
	"yora/internal/event"
)

var _ event.Message = Message{}

// Message 表示一个消息，由多个 Segment 组成
type Message []Segment

// GetSegmentsByType 根据类型获取所有匹配的片段
func (m Message) GetSegmentsByType(segmentType string) []event.Segment {
	var result []event.Segment
	for i := range m {
		if m[i].IsType(segmentType) {
			// 拷贝值到新变量，防止闭包问题
			seg := m[i]
			result = append(result, &seg)
		}
	}
	return result
}

// HasType 检查消息中是否包含指定类型的片段
func (m Message) HasType(segmentType string) bool {
	for _, segment := range m {
		if segment.IsType(segmentType) {
			return true
		}
	}
	return false
}

// IsEmpty 检查消息是否为空
func (m Message) IsEmpty() bool {
	return len(m) == 0
}

// PlainText 将消息转换为纯文本
func (m Message) PlainText() string {
	var parts []string
	for _, segment := range m {
		parts = append(parts, segment.String())
	}
	return strings.Join(parts, "")
}

// Segments 返回消息中的所有片段
func (m Message) Segments() []event.Segment {
	result := make([]event.Segment, len(m))
	for i := range m {
		seg := m[i]
		result[i] = &seg
	}
	return result
}

// String 将消息转换为字符串表示
func (m Message) String() string {
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("Message[%d segments]", len(m))
	}
	return string(data)
}

// NewMessage 创建新消息
func NewMessage(segments ...Segment) Message {
	return Message(segments)
}

// Append 追加消息段（返回新 Message）
func (m Message) Append(segment event.Segment) event.Message {
	if s, ok := segment.(Segment); ok {
		return append(m, s)
	}
	panic("segment is not of type Segment")
}

// New 创建消息对象，可以是 string、Segment、[]Segment、map、[]any、Message
func New(data any) Message {
	switch v := data.(type) {
	case string:
		return Message{{TypeStr: "text", DataMap: map[string]any{"text": v}}}

	case Message:
		return v

	case []Segment:
		return Message(v)

	case Segment:
		return Message{v}

	case []any:
		var msg Message
		for _, item := range v {
			seg := New(item)
			msg = append(msg, seg...)
		}
		return msg

	case map[string]any:
		typ, ok := v["type"].(string)
		if !ok {
			return nil
		}
		dataMap, ok := v["data"].(map[string]any)
		if !ok {
			return nil
		}
		return Message{{TypeStr: typ, DataMap: dataMap}}

	default:
		return nil
	}
}
