package event

type Sender interface {
	// ID 唯一用户ID
	ID() string

	// Username 用户名（如果有）
	Username() string

	// DisplayName 显示名称/昵称
	DisplayName() string

	// AvatarURL 头像URL（如果有）
	AvatarURL() string

	// IsBot 是否为机器人
	IsBot() bool

	// IsAnonymous 是否为匿名用户
	IsAnonymous() bool

	// Raw 原始结构体
	Raw() any

	// Extra 平台特定的额外信息
	Extra() map[string]any

	Role() string
}

type GroupSender interface {
	Sender

	// ChatID 群聊ID
	ChatID() string

	// Role 群内角色：owner、admin、member 等
	Role() string

	// Extra 群内特定信息（如群名片、头衔等）
	Extra() map[string]any
}
