package onebot

import (
	"context"
	"fmt"
	"reflect"
	"yora/internal/event"
	"yora/internal/matcher"
)

type Handler struct {
	handlerFunc  reflect.Value
	dependencies []matcher.Dependent
	typeCache    map[reflect.Type]matcher.Dependent // 类型缓存
}

// 依赖注入的处理器
func NewHandler(f any, deps ...matcher.Dependent) *Handler {

	hds := make([]matcher.Dependent, 0, len(deps)+1)
	hds = append(hds, OnebotEvent(), OneBot())

	h := &Handler{
		handlerFunc:  reflect.ValueOf(f),
		dependencies: hds,
		typeCache:    make(map[reflect.Type]matcher.Dependent),
	}
	if err := h.validateAndCacheParameters(); err != nil {
		panic(err)
	}
	return h
}

func (h *Handler) Match(ctx context.Context, e event.Event) bool {
	for _, dep := range h.dependencies {
		if dep.Match(ctx, e) {
			return true
		}
	}
	return false
}

func (h *Handler) Type() reflect.Type {
	return h.handlerFunc.Type()
}

func (h *Handler) Resolve(ctx context.Context, e event.Event) (any, error) {
	funcType := h.handlerFunc.Type()
	numParams := funcType.NumIn()
	args := make([]reflect.Value, numParams)

	for i := 0; i < numParams; i++ {
		paramType := funcType.In(i)

		// 从缓存中获取依赖
		dep, exists := h.typeCache[paramType]
		if !exists {
			return nil, fmt.Errorf("no cached dependency found for parameter %d (type: %v)", i, paramType)
		}

		// 调用依赖的 Resolve 获取参数值
		val, err := dep.Resolve(ctx, e)
		if err != nil {
			return nil, fmt.Errorf("dependency.Resolve failed for parameter %d (type: %v): %w", i, paramType, err)
		}

		if val == nil {
			return nil, fmt.Errorf("dependency.Resolve returned nil for parameter %d (type: %v)", i, paramType)
		}

		valueReflect := reflect.ValueOf(val)

		// 确保返回的值类型能赋值给参数类型
		if !valueReflect.Type().AssignableTo(paramType) {
			// 允许指针和值类型的转换
			if valueReflect.Kind() == reflect.Ptr && valueReflect.Elem().Type() == paramType {
				valueReflect = valueReflect.Elem()
			} else if paramType.Kind() == reflect.Ptr && valueReflect.Type() == paramType.Elem() {
				// 这里可以考虑包装成指针
				ptr := reflect.New(valueReflect.Type())
				ptr.Elem().Set(valueReflect)
				valueReflect = ptr
			} else {
				return nil, fmt.Errorf("dependency.Resolve returned incompatible type for parameter %d: expected %v, got %v",
					i, paramType, valueReflect.Type())
			}
		}

		args[i] = valueReflect
	}

	// 调用处理函数，传入组装好的参数
	results := h.handlerFunc.Call(args)

	// 处理返回值中可能的 error
	if len(results) > 0 {
		lastResult := results[len(results)-1]
		if err, ok := lastResult.Interface().(error); ok && err != nil {
			return nil, err
		}
	}

	// 如果有返回值，返回第一个（除了最后的 error）
	if len(results) > 0 {
		if len(results) == 1 {
			return results[0].Interface(), nil
		}
		return results[0].Interface(), nil
	}

	return nil, nil
}

// findMatchingDependency 为指定类型找到匹配的依赖 (支持降级匹配)
func (h *Handler) findMatchingDependency(targetType reflect.Type) (matcher.Dependent, error) {
	var fallbackDependencies []matcher.Dependent

	// 注入基础依赖

	// 第一轮：寻找精确类型匹配的依赖
	for _, dep := range h.dependencies {
		depType := dep.Type()

		// 如果是动态依赖（Type() 返回 nil），收集起来作为备选
		if depType == nil {
			fallbackDependencies = append(fallbackDependencies, dep)
			continue
		}

		// 精确类型匹配
		if h.isTypeCompatible(depType, targetType) {
			return dep, nil
		}
	}

	// 第二轮：如果没找到精确匹配，按顺序尝试使用动态依赖
	if len(fallbackDependencies) > 0 {
		// 返回第一个动态依赖，运行时会尝试从 context 获取
		return fallbackDependencies[0], nil
	}

	return nil, fmt.Errorf("no compatible dependency found for type: %v", targetType)
}

// isTypeCompatible 检查类型兼容性
func (h *Handler) isTypeCompatible(sourceType, targetType reflect.Type) bool {
	// nil 检查 - 防止动态依赖导致的 panic
	if sourceType == nil || targetType == nil {
		return false
	}

	// 直接类型匹配
	if sourceType == targetType {
		return true
	}

	// 可赋值性检查
	if sourceType.AssignableTo(targetType) {
		return true
	}

	// 指针和值类型的转换
	if sourceType.Kind() == reflect.Ptr && targetType.Kind() != reflect.Ptr {
		return sourceType.Elem() == targetType
	}
	if sourceType.Kind() != reflect.Ptr && targetType.Kind() == reflect.Ptr {
		return sourceType == targetType.Elem()
	}

	// 接口实现检查
	if targetType.Kind() == reflect.Interface {
		return sourceType.Implements(targetType)
	}

	return false
}

// validateAndCacheParameters 校验参数并建立类型缓存
func (h *Handler) validateAndCacheParameters() error {
	funcType := h.handlerFunc.Type()
	numParams := funcType.NumIn()

	// 为每个参数类型找到匹配的依赖并缓存
	for i := 0; i < numParams; i++ {
		paramType := funcType.In(i)

		matchedDep, err := h.findMatchingDependency(paramType)
		if err != nil {
			return fmt.Errorf("no dependency found for parameter %d (type: %v): %w", i, paramType, err)
		}

		// 缓存类型映射
		h.typeCache[paramType] = matchedDep
	}

	return nil
}
