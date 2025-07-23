package onebot

import (
	"context"
	"fmt"
	"sync"
	"yora/internal/event"
	"yora/internal/matcher"
)

type Dispatcher struct {
	mu       sync.RWMutex
	matchers []matcher.Matcher
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Subscribe(m matcher.Matcher) error {
	if m == nil {
		return fmt.Errorf("matcher cannot be nil")
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.matchers = append(d.matchers, m)
	return nil
}

func (d *Dispatcher) Unsubscribe(m matcher.Matcher) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, matcher := range d.matchers {
		if matcher == m {
			d.matchers = append(d.matchers[:i], d.matchers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("matcher not found")
}

func (d *Dispatcher) Dispatch(ctx context.Context, e event.Event) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var lastErr error
	for _, matcher := range d.matchers {
		if !matcher.CanHandle(ctx, e) {
			continue
		}
		if err := matcher.Handle(ctx, e); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
