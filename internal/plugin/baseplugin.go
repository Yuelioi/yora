package plugin

import (
	"context"
	"yora/internal/hook"
	"yora/internal/matcher"
)

var _ Plugin = (*BasePlugin)(nil)

type BasePlugin struct {
	matchers    []*matcher.Matcher
	metadata    *Metadata
	config      *BaseConfig
	hookManager *hook.HookManager
}

func NewBasePlugin() BasePlugin {
	return BasePlugin{
		hookManager: hook.NewHookManager(),
		matchers:    make([]*matcher.Matcher, 0),
		metadata:    &Metadata{},
		config:      NewBaseConfig(),
	}
}

// Load implements Plugin.
func (p *BasePlugin) Load() error {
	panic("unimplemented")
}

func (p *BasePlugin) Init() error {

	ctx := hook.NewPluginHookContext(context.Background(), hook.PluginBeforeInit, p)
	if err := p.TriggerHook(hook.PluginBeforeInit, ctx.HookContext); err != nil {
		return err
	}

	p.matchers = make([]*matcher.Matcher, 0)
	p.config = NewBaseConfig()

	ctx = hook.NewPluginHookContext(context.Background(), hook.PluginAfterInit, p)
	return p.TriggerHook(hook.PluginAfterInit, ctx.HookContext)
}

func (p *BasePlugin) BeforeLoad() error {
	ctx := hook.NewPluginHookContext(context.Background(), hook.PluginBeforeLoad, p)
	return p.TriggerHook(hook.PluginBeforeLoad, ctx.HookContext)
}

func (p *BasePlugin) AfterLoad() error {
	ctx := hook.NewPluginHookContext(context.Background(), hook.PluginAfterLoad, p)
	return p.TriggerHook(hook.PluginAfterLoad, ctx.HookContext)
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

func (p *BasePlugin) Config() Config {
	return p.config
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
