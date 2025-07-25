package plugin

import "yora/internal/matcher"

// Plugin 插件接口
type Plugin interface {
	Metadata() *Metadata

	SetMetadata(*Metadata)

	RegisterMatcher(m *matcher.Matcher)

	Load() error

	Unload() error

	Matchers() []*matcher.Matcher
}

// Metadata 插件元数据
type Metadata struct {
	ID          string // 插件标识(必填)
	Name        string // 插件显示名称(必填)
	Description string // 插件描述
	Version     string // 插件版本
	Author      string // 插件作者
	Usage       string // 插件用途
	Group       string // 插件分组
	Extra       any    // 额外信息
}

type BasePlugin struct {
	matchers []*matcher.Matcher
	metadata *Metadata
}

func (p *BasePlugin) Matchers() []*matcher.Matcher {
	return p.matchers
}

func (p *BasePlugin) Unload() error {
	return nil
}

func (p *BasePlugin) RegisterMatcher(m *matcher.Matcher) {
	p.matchers = append(p.matchers, m)
}

func (p *BasePlugin) Metadata() *Metadata {
	return p.metadata
}
func (p *BasePlugin) SetMetadata(m *Metadata) {
	p.metadata = m
}
