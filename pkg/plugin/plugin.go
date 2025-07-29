package plugin

import (
	"yora/pkg/conf"
	"yora/pkg/hook"
)

// Plugin 插件接口
type Plugin interface {
	// 插件元数据
	Metadata() *Metadata
	SetMetadata(*Metadata)

	// 插件配置
	Config() *conf.PluginConfig
	SetConfig(c *conf.PluginConfig)

	// 匹配器
	RegisterMatcher(m *Matcher)
	Matchers() []*Matcher

	// 生命周期
	Init() error
	Load() error
	Unload() error

	hook.Hookable
}
