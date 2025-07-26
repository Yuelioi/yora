package middleware

import (
	"context"
	"yora/internal/event"
	"yora/internal/log"
)

var logger = log.NewMiddleware("中间件")

type middlewareFuncWrapper struct {
	name string
	fn   func(ctx context.Context, event event.Event, next func(ctx context.Context, event event.Event) error) error
}

func (m middlewareFuncWrapper) Process(ctx context.Context, event event.Event, next func(ctx context.Context, event event.Event) error) error {
	err := m.fn(ctx, event, next)
	if err != nil {
		logger.Error().Err(err).Str("middleware", m.name).Msg("中间件处理失败")
	}
	return err
}

func (m middlewareFuncWrapper) Name() string {
	return m.name
}

// 将普通函数转为中间件接口
func MiddlewareFunc(name string, fn func(ctx context.Context, event event.Event, next func(ctx context.Context, event event.Event) error) error) Middleware {
	return middlewareFuncWrapper{
		name: name,
		fn:   fn,
	}
}
