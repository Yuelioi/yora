package matcher

import (
	"yora/internal/rule"
)

func On(rule rule.Rule, handlers ...*Handler) *Matcher {
	return NewMatcher(rule, handlers...)
}

// 事件
func OnMetaEvent(handlers ...*Handler) *Matcher {
	return NewMatcher(rule.IsMetaEvent(), handlers...)
}

func OnMessage(handlers ...*Handler) *Matcher {
	return NewMatcher(rule.IsMessageEvent(), handlers...)
}

func OnNotice(handlers ...*Handler) *Matcher {
	return NewMatcher(rule.IsNoticeEvent(), handlers...)
}

func OnRequest(handlers ...*Handler) *Matcher {
	return NewMatcher(rule.IsRequestEvent(), handlers...)
}

func OnCustomEvent(eventName string, handlers ...*Handler) *Matcher {
	return NewMatcher(rule.IsCustomEvent(eventName), handlers...)
}

// 命令
func OnStartsWith(prefix string, handlers ...*Handler) *Matcher {
	return NewMatcher(rule.StartsWith(prefix), handlers...)
}

func OnEndsWith(suffix string, handlers ...*Handler) *Matcher {
	return NewMatcher(rule.EndsWith(suffix), handlers...)
}

func OnFullMatch(pattern string, handlers ...*Handler) *Matcher {
	return NewMatcher(rule.FullMatch(pattern), handlers...)
}

func OnKeyword(keyword string, handlers ...*Handler) *Matcher {
	return NewMatcher(rule.Keyword(keyword), handlers...)
}

func OnCommand(cmds []string, caseSensitive bool, handlers ...*Handler) *Matcher {
	return NewMatcher(rule.Command(caseSensitive, cmds...), handlers...)
}

func OnRegex(pattern string, handlers ...*Handler) *Matcher {
	return NewMatcher(rule.Regex(pattern), handlers...)
}
