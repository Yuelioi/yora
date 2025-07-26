package depends

import (
	"context"
	"yora/internal/event"
)

func CommandArgs(cmds []string) Dependent {
	return DependentFunc(func(ctx context.Context, e event.Event) any {
		return cmds
	})
}
