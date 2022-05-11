package imUser

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type MuteParam struct {
	Accid string `json:"accid" validate:"required"` // 用户帐号
	Mute  bool   `json:"mute" validate:"required"`  // 是否全局禁言： true：全局禁言，false:取消全局禁言
}

// Mute 账号全局禁言
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/TU5Mjc5MTg?platformId=60353
// 设置或取消账号的全局禁言状态；
// 账号被设置为全局禁言后，不能发送“点对点”、“群”、“聊天室”消息
func (api *Api) Mute(param *CreateParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post("/user/mute.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
