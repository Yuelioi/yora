package adapter

import (
	"context"
	"yora/internal/event"
	"yora/internal/log"
)

var logger = log.NewMiddlewareLogger("中间件")

// Middleware 中间件接口
type Middleware interface {
	Process(ctx context.Context, event event.Event, next func(ctx context.Context, event event.Event) error) error
}

// MiddlewareFunc 中间件函数类型
type MiddlewareFunc func(ctx context.Context, event event.Event, next func(ctx context.Context, event event.Event) error) error

func (f MiddlewareFunc) Process(ctx context.Context, event event.Event, next func(ctx context.Context, event event.Event) error) error {
	err := f(ctx, event, next)
	if err != nil {
		logger.Error().Err(err).Msg("中间件处理失败")
	}
	return err
}
