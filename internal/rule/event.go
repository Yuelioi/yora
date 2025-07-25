package rule

import (
	"context"
	"yora/internal/event"
)

func IsMetaEvent() RuleFunc {
	return func(ctx context.Context, event event.Event) bool {
		return event.Type() == "meta_event"
	}
}

func IsMessageEvent() RuleFunc {
	return func(ctx context.Context, event event.Event) bool {
		return event.Type() == "message"
	}
}

func IsNoticeEvent() RuleFunc {
	return func(ctx context.Context, event event.Event) bool {
		return event.Type() == "notice"
	}
}

func IsRequestEvent() RuleFunc {
	return func(ctx context.Context, event event.Event) bool {
		return event.Type() == "request"
	}
}

func IsCustomEvent(eventName string) RuleFunc {
	return func(ctx context.Context, event event.Event) bool {
		return event.Type() == eventName
	}
}
