package bot

import (
	"yora/pkg/adapter"
	"yora/pkg/conf"
	"yora/pkg/message"
	"yora/pkg/middleware"
	"yora/pkg/plugin"
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
	RegisterAdapters(adapters ...adapter.Adapter) error

	// 注册中间件
	RegisterMiddlewares(middlewares ...middleware.Middleware) error

	// 注册插件
	RegisterPlugins(plugins ...plugin.Plugin) error

	// 获取配置
	Config() *conf.BotConfig

	// 运行Bot(执行初始化等操作)
	Run() error

	// 检查机器人是否正在运行
	IsRunning() bool

	// 关闭Bot
	ShutDown() error

	// 获取插件
	Plugins() []plugin.Plugin
}
