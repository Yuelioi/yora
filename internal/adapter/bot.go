package adapter

type Bot interface {
	SelfID() string

	// 所属平台名，如 "onebot"
	Platform() string

	// 发送消息（通用格式）
	SendMessage(params ...any) error

	CallAPI(params ...any) (any, error)

	// 获取原始 bot 对象（如 OneBot 的详细结构）
	Raw() any

	// 添加中间件
	AddMiddleware(middleware ...Middleware)
}
