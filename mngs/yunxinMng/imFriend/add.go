package imFriend

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

// 用户关系托管（加好友）
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DQ0MTY1NzI?platformId=60353
// 两人保持好友关系

type AddParam struct {
	Accid    string `json:"accid" validate:"required"`  // Accid 加好友发起者accid
	FAccid   string `json:"faccid" validate:"required"` // FAccid 加好友接收者accid
	Type     int    `json:"type" validate:"required"`   // 1直接加好友，2请求加好友，3同意加好友，4拒绝加好友
	Msg      string `json:"msg"`                        // 加好友对应的请求消息，第三方组装，最长256字符
	ServerEx string `json:"serverex"`                   // 服务器端扩展字段，限制长度256 此字段client端只读，server端读写
}

func (api *Api) Add(param *AddParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"add.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
