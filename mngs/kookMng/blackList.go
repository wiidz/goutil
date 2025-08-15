package kookMng

import (
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/networkStruct"
)

// AddBlockList 给用户添加到黑名单
func (m *KookMng) AddBlockList(guildID, kookID, reasonStr string, delMsgDays int) (resp *ApiResponse, err error) {
	data, err := typeHelper.JsonEncodeDecodeMap(&AddUserBlackListParam{
		GuildID:    guildID,
		TargetID:   kookID,
		Remark:     reasonStr,
		DelMsgDays: delMsgDays,
	})
	if err != nil {
		return
	}
	return m.sendRequest(suffixs.AddUserBlackList, data, networkStruct.Post)
}

// DelBlockList 移除用户添加到黑名单
func (m *KookMng) DelBlockList(guildID, kookID string) (resp *ApiResponse, err error) {
	data, err := typeHelper.JsonEncodeDecodeMap(&DelUserBlackListParam{
		GuildID:  guildID,
		TargetID: kookID,
	})
	if err != nil {
		return
	}
	return m.sendRequest(suffixs.DelUserBlackList, data, networkStruct.Post)
}
