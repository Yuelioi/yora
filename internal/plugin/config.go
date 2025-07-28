package plugin

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"
)

var _ Config = (*BaseConfig)(nil)

// BaseConfig 默认配置实现
type BaseConfig struct {
	data      map[string]any
	defaults  map[string]any
	callbacks []func(string, any, any)
	mutex     sync.RWMutex
}

func NewBaseConfig() *BaseConfig {
	return &BaseConfig{
		data:     make(map[string]any),
		defaults: make(map[string]any),
	}
}

// 字段操作实现
func (c *BaseConfig) Set(key string, value any) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	oldValue, exists := c.data[key]
	c.data[key] = value

	// 触发变更回调
	for _, callback := range c.callbacks {
		if exists {
			callback(key, oldValue, value)
		} else {
			callback(key, nil, value)
		}
	}

	return nil
}

func (c *BaseConfig) Get(key string) (any, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if value, exists := c.data[key]; exists {
		return value, true
	}

	// 如果没有值，返回默认值
	if defaultValue, exists := c.defaults[key]; exists {
		return defaultValue, true
	}

	return nil, false
}

func (c *BaseConfig) GetString(key string, defaultValue string) string {
	value, exists := c.Get(key)
	if !exists {
		return defaultValue
	}

	if str, ok := value.(string); ok {
		return str
	}

	return defaultValue
}

func (c *BaseConfig) GetInt(key string, defaultValue int) int {
	value, exists := c.Get(key)
	if !exists {
		return defaultValue
	}

	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if v, err := strconv.Atoi(v); err == nil {
			return v
		}
		return defaultValue
	default:
		return defaultValue
	}
}

func (c *BaseConfig) GetBool(key string, defaultValue bool) bool {
	value, exists := c.Get(key)
	if !exists {
		return defaultValue
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		return defaultValue
	default:
		return defaultValue
	}
}

func (c *BaseConfig) GetFloat64(key string, defaultValue float64) float64 {
	value, exists := c.Get(key)
	if !exists {
		return defaultValue
	}

	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		if v, err := strconv.ParseFloat(v, 64); err == nil {
			return v
		}
		return defaultValue
	default:
		return defaultValue
	}
}

// 批量操作
func (c *BaseConfig) SetAll(data map[string]any) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, value := range data {
		c.data[key] = value
	}
	return nil
}

func (c *BaseConfig) GetAll() map[string]any {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	result := make(map[string]any)
	for key, value := range c.data {
		result[key] = value
	}
	return result
}

// 持久化
func (c *BaseConfig) SaveToJSON(filepath string) error {
	data := c.GetAll()
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, jsonData, 0644)
}

func (c *BaseConfig) LoadFromJSON(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	return c.FromJSON(data)
}

func (c *BaseConfig) ToJSON() ([]byte, error) {
	data := c.GetAll()
	return json.Marshal(data)
}

func (c *BaseConfig) FromJSON(data []byte) error {
	var configData map[string]any
	if err := json.Unmarshal(data, &configData); err != nil {
		return err
	}

	return c.SetAll(configData)
}

// 配置管理
func (c *BaseConfig) Has(key string) bool {
	_, exists := c.Get(key)
	return exists
}

func (c *BaseConfig) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.data[key]; exists {
		delete(c.data, key)
		return true
	}
	return false
}

func (c *BaseConfig) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]any)
}

func (c *BaseConfig) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}

// 验证和默认值
func (c *BaseConfig) SetDefault(key string, value any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.defaults[key] = value
}

func (c *BaseConfig) Validate() error {
	// 添加具体的验证逻辑

	return nil
}

// 配置变更监听
func (c *BaseConfig) OnChange(callback func(key string, oldValue, newValue any)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.callbacks = append(c.callbacks, callback)
}
