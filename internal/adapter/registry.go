package adapter

type Registry interface {

	// Register 注册协议适配器
	Register(adapter Adapter) error

	// Unregister 注销协议适配器
	Unregister(protocol Protocol) error

	// GetAdapter 获取协议适配器
	GetAdapter(protocol Protocol) (Adapter, error)
}
