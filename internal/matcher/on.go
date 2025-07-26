package matcher

import (
	"yora/internal/depends"
	"yora/internal/rule"
)

func On(rule rule.Rule, handler *Handler) *Matcher {
	return NewMatcher(rule, handler)
}

// 事件
func OnMetaEvent(handler *Handler) *Matcher {
	return NewMatcher(rule.IsMetaEvent(), handler)
}

func OnMessage(handler *Handler) *Matcher {
	return NewMatcher(rule.IsMessageEvent(), handler)
}

func OnNotice(handler *Handler) *Matcher {
	return NewMatcher(rule.IsNoticeEvent(), handler)
}

func OnRequest(handler *Handler) *Matcher {
	return NewMatcher(rule.IsRequestEvent(), handler)
}

func OnCustomEvent(eventName string, handler *Handler) *Matcher {
	return NewMatcher(rule.IsCustomEvent(eventName), handler)
}

// 命令
func OnStartsWith(prefix string, handler *Handler) *Matcher {
	return NewMatcher(rule.StartsWith(prefix), handler)
}

func OnEndsWith(suffix string, handler *Handler) *Matcher {
	return NewMatcher(rule.EndsWith(suffix), handler)
}

func OnFullMatch(pattern string, handler *Handler) *Matcher {
	return NewMatcher(rule.FullMatch(pattern), handler)
}

func OnKeyword(keyword string, handler *Handler) *Matcher {
	return NewMatcher(rule.Keyword(keyword), handler)
}

func OnCommand(cmds []string, caseSensitive bool, handler *Handler) *Matcher {
	handler.RegisterDependent(depends.CommandArgs(cmds))
	return NewMatcher(rule.Command(caseSensitive, cmds...), handler)
}

func OnRegex(pattern string, handler *Handler) *Matcher {
	return NewMatcher(rule.Regex(pattern), handler)
}
