package api

import (
	"context"
	"yora/protocols/onebot/client"
	"yora/protocols/onebot/client/models"
)

var c = client.GetClient(context.Background())

// GetPrivateFile 获取私聊文件
func GetPrivateFile(userID string, fileID string, fileHash string) (*models.GetPrivateFileResponse, error) {
	req := models.GetPrivateFileRequest{
		UserID:   userID,
		FileID:   fileID,
		FileHash: fileHash,
	}
	resp, err := client.Call[models.GetPrivateFileRequest, models.GetPrivateFileResponse](c, "get_private_file", req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
