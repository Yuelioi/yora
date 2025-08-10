package bot

import (
	"context"
	"yora/pkg/event"
	"yora/pkg/provider"
)

// 获取配置
func Config() provider.Provider {
	return provider.StaticProvider(func(ctx context.Context, e event.Event) any {
		return GetBot().Config()
	})
}

// 获取Bot
func BotProvider() provider.Provider {
	return provider.StaticProvider(func(ctx context.Context, e event.Event) any {
		return GetBot()
	})
}
