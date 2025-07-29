package adapter

// 协议能力集
type Capabilities struct {
	// 是否支持群聊
	SupportsGroupChat bool

	// 是否支持私聊
	SupportsPrivateChat bool

	// 是否支持文件上传
	SupportsFileUpload bool

	// 是否支持富文本
	SupportsRichText bool

	// 是否支持回复消息
	SupportsReply bool

	// 是否支持转发消息
	SupportsForward bool

	// 是否支持编辑消息
	SupportsEdit bool

	// 是否支持删除消息
	SupportsDelete bool

	// 支持的消息段类型
	SupportedSegmentTypes []string

	// 最大消息长度
	MaxMessageLength int

	// 最大文件大小（字节）
	MaxFileSize int64

	// 协议特定的其他能力
	Extra map[string]any
}
