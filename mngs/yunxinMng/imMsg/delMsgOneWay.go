package imMsg

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type DelMsgOneWayParam struct {
	DeleteMsgID string `json:"deleteMsgid" validate:"required"` // 要撤回消息的msgid
	TimeTag     int64  `json:"timetag" validate:"required"`     // 要撤回消息的创建时间
	Type        int    `json:"type" validate:"required"`        // 13:表示点对点消息撤回，14:表示群消息撤回，其它为参数错误
	From        string `json:"from" validate:"required"`        // 发消息的accid
	To          string `json:"to" validate:"required"`          // 如果点对点消息，为接收消息的accid,如果群消息，为对应群的tid
	Msg         string `json:"msg,omitempty"`                   // 可以带上对应的描述
}

// DelMsgOneWay 单向撤回消息
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DEwMTE3NzQ?platformId=60353#单向撤回消息
// 1、可以单向撤回点对点消息和群消息，撤回之后，消息接收者会收到一条单向撤回的通知，并删除对应的离线消息、漫游消息、历史消息
// 2、撤回之后，消息发送者无感知，可以正常使用漫游消息、历史消息
// 3、客户端要求至少v7.2.0版本，否则无法收到撤回通知（但是历史消息依然会单向删除）
func (api *Api) DelMsgOneWay(param *DelMsgOneWayParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"delMsgOneWay.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
