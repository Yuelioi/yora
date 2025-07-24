package adapter

import "yora/protocols/onebot/message"

type Bot interface {
	SelfID() string

	// 所属平台名，如 "onebot"
	Platform() string

	// 发送消息（通用格式）
	Send(messageType string, userId string, groupId string, message message.Message) (any, error)

	// 调用 API（通用格式）
	CallAPI(params ...any) (any, error)

	// 添加中间件
	AddMiddleware(middleware Middleware)

	// 运行Bot
	Run() error

	// 关闭Bot
	ShutDown() error
}
