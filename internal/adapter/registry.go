package adapter

type Registry interface {

	// Register 注册协议适配器
	Register(adapter Adapter) error

	// Unregister 注销协议适配器
	Unregister(protocol Protocol) error

	// GetAdapter 获取协议适配器
	GetAdapter(protocol Protocol) (Adapter, error)

	// 获取协议适配器
	Adapters() []Adapter

	// 获取所有中间件
	Middlewares() []Middleware
}
