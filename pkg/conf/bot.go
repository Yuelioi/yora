package conf

type BotConfig struct {
	*BaseConfig
	Port        string `json:"port"`         // Port to listen on
	LoggerLevel string `json:"logger_level"` // Logger level
	SelfID      string `json:"self_id"`      // Bot's ID
}

func NewBotConfig() *BotConfig {
	return &BotConfig{
		BaseConfig: NewBaseConfig(),
	}
}
