package event

import "yora/pkg/message"

// MessageEvent 消息事件接口
type MessageEvent interface {
	Event

	// UserID 发言者用户ID
	UserID() string

	// ChatID 聊天会话ID（群聊ID、频道ID等）
	ChatID() string

	// Message 返回结构化消息
	Message() message.Message

	// RawMessage 返回原始文本内容
	RawMessage() string

	// Sender 返回发送者信息
	Sender() message.Sender

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

type GroupMessageEvent interface {
	Event

	// UserID 发言者用户ID
	UserID() string

	// ChatID 聊天会话ID（群聊ID、频道ID等）
	ChatID() string

	// Message 返回结构化消息
	Message() message.Message

	// RawMessage 返回原始文本内容
	RawMessage() string

	// Sender 返回发送者信息
	Sender() message.Sender

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

type PrimaryMessageEvent interface {
	Event

	// UserID 发言者用户ID
	UserID() string

	// Message 返回结构化消息
	Message() message.Message

	// RawMessage 返回原始文本内容
	RawMessage() string

	// Sender 返回发送者信息
	Sender() message.Sender

	// MessageID 消息ID（字符串格式，适配所有平台）
	MessageID() string

	// ReplyTo 回复的消息ID（如果是回复消息）
	ReplyTo() string

	// Extra 额外数据
	Extra() map[string]any
}
