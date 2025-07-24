package adapter

import (
	"yora/internal/event"
	"yora/protocols/onebot/message"
)

// ProtocolAdapter 协议适配器接口
type Adapter interface {
	// Protocol 返回适配的协议类型
	Protocol() Protocol

	// ParseEvent 解析原始事件数据为通用事件接口
	ParseEvent(raw any) (event.Event, error)

	// ParseMessage 解析协议特定的消息为通用消息段
	ParseMessage(raw string) ([]event.Segment, error)

	// ValidateEvent 验证事件数据是否有效
	ValidateEvent(event event.Event) error

	// GetCapabilities 获取协议能力集
	GetCapabilities() Capabilities

	CallAPI(action string, params any) (any, error)

	Send(messageType string, userId string, groupId string, message message.Message) (any, error)
}

type PlatformInfo interface {
	// Name 平台名称
	Name() string

	// Version 平台/协议版本
	Version() string

	// Extra 平台特定的额外信息
	Extra() map[string]any
}

// Capabilities 协议能力集
type Capabilities struct {
	// SupportsGroupChat 是否支持群聊
	SupportsGroupChat bool

	// SupportsPrivateChat 是否支持私聊
	SupportsPrivateChat bool

	// SupportsFileUpload 是否支持文件上传
	SupportsFileUpload bool

	// SupportsRichText 是否支持富文本
	SupportsRichText bool

	// SupportsReply 是否支持回复消息
	SupportsReply bool

	// SupportsForward 是否支持转发消息
	SupportsForward bool

	// SupportsEdit 是否支持编辑消息
	SupportsEdit bool

	// SupportsDelete 是否支持删除消息
	SupportsDelete bool

	// SupportedSegmentTypes 支持的消息段类型
	SupportedSegmentTypes []string

	// MaxMessageLength 最大消息长度
	MaxMessageLength int

	// MaxFileSize 最大文件大小（字节）
	MaxFileSize int64

	// Extra 协议特定的其他能力
	Extra map[string]any
}
