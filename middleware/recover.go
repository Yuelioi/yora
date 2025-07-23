package middleware

import (
	"context"
	"fmt"
	"yora/internal/adapter"
	"yora/internal/event"
)

func RecoveryMiddleware() adapter.Middleware {
	return adapter.MiddlewareFunc(func(ctx context.Context, event event.Event, next func(ctx context.Context, event event.Event) error) error {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[Recovery] 捕获到 panic: %v\n", r)
			}
		}()

		return next(ctx, event)
	})
}
