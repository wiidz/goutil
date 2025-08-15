package kookMng

import (
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/networkStruct"
)

// GetUserGroupList 获取服务器角色列表
func (m *KookMng) GetUserGroupList(guildID string, page, pageSize int) (resp *ApiResponse, err error) {
	data, err := typeHelper.JsonEncodeDecodeMap(&GetUserGroupListParam{
		GuildID:  guildID,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return
	}
	return m.sendRequest(suffixs.GetUserGroupList, data, networkStruct.Get)
}

// AddUserGroup 给用户添加分组
func (m *KookMng) AddUserGroup(guildID, targetID string, roleID uint64) (resp *ApiResponse, err error) {
	data, err := typeHelper.JsonEncodeDecodeMap(&UpdateUserGroupParam{
		GuildID: guildID,
		UserID:  targetID,
		RoleID:  roleID,
	})
	if err != nil {
		return
	}
	return m.sendRequest(suffixs.AddUserGroup, data, networkStruct.Post)
}

// DelUserGroup 移除用户分组
func (m *KookMng) DelUserGroup(guildID, targetID string, roleID uint64) (resp *ApiResponse, err error) {
	data, err := typeHelper.JsonEncodeDecodeMap(&UpdateUserGroupParam{
		GuildID: guildID,
		UserID:  targetID,
		RoleID:  roleID,
	})
	if err != nil {
		return
	}
	return m.sendRequest(suffixs.DelUserGroup, data, networkStruct.Post)
}
