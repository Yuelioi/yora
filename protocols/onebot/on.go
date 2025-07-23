package onebot

import (
	"yora/internal/matcher"
)

func On(typ string, rule matcher.Rule, permission matcher.Permission, priority int, block bool, handlers ...matcher.Dependent) matcher.Matcher {
	return &Matcher{
		typ:        typ,
		priority:   priority,
		block:      block,
		rule:       rule,
		permission: permission,
		Handlers:   handlers,
	}

}

func OnCommand(cmd string, permission matcher.Permission, priority int, block bool, handlers ...matcher.Dependent) matcher.Matcher {

	return On("message", Command(cmd), permission, priority, block, handlers...)
}
