package middleware

import (
	"context"
	"yora/internal/event"
)

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
