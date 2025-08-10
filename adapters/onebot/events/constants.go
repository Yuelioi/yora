package events

type PostType string
type MessageType string
type SubType string

const (
	PostTypeMessage   PostType = "message"
	PostTypeNotice    PostType = "notice"
	PostTypeMetaEvent PostType = "meta_event"
	PostTypeRequest   PostType = "request"
)

const (
	MessageTypeGroup   MessageType = "group"
	MessageTypePrivate MessageType = "private"
)
const (
	SubTypeNormal  SubType = "normal"
	SubTypeConnect SubType = "connect"
)
