package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type JoinTeamsParam struct {
	TeamID string `json:"tid"` // 网易云信服务器产生，群唯一标识，创建群时会返回
}

type JoinTeamsResp struct {
	*imClient.CommonResp
	Count int            `json:"count"`
	Infos []*TeamPreview `json:"infos"`
}

type TeamPreview struct {
	Owner    string `json:"owner"`
	TeamName string `json:"tname"`
	MaxUsers int    `json:"maxusers"`
	Tid      int    `json:"tid"`
	Size     int    `json:"size"`
	Custom   string `json:"custom"`
}

// JoinTeams 获取某用户所加入的群信息
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#获取某用户所加入的群信息
// 获取某个用户所加入高级群的群信息
func (api *Api) JoinTeams(param *JoinTeamsParam) (*JoinTeamsResp, error) {
	res, err := api.Client.Post(SubDomain+"joinTeams.action", param, &JoinTeamsResp{})
	return res.(*JoinTeamsResp), err
}
