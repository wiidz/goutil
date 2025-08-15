package kookMng

import (
	"errors"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/networkStruct"
)

// SendDirectMsg 发送私聊消息
func (m *KookMng) SendDirectMsg(targetID, content, quote, nonce string) (resp *ApiResponse, err error) {
	// 判断一下不能给自己发
	if targetID == m.Config.RobotID {
		err = errors.New("不能给自己发送消息")
		return
	}

	data, err := typeHelper.JsonEncodeDecodeMap(&CreateDirectMessageParam{
		Type:     9,
		TargetID: targetID,
		//ChatCode:   "",
		Content: content,
		Quote:   quote,
		Nonce:   nonce,
		//TemplateID: "",
	})
	if err != nil {
		return
	}
	return m.sendRequest(suffixs.DirectMsgCreate, data, networkStruct.Post)
}

// CreateUserChat 创建私聊（发送私聊之前要创建）
func (m *KookMng) CreateUserChat(targetID string) (resp *ApiResponse, err error) {
	// 判断一下不能给自己发
	if targetID == m.Config.RobotID {
		err = errors.New("不能给自己创建聊天")
		return
	}

	data, err := typeHelper.JsonEncodeDecodeMap(&CreateUserChatParam{
		TargetID: targetID,
	})
	if err != nil {
		return
	}
	return m.sendRequest(suffixs.UserChatCreate, data, networkStruct.Post)
}
