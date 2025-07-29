package adapter

import (
	"net/http"
	"yora/pkg/event"
	"yora/pkg/message"
)

// 协议适配器接口
type Adapter interface {
	// 返回适配的协议类型
	Protocol() Protocol

	// 解析原始事件数据为通用事件接口
	ParseEvent(raw any) (event.Event, error)

	// 解析协议特定的消息为通用消息段
	ParseMessage(raw string) ([]message.Segment, error)

	// 验证事件数据是否有效
	ValidateEvent(event event.Event) error

	// 获取协议能力集
	GetCapabilities() Capabilities

	// 调用协议API
	CallAPI(action string, params any) (any, error)

	// // 处理HTTP请求
	// HandleHTTP(w http.ResponseWriter, r *http.Request) error

	// 处理WebSocket请求
	HandleWebSocket(w http.ResponseWriter, r *http.Request, f func(message []byte)) error

	// 发送消息
	Send(userId string, groupId string, message message.Message) (any, error)
}

type Registry interface {
	// 注册协议适配器
	Register(adapter Adapter) error

	// 注销协议适配器
	Unregister(protocol Protocol) error

	// 获取协议适配器
	Adapters() map[Protocol]Adapter
}
