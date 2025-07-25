package event

import (
	"yora/internal/event"
)

var _ event.MetaEvent = (*MetaEvent)(nil)

type MetaEvent struct {
	Event
}

// Extra implements event.MetaEvent.
func (m *MetaEvent) Extra() map[string]any {
	panic("unimplemented")
}

// Status implements event.MetaEvent.
func (m *MetaEvent) Status() map[string]any {
	panic("unimplemented")
}
