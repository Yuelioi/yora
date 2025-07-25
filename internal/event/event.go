package event

import (
	"time"
)

// Event 基础事件接口
type Event interface {

	// Type 返回事件类型，比如 message、notice、meta_event
	Type() string

	// SubType 返回事件子类型，比如 group、private、poke 等
	SubType() string

	// Time 返回事件时间戳
	Time() time.Time

	// SelfID 是 bot 自身 ID
	SelfID() string

	// Raw 返回原始事件数据
	Raw() any
}

// MessageEvent 消息事件接口
type MessageEvent interface {
	Event

	// UserID 发言者用户ID
	UserID() string

	// ChatID 聊天会话ID（群聊ID、频道ID等）
	ChatID() string

	// Message 返回结构化消息
	Message() Message

	// RawMessage 返回原始文本内容
	RawMessage() string

	// Sender 返回发送者信息
	Sender() Sender

	// IsGroup 判断是否为群聊消息
	IsGroup() bool

	// IsPrivate 判断是否为私聊消息
	IsPrivate() bool

	// MessageID 消息ID（字符串格式，适配所有平台）
	MessageID() string

	// ReplyTo 回复的消息ID（如果是回复消息）
	ReplyTo() string

	// Extra 额外数据
	Extra() map[string]any
}

// MetaEvent 元事件接口（如心跳、状态变更等）
type MetaEvent interface {
	Event

	// Status 状态信息
	Status() map[string]any

	// Extra 额外数据
	Extra() map[string]any
}

type NoticeEvent interface {
	Event

	// UserID 相关用户ID
	UserID() string

	// ChatID 相关聊天会话ID
	ChatID() string

	// OperatorID 操作者ID（如果有）
	OperatorID() string

	// Extra 额外数据
	Extra() map[string]any
}

// RequestEvent 请求事件接口（如好友申请、入群申请等）
type RequestEvent interface {
	Event

	// UserID 请求用户ID
	UserID() string

	// ChatID 相关聊天会话ID（如果有）
	ChatID() string

	// Comment 请求备注信息
	Comment() string

	// Flag 请求标识符，用于响应请求
	Flag() string

	// Extra 额外数据
	Extra() map[string]any
}
