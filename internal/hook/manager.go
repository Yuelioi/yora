package hook

import (
	"fmt"
	"sync"
	"time"
)

// HookInfo Hook信息
type HookInfo struct {
	ID       string       // Hook唯一标识
	Handler  HookHandler  // 处理函数
	Priority HookPriority // 优先级
	Once     bool         // 是否只执行一次
	executed bool         // 是否已执行（用于Once）
}

// HookManager Hook管理器
type HookManager struct {
	hooks    map[HookType][]*HookInfo
	disabled map[HookType]bool // 禁用的Hook类型

	mutex sync.RWMutex
}

// NewHookManager 创建新的Hook管理器
func NewHookManager() *HookManager {
	return &HookManager{
		hooks:    make(map[HookType][]*HookInfo),
		disabled: make(map[HookType]bool),
	}
}

// AddHook 添加Hook
func (hm *HookManager) AddHook(hookType HookType, handler HookHandler) string {
	return hm.AddHookWithOptions(hookType, handler, PriorityNormal, "", false)
}

// AddHookWithPriority 添加带优先级的Hook
func (hm *HookManager) AddHookWithPriority(hookType HookType, handler HookHandler, priority HookPriority) string {
	return hm.AddHookWithOptions(hookType, handler, priority, "", false)
}

// AddOnceHook 添加只执行一次的Hook
func (hm *HookManager) AddOnceHook(hookType HookType, handler HookHandler) string {
	return hm.AddHookWithOptions(hookType, handler, PriorityNormal, "", true)
}

// AddHookWithOptions 添加Hook（完整选项）
func (hm *HookManager) AddHookWithOptions(hookType HookType, handler HookHandler, priority HookPriority, id string, once bool) string {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	if id == "" {
		id = fmt.Sprintf("hook_%d", time.Now().UnixNano())
	}

	hookInfo := &HookInfo{
		ID:       id,
		Handler:  handler,
		Priority: priority,
		Once:     once,
	}

	// 按优先级插入
	hooks := hm.hooks[hookType]
	inserted := false
	for i, existing := range hooks {
		if hookInfo.Priority > existing.Priority {
			// 插入到当前位置
			hm.hooks[hookType] = append(hooks[:i], append([]*HookInfo{hookInfo}, hooks[i:]...)...)
			inserted = true
			break
		}
	}

	if !inserted {
		hm.hooks[hookType] = append(hooks, hookInfo)
	}

	return id
}

// RemoveHook 移除Hook
func (hm *HookManager) RemoveHook(hookType HookType, id string) bool {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	hooks := hm.hooks[hookType]
	for i, hook := range hooks {
		if hook.ID == id {
			hm.hooks[hookType] = append(hooks[:i], hooks[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveAllHooks 移除所有Hook
func (hm *HookManager) RemoveAllHooks(hookType HookType) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	delete(hm.hooks, hookType)
}

// DisableHook 禁用Hook类型
func (hm *HookManager) DisableHook(hookType HookType) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	hm.disabled[hookType] = true
}

// EnableHook 启用Hook类型
func (hm *HookManager) EnableHook(hookType HookType) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	delete(hm.disabled, hookType)
}

// IsDisabled 检查Hook类型是否被禁用
func (hm *HookManager) IsDisabled(hookType HookType) bool {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	return hm.disabled[hookType]
}

// TriggerHook 触发Hook
func (hm *HookManager) TriggerHook(hookType HookType, ctx *HookContext) error {
	if hm.IsDisabled(hookType) {
		return nil
	}

	hm.mutex.RLock()
	hooks := make([]*HookInfo, len(hm.hooks[hookType]))
	copy(hooks, hm.hooks[hookType])
	hm.mutex.RUnlock()

	var toRemove []string

	for _, hookInfo := range hooks {
		// 检查是否已执行过（用于Once Hook）
		if hookInfo.Once && hookInfo.executed {
			continue
		}

		// 执行Hook
		if err := hookInfo.Handler(ctx); err != nil {
			// 如果不是错误Hook，可以选择继续执行其他Hook或者停止
			return err
		}

		// 标记Once Hook为已执行
		if hookInfo.Once {
			hookInfo.executed = true
			toRemove = append(toRemove, hookInfo.ID)
		}
	}

	// 移除已执行的Once Hook
	for _, id := range toRemove {
		hm.RemoveHook(hookType, id)
	}

	return nil
}

// TriggerHookAsync 异步触发Hook
func (hm *HookManager) TriggerHookAsync(hookType HookType, ctx *HookContext) {
	go func() {
		if err := hm.TriggerHook(hookType, ctx); err != nil {
			// 可以在这里记录日志或处理错误
		}
	}()
}

// ListHooks 列出所有Hook
func (hm *HookManager) ListHooks(hookType HookType) []string {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	var ids []string
	for _, hook := range hm.hooks[hookType] {
		ids = append(ids, hook.ID)
	}
	return ids
}

// HookCount 获取Hook数量
func (hm *HookManager) HookCount(hookType HookType) int {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	return len(hm.hooks[hookType])
}
