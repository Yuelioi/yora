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

// 插件
const (
	PluginBeforeInit   HookType = "plugin.before_init"
	PluginAfterInit    HookType = "plugin.after_init"
	PluginBeforeLoad   HookType = "plugin.before_load"
	PluginAfterLoad    HookType = "plugin.after_load"
	PluginBeforeUnload HookType = "plugin.before_unload"
	PluginAfterUnload  HookType = "plugin.after_unload"
	PluginOnError      HookType = "plugin.on_error"
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
