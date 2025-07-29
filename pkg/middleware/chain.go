package middleware

import (
	"context"
	"yora/pkg/event"
)

// 构建中间件 并返回最终处理函数
func Chain(middlewares []Middleware, final func(ctx context.Context, event event.Event) error) func(ctx context.Context, event event.Event) error {
	if len(middlewares) == 0 {
		return final
	}
	current := final
	for i := len(middlewares) - 1; i >= 0; i-- {
		m := middlewares[i]
		next := current
		current = func(ctx context.Context, evt event.Event) error {
			return m.Process(ctx, evt, next)
		}
	}
	return current
}
