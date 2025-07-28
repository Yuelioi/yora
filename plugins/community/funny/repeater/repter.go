package repeater

import (
	"reflect"
	"yora/internal/bot"
	"yora/internal/matcher"
	"yora/internal/message"
	"yora/internal/plugin"
	"yora/protocols/onebot/event"
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
	plugin.BasePlugin // 继承基础插件

	records map[string]*repeatRecord
}

func (r *repeater) Load() error {
	// 创建匹配器
	rp := matcher.NewHandler(r.record)
	rm := matcher.OnMessage(rp)

	// 设置元数据
	r.SetMetadata(
		&plugin.Metadata{
			ID:          "repeater",
			Name:        "Repeater",
			Description: "Repeats the message",
			Version:     "0.1.0",
			Author:      "Yora",
			Usage:       "复读机",
			Examples:    []string{"!repeat", "!repeat 10", "!repeat 100"},
			Group:       "funny",
			Extra:       map[string]any{},
		},
	)

	// 注册匹配器
	r.RegisterMatcher(rm)
	return nil
}

func (r *repeater) record(evt *event.MessageEvent, bot bot.Bot) {
	maxRepeat := r.Config().GetInt(MaxRepeat, 3)

	groupID := evt.ChatID()
	msg := evt.Message()

	rec, exists := r.records[groupID]
	if !exists {
		r.records[groupID] = &repeatRecord{lastMsg: msg, count: 1}
		return
	}

	if reflect.DeepEqual(rec.lastMsg, msg) {
		rec.count++
		if rec.count > maxRepeat {
			bot.Send(evt.UserID(), groupID, msg)
			rec.count = 1 // 重置计数
		}
	} else {
		// 消息不同，重置记录
		rec.lastMsg = msg
		rec.count = 1
	}

}
