package middleware

import (
	"context"
	"time"
	"yora/internal/adapter"
	"yora/internal/event"
	"yora/internal/log"
)

var logger = log.NewMiddleware("事件耗时统计")

// LoggingMiddleware 日志中间件
func LoggingMiddleware() adapter.Middleware {
	return adapter.MiddlewareFunc("日志中间件", func(ctx context.Context, event event.Event, next func(ctx context.Context, event event.Event) error) error {
		start := time.Now()

		logger.Debug().Msgf("处理事件: %s, 类型: %s", event.SelfID(), event.Type())

		err := next(ctx, event)

		duration := time.Since(start)
		if err != nil {
			return err

		} else {
			logger.Info().Msgf("事件处理成功 BOT: %s, 类型: %s, 耗时: %v", event.SelfID(), event.Type(), duration)
		}

		return err
	})
}
