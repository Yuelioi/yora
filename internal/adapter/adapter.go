package adapter

import (
	"yora/internal/event"
	"yora/internal/message"
	"yora/internal/middleware"
)

// ProtocolAdapter 协议适配器接口
type Adapter interface {
	// Protocol 返回适配的协议类型
	Protocol() Protocol

	// ParseEvent 解析原始事件数据为通用事件接口
	ParseEvent(raw any) (event.Event, error)

	// ParseMessage 解析协议特定的消息为通用消息段
	ParseMessage(raw string) ([]message.Segment, error)

	// ValidateEvent 验证事件数据是否有效
	ValidateEvent(event event.Event) error

	// GetCapabilities 获取协议能力集
	GetCapabilities() Capabilities

	// CallAPI 调用协议API
	CallAPI(action string, params any) (any, error)

	// Send 发送消息
	Send(messageType string, userId string, groupId string, message message.Message) (any, error)
}

type Registry interface {
	// Register 注册协议适配器
	Register(adapter Adapter) error

	// Unregister 注销协议适配器
	Unregister(protocol Protocol) error

	// GetAdapter 获取协议适配器
	GetAdapter(protocol Protocol) (Adapter, error)

	// 获取协议适配器
	Adapters() map[Protocol]Adapter

	// 获取所有中间件
	Middlewares() []middleware.Middleware
}
