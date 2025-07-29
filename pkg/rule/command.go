package rule

import (
	"context"
	"regexp"
	"strings"
	"yora/pkg/event"
)

func getRawMessage(e event.Event) string {
	if msgEvent, ok := e.(event.MessageEvent); ok {
		return msgEvent.RawMessage()
	}
	return ""
}

// StartsWith 前缀规则
func StartsWith(prefix ...string) Rule {
	return RuleFunc(func(ctx context.Context, e event.Event) bool {
		if msgEvent, ok := e.(event.MessageEvent); ok {
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
func EndsWith(suffix ...string) Rule {
	return RuleFunc(func(ctx context.Context, e event.Event) bool {
		if msgEvent, ok := e.(event.MessageEvent); ok {
			for _, s := range suffix {
				if strings.HasSuffix(msgEvent.RawMessage(), s) {
					return true
				}
			}
		}
		return false
	})
}

func FullMatch(s string) Rule {
	return RuleFunc(func(ctx context.Context, e event.Event) bool {
		if msgEvent, ok := e.(event.MessageEvent); ok {
			return msgEvent.RawMessage() == s
		}
		return false
	})
}

// Keyword 关键词规则
func Keyword(keywords ...string) Rule {
	return RuleFunc(func(ctx context.Context, e event.Event) bool {
		if msgEvent, ok := e.(event.MessageEvent); ok {
			for _, keyword := range keywords {
				if strings.Contains(msgEvent.RawMessage(), keyword) {
					return true
				}
			}
		}
		return false
	})
}

// Command 命令规则
func Command(caseSensitive bool, cmds ...string) Rule {
	return RuleFunc(func(ctx context.Context, e event.Event) bool {

		if msgEvent, ok := e.(event.MessageEvent); ok {
			message := strings.TrimSpace(msgEvent.RawMessage())
			if !caseSensitive {
				message = strings.ToLower(message)
			}
			for _, cmd := range cmds {
				match := cmd
				if !caseSensitive {
				}
				if strings.HasPrefix(message, match) {
					return true
				}
			}
		}
		return false
	})
}

// Regex 正则表达式规则
func Regex(pattern string) Rule {
	re := regexp.MustCompile(pattern)
	return RuleFunc(func(ctx context.Context, e event.Event) bool {
		if msgEvent, ok := e.(event.MessageEvent); ok {
			return re.MatchString(msgEvent.RawMessage())
		}
		return false
	})
}
