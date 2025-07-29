package adapter

// Protocol 协议类型
type Protocol string

const (
	ProtocolOneBot   Protocol = "onebot"
	ProtocolTelegram Protocol = "telegram"
	ProtocolQQ       Protocol = "qq"
	ProtocolDiscord  Protocol = "discord"
	ProtocolWechat   Protocol = "wechat"
	ProtocolFeishu   Protocol = "feishu"
)

// 标准消息段类型常量
const (
	SegmentTypeText     = "text"
	SegmentTypeImage    = "image"
	SegmentTypeAudio    = "audio"
	SegmentTypeVideo    = "video"
	SegmentTypeFile     = "file"
	SegmentTypeAt       = "at"
	SegmentTypeEmoji    = "emoji"
	SegmentTypeReply    = "reply"
	SegmentTypeForward  = "forward"
	SegmentTypeLocation = "location"
	SegmentTypeContact  = "contact"
	SegmentTypeLink     = "link"
	SegmentTypeCode     = "code"
	SegmentTypeQuote    = "quote"
)

// 标准事件类型常量
const (
	EventTypeMessage = "message"
	EventTypeNotice  = "notice"
	EventTypeRequest = "request"
	EventTypeMeta    = "meta"
)

// 标准聊天类型常量
const (
	ChatTypePrivate    = "private"
	ChatTypeGroup      = "group"
	ChatTypeChannel    = "channel"
	ChatTypeSuperGroup = "supergroup"
	ChatTypeForum      = "forum"
)
