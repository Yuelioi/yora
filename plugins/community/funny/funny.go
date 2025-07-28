package funny

import (
	"yora/internal/plugin"
	"yora/plugins/community/funny/chat"
	"yora/plugins/community/funny/repeater"
)

var Plugins = []plugin.Plugin{}

func init() {

	Plugins = append(Plugins, repeater.New(), chat.New())

}
