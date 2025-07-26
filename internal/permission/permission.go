package permission

import (
	"context"
	"yora/internal/condition"
	"yora/internal/event"
)

// 权限
type Permission condition.Condition

// PermissionFunc 权限函数类型
type PermissionFunc func(ctx context.Context, event event.Event) bool

func (f PermissionFunc) Match(ctx context.Context, event event.Event) bool {
	return f(ctx, event)
}
