package middleware

import (
	"context"
	"fmt"
	"time"
	"yora/pkg/event"
	"yora/pkg/middleware"
)

func TimeoutMiddleware(timeout time.Duration) middleware.Middleware {
	return middleware.MiddlewareFunc("超时中间件", func(ctx context.Context, event event.Event, next middleware.HandlerFunc) error {

		ec, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- next(ctx, event)
		}()

		select {
		case err := <-done:
			return err
		case <-ec.Done():
			return fmt.Errorf("处理超时: %v", timeout)
		}
	})
}
