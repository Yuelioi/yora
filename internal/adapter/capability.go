package adapter

// Capabilities 协议能力集
type Capabilities struct {
	// SupportsGroupChat 是否支持群聊
	SupportsGroupChat bool

	// SupportsPrivateChat 是否支持私聊
	SupportsPrivateChat bool

	// SupportsFileUpload 是否支持文件上传
	SupportsFileUpload bool

	// SupportsRichText 是否支持富文本
	SupportsRichText bool

	// SupportsReply 是否支持回复消息
	SupportsReply bool

	// SupportsForward 是否支持转发消息
	SupportsForward bool

	// SupportsEdit 是否支持编辑消息
	SupportsEdit bool

	// SupportsDelete 是否支持删除消息
	SupportsDelete bool

	// SupportedSegmentTypes 支持的消息段类型
	SupportedSegmentTypes []string

	// MaxMessageLength 最大消息长度
	MaxMessageLength int

	// MaxFileSize 最大文件大小（字节）
	MaxFileSize int64

	// Extra 协议特定的其他能力
	Extra map[string]any
}
