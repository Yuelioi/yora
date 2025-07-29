package conf

var _ Config = (*PluginConfig)(nil)

type PluginConfig struct {
	*BaseConfig
}

func NewPluginConfig() *PluginConfig {
	return &PluginConfig{
		BaseConfig: NewBaseConfig(),
	}
}
