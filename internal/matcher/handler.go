package matcher

import (
	"context"
	"fmt"
	"reflect"
	"yora/internal/event"
)

type Handler struct {
	fn        any
	deps      []Dependent                // 依赖容器 (仅注册使用)
	typedDeps map[reflect.Type]Dependent // 类型缓存 (在调用前使用)
}

// RegisterDependent 注册依赖
func (h *Handler) RegisterDependent(provider Dependent) *Handler {
	if h.deps == nil {
		h.deps = make([]Dependent, 0)
	}

	h.deps = append(h.deps, provider)
	return h
}

func NewHandler(fn any) *Handler {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	if fnType.Kind() != reflect.Func {
		panic("not a function")
	}

	h := &Handler{
		fn:        fn,
		deps:      make([]Dependent, 0),
		typedDeps: make(map[reflect.Type]Dependent),
	}

	return h
}

func (h *Handler) Call(ctx context.Context, e event.Event) error {
	fnValue := reflect.ValueOf(h.fn)
	fnType := fnValue.Type()

	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("not a function")
	}

	// 准备参数
	args := make([]reflect.Value, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		paramType := fnType.In(i)

		// 特殊处理context和event
		switch {
		case paramType == reflect.TypeOf((*context.Context)(nil)).Elem():
			args[i] = reflect.ValueOf(ctx)
		case paramType == reflect.TypeOf((*event.Event)(nil)).Elem():
			args[i] = reflect.ValueOf(e)
		default:
			// 从依赖容器中获取
			if provider, ok := h.typedDeps[paramType]; ok {
				args[i] = reflect.ValueOf(provider)
			} else {
				return fmt.Errorf("no provider for type %v", paramType)
			}
		}
	}

	// 调用函数
	results := fnValue.Call(args)

	// 处理返回值
	if len(results) > 0 {
		if err, ok := results[len(results)-1].Interface().(error); ok {
			return err
		}
	}

	return nil
}

// validate校验参数并建立类型缓存
func (h *Handler) Validate() error {
	// funcType := h.handlerFunc.Type()
	// numParams := funcType.NumIn()

	// // 为每个参数类型找到匹配的依赖并缓存
	// for i := 0; i < numParams; i++ {
	// 	paramType := funcType.In(i)

	// 	matchedDep, err := h.findMatchingDependency(paramType)
	// 	if err != nil {
	// 		return fmt.Errorf("no dependency found for parameter %d (type: %v): %w", i, paramType, err)
	// 	}

	// 	// 缓存类型映射
	// 	h.typeCache[paramType] = matchedDep
	// }

	return nil
}

// findMatchingDependency 为指定类型找到匹配的依赖 (支持降级匹配)
func (h *Handler) findMatchingDependency(targetType reflect.Type) (Dependent, error) {
	// var fallbackDependencies []Dependent

	// // 注入基础依赖

	// // 第一轮：寻找精确类型匹配的依赖
	// for _, dep := range h.deps {
	// 	depType := dep.Provide(ctx context.Context, e event.Event)

	// 	// 如果是动态依赖（Type() 返回 nil），收集起来作为备选
	// 	if depType == nil {
	// 		fallbackDependencies = append(fallbackDependencies, dep)
	// 		continue
	// 	}

	// 	// 精确类型匹配
	// 	if isTypeCompatible(depType, targetType) {
	// 		return dep, nil
	// 	}
	// }

	// // 第二轮：如果没找到精确匹配，按顺序尝试使用动态依赖
	// if len(fallbackDependencies) > 0 {
	// 	// 返回第一个动态依赖，运行时会尝试从 context 获取
	// 	return fallbackDependencies[0], nil
	// }

	return nil, fmt.Errorf("no compatible dependency found for type: %v", targetType)
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
