package imFriend

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

// 查看指定用户的黑名单和静音列表
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DQ0MTY1NzI?platformId=60353#查看指定用户的黑名单和静音列表
// 查看用户的黑名单和静音列表

type GetSpecialRelationParam struct {
	Accid string `json:"accid" validate:"required"` // 用户帐号
}

type GetSpecialRelationResp struct {
	*imClient.CommonResp
	MuteList  []string `json:"mutelist"`  // 被静音的帐号列表
	BlackList []string `json:"BlackList"` // //加黑的帐号列表
}

func (api *Api) GetSpecialRelations(param *DeleteParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"listBlackAndMuteList.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
