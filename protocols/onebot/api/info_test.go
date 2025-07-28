package api

import (
	"testing"
)

// 信息获取测试辅助类
type InfoTestHelper struct {
	*TestHelper
}

func NewInfoTestHelper(t *testing.T) *InfoTestHelper {
	h := NewTestHelper(t)
	return &InfoTestHelper{
		TestHelper: h,
	}
}

// 1. 获取好友列表测试
func TestGetFriendList(t *testing.T) {
	h := NewInfoTestHelper(t)

	resp, err := h.api.GetFriendList()
	h.StatusOk(resp, err, "获取好友列表")

	t.Logf("获取好友列表成功，好友数量: %d", len(resp.Data))
	if len(resp.Data) > 0 {
		firstFriend := resp.Data[0]
		t.Logf("第一个好友信息: ID=%d, 昵称=%s, 备注=%s",
			firstFriend.UserID, firstFriend.Nickname, firstFriend.Remark)
	}

}

// 2. 获取群信息测试
func TestGetGroupInfo(t *testing.T) {
	h := NewInfoTestHelper(t)

	groupID := GID

	// 测试使用缓存
	resp, err := h.api.GetGroupInfo(groupID, false)
	h.StatusOk(resp, err, "获取群信息（使用缓存）")

	h.t.Logf("获取群信息成功, 群名 = %s", resp.Data.GroupName)

}

// 3. 获取群成员列表测试
func TestGetGroupMemberList(t *testing.T) {
	h := NewInfoTestHelper(t)

	groupID := GID

	resp, err := h.api.GetGroupMemberList(groupID)
	h.StatusOk(resp, err, "获取群成员列表")

	if len(resp.Data) > 0 {
		firstMember := resp.Data[0]
		t.Logf("第一个成员信息: ID=%d, 昵称=%s, 群名片=%s, 角色=%s",
			firstMember.UserID, firstMember.Nickname, firstMember.Card, firstMember.Role)
	}

}

// 4. 获取群成员信息测试
func TestGetGroupMemberInfo(t *testing.T) {
	h := NewInfoTestHelper(t)

	groupID := GID

	// 测试使用缓存
	resp, err := h.api.GetGroupMemberInfo(groupID, TID, false)
	h.StatusOk(resp, err, "获取群成员信息（使用缓存）")

	t.Logf("获取群成员信息成功（使用缓存），群ID: %d, 用户ID: %d, 昵称: %s, 角色: %s",
		groupID, TID, resp.Data.Nickname, resp.Data.Role)

	// 测试不使用缓存
	resp2, err2 := h.api.GetGroupMemberInfo(groupID, TID, true)
	h.StatusOk(resp2, err2, "获取群成员信息（不使用缓存）")

	t.Logf("获取群成员信息成功（不使用缓存），群ID: %d, 用户ID: %d, 昵称: %s, 角色: %s",
		groupID, TID, resp2.Data.Nickname, resp2.Data.Role)

}

// 5. 获取群列表测试
func TestGetGroupList(t *testing.T) {
	h := NewInfoTestHelper(t)

	// 测试使用缓存
	resp, err := h.api.GetGroupList(false)
	h.StatusOk(resp, err, "获取群列表（使用缓存）")

	t.Logf("获取群列表成功（使用缓存），群数量: %d", len(resp.Data))
	if len(resp.Data) > 0 {
		firstGroup := resp.Data[0]
		t.Logf("第一个群信息: ID=%d, 群名=%s, 成员数=%d",
			firstGroup.GroupID, firstGroup.GroupName, firstGroup.MemberCount)
	}

	// 测试不使用缓存
	resp2, err2 := h.api.GetGroupList(true)
	h.StatusOk(resp2, err2, "获取群列表（不使用缓存）")

	t.Logf("获取群列表成功（不使用缓存），群数量: %d", len(resp2.Data))
}

// 6. 获取当前登录账号信息测试
func TestGetLoginInfo(t *testing.T) {
	h := NewInfoTestHelper(t)

	resp, err := h.api.GetLoginInfo()
	h.StatusOk(resp, err, "获取当前登录账号信息")

	t.Logf("获取当前登录账号信息成功，用户ID: %d, 昵称: %s",
		resp.Data.UserID, resp.Data.Nickname)
}

// 7. 获取状态信息测试
func TestGetStatus(t *testing.T) {
	h := NewInfoTestHelper(t)

	resp, err := h.api.GetStatus()
	h.StatusOk(resp, err, "获取状态信息")

	t.Logf("获取状态信息成功，在线状态: %v, 运行状态: %v",
		resp.Data.Online, resp.Data.Good)

}

// 8. 获取陌生人信息测试
func TestGetStrangerInfo(t *testing.T) {
	h := NewInfoTestHelper(t)

	userID := TID

	// 测试使用缓存
	resp, err := h.api.GetStrangerInfo(userID, false)
	h.StatusOk(resp, err, "获取陌生人信息（使用缓存）")

	t.Logf("获取陌生人信息成功（使用缓存），用户ID: %d, 昵称: %s, 性别: %s, 年龄: %d",
		userID, resp.Data.Nickname, resp.Data.Sex, resp.Data.Age)

	// 测试不使用缓存
	resp2, err2 := h.api.GetStrangerInfo(userID, true)
	h.StatusOk(resp2, err2, "获取陌生人信息（不使用缓存）")

	t.Logf("获取陌生人信息成功（不使用缓存），用户ID: %d, 昵称: %s, 性别: %s, 年龄: %d",
		userID, resp2.Data.Nickname, resp2.Data.Sex, resp2.Data.Age)
}

// 9. 获取版本信息测试
func TestGetVersionInfo(t *testing.T) {
	h := NewInfoTestHelper(t)

	resp, err := h.api.GetVersionInfo()
	h.StatusOk(resp, err, "获取版本信息")

	t.Logf("获取版本信息成功，应用名: %s, 版本: %s, 协议版本: %s",
		resp.Data.AppName, resp.Data.AppVersion, resp.Data.ProtocolVersion)
}
