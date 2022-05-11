package imFriend

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type DeleteParam struct {
	Accid         string `json:"accid" validate:"required"`  // 发起者accid
	FAccid        string `json:"faccid" validate:"required"` // 要删除朋友的accid
	IsDeleteAlias int    `json:"isDeleteAlias"`              // 是否需要删除备注信息 默认false:不需要，true:需要
}

// Delete 删除好友
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DQ0MTY1NzI?platformId=60353
// 删除好友关系
func (api *Api) Delete(param *DeleteParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"delete.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
