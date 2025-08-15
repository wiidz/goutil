package kookMng

import (
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/networkStruct"
)

// SendChanelMsg 向频道发送消息
func (m *KookMng) SendChanelMsg(chanelID, content, quote, nonce string) (resp *ApiResponse, err error) {
	data, err := typeHelper.JsonEncodeDecodeMap(&CreateChanelMsgParam{
		Type:         9,
		TargetID:     chanelID,
		Content:      content,
		Quote:        quote,
		Nonce:        nonce,
		TempTargetID: "",
		TemplateID:   "",
	})
	if err != nil {
		return
	}
	return m.sendRequest(suffixs.ChanelMessageCreate, data, networkStruct.Post)
}

// DelChanelMsg 删除频道中的消息
// 普通用户只能删除自己的消息，有权限的用户可以删除权限范围内他人的消息
func (m *KookMng) DelChanelMsg(msgID string) (resp *ApiResponse, err error) {
	data, err := typeHelper.JsonEncodeDecodeMap(&DeleteChanelMsgParam{
		MsgID: msgID,
	})
	if err != nil {
		return
	}
	return m.sendRequest(suffixs.ChanelMessageDelete, data, networkStruct.Post)
}
