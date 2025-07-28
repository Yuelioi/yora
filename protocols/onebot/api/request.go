package api

import (
	"yora/protocols/onebot/client"
	"yora/protocols/onebot/client/models"
)

// 处理加好友请求
func (api *API) SetFriendAdd(flag string, approve bool, reason string) error {
	req := models.SetFriendAddRequest{
		Flag:    flag,
		Approve: approve,
		Reason:  reason,
	}
	_, err := client.Call[models.SetFriendAddRequest, interface{}](api.client, "set_friend_add", req)
	return err
}

// 处理加群请求/邀请
func (api *API) SetGroupAdd(flag string, approve bool, reason string) error {
	req := models.SetGroupAddRequest{
		Flag:    flag,
		Approve: approve,
		Reason:  reason,
	}
	_, err := client.Call[models.SetGroupAddRequest, interface{}](api.client, "set_group_add", req)
	return err
}
