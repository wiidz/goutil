package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type QueryDetailParam struct {
	TeamID string `json:"tid"` // 网易云信服务器产生，群唯一标识，创建群时会返回
}

type QueryDetailResp struct {
	*imClient.CommonResp
	TeamInfo struct {
		Icon               interface{} `json:"icon"`
		Announcement       interface{} `json:"announcement"`
		UpdateTeamInfoMode int         `json:"uptinfomode"`
		MaxUsers           int         `json:"maxusers"`
		Intro              interface{} `json:"intro"`
		UpdateCustomMode   int         `json:"upcustommode"`
		TeamName           string      `json:"tname"`
		BeInviteMode       int         `json:"beinvitemode"`
		JoinMode           int         `json:"joinmode"`
		TeamID             int         `json:"tid"`
		InviteMode         int         `json:"invitemode"`
		Mute               bool        `json:"mute"`
		Custom             string      `json:"custom"`
		ClientCustom       string      `json:"clientCustom"`
		CreateTime         int64       `json:"createtime"`
		UpdateTime         int64       `json:"updatetime"`
		Owner              struct {
			CreateTime int64       `json:"createtime"`
			UpdateTime int64       `json:"updatetime"`
			Nick       string      `json:"nick"`
			Accid      string      `json:"accid"`
			Mute       bool        `json:"mute"`
			Custom     interface{} `json:"custom"`
		} `json:"owner"`
		Admins []struct {
			CreateTime int64  `json:"createtime"`
			UpdateTime int64  `json:"updatetime"`
			Nick       string `json:"nick"`
			Accid      string `json:"accid"`
			Mute       bool   `json:"mute"`
			Custom     string `json:"custom"`
		} `json:"admins"`
		Members []struct {
			CreateTime int64       `json:"createtime"`
			UpdateTime int64       `json:"updatetime"`
			Nick       string      `json:"nick"`
			Accid      string      `json:"accid"`
			Mute       bool        `json:"mute"`
			Custom     interface{} `json:"custom"`
		} `json:"members"`
	} `json:"tinfo"`
}

// QueryDetail 获取群组详细信息
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#获取群组详细信息
// 查询指定群的详细信息（群信息+成员详细信息）
func (api *Api) QueryDetail(param *QueryDetailParam) (*QueryDetailResp, error) {
	res, err := api.Client.Post(SubDomain+"queryDetail.action", param, &QueryDetailResp{})
	return res.(*QueryDetailResp), err
}
