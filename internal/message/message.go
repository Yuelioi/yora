package message

// Message 消息接口
type Message interface {
	// Segments 获取所有消息段
	Segments() []Segment

	// String 转为用户可读的消息文本
	String() string

	// PlainText 提取纯文本内容
	PlainText() string

	// IsEmpty 判断消息是否为空
	IsEmpty() bool

	// HasType 判断是否包含指定类型的消息段
	HasType(segmentType string) bool

	// GetSegmentsByType 获取指定类型的所有消息段
	GetSegmentsByType(segmentType string) []Segment

	// Append 添加消息段
	Append(seg Segment) Message
}

// Segment 消息段接口
type Segment interface {
	// Type 类型，如 "text", "image", "at"
	Type() string

	// Data 原始数据，用于构造或序列化
	Data() map[string]any

	// String 可读表示（用于日志或调试）
	String() string

	// IsType 判断是否为指定类型
	IsType(segmentType string) bool

	// GetData 获取指定键的数据
	GetData(key string) (any, bool)
}

type Sender interface {
	// ID 唯一用户ID
	ID() string

	// Username 用户名（如果有）
	Username() string

	// DisplayName 显示名称/昵称
	DisplayName() string

	// AvatarURL 头像URL（如果有）
	AvatarURL() string

	// IsAnonymous 是否为匿名用户
	IsAnonymous() bool

	// Raw 原始结构体
	Raw() any

	// Role 群内角色：owner、admin、member...
	Role() string

	// Extra 平台特定的额外信息
	Extra() map[string]any
}
