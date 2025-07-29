package middleware

import (
	"context"
	"fmt"
	"runtime/debug"
	"yora/pkg/event"
	"yora/pkg/log"
	"yora/pkg/middleware"
)

var loggerRecover = log.NewMiddleware("PANIC 中间件")

func RecoveryMiddleware() middleware.Middleware {
	return middleware.MiddlewareFunc("panic中间件", func(ctx context.Context, event event.Event, next middleware.HandlerFunc) error {
		defer func() {
			if r := recover(); r != nil {
				var errMsg string

				switch v := r.(type) {
				case error:
					errMsg = v.Error()
				default:
					errMsg = fmt.Sprintf("%v", v)
				}

				stack := string(debug.Stack())
				loggerRecover.Error().
					Str("event_type", event.Type()).
					Str("panic", errMsg).
					Str("stack", stack).
					Msg("捕获 panic")
			}
		}()

		return next(ctx, event)
	})
}
