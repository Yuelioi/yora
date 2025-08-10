package provider

import (
	"context"
	"yora/pkg/event"
)

func Event() Provider {
	return DynamicProvider(func(ctx context.Context, e event.Event) any {
		return e
	})
}

func Ctx() Provider {
	return DynamicProvider(func(ctx context.Context, e event.Event) any {
		return ctx
	})
}

func MessageEvent() Provider {
	return DynamicProvider(func(ctx context.Context, e event.Event) any {
		if e.Type() == "message" {
			e2, ok := e.(event.MessageEvent)
			if ok {
				return e2
			}
			return e
		}
		return nil
	})
}

func MetaEvent() Provider {
	return DynamicProvider(func(ctx context.Context, e event.Event) any {
		if e.Type() == "meta_event" {
			e2, ok := e.(event.MetaEvent)
			if ok {
				return e2
			}
			return e
		}
		return nil
	})
}

func NoticeEvent() Provider {
	return DynamicProvider(func(ctx context.Context, e event.Event) any {
		if e.Type() == "notice" {
			e2, ok := e.(event.NoticeEvent)
			if ok {
				return e2
			}
			return e
		}
		return nil
	})
}

func RequestEvent() Provider {
	return DynamicProvider(func(ctx context.Context, e event.Event) any {
		if e.Type() == "request" {
			e2, ok := e.(event.RequestEvent)
			if ok {
				return e2
			}
			return e
		}
		return nil
	})
}
