package plugin

// Plugin 插件接口
type Plugin interface {
	Load() error
	Unload() error
	Metadata() *Metadata
}

// Metadata 插件元数据
type Metadata struct {
	ID          string // 插件ID
	Name        string // 插件名称
	Description string // 插件描述
	Version     string // 插件版本
	Author      string // 插件作者
	Usage       string // 插件用途
	Extra       any    // 额外信息
}
