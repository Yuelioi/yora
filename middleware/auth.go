package middleware

import (
	"context"
	"fmt"
	"yora/internal/adapter"
	"yora/internal/event"
)

func AuthMiddleware() adapter.Middleware {
	return adapter.MiddlewareFunc(func(ctx context.Context, e event.Event, next func(ctx context.Context, event event.Event) error) error {
		// 权限检查逻辑
		if msgEvent, ok := e.(event.MessageEvent); ok {
			// 检查用户权限
			if msgEvent.Sender().ID() != "" {
				// 根据用户角色进行权限检查
				switch msgEvent.Sender().Role() {
				case "admin", "owner":
					// 管理员和群主有所有权限
					return next(ctx, e)
				case "member":
					// 普通成员的权限检查
					return next(ctx, e)
				default:
					return fmt.Errorf("权限不足: 未知用户角色 %s", msgEvent.Sender().Role())
				}
			}
		}

		// 其他类型事件直接通过
		return next(ctx, e)
	})
}
