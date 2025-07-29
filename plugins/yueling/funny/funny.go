package funny

import (
	"yora/pkg/plugin"
	"yora/plugins/yueling/funny/chat"
	"yora/plugins/yueling/funny/repeater"
)

var Plugins = []plugin.Plugin{}

func init() {

	Plugins = append(Plugins, repeater.New(), chat.New())

}
