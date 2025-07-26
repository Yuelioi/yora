package models

// Friend 好友相关结构体

// GetFriendListResponse 获取好友列表响应
type GetFriendListResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    []struct {
		UserID   int    `json:"user_id"`  // 用户ID
		QID      string `json:"q_id"`     // QQ号
		Nickname string `json:"nickname"` // 昵称
		Remark   string `json:"remark"`   // 备注名
		Group    struct {
			GroupID   int    `json:"group_id"`   // 分组ID
			GroupName string `json:"group_name"` // 分组名称
		} `json:"group"` // 好友分组信息
	} `json:"data"` // 好友列表
}

// Group 群组相关结构体

// GetGroupInfoRequest 获取群信息请求
type GetGroupInfoRequest struct {
	GroupID int  `json:"group_id"` // 群ID
	NoCache bool `json:"no_cache"` // 是否不使用缓存
}

// GetGroupInfoResponse 获取群信息响应
type GetGroupInfoResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		// TODO: 添加群信息具体字段
	} `json:"data"` // 群信息数据
}

// GetGroupListRequest 获取群列表请求
type GetGroupListRequest struct {
	NoCache bool `json:"no_cache"` // 是否不使用缓存
}

// GetGroupListResponse 获取群列表响应
type GetGroupListResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    []struct {
		// TODO: 添加群列表具体字段
	} `json:"data"` // 群列表数据
}

// GetGroupMemberInfoRequest 获取群成员信息请求
type GetGroupMemberInfoRequest struct {
	GroupID int  `json:"group_id"` // 群ID
	UserID  int  `json:"user_id"`  // 用户ID
	NoCache bool `json:"no_cache"` // 是否不使用缓存
}

// GetGroupMemberInfoResponse 获取群成员信息响应
type GetGroupMemberInfoResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		// TODO: 添加群成员信息具体字段
	} `json:"data"` // 群成员信息数据
}

// GetGroupMemberListRequest 获取群成员列表请求
type GetGroupMemberListRequest struct {
	GroupID int `json:"group_id"` // 群ID
}

// GetGroupMemberListResponse 获取群成员列表响应
type GetGroupMemberListResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    []struct {
		// TODO: 添加群成员列表具体字段
	} `json:"data"` // 群成员列表数据
}

// User 用户相关结构体

// GetLoginInfoResponse 获取登录信息响应
type GetLoginInfoResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		UserID   int    `json:"user_id"`  // 用户ID
		Nickname string `json:"nickname"` // 昵称
	} `json:"data"` // 登录用户信息
}

// GetStrangerInfoResponse 获取陌生人信息响应
type GetStrangerInfoResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		UserID   int    `json:"user_id"`  // 用户ID
		Avatar   string `json:"avatar"`   // 头像链接
		QID      string `json:"q_id"`     // QQ号
		Nickname string `json:"nickname"` // 昵称
		Sign     string `json:"sign"`     // 个性签名
		Sex      string `json:"sex"`      // 性别
		Age      int    `json:"age"`      // 年龄
		Level    int    `json:"level"`    // 等级
		Status   struct {
			StatusID int    `json:"status_id"` // 状态ID
			FaceID   int    `json:"face_id"`   // 表情ID
			Message  string `json:"message"`   // 状态消息
		} `json:"status"` // 用户状态信息
		RegisterTime string `json:"RegisterTime"` // 注册时间
		Business     []struct {
			Type   int    `json:"type"`   // 业务类型
			Name   string `json:"name"`   // 业务名称
			Level  int    `json:"level"`  // 业务等级
			Icon   string `json:"icon"`   // 业务图标
			Ispro  int    `json:"ispro"`  // 是否专业版
			Isyear int    `json:"isyear"` // 是否年费
		} `json:"Business"` // 业务信息列表
	} `json:"data"` // 陌生人信息数据
}

// System 系统相关结构体

// GetStatusResponse 获取状态响应
type GetStatusResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		AppInitialized bool `json:"app_initialized"` // 应用是否已初始化
		AppEnabled     bool `json:"app_enabled"`     // 应用是否启用
		PluginsGood    bool `json:"plugins_good"`    // 插件是否正常
		AppGood        bool `json:"app_good"`        // 应用是否正常
		Online         bool `json:"online"`          // 是否在线
		Good           bool `json:"good"`            // 整体状态是否正常
		Memory         int  `json:"memory"`          // 内存使用量
	} `json:"data"` // 状态信息数据
}

// GetVersionInfoResponse 获取版本信息响应
type GetVersionInfoResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		AppName         string `json:"app_name"`         // 应用名称
		AppVersion      string `json:"app_version"`      // 应用版本
		ProtocolVersion string `json:"protocol_version"` // 协议版本
		NtProtocol      string `json:"nt_protocol"`      // NT协议版本
	} `json:"data"` // 版本信息数据
}
