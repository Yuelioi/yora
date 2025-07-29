package plugin

import (
	"yora/pkg/conf"
	"yora/pkg/hook"
)

var _ Plugin = (*BasePlugin)(nil)

type BasePlugin struct {
	matchers    []*Matcher
	metadata    *Metadata
	config      *conf.PluginConfig
	hookManager *hook.HookManager
}

// !需要插件自己实现
func (p *BasePlugin) Load() error {
	panic("unimplemented")
}

// !需要插件自己实现
func (p *BasePlugin) Unload() error {
	return nil
}

// 初始化
func (p *BasePlugin) Init() error {
	p.hookManager = hook.NewHookManager()

	p.matchers = make([]*Matcher, 0)
	p.config = conf.NewPluginConfig()
	p.metadata = &Metadata{}

	return nil
}

func (p *BasePlugin) RegisterMatcher(m *Matcher) {
	p.matchers = append(p.matchers, m)
}

func (p *BasePlugin) Matchers() []*Matcher {
	return p.matchers
}

func (p *BasePlugin) Metadata() *Metadata {
	return p.metadata
}
func (p *BasePlugin) SetMetadata(m *Metadata) {
	p.metadata = m
}

func (p *BasePlugin) Config() *conf.PluginConfig {
	return p.config
}
func (p *BasePlugin) SetConfig(c *conf.PluginConfig) {
	p.config = c
}

func (p *BasePlugin) HookManager() *hook.HookManager {
	return p.hookManager
}

func (p *BasePlugin) RegisterHook(hookType hook.HookType, handler hook.HookHandler) string {
	return p.hookManager.AddHook(hookType, handler)
}

func (p *BasePlugin) TriggerHook(hookType hook.HookType, ctx *hook.HookContext) error {
	return p.hookManager.TriggerHook(hookType, ctx)
}
