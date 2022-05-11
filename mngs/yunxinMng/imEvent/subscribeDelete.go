package imEvent

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type SubscribeDeleteParam struct {
	Accid           string `json:"accid" validate:"required"`           // 事件订阅人账号
	EventType       int    `json:"eventType" validate:"required"`       // 事件类型，固定设置为1，即 eventType=1
	PublisherAccids string `json:"publisherAccids" validate:"required"` // 取消被订阅人的账号列表，最多100个账号，JSONArray格式。示例：["pub_user1","pub_user2"]
}
type SubscribeDeleteResp struct {
	*imClient.CommonResp
	FailedAccid []string `json:"failed_accid"` // 订阅失败的账号数组
}

// SubscribeDelete 取消在线状态事件订阅
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc5NDQwODk?platformId=60353#取消在线状态事件订阅
// 取消订阅指定人员的在线状态事件
func (api *Api) SubscribeDelete(param *SubscribeDeleteParam) (*SubscribeDeleteResp, error) {
	res, err := api.Client.Post(SubDomain+"subscribe/delete.action", param, &SubscribeDeleteResp{})
	return res.(*SubscribeDeleteResp), err
}
