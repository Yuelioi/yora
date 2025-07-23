package matcher

import (
	"context"
	"yora/internal/event"
)

// Permission 权限接口
type Permission interface {
	Check(ctx context.Context, event event.Event) bool
}

// PermissionFunc 权限函数类型
type PermissionFunc func(ctx context.Context, event event.Event) bool

func (f PermissionFunc) Check(ctx context.Context, event event.Event) bool {
	return f(ctx, event)
}

func AnyPermission(permissions ...Permission) Permission {
	return PermissionFunc(func(ctx context.Context, event event.Event) bool {
		for _, p := range permissions {
			if p.Check(ctx, event) {
				return true
			}
		}
		return false
	})
}

func AllPermissions(permissions ...Permission) Permission {
	return PermissionFunc(func(ctx context.Context, event event.Event) bool {
		for _, p := range permissions {
			if !p.Check(ctx, event) {
				return false
			}
		}
		return true
	})
}

func NotPermission(permission Permission) Permission {
	return PermissionFunc(func(ctx context.Context, event event.Event) bool {
		return !permission.Check(ctx, event)
	})
}
