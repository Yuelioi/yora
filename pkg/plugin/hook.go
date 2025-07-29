package plugin

import "yora/pkg/hook"

type Wrapper[T any] struct {
	Plugin T
}

// 泛型链式注册器
func Register[T Plugin](p T) *Wrapper[T] {
	if base, ok := any(p).(interface{ Init() error }); ok {
		base.Init()
	}
	return &Wrapper[T]{Plugin: p}
}

func (w *Wrapper[T]) WithHook(event hook.HookType, fn hook.HookHandler) *Wrapper[T] {
	if h, ok := any(w.Plugin).(interface {
		RegisterHook(hook.HookType, hook.HookHandler)
	}); ok {
		h.RegisterHook(event, fn)
	}
	return w
}
