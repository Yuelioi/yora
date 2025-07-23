package models

// Response 通用响应结构体
type Response struct {
	Status  string `json:"status"`
	Retcode int    `json:"retcode"`
	Data    any    `json:"data,omitempty"`
}

// PrivateFile 私聊文件相关结构体

// GetPrivateFileRequest 获取私聊文件资源链接请求
type GetPrivateFileRequest struct {
	UserID   string `json:"user_id"`   // 用户ID
	FileID   string `json:"file_id"`   // 文件ID
	FileHash string `json:"file_hash"` // 文件哈希值
}

// GetPrivateFileResponse 获取私聊文件资源链接响应
type GetPrivateFileResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		URL string `json:"url"` // 文件资源链接
	} `json:"data"`
}

// UploadPrivateFileRequest 上传私聊文件请求
type UploadPrivateFileRequest struct {
	UserID int    `json:"user_id"` // 用户ID
	File   string `json:"file"`    // 文件内容/路径
	Name   string `json:"name"`    // 文件名
}

// GroupFile 群文件相关结构体

// GetGroupFileRequest 获取群文件资源链接请求
type GetGroupFileRequest struct {
	GroupID int    `json:"group_id"` // 群ID
	FileID  string `json:"file_id"`  // 文件ID
}

// GetGroupFileResponse 获取群文件资源链接响应
type GetGroupFileResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		URL string `json:"url"` // 文件资源链接
	} `json:"data"`
}

// GetGroupRootFilesRequest 获取群根目录文件列表请求
type GetGroupRootFilesRequest struct {
	GroupID int `json:"group_id"` // 群ID
}

// GetGroupRootFilesResponse 获取群根目录文件列表响应
type GetGroupRootFilesResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		// TODO: 添加具体的文件列表字段
	} `json:"data"`
}

// GetGroupSubFilesRequest 获取群子目录文件列表请求
type GetGroupSubFilesRequest struct {
	GroupID  int    `json:"group_id"`  // 群ID
	FolderID string `json:"folder_id"` // 文件夹ID
}

// GetGroupSubFilesResponse 获取群子目录文件列表响应
type GetGroupSubFilesResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		// TODO: 添加具体的文件列表字段
	} `json:"data"`
}

// MoveGroupFileRequest 移动群文件请求
type MoveGroupFileRequest struct {
	GroupID         int    `json:"group_id"`         // 群ID
	FileID          string `json:"file_id"`          // 文件ID
	ParentDirectory string `json:"parent_directory"` // 原目录
	TargetDirectory string `json:"target_directory"` // 目标目录
}

// DeleteGroupFileRequest 删除群文件请求
type DeleteGroupFileRequest struct {
	GroupID int    `json:"group_id"` // 群ID
	FileID  string `json:"file_id"`  // 文件ID
}

// UploadGroupFileRequest 上传群文件请求
type UploadGroupFileRequest struct {
	GroupID int    `json:"group_id"` // 群ID
	File    string `json:"file"`     // 文件内容/路径
	Name    string `json:"name"`     // 文件名
	Folder  string `json:"folder"`   // 文件夹路径
}

// GroupFolder 群文件夹相关结构体

// CreateGroupFolderRequest 创建群文件文件夹请求
type CreateGroupFolderRequest struct {
	GroupID  int    `json:"group_id"`  // 群ID
	Name     string `json:"name"`      // 文件夹名称
	ParentID string `json:"parent_id"` // 父文件夹ID
}

// DeleteGroupFolderRequest 删除群文件文件夹请求
type DeleteGroupFolderRequest struct {
	GroupID  int    `json:"group_id"`  // 群ID
	FolderID string `json:"folder_id"` // 文件夹ID
}

// RenameGroupFolderRequest 重命名群文件文件夹请求
type RenameGroupFolderRequest struct {
	GroupID       int    `json:"group_id"`        // 群ID
	FolderID      string `json:"folder_id"`       // 文件夹ID
	NewFolderName string `json:"new_folder_name"` // 新文件夹名称
}
