package message

import (
	"strings"
	basemsg "yora/internal/message"
)

var _ basemsg.Message = Message{}

// Message 表示一个消息，由多个 Segment 组成
type Message []basemsg.Segment

// GetSegmentsByType 根据类型获取所有匹配的片段
func (m Message) GetSegmentsByType(segmentType string) []basemsg.Segment {
	var result []basemsg.Segment
	for _, segment := range m {
		if segment.IsType(segmentType) {
			result = append(result, segment)
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
		if segment.IsType("text") {
			parts = append(parts, segment.String())
		}
	}
	return strings.Join(parts, "")
}

// Segments 返回消息中的所有片段
func (m Message) Segments() []basemsg.Segment {
	result := make([]basemsg.Segment, len(m))
	copy(result, m)
	return result
}

// String 将消息转换为字符串表示
func (m Message) String() string {
	if len(m) == 0 {
		return ""
	}

	var parts []string
	for _, segment := range m {
		str := segment.String()
		if str != "" {
			parts = append(parts, str)
		}
	}
	return strings.Join(parts, "")
}

// NewMessage 创建新消息
func NewMessage(segments ...basemsg.Segment) Message {
	return Message(segments)
}

// Append 追加消息段（返回新 Message）
func (m Message) Append(seg basemsg.Segment) basemsg.Message {
	return append(m, seg)
}

// New 创建消息对象，可以是 string、Segment、[]Segment、map、[]any、Message
func New(data any) Message {
	switch v := data.(type) {
	case string:
		if v == "" {
			return Message{}
		}
		return Message{NewTextSegment(v)}

	case Message:
		return v

	case []basemsg.Segment:
		return Message(v)

	case []*Segment:
		result := make(Message, len(v))
		for i, seg := range v {
			result[i] = seg
		}
		return result

	case *Segment:
		if v == nil {
			return Message{}
		}
		return Message{v}

	case Segment:
		return Message{&v}

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
			return Message{}
		}

		var dataMap map[string]any
		if data, exists := v["data"]; exists {
			if dm, ok := data.(map[string]any); ok {
				dataMap = dm
			} else {
				dataMap = make(map[string]any)
			}
		} else {
			dataMap = make(map[string]any)
		}

		return Message{NewSegment(typ, dataMap)}

	default:
		return Message{}
	}
}
