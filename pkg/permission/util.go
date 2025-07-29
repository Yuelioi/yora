package permission

import "yora/pkg/event"

func getRole(e event.Event) string {
	if msgEvent, ok := e.(event.GroupMessageEvent); ok {
		return msgEvent.Sender().Role()
	}
	return ""
}
