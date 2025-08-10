package events

import (
	"time"
	"yora/pkg/event"
)

var _ event.RequestEvent = (*RequestEvent)(nil)

type RequestEvent struct {
	Event
}

// Comment implements event.RequestEvent.
func (r *RequestEvent) Comment() string {
	panic("unimplemented")
}

// Extra implements event.RequestEvent.
func (r *RequestEvent) Extra() map[string]any {
	panic("unimplemented")
}

// Flag implements event.RequestEvent.
func (r *RequestEvent) Flag() string {
	panic("unimplemented")
}

// Time implements event.RequestEvent.
// Subtle: this method shadows the method (Event).Time of RequestEvent.Event.
func (r *RequestEvent) Time() time.Time {
	panic("unimplemented")
}

// UserID implements event.RequestEvent.
// Subtle: this method shadows the method (Event).UserID of RequestEvent.Event.
func (r *RequestEvent) UserID() string {
	panic("unimplemented")
}
