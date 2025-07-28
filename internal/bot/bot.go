package bot

import (
	"yora/internal/adapter"
	"yora/internal/message"
	"yora/internal/middleware"
	"yora/internal/plugin"
)

type Bot interface {

	// 机器人自身ID
	SelfID() string

	// 所属平台名，如 "onebot"
	Platform() string

	// 发送消息（通用格式）
	Send(userId string, groupId string, message message.Message) (any, error)

	// 调用 API（通用格式）
	CallAPI(params ...any) (any, error)

	// 注册适配器
	RegisterAdapter(adapter adapter.Adapter) error

	// 添加中间件
	AddMiddleware(middleware middleware.Middleware)

	// 运行Bot(执行初始化等操作)
	Run() error

	// 关闭Bot
	ShutDown() error

	// 获取插件
	Plugins() []plugin.Plugin
}
