package bot

import (
	"context"
	"yora/internal/depends"
	"yora/internal/event"
)

// 基础依赖 用于注入
func OneBot() depends.Dependent {
	return depends.DependentFunc(func(ctx context.Context, e event.Event) any {
		return GetBot()
	})
}
