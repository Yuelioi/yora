package permission

import (
	"context"
	"yora/internal/condition"
	"yora/internal/event"
)

// Everyone 所有人权限
func Everyone() Permission {
	return PermissionFunc(func(ctx context.Context, e event.Event) bool {
		return true
	})
}

// 仅超级用户权限
func SuperUser(superUsers ...string) Permission {
	userSet := make(map[string]bool)
	for _, user := range superUsers {
		userSet[user] = true
	}

	return PermissionFunc(func(ctx context.Context, e event.Event) bool {
		if msgEvent, ok := e.(event.MessageEvent); ok {
			return userSet[msgEvent.UserID()]
		}
		return false
	})
}

// GroupOwner 仅群主权限
func GroupOwner() Permission {
	return PermissionFunc(func(ctx context.Context, e event.Event) bool {
		return getRole(e) == "owner"
	})
}

// GroupAdmin 仅群管理员权限
func GroupAdmin() Permission {
	return PermissionFunc(func(ctx context.Context, e event.Event) bool {
		return getRole(e) == "admin"
	})
}

// GroupMember 仅群成员权限
func GroupMember() Permission {
	return PermissionFunc(func(ctx context.Context, e event.Event) bool {
		return getRole(e) == "member"
	})
}

// 超级用户、群主、管理员
func GroupAdminOrOwner() Permission {
	return condition.Any(SuperUser(), GroupOwner(), GroupAdmin())
}
