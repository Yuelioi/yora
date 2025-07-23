package middleware

import (
	"context"
	"time"
	"yora/internal/adapter"
	"yora/internal/event"

	"github.com/rs/zerolog/log"
)

// LoggingMiddleware 日志中间件
func LoggingMiddleware() adapter.Middleware {
	return adapter.MiddlewareFunc(func(ctx context.Context, event event.Event, next func(ctx context.Context, event event.Event) error) error {
		start := time.Now()

		log.Debug().Msgf("开始处理事件: %s, 类型: %s", event.SelfID(), event.Type())

		err := next(ctx, event)

		duration := time.Since(start)
		if err != nil {
			return err

		} else {
			log.Info().Msgf("事件处理成功: %s, 类型: %s, 耗时: %v", event.SelfID(), event.Type(), duration)
		}

		return err
	})
}
