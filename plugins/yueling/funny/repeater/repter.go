package repeater

import (
	"yora/adapters/onebot/events"
	"yora/pkg/bot"
	"yora/pkg/message"
	"yora/pkg/plugin"
)

var _ plugin.Plugin = (*repeater)(nil)

func New() plugin.Plugin {
	return &repeater{
		records: make(map[string]*repeatRecord),
	}
}

// 配置信息
var (
	MaxRepeat = "max_repeat"
)

// 自定义记录
type repeatRecord struct {
	lastMsg message.Message
	count   int
}

// repeater 插件
type repeater struct {
	records map[string]*repeatRecord
}

// Matchers implements plugin.Plugin.
func (r *repeater) Matchers() []*plugin.Matcher {
	panic("unimplemented")
}

func (r *repeater) Load() error {
	// 创建匹配器
	// rp := handler.NewHandler(r.record)
	// rm := on.OnMessage(rp)

	// 注册匹配器
	return nil
}

func (r *repeater) PluginInfo() *plugin.PluginInfo {
	return &plugin.PluginInfo{
		ID:          "repeater",
		Name:        "Repeater",
		Description: "Repeats the message",
		Version:     "0.1.0",
		Author:      "Yora",
		Usage:       "复读机",
		Examples:    []string{"!repeat", "!repeat 10", "!repeat 100"},
		Group:       "funny",
		Extra:       map[string]any{},
	}
}

func (r *repeater) record(evt *events.MessageEvent, bot bot.Bot) {
	// maxRepeat := r.Config().GetInt(MaxRepeat, 3)

	// groupID := evt.ChatID()
	// msg := evt.Message()

	// rec, exists := r.records[groupID]
	// if !exists {
	// 	r.records[groupID] = &repeatRecord{lastMsg: msg, count: 1}
	// 	return
	// }

	// if reflect.DeepEqual(rec.lastMsg, msg) {
	// 	rec.count++
	// 	if rec.count > maxRepeat {
	// 		bot.Send(evt.UserID(), groupID, msg)
	// 		rec.count = 1 // 重置计数
	// 	}
	// } else {
	// 	// 消息不同，重置记录
	// 	rec.lastMsg = msg
	// 	rec.count = 1
	// }

}
