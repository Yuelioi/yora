package matcher

import (
	"context"
	"reflect"
	"yora/internal/event"
)

type DepCTXKey string

type Dependent interface {
	// Resolve 基于事件解析依赖项 - 给Handler用的主要方法
	Resolve(ctx context.Context, e event.Event) (any, error)

	// Match 判断此依赖是否能处理给定事件 - 用于Handler选择合适的依赖
	Match(ctx context.Context, e event.Event) bool

	// Type 返回当前依赖的类型
	Type() reflect.Type
}

type DependencyFunc struct {
	match   func(ctx context.Context, e event.Event) bool
	resolve func(ctx context.Context, e event.Event) (any, error)
}

func (f DependencyFunc) Match(ctx context.Context, e event.Event) bool {
	if f.match == nil {
		return true
	}
	return f.match(ctx, e)
}

func (f DependencyFunc) Resolve(ctx context.Context, e event.Event) (any, error) {
	return f.resolve(ctx, e)
}

func (f DependencyFunc) Type() reflect.Type {
	var zero any
	return reflect.TypeOf(zero)
}

func NewDependency(
	resolve func(ctx context.Context, e event.Event) (any, error),
	match ...func(ctx context.Context, e event.Event) bool,
) Dependent {
	var m func(ctx context.Context, e event.Event) bool
	if len(match) > 0 {
		m = match[0]
	}
	return DependencyFunc{match: m, resolve: resolve}
}
