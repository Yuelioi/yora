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
