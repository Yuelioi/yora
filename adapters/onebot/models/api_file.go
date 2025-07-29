package models

type FileInfo struct {
	BusID         int     `json:"busid"`
	DeadTime      int     `json:"dead_time"`
	DownloadTimes int     `json:"download_times"`
	FileID        string  `json:"file_id"`
	FileName      string  `json:"file_name"`
	FileSize      int     `json:"file_size"`
	GroupID       int64   `json:"group_id"`
	ModifyTime    float64 `json:"modify_time"`
	UploadTime    float64 `json:"upload_time"`
	Uploader      int64   `json:"uploader"`
	UploaderName  string  `json:"uploader_name"`
}

// 文件夹信息结构体
type FolderInfo struct {
	CreateName     string  `json:"create_name"`
	CreateTime     float64 `json:"create_time"`
	Creator        int64   `json:"creator"`
	FolderID       string  `json:"folder_id"`
	FolderName     string  `json:"folder_name"`
	GroupID        int64   `json:"group_id"`
	TotalFileCount int     `json:"total_file_count"`
}

// 根目录响应数据结构体
type DirectoryData struct {
	Files   []FileInfo   `json:"files"`
	Folders []FolderInfo `json:"folders"`
}

// URL 表示一个带有 URL 字段的结构体
type URL = struct {
	URL string `json:"url"`
}

// 获取私聊文件资源链接请求
type GetPrivateFileRequest struct {
	UserID   int    `json:"user_id"`             // 用户ID
	FileID   string `json:"file_id"`             // 文件ID
	FileHash string `json:"file_hash,omitempty"` // 文件哈希值(可选)
}

// 获取私聊文件资源链接响应
type GetPrivateFileResponse = Response[URL]

// 获取群文件资源链接请求
type GetGroupFileURLRequest struct {
	GroupID int    `json:"group_id"` // 群ID
	FileID  string `json:"file_id"`  // 文件ID
}

// 获取群文件资源链接响应
type GetGroupFileURLResponse = Response[URL]

// 获取群根目录文件列表请求
type GetGroupRootFilesRequest struct {
	GroupID int `json:"group_id"` // 群ID
}

// 获取群目录文件列表响应
type GetGroupFilesResponse = Response[DirectoryData]

// 获取群子目录文件列表请求
type GetGroupSubFilesRequest struct {
	GroupID  int    `json:"group_id"`  // 群ID
	FolderID string `json:"folder_id"` // 文件夹ID
}

// 移动群文件请求
type MoveGroupFileRequest struct {
	GroupID         int    `json:"group_id"`         // 群ID
	FileID          string `json:"file_id"`          // 文件ID
	ParentDirectory string `json:"parent_directory"` // 原目录
	TargetDirectory string `json:"target_directory"` // 目标目录
}

// 移动群文件响应
type MoveGroupFileResponse = Response[struct {
	// TODO: 添加具体的移动结果字段
}]

// 删除群文件请求
type DeleteGroupFileRequest struct {
	GroupID int    `json:"group_id"` // 群ID
	FileID  string `json:"file_id"`  // 文件ID
}

// 删除群文件响应
type DeleteGroupFileResponse = Response[any]

// 创建群文件文件夹请求
type CreateGroupFolderRequest struct {
	GroupID  int    `json:"group_id"`  // 群ID
	Name     string `json:"name"`      // 文件夹名称
	ParentID string `json:"parent_id"` // 父文件夹ID，tx不再允许在非根目录创建文件夹了，该值废弃，请直接传递"/"
}

// 创建群文件文件夹响应
type CreateGroupFolderResponse = Response[any]

// 删除群文件文件夹请求
type DeleteGroupFolderRequest struct {
	GroupID  int    `json:"group_id"`  // 群ID
	FolderID string `json:"folder_id"` // 文件夹ID
}

// 删除群文件文件夹响应
type DeleteGroupFolderResponse = Response[any]

// 重命名群文件文件夹请求
type RenameGroupFolderRequest struct {
	GroupID       int    `json:"group_id"`        // 群ID
	FolderID      string `json:"folder_id"`       // 文件夹ID
	NewFolderName string `json:"new_folder_name"` // 新文件夹名称
}

// 重命名群文件文件夹响应
type RenameGroupFolderResponse = Response[any]

// 上传群文件请求
type UploadGroupFileRequest struct {
	GroupID int    `json:"group_id"` // 群ID
	File    string `json:"file"`     // 本地文件路径
	Name    string `json:"name"`     // 文件名
	Folder  string `json:"folder"`   // 文件夹ID 默认 "/"
}

// 上传群文件响应
type UploadGroupFileResponse = Response[any]

// 上传私聊文件请求
type UploadPrivateFileRequest struct {
	UserID int    `json:"user_id"` // 用户ID
	File   string `json:"file"`    // 本地文件路径
	Name   string `json:"name"`    // 文件名
}

// 上传私聊文件响应
type UploadPrivateFileResponse = Response[any]
