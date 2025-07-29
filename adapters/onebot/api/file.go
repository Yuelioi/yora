package api

import (
	"yora/adapters/onebot/client"
	"yora/adapters/onebot/models"
)

// 获取私聊文件资源链接
//
// 参数：
//   - userID: 用户ID
//   - fileID: 文件ID
//   - fileHash: 文件哈希值（可选）
func (api *API) GetPrivateFile(userID int, fileID string, fileHash string) (*models.GetPrivateFileResponse, error) {
	req := models.GetPrivateFileRequest{
		UserID:   userID,
		FileID:   fileID,
		FileHash: fileHash,
	}
	return client.Call[models.GetPrivateFileRequest, models.GetPrivateFileResponse](api.client, "get_private_file", req)
}

// 获取群文件资源链接
//
// 参数：
//   - groupID: 群ID
//   - fileID: 文件ID
func (api *API) GetGroupFileURL(groupID int, fileID string) (*models.GetGroupFileURLResponse, error) {
	req := models.GetGroupFileURLRequest{
		GroupID: groupID,
		FileID:  fileID,
	}
	return client.Call[models.GetGroupFileURLRequest, models.GetGroupFileURLResponse](api.client, "get_group_file_url", req)
}

// 获取群根目录文件列表
//
// 参数：
//   - groupID: 群ID
func (api *API) GetGroupRootFiles(groupID int) (*models.GetGroupFilesResponse, error) {
	req := models.GetGroupRootFilesRequest{
		GroupID: groupID,
	}
	return client.Call[models.GetGroupRootFilesRequest, models.GetGroupFilesResponse](api.client, "get_group_root_files", req)
}

// 获取群子目录文件列表
//
// 参数：
//   - groupID: 群ID
//   - folderID: 文件夹ID
func (api *API) GetGroupSubFiles(groupID int, folderID string) (*models.GetGroupFilesResponse, error) {
	req := models.GetGroupSubFilesRequest{
		GroupID:  groupID,
		FolderID: folderID,
	}
	return client.Call[models.GetGroupSubFilesRequest, models.GetGroupFilesResponse](api.client, "get_group_files_by_folder", req)
}

// 移动群文件
//
// 参数：
//   - groupID: 群ID
//   - fileID: 文件ID
//   - parentDirectory: 当前文件所在目录ID
//   - targetDirectory: 目标目录ID
func (api *API) MoveGroupFile(groupID int, fileID string, parentDirectory string, targetDirectory string) (*models.MoveGroupFileResponse, error) {
	req := models.MoveGroupFileRequest{
		GroupID:         groupID,
		FileID:          fileID,
		ParentDirectory: parentDirectory,
		TargetDirectory: targetDirectory,
	}
	return client.Call[models.MoveGroupFileRequest, models.MoveGroupFileResponse](api.client, "move_group_file", req)
}

// 删除群文件
//
// 参数：
//   - groupID: 群ID
//   - fileID: 文件ID
func (api *API) DeleteGroupFile(groupID int, fileID string) (*models.DeleteGroupFileResponse, error) {
	req := models.DeleteGroupFileRequest{
		GroupID: groupID,
		FileID:  fileID,
	}
	return client.Call[models.DeleteGroupFileRequest, models.DeleteGroupFileResponse](api.client, "delete_group_file", req)
}

// 创建群文件夹
//
// 参数：
//   - groupID: 群ID
//   - name: 文件夹名称
func (api *API) CreateGroupFolder(groupID int, name string) (*models.CreateGroupFolderResponse, error) {
	req := models.CreateGroupFolderRequest{
		GroupID:  groupID,
		Name:     name,
		ParentID: "/",
	}
	return client.Call[models.CreateGroupFolderRequest, models.CreateGroupFolderResponse](api.client, "create_group_file_folder", req)
}

// 删除群文件夹
//
// 参数：
//   - groupID: 群ID
//   - folderID: 文件夹ID
func (api *API) DeleteGroupFolder(groupID int, folderID string) (*models.DeleteGroupFolderResponse, error) {
	req := models.DeleteGroupFolderRequest{
		GroupID:  groupID,
		FolderID: folderID,
	}
	return client.Call[models.DeleteGroupFolderRequest, models.DeleteGroupFolderResponse](api.client, "delete_group_file_folder", req)
}

// 重命名群文件夹
//
// 参数：
//   - groupID: 群ID
//   - folderID: 文件夹ID
//   - newFolderName: 新文件夹名称
func (api *API) RenameGroupFolder(groupID int, folderID string, newFolderName string) (*models.RenameGroupFolderResponse, error) {
	req := models.RenameGroupFolderRequest{
		GroupID:       groupID,
		FolderID:      folderID,
		NewFolderName: newFolderName,
	}
	return client.Call[models.RenameGroupFolderRequest, models.RenameGroupFolderResponse](api.client, "rename_group_file_folder", req)
}

// 上传群文件
//
// 参数：
//   - groupID: 群ID
//   - file: 文件链接, 本地绝对文件路径
//   - name: 文件名
//   - folder: 文件夹ID（默认值 "/" 表示根目录）
func (api *API) UploadGroupFile(groupID int, file string, name string, folder string) (*models.UploadGroupFileResponse, error) {
	req := models.UploadGroupFileRequest{
		GroupID: groupID,
		File:    file,
		Name:    name,
		Folder:  folder,
	}
	return client.Call[models.UploadGroupFileRequest, models.UploadGroupFileResponse](api.client, "upload_group_file", req)
}

// 上传私聊文件
//
// 参数：
//   - userID: 用户ID
//   - file: 文件链接, 本地绝对文件路径
//   - name: 文件名
func (api *API) UploadPrivateFile(userID int, file string, name string) (*models.UploadPrivateFileResponse, error) {
	req := models.UploadPrivateFileRequest{
		UserID: userID,
		File:   file,
		Name:   name,
	}
	return client.Call[models.UploadPrivateFileRequest, models.UploadPrivateFileResponse](api.client, "upload_private_file", req)
}
