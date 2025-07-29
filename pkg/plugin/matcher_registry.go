package plugin

import (
	"context"
	"sort"
	"sync"
	"yora/pkg/event"
	"yora/pkg/log"

	"github.com/rs/zerolog"
)

var (
	once            sync.Once
	matcherRegistry *MatcherRegistry
)

// 匹配器注册中心
type MatcherRegistry struct {
	matchers []*Matcher
	logger   zerolog.Logger
	mu       sync.RWMutex
}

func GetMatcherRegistry() *MatcherRegistry {
	once.Do(func() {
		matcherRegistry = newMatcherRegistry()
	})
	return matcherRegistry
}

func newMatcherRegistry() *MatcherRegistry {
	return &MatcherRegistry{
		matchers: make([]*Matcher, 0),
		logger:   log.NewHandler("MatcherRegistry"),
	}
}

func (mr *MatcherRegistry) RegisterMatchers(ms ...*Matcher) {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	mr.matchers = append(mr.matchers, ms...)

	sort.Slice(mr.matchers, func(i, j int) bool {
		return mr.matchers[i].Priority < mr.matchers[j].Priority
	})
}

// 匹配事件并缓存匹配到的matcher
func (mr *MatcherRegistry) MatchedMatchers(ctx context.Context, evt event.Event) []*Matcher {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	matched := make([]*Matcher, 0, len(mr.matchers))

	for _, m := range mr.matchers {
		if m.Rule.Match(ctx, evt) {
			matched = append(matched, m)
		}
		if m.Block {
			break
		}
	}

	return matched
}
