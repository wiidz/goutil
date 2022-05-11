package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type GetMuteMembersParam struct {
	Tid   string `json:"tid" validate:"required"`   // 网易云信服务器产生，群唯一标识，创建群时会返回
	Owner string `json:"owner" validate:"required"` // 群主accid
}
type GetMuteMembersResp struct {
	*imClient.CommonResp
	Mutes []struct {
		Nick  string `json:"nick"`
		Accid string `json:"accid"`
		Tid   int    `json:"tid"`
		Type  int    `json:"type"`
	}
}

// GetMuteMembers 获取群组禁言列表
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#获取群组禁言列表
// 获取群组禁言的成员列表
func (api *Api) GetMuteMembers(param *GetMuteMembersParam) (*GetMuteMembersResp, error) {
	res, err := api.Client.Post(SubDomain+"listTeamMute.action", param, &GetMuteMembersResp{})
	return res.(*GetMuteMembersResp), err
}
