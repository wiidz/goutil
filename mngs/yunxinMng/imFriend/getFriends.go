package imFriend

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type GetFriendsParam struct {
	Accid      string `json:"accid" validate:"required"`       // 发起者accid
	UpdateTime int    `json:"update_time" validate:"required"` // 更新时间戳，接口返回该时间戳之后有更新的好友列表
	CreateTime int    `json:"create_time"`                     // 【Deprecated】定义同updatetime
}

type GetFriendsResp struct {
	*imClient.CommonResp
	Size    int       `json:"size"`    // 数量
	Friends []*Friend `json:"friends"` // 好友
}

type Friend struct {
	CreateTime  int64  `json:"createtime"`
	Bidirection bool   `json:"bidirection"`
	FAccid      string `json:"faccid"`
}

// GetFriends 获取好友关系
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DQ0MTY1NzI?platformId=60353
// 查询某时间点起到现在有更新的双向好友
func (api *Api) GetFriends(param *GetFriendsParam) (*GetFriendsResp, error) {
	res, err := api.Client.Post(SubDomain+"get.action", param, &GetFriendsResp{})
	return res.(*GetFriendsResp), err
}
