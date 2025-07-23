package onebot

import (
	"context"
	"strconv"
	"yora/internal/event"
	"yora/internal/matcher"
)

// Everyone 所有人权限
func Everyone() matcher.Permission {
	return matcher.PermissionFunc(func(ctx context.Context, e event.Event) bool {
		return true
	})
}

// SuperUser 仅超级用户权限
func SuperUser(superUsers ...string) matcher.Permission {
	userSet := make(map[string]bool)
	for _, user := range superUsers {
		userSet[user] = true
	}

	return matcher.PermissionFunc(func(ctx context.Context, e event.Event) bool {
		if msgEvent, ok := e.(*Event); ok {
			return userSet[strconv.Itoa(msgEvent.UserIDInt)]
		}
		return false
	})
}

func getRole(event event.Event) string {
	if msgEvent, ok := event.(*Event); ok {
		return msgEvent.Sender().Role()
	}
	return ""
}

// GroupOwner 仅群主权限
func GroupOwner() matcher.Permission {
	return matcher.PermissionFunc(func(ctx context.Context, e event.Event) bool {
		return getRole(e) == "owner"
	})
}

// GroupAdmin 仅群管理员权限
func GroupAdmin() matcher.Permission {
	return matcher.PermissionFunc(func(ctx context.Context, e event.Event) bool {
		return getRole(e) == "admin"
	})
}

// GroupMember 仅群成员权限
func GroupMember() matcher.Permission {
	return matcher.PermissionFunc(func(ctx context.Context, e event.Event) bool {
		return getRole(e) == "member"
	})
}

// 超级用户、群主、管理员
func GroupAdminOrOwner() matcher.Permission {
	return matcher.AnyPermission(SuperUser(), GroupOwner(), GroupAdmin())
}
