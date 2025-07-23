package adapter

import (
	"yora/internal/event"
	"yora/internal/matcher"
)

type Dispatcher interface {
	Dispatch(event event.Event)

	Subscribe(matcher matcher.Matcher) error

	Unsubscribe(matcher matcher.Matcher) error
}
