package hook

// HookType Hook类型定义
type HookType string

// BOT
const (
	BotOnStart     HookType = "bot.on_start"
	BotOnStop      HookType = "bot.on_stop"
	BotOnReload    HookType = "bot.on_reload"
	BotHealthCheck HookType = "bot.health_check"
)

// Plugin
const (
	PluginOnStart        HookType = "plugin.on_start"
	PluginOnStop         HookType = "plugin.on_stop"
	PluginOnReload       HookType = "plugin.on_reload"
	PluginOnConfigLoad   HookType = "plugin.on_config_load"
	PluginOnConfigChange HookType = "plugin.on_config_change"
)

// Matcher
const (
	MatcherBeforeMatch  HookType = "matcher.before_match"
	MatcherAfterMatch   HookType = "matcher.after_match"
	MatcherBeforeHandle HookType = "matcher.before_handle"
	MatcherAfterHandle  HookType = "matcher.after_handle"
	MatcherOnError      HookType = "matcher.on_error"
)

// 消息处理
const (
	MessageBeforeProcess HookType = "message.before_process"
	MessageAfterProcess  HookType = "message.after_process"
	MessageOnReceive     HookType = "message.on_receive"
	MessageOnSend        HookType = "message.on_send"
	MessageOnMatch       HookType = "message.on_match"
	MessageOnNoMatch     HookType = "message.on_no_match"
)

// HookPriority Hook优先级
type HookPriority int

const (
	PriorityHighest HookPriority = 100
	PriorityHigh    HookPriority = 75
	PriorityNormal  HookPriority = 50
	PriorityLow     HookPriority = 25
	PriorityLowest  HookPriority = 0
)
