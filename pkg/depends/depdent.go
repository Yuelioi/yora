package depends

import (
	"context"
	"yora/pkg/event"
)

// Dependent 接口用于提供依赖
type Dependent interface {
	Provide(ctx context.Context, e event.Event) any
}

// 函数式 Dependent
var _ Dependent = DependentFunc(nil)

type DependentFunc func(ctx context.Context, e event.Event) any

func (f DependentFunc) Provide(ctx context.Context, e event.Event) any {
	return f(ctx, e)
}
