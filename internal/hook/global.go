package hook

// 全局Hook管理器
var globalHookManager = NewHookManager()

// GlobalHookManager 获取全局Hook管理器
func GlobalHookManager() *HookManager {
	return globalHookManager
}

// RegisterGlobalHook 注册全局Hook
func RegisterGlobalHook(hookType HookType, handler HookHandler) string {
	return globalHookManager.AddHook(hookType, handler)
}

// TriggerGlobalHook 触发全局Hook
func TriggerGlobalHook(hookType HookType, ctx *HookContext) error {
	return globalHookManager.TriggerHook(hookType, ctx)
}
