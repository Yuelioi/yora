package api

import (
	"errors"
	"strings"
	"testing"
	"yora/protocols/onebot/client/models"

	"github.com/stretchr/testify/assert"
)

// ! 不知道如何拿到私有文件id
func TestGetPrivateFile(t *testing.T) {
	h := NewTestHelper(t)

	h.t.Skip("私有文件API未实现")

	resp, err := h.api.GetPrivateFile(UID, "", "")
	h.StatusOk(resp, err, "获取私有文件")

	assert.True(t, strings.HasPrefix(resp.Data.URL, "http"), "URL格式验证")

}

type TestFileHelper struct {
	*TestHelper
}

func NewTestFileHelper(t *testing.T) *TestFileHelper {
	th := NewTestHelper(t)
	return &TestFileHelper{
		TestHelper: th,
	}
}

// 创建临时测试文件夹并返回清理函数
func (h *TestFileHelper) createTempFolder(folderName string) (string, func()) {
	_, err := h.api.CreateGroupFolder(GID, folderName)
	if err != nil {
		h.t.Logf("创建临时文件夹失败，可能已存在: %v", err)
	}

	folderID, err := h.getFolderIDByName(GID, folderName)
	if err != nil {
		h.t.Fatalf("获取临时文件夹ID失败: %v", err)
	}

	cleanup := func() {
		if _, err := h.api.DeleteGroupFolder(GID, folderID); err != nil {
			h.t.Logf("清理临时文件夹失败: %v", err)
		}
	}

	return folderID, cleanup
}

// 上传临时测试文件并返回清理函数
func (h *TestFileHelper) uploadTempFile(fileName, folder string) (string, func()) {
	absFile, err := h.getTestFilePath()
	if err != nil {
		h.t.Fatalf("获取文件绝对路径失败: %v", err)
	}

	_, err = h.api.UploadGroupFile(GID, absFile, fileName, folder)
	if err != nil {
		h.t.Fatalf("上传临时文件失败: %v", err)
	}

	fileID, err := h.getFileIDByName(GID, fileName)
	if err != nil {
		h.t.Fatalf("获取临时文件ID失败: %v", err)
	}

	cleanup := func() {
		if _, err := h.api.DeleteGroupFile(GID, fileID); err != nil {
			h.t.Logf("清理临时文件失败: %v", err)
		}
	}

	return fileID, cleanup
}

// 递归搜索文件ID（优化后的版本）
func (h *TestFileHelper) getFileIDByName(groupID int, filename string) (string, error) {
	return h.searchFileInDirectory(groupID, "", filename)
}

// 在指定目录中搜索文件
func (h *TestFileHelper) searchFileInDirectory(groupID int, folderID, filename string) (string, error) {
	var resp interface{}

	// 根据是否有folderID决定调用哪个API
	if folderID == "" {
		respRoot, err := h.api.GetGroupRootFiles(groupID)
		if err != nil {
			return "", err
		}
		resp = respRoot
	} else {
		respSub, err := h.api.GetGroupSubFiles(groupID, folderID)
		if err != nil {
			return "", err
		}
		resp = respSub
	}

	// 类型断言获取文件和文件夹列表
	var files []models.FileInfo
	var folders []models.FolderInfo

	switch r := resp.(type) {
	case *models.GetGroupFilesResponse:
		files = r.Data.Files
		folders = r.Data.Folders

	default:
		return "", errors.New("未知的响应类型")
	}

	// 在当前目录的文件中搜索
	for _, file := range files {
		if file.FileName == filename {
			return file.FileID, nil
		}
	}

	// 递归搜索子文件夹
	for _, folder := range folders {
		if fileID, err := h.searchFileInDirectory(groupID, folder.FolderID, filename); err == nil {
			return fileID, nil
		}
	}

	return "", errors.New("文件不存在")
}

// 获取文件夹ID
func (h *TestFileHelper) getFolderIDByName(groupID int, folderName string) (string, error) {
	resp, err := h.api.GetGroupRootFiles(groupID)
	if err != nil {
		return "", err
	}

	for _, folder := range resp.Data.Folders {
		if folder.FolderName == folderName {
			return folder.FolderID, nil
		}
	}
	return "", errors.New("文件夹不存在")
}

// 测试用例

func TestGetGroupFileURL(t *testing.T) {
	h := NewTestFileHelper(t)

	// 上传临时文件
	fileName := "test_get_url.jpg"
	fileID, cleanup := h.uploadTempFile(fileName, "/")
	defer cleanup()

	// 测试获取文件URL
	resp, err := h.api.GetGroupFileURL(GID, fileID)
	h.StatusOk(resp, err, "获取群文件URL")
	assert.True(t, strings.HasPrefix(resp.Data.URL, "http"), "URL格式验证")

}

func TestGetGroupRootFiles(t *testing.T) {
	h := NewTestFileHelper(t)

	resp, err := h.api.GetGroupRootFiles(GID)
	h.StatusOk(resp, err, "获取群根文件")

}

func TestGetGroupSubFiles(t *testing.T) {
	h := NewTestFileHelper(t)

	// 创建临时文件夹
	folderName := "test_sub_files"
	folderID, cleanup := h.createTempFolder(folderName)
	defer cleanup()

	// 测试获取子文件
	resp, err := h.api.GetGroupSubFiles(GID, folderID)
	h.StatusOk(resp, err, "获取群子文件")

}

func TestMoveGroupFile(t *testing.T) {
	h := NewTestFileHelper(t)

	// 创建目标文件夹
	folderName := "test_move_target"
	folderID, cleanupFolder := h.createTempFolder(folderName)
	defer cleanupFolder()

	// 上传文件到根目录
	fileName := "test_move.jpg"
	fileID, cleanupFile := h.uploadTempFile(fileName, "/")
	defer cleanupFile()

	// 移动文件
	resp, err := h.api.MoveGroupFile(GID, fileID, "", folderID)

	h.StatusOk(resp, err, "移动群文件")

}

func TestDeleteGroupFile(t *testing.T) {
	h := NewTestFileHelper(t)

	// 上传临时文件
	fileName := "test_delete.jpg"
	fileID, _ := h.uploadTempFile(fileName, "/") // 不使用cleanup，因为我们要手动删除

	// 删除文件
	resp, err := h.api.DeleteGroupFile(GID, fileID)
	h.StatusOk(resp, err, "删除群文件")

}

func TestCreateGroupFolder(t *testing.T) {
	h := NewTestFileHelper(t)

	folderName := "test_create_folder"

	// 创建文件夹
	resp, err := h.api.CreateGroupFolder(GID, folderName)
	h.StatusOk(resp, err, "创建群文件夹")

	// 清理：删除创建的文件夹
	if folderID, err := h.getFolderIDByName(GID, folderName); err == nil {
		rp, err := h.api.DeleteGroupFolder(GID, folderID)
		h.StatusOk(rp, err, "删除群文件夹")
	}
}

func TestDeleteGroupFolder(t *testing.T) {
	h := NewTestFileHelper(t)

	// 创建临时文件夹
	folderName := "test_delete_folder"
	folderID, _ := h.createTempFolder(folderName) // 不使用cleanup，因为我们要手动删除

	// 删除文件夹
	resp, err := h.api.DeleteGroupFolder(GID, folderID)
	h.StatusOk(resp, err, "删除群文件夹")

}

func TestRenameGroupFolder(t *testing.T) {
	h := NewTestFileHelper(t)

	// 创建临时文件夹
	oldName := "test_rename_old"
	newName := "test_rename_new"
	folderID, cleanup := h.createTempFolder(oldName)
	defer cleanup()

	// 重命名文件夹
	resp, err := h.api.RenameGroupFolder(GID, folderID, newName)
	h.StatusOk(resp, err, "重命名群文件夹")

}

func TestUploadGroupFile(t *testing.T) {
	h := NewTestFileHelper(t)

	absFile, err := h.getTestFilePath()
	if err != nil {
		t.Fatalf("获取文件绝对路径失败: %v", err)
	}

	fileName := "test_upload.jpg"

	// 上传文件
	resp, err := h.api.UploadGroupFile(GID, absFile, fileName, "/")
	h.StatusOk(resp, err, "上传群文件")

	// 清理：删除上传的文件
	if fileID, err := h.getFileIDByName(GID, fileName); err == nil {
		h.api.DeleteGroupFile(GID, fileID)
	}

}

func TestUploadPrivateFile(t *testing.T) {

	h := NewTestFileHelper(t)
	h.t.Skip("私有文件API未实现")

	absFile, err := h.getTestFilePath()
	if err != nil {
		t.Fatalf("获取文件绝对路径失败: %v", err)
	}

	fileName := "test_private_upload.jpg"

	// 上传私有文件
	resp, err := h.api.UploadPrivateFile(UID, absFile, fileName)

	h.StatusOk(resp, err, "上传私有文件")
	if err == nil {
		t.Logf("上传私有文件成功: %+v", resp)
		// !注意：私有文件的清理不了
		// !可能需要撤回消息?
	}
}
