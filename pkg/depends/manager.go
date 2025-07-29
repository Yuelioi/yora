package depends

import (
	"context"
	"reflect"
	"yora/pkg/event"
)

type DependencyManager struct {
	// 基础依赖（全局单例）
	baseDeps map[reflect.Type]interface{}
	// 工厂函数（每次创建新实例）
	factories map[reflect.Type]func(ctx context.Context, e event.Event) interface{}
	// 动态依赖提供者
	dynamicProviders []Dependent
}

// 在系统启动时注册
func (gdm *DependencyManager) RegisterBaseDep(dep interface{}) {
	gdm.baseDeps[reflect.TypeOf(dep)] = dep
}

func (gdm *DependencyManager) RegisterFactory(targetType reflect.Type, factory func(ctx context.Context, e event.Event) interface{}) {
	gdm.factories[targetType] = factory
}

func (gdm *DependencyManager) BaseDepends() map[reflect.Type]interface{} {
	return gdm.baseDeps
}

func (gdm *DependencyManager) Factories() map[reflect.Type]func(ctx context.Context, e event.Event) interface{} {
	return gdm.factories
}

func (gdm *DependencyManager) DynamicProviders() []Dependent {
	return gdm.dynamicProviders
}
