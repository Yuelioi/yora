package hook

import "context"

type PluginHookContext struct {
	*HookContext
	Plugin any // 插件实例
}

func NewPluginHookContext(ctx context.Context, hookType HookType, plugin any) *PluginHookContext {
	hc := NewHookContext(ctx, hookType)
	return &PluginHookContext{
		HookContext: hc,
		Plugin:      plugin,
	}
}

type BotHookContext struct {
	*HookContext
	Bot any // 机器人实例
}

func NewBotHookContext(ctx context.Context, hookType HookType, bot any) *BotHookContext {
	hc := NewHookContext(ctx, hookType)
	return &BotHookContext{
		HookContext: hc,
		Bot:         bot,
	}
}

type MessageHookContext struct {
	*HookContext
	Message any // 消息实例
}

func NewMessageHookContext(ctx context.Context, hookType HookType, message any) *MessageHookContext {
	hc := NewHookContext(ctx, hookType)
	return &MessageHookContext{
		HookContext: hc,
		Message:     message,
	}
}
