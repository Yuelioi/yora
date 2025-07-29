// Package handler 封装了事件处理函数的逻辑，支持依赖注入与自动调用。
//
// 核心功能：
//   - 使用 NewHandler 构造并注册事件处理器
//   - 使用 RegisterDependent 方法注入函数依赖
//   - 使用 Call 方法执行并注入上下文和依赖参数
//
// 本包适用于构建具有自动依赖注入能力的事件驱动系统。

package handler

import (
	"context"
	"reflect"
	"yora/pkg/depends"
	"yora/pkg/event"
)

// 事件处理函数
type Handler struct {
	id         uintptr        // 函数的唯一标识（指针地址）
	fnType     reflect.Type   // 函数类型（签名）
	fnValue    reflect.Value  // 函数值（用于调用）
	numParams  int            // 参数数量
	paramTypes []reflect.Type // 参数类型（用于依赖匹配）
}

// 创建一个新的 Handler 实例并注册到 Registry
func NewHandler(fn any) *Handler {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	if fnType.Kind() != reflect.Func {
		panic("NewHandler: argument is not a function")
	}

	numParams := fnType.NumIn()
	paramTypes := make([]reflect.Type, numParams)
	for i := 0; i < numParams; i++ {
		paramTypes[i] = fnType.In(i)
	}

	ptr := fnValue.Pointer()
	handler := &Handler{
		id:         ptr,
		fnValue:    fnValue,
		fnType:     fnType,
		numParams:  numParams,
		paramTypes: paramTypes,
	}

	// 将 handler 注册到全局的依赖注册器
	GetHandlerRegistry().RegisterHandler(handler)
	return handler
}

// 注册依赖到注册中心
func (h *Handler) RegisterDependent(deps ...depends.Dependent) *Handler {
	GetHandlerRegistry().RegisterDependent(deps...)
	return h
}

// 执行 handler 函数
func (h *Handler) Call(ctx context.Context, e event.Event) error {
	reg := GetHandlerRegistry()

	// 为本次调用构建参数依赖（全局缓存中如已存在则跳过）
	if err := reg.BuildTypedValues(h.id, h.paramTypes, ctx, e); err != nil {
		return err
	}

	// 获取参数值并填充
	args := make([]reflect.Value, h.numParams)
	for i, t := range h.paramTypes {
		v, err := reg.GetTypedDependency(t)
		if err != nil {
			return err
		}
		args[i] = v
	}

	// 执行函数
	results := h.fnValue.Call(args)

	// 处理最后一个返回值（如果为 error 且非 nil，则返回）
	if len(results) > 0 {
		if last := results[len(results)-1]; !last.IsNil() {
			if err, ok := last.Interface().(error); ok {
				return err
			}
		}
	}

	return nil
}
