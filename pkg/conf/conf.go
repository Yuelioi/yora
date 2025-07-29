package conf

// 配置接口
type Config interface {
	// 字段操作
	Set(key string, value any) error
	Get(key string) (any, bool)

	GetString(key string, defaultValue string) string
	GetInt(key string, defaultValue int) int
	GetBool(key string, defaultValue bool) bool
	GetFloat64(key string, defaultValue float64) float64

	// 批量操作
	SetAll(data map[string]any) error
	GetAll() map[string]any

	// 持久化
	SaveToJSON(filepath string) error
	LoadFromJSON(filepath string) error
	ToJSON() ([]byte, error)
	FromJSON(data []byte) error

	// 配置管理
	Has(key string) bool
	Delete(key string) bool
	Clear()
	Keys() []string

	// 验证和默认值
	SetDefault(key string, value any)
	Validate() error

	// 配置变更监听
	OnChange(callback func(key string, oldValue, newValue any))
}
