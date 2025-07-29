package hook

import (
	"context"
	"sync"
	"time"
)

// HookContext Hook执行上下文
type HookContext struct {
	ctx      context.Context
	hookType HookType

	data     map[string]any
	metadata map[string]string

	timestamp time.Time
	err       error

	mutex sync.RWMutex
}

// NewHookContext 创建新的Hook上下文
func NewHookContext(ctx context.Context, hookType HookType) *HookContext {
	if ctx == nil {
		ctx = context.Background()
	}

	return &HookContext{
		ctx:       ctx,
		hookType:  hookType,
		data:      make(map[string]any),
		metadata:  make(map[string]string),
		timestamp: time.Now(),
	}
}

// Context 获取上下文
func (hc *HookContext) Context() context.Context {
	return hc.ctx
}

// HookType 获取Hook类型
func (hc *HookContext) HookType() HookType {
	return hc.hookType
}

// Timestamp 获取创建时间
func (hc *HookContext) Timestamp() time.Time {
	return hc.timestamp
}

// Set 设置数据
func (hc *HookContext) Set(key string, value any) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	hc.data[key] = value
}

// Get 获取数据
func (hc *HookContext) Get(key string) (any, bool) {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()
	value, exists := hc.data[key]
	return value, exists
}

// GetString 获取字符串数据
func (hc *HookContext) GetString(key string) string {
	if value, exists := hc.Get(key); exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetInt 获取整数数据
func (hc *HookContext) GetInt(key string) int {
	if value, exists := hc.Get(key); exists {
		if i, ok := value.(int); ok {
			return i
		}
	}
	return 0
}

// GetBool 获取布尔数据
func (hc *HookContext) GetBool(key string) bool {
	if value, exists := hc.Get(key); exists {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return false
}

// SetError 设置错误
func (hc *HookContext) SetError(err error) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	hc.err = err
}

// GetError 获取错误
func (hc *HookContext) GetError() error {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()
	return hc.err
}

// SetMetadata 设置元数据
func (hc *HookContext) SetMetadata(key, value string) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	hc.metadata[key] = value
}

// GetMetadata 获取元数据
func (hc *HookContext) GetMetadata(key string) string {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()
	return hc.metadata[key]
}

// GetAllData 获取所有数据
func (hc *HookContext) GetAllData() map[string]any {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()

	result := make(map[string]any)
	for k, v := range hc.data {
		result[k] = v
	}
	return result
}
