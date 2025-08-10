package handler

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"yora/pkg/event"
	"yora/pkg/provider"
)

// TODO
// 依赖有的是动态的 有的是静态的

// 是全局依赖注入管理器（负责依赖匹配与缓存）
type HandlerRegistry struct {
	mu sync.RWMutex

	// 静态依赖（全局单例）
	staticDeps map[reflect.Type]any

	// 动态依赖
	dynamicProviders []provider.Provider // （插件/系统注册）

	// 工厂函数（每次创建新实例）
	factories map[reflect.Type]func(ctx context.Context, e event.Event) any

	typedDepsValues map[reflect.Type]reflect.Value // 缓存：类型 -> 构造出的值（当前事件周期有效）

	paramTypesMap map[uintptr][]reflect.Type // handlerID -> 参数类型列表（用于调试/辅助）

}

var (
	once sync.Once
	h    *HandlerRegistry
)

// 获取单例全局依赖注入注册器
func GetHandlerRegistry() *HandlerRegistry {
	once.Do(func() {
		h = &HandlerRegistry{
			paramTypesMap:    make(map[uintptr][]reflect.Type),
			typedDepsValues:  make(map[reflect.Type]reflect.Value),
			staticDeps:       make(map[reflect.Type]any),
			factories:        make(map[reflect.Type]func(ctx context.Context, e event.Event) any),
			dynamicProviders: make([]provider.Provider, 0),
		}
	})

	return h
}

// ! 清空依赖注入缓存（每次事件调用前执行一次）
func (r *HandlerRegistry) ResetCache() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.typedDepsValues = make(map[reflect.Type]reflect.Value)
}

func (r *HandlerRegistry) RegisterProviders(providers ...provider.Provider) *HandlerRegistry {
	for _, pro := range providers {
		switch p := pro.(type) {
		case provider.StaticProvider:

			v := pro.Provide(context.Background(), nil)
			if v == nil {
				continue
			}
			vt := reflect.TypeOf(v)
			vv := reflect.ValueOf(v)
			r.staticDeps[vt] = vv
		case provider.DynamicProvider:
			r.dynamicProviders = append(r.dynamicProviders, p)
		default:
			// 默认当作动态处理
			r.dynamicProviders = append(r.dynamicProviders, p)
		}
	}
	return r
}

// 注册 Handler 的元信息（参数类型等）
func (r *HandlerRegistry) RegisterHandler(handler *Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.paramTypesMap[handler.id] = handler.paramTypes
}

// 构建并缓存当前调用所需的依赖（只构建未缓存的）
func (r *HandlerRegistry) BuildTypedValues(id uintptr, paramTypes []reflect.Type, ctx context.Context, e event.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.typedDepsValues == nil {
		r.typedDepsValues = make(map[reflect.Type]reflect.Value)
	}

	// 注入静态依赖
	for t, pro := range r.staticDeps {
		if _, exists := r.typedDepsValues[t]; exists {
			continue // 已构建，跳过
		}
		v := reflect.ValueOf(pro)
		r.typedDepsValues[t] = v
	}

	// 注入动态依赖
	for _, t := range paramTypes {
		if _, exists := r.typedDepsValues[t]; exists {
			continue // 已构建，跳过
		}

		v, err := r.findMatchingDependencyValue(t, ctx, e)
		if err != nil {
			return fmt.Errorf("构建依赖失败 [%v]: %w", t, err)
		}

		r.typedDepsValues[t] = v
	}

	return nil
}

// 获取已缓存的依赖值（用于注入参数）
func (r *HandlerRegistry) GetTypedDependency(t reflect.Type) (reflect.Value, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if v, ok := r.typedDepsValues[t]; ok {
		return v, nil
	}
	return reflect.Value{}, fmt.Errorf("找不到类型 [%v] 的依赖", t)
}

// 根据类型查找并返回匹配的依赖
func (r *HandlerRegistry) findMatchingDependencyValue(t reflect.Type, ctx context.Context, e event.Event) (reflect.Value, error) {
	for _, pro := range r.dynamicProviders {
		v := pro.Provide(ctx, e)
		if v == nil {
			continue
		}
		vt := reflect.TypeOf(v)
		vv := reflect.ValueOf(v)
		if isTypeCompatible(vt, t) {
			return vv, nil
		}
	}
	return reflect.Value{}, fmt.Errorf("未匹配到类型 [%v] 的依赖", t)
}

// 判断两个类型是否兼容（用于依赖匹配）
func isTypeCompatible(src, tgt reflect.Type) bool {
	if src == nil || tgt == nil {
		return false
	}
	if src == tgt || src.AssignableTo(tgt) {
		return true
	}
	if src.Kind() == reflect.Ptr && tgt.Kind() != reflect.Ptr {
		return src.Elem() == tgt
	}
	if src.Kind() != reflect.Ptr && tgt.Kind() == reflect.Ptr {
		return src == tgt.Elem()
	}
	if tgt.Kind() == reflect.Interface {
		return src.Implements(tgt)
	}
	return false
}
