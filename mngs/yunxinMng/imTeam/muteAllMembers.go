package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type MuteAllMembersParam struct {
	Tid      string `json:"tid" validate:"required"`   // 网易云信服务器产生，群唯一标识，创建群时会返回
	Owner    string `json:"owner" validate:"required"` // 群主accid
	Mute     string `json:"mute"`                      // true:禁言，false:解禁(mute和muteType至少提供一个，都提供时按mute处理)
	MuteType int    `json:"muteType"`                  // 禁言类型 0:解除禁言，1:禁言普通成员 3:禁言整个群(包括群主)
	Attach   string `json:"attach"`                    // 自定义扩展字段，最大长度512
}

// MuteAllMembers 将群组整体禁言
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#将群组整体禁言
// 禁言群组，普通成员不能发送消息，创建者和管理员可以发送消息
func (api *Api) MuteAllMembers(param *MuteAllMembersParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"muteTlistAll.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
