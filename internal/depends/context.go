package depends

import (
	"context"
	"yora/internal/event"
)

func Ctx() Dependent {
	return DependentFunc(func(ctx context.Context, e event.Event) any {
		return ctx
	})
}
