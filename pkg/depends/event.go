package depends

import (
	"context"
	"yora/pkg/event"
)

func Event() Dependent {
	return DependentFunc(func(ctx context.Context, e event.Event) any {
		return e
	})
}

func MessageEvent() Dependent {
	return DependentFunc(func(ctx context.Context, e event.Event) any {
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

func MetaEvent() Dependent {
	return DependentFunc(func(ctx context.Context, e event.Event) any {
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

func NoticeEvent() Dependent {
	return DependentFunc(func(ctx context.Context, e event.Event) any {
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

func RequestEvent() Dependent {
	return DependentFunc(func(ctx context.Context, e event.Event) any {
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
