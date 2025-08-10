package provider

import (
	"context"
	"yora/pkg/event"
)

func CommandArgs(cmds []string) Provider {
	return DynamicProvider(func(ctx context.Context, e event.Event) any {
		return cmds
	})
}
