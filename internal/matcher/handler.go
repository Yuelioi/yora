package matcher

import (
	"context"
	"fmt"
	"reflect"
	"yora/internal/depends"
	"yora/internal/event"
)

type Handler struct {
	fnType         reflect.Type                   // 函数类型
	fnValue        reflect.Value                  // 函数
	deps           []depends.Dependent            // 依赖容器 (仅注册使用)
	typedDepsValue map[reflect.Type]reflect.Value // 类型缓存 (在调用前使用)
	numParams      int                            // 参数数量缓存
	paramTypes     []reflect.Type                 // 参数
}

// RegisterDependent 注册依赖
func (h *Handler) RegisterDependent(providers ...depends.Dependent) *Handler {
	h.deps = append(h.deps, providers...)
	return h
}

// BuildDependentType 建立依赖类型缓存
func (h *Handler) BuildDependentType(ctx context.Context, e event.Event) error {
	for i := 0; i < h.numParams; i++ {
		paramType := h.paramTypes[i]

		matchedDep, err := h.findMatchingDependencyValue(paramType, ctx, e)
		if err != nil {
			return fmt.Errorf("no dependency found for parameter %d (type: %v): %w", i, paramType, err)
		}

		// 缓存类型映射
		h.typedDepsValue[paramType] = matchedDep
	}

	return nil
}

func NewHandler(fn any) *Handler {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	if fnType.Kind() != reflect.Func {
		panic("new handler: not a function")
	}

	numParams := fnType.NumIn()
	paramTypes := make([]reflect.Type, numParams)

	for i := 0; i < numParams; i++ {
		paramTypes[i] = fnType.In(i)
	}

	return &Handler{
		fnValue:        fnValue,
		fnType:         fnType,
		deps:           make([]depends.Dependent, 0),
		typedDepsValue: make(map[reflect.Type]reflect.Value),
		numParams:      numParams,
		paramTypes:     paramTypes,
	}
}

// Call 执行函数
func (h *Handler) Call(ctx context.Context, e event.Event) error {
	// 准备参数
	args := make([]reflect.Value, h.numParams)
	for i := 0; i < h.numParams; i++ {
		paramType := h.paramTypes[i]

		// 从依赖容器中获取
		if v, ok := h.typedDepsValue[paramType]; ok {
			args[i] = v
		} else {
			return fmt.Errorf("no provider for type %v", paramType)
		}
	}

	// 调用函数
	results := h.fnValue.Call(args)

	// 处理返回值 - 只检查最后一个返回值是否为 error
	if len(results) > 0 {
		if lastResult := results[len(results)-1]; !lastResult.IsNil() {
			if err, ok := lastResult.Interface().(error); ok && err != nil {
				return err
			}
		}
	}

	return nil
}

// findMatchingDependency 为指定类型找到匹配的依赖
func (h *Handler) findMatchingDependencyValue(targetType reflect.Type, ctx context.Context, e event.Event) (reflect.Value, error) {
	for _, dep := range h.deps {
		v := dep.Provide(ctx, e)
		if v == nil {
			continue // 跳过 nil 依赖
		}

		depType := reflect.TypeOf(v)
		depValue := reflect.ValueOf(v)

		if isTypeCompatible(depType, targetType) {
			return depValue, nil
		}
	}

	return reflect.Value{}, fmt.Errorf("no compatible dependency found for type: %v", targetType)
}

// isTypeCompatible 检查类型兼容性
func isTypeCompatible(sourceType, targetType reflect.Type) bool {
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
