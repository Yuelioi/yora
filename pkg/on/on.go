package on

import (
	"yora/pkg/handler"
	"yora/pkg/plugin"
	"yora/pkg/provider"
	"yora/pkg/rule"
)

func On(rule rule.Rule, handler *handler.Handler) *plugin.Matcher {
	return plugin.NewMatcher(rule, handler)
}

// 事件
func OnMetaEvent(handler *handler.Handler) *plugin.Matcher {
	return plugin.NewMatcher(rule.IsMetaEvent(), handler)
}

func OnMessage(handler *handler.Handler) *plugin.Matcher {
	return plugin.NewMatcher(rule.IsMessageEvent(), handler)
}

func OnNotice(handler *handler.Handler) *plugin.Matcher {
	return plugin.NewMatcher(rule.IsNoticeEvent(), handler)
}

func OnRequest(handler *handler.Handler) *plugin.Matcher {
	return plugin.NewMatcher(rule.IsRequestEvent(), handler)
}

func OnCustomEvent(eventName string, handler *handler.Handler) *plugin.Matcher {
	return plugin.NewMatcher(rule.IsCustomEvent(eventName), handler)
}

// 命令
func OnStartsWith(prefix string, handler *handler.Handler) *plugin.Matcher {
	return plugin.NewMatcher(rule.StartsWith(prefix), handler)
}

func OnEndsWith(suffix string, handler *handler.Handler) *plugin.Matcher {
	return plugin.NewMatcher(rule.EndsWith(suffix), handler)
}

func OnFullMatch(pattern string, handler *handler.Handler) *plugin.Matcher {
	return plugin.NewMatcher(rule.FullMatch(pattern), handler)
}

func OnKeyword(keyword string, handler *handler.Handler) *plugin.Matcher {
	return plugin.NewMatcher(rule.Keyword(keyword), handler)
}

func OnCommand(cmds []string, caseSensitive bool, handler *handler.Handler) *plugin.Matcher {
	handler.RegisterProviders(provider.CommandArgs(cmds))
	return plugin.NewMatcher(rule.Command(caseSensitive, cmds...), handler)
}

func OnRegex(pattern string, handler *handler.Handler) *plugin.Matcher {
	return plugin.NewMatcher(rule.Regex(pattern), handler)
}
