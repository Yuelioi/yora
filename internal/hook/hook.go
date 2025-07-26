// hook/types.go
package hook

// Hookable 可Hook的接口
type Hookable interface {
	// HookManager 返回Hook管理器
	HookManager() *HookManager

	// RegisterHook 注册Hook
	RegisterHook(hookType HookType, handler HookHandler) string

	// TriggerHook 触发Hook
	TriggerHook(hookType HookType, ctx *HookContext) error
}
