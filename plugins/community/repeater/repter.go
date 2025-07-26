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

var Repeater = &repeater{
	records: map[string][]record{},
}

type record struct {
	msg      message.Message
	repeated int
}

type repeater struct {
	plugin.BasePlugin
	records map[string][]record
}

type config struct {
	MaxRepeated int `json:"max_repeated"`
}

// Load implements plugin.Plugin.
func (r *repeater) Load() error {
	rp := matcher.NewHandler(r.record)
	rm := matcher.OnMessage(rp)

	r.SetMetadata(
		&plugin.Metadata{
			ID:          "repeater",
			Name:        "Repeater",
			Description: "Repeats the message",
			Version:     "0.1.0",
			Author:      "Yora",
			Usage:       "复读机",
			Examples:    []string{"!repeat", "!repeat 10", "!repeat 100"},
			Group:       "Community",
			Extra:       map[string]any{},
		},
	)
	r.RegisterMatcher(rm)
	return nil
}

func (r *repeater) record(evt *event.MessageEvent, bot bot.Bot) {
	groupID := evt.ChatID()

	msg := evt.Message()

	if recs, ok := r.records[groupID]; ok {
		if len(recs) > 0 {
			lastMsg := recs[len(recs)-1]
			if reflect.DeepEqual(lastMsg, msg) {
				bot.Send("message", evt.UserID(), groupID, msg)
			}
		}
		r.records[groupID] = append(recs, record{msg: msg})

	}
	r.records[groupID] = []record{{msg: msg, repeated: 1}}
}
