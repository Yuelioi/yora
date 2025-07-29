package api

import (
	"yora/adapters/onebot/client"
	"yora/adapters/onebot/models"
)

func (api *API) GetFriendList() (*models.GetFriendListResponse, error) {
	req := models.GetFriendListRequest{}
	return client.Call[models.GetFriendListRequest, models.GetFriendListResponse](api.client, "get_friend_list", req)

}

// GetGroupInfo 获取群信息
//
// 参数：
//   - groupID: 群号
//   - noCache: 是否不使用缓存（true 表示跳过缓存，实时获取）
func (api *API) GetGroupInfo(groupID int, noCache bool) (*models.GetGroupInfoResponse, error) {
	req := models.GetGroupInfoRequest{
		GroupID: groupID,
		NoCache: noCache,
	}
	return client.Call[models.GetGroupInfoRequest, models.GetGroupInfoResponse](api.client, "get_group_info", req)

}

// GetGroupMemberList 获取群成员列表
//
// 参数：
//   - groupID: 群号
func (api *API) GetGroupMemberList(groupID int) (*models.GetGroupMemberListResponse, error) {
	req := models.GetGroupMemberListRequest{
		GroupID: groupID,
	}
	return client.Call[models.GetGroupMemberListRequest, models.GetGroupMemberListResponse](api.client, "get_group_member_list", req)

}

// GetGroupMemberInfo 获取群成员信息
//
// 参数：
//   - groupID: 群号
//   - userID: 用户 QQ 号
//   - noCache: 是否不使用缓存（true 表示跳过缓存，实时获取）
func (api *API) GetGroupMemberInfo(groupID int, userID int, noCache bool) (*models.GetGroupMemberInfoResponse, error) {
	req := models.GetGroupMemberInfoRequest{
		GroupID: groupID,
		UserID:  userID,
		NoCache: noCache,
	}
	return client.Call[models.GetGroupMemberInfoRequest, models.GetGroupMemberInfoResponse](api.client, "get_group_member_info", req)

}

// GetGroupList 获取群列表
//
// 参数：
//   - noCache: 是否不使用缓存（true 表示跳过缓存，实时获取）
func (api *API) GetGroupList(noCache bool) (*models.GetGroupListResponse, error) {
	req := models.GetGroupListRequest{
		NoCache: noCache,
	}
	return client.Call[models.GetGroupListRequest, models.GetGroupListResponse](api.client, "get_group_list", req)

}

// GetLoginInfo 获取当前登录账号信息
func (api *API) GetLoginInfo() (*models.GetLoginInfoResponse, error) {
	req := models.GetLoginInfoRequest{}
	return client.Call[models.GetLoginInfoRequest, models.GetLoginInfoResponse](api.client, "get_login_info", req)

}

// GetStatus 获取状态信息（包括在线情况、运行时间、插件状态等）
func (api *API) GetStatus() (*models.GetStatusResponse, error) {
	req := models.GetStatusRequest{}
	return client.Call[models.GetStatusRequest, models.GetStatusResponse](api.client, "get_status", req)

}

// GetStrangerInfo 获取陌生人信息
//
// 参数：
//   - userID: 用户 QQ 号
//   - noCache: 是否不使用缓存（true 表示跳过缓存，实时获取）
func (api *API) GetStrangerInfo(userID int, noCache bool) (*models.GetStrangerInfoResponse, error) {
	req := models.GetStrangerInfoRequest{
		UserID:  userID,
		NoCache: noCache,
	}
	return client.Call[models.GetStrangerInfoRequest, models.GetStrangerInfoResponse](api.client, "get_stranger_info", req)

}

// 获取版本信息
func (api *API) GetVersionInfo() (*models.GetVersionInfoResponse, error) {
	req := models.GetVersionInfoRequest{}
	return client.Call[models.GetVersionInfoRequest, models.GetVersionInfoResponse](api.client, "get_version_info", req)

}
