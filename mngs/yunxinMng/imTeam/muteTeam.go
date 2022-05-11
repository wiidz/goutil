package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type MuteTeamParam struct {
	TeamID string `json:"tid" validate:"required"`   // 网易云信服务器产生，群唯一标识，创建群时会返回
	Accid  string `json:"accid" validate:"required"` // 要操作的群成员accid
	Ope    int    `json:"ope" validate:"required"`   // 1：关闭消息提醒，2：打开消息提醒，其他值无效
}

// MuteTeam 修改消息提醒开关
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#修改消息提醒开关
// 高级群修改消息提醒开关
func (api *Api) MuteTeam(param *MuteTeamParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"muteTeam.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
