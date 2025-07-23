package onebot

import (
	"context"
	"regexp"
	"strings"
	"yora/internal/event"
	"yora/internal/matcher"
)

// Command 命令规则
func Command(cmd ...string) matcher.Rule {
	return matcher.RuleFunc(func(ctx context.Context, e event.Event) bool {
		if msgEvent, ok := e.(*Event); ok {
			message := strings.TrimSpace(msgEvent.RawMessage())
			for _, c := range cmd {
				if strings.HasPrefix(message, c) {
					return true
				}
			}
		}
		return false
	})
}

// Keyword 关键词规则
func Keyword(keywords ...string) matcher.Rule {
	return matcher.RuleFunc(func(ctx context.Context, e event.Event) bool {
		if msgEvent, ok := e.(*Event); ok {
			for _, keyword := range keywords {
				if strings.Contains(msgEvent.RawMessage(), keyword) {
					return true
				}
			}
		}
		return false
	})
}

// Regex 正则表达式规则
func Regex(pattern string) matcher.Rule {
	re := regexp.MustCompile(pattern)
	return matcher.RuleFunc(func(ctx context.Context, e event.Event) bool {
		if msgEvent, ok := e.(*Event); ok {
			return re.MatchString(msgEvent.RawMessage())
		}
		return false
	})
}

// StartsWith 前缀规则
func StartsWith(prefix ...string) matcher.Rule {
	return matcher.RuleFunc(func(ctx context.Context, e event.Event) bool {
		if msgEvent, ok := e.(*Event); ok {
			for _, p := range prefix {
				if strings.HasPrefix(msgEvent.RawMessage(), p) {
					return true
				}
			}
		}
		return false
	})
}

// EndsWith 后缀规则
func EndsWith(suffix ...string) matcher.Rule {
	return matcher.RuleFunc(func(ctx context.Context, e event.Event) bool {
		if msgEvent, ok := e.(*Event); ok {
			for _, s := range suffix {
				if strings.HasSuffix(msgEvent.RawMessage(), s) {
					return true
				}
			}
		}
		return false
	})
}
