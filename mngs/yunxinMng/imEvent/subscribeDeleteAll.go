package imEvent

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type SubscribeDeleteAllParam struct {
	Accid           string `json:"accid" validate:"required"`           // 事件订阅人账号
	EventType       int    `json:"eventType" validate:"required"`       // 事件类型，固定设置为1，即 eventType=1
	PublisherAccids string `json:"publisherAccids" validate:"required"` // 取消被订阅人的账号列表，最多100个账号，JSONArray格式。示例：["pub_user1","pub_user2"]
}

// SubscribeDeleteAll 取消全部在线状态事件订阅
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc5NDQwODk?platformId=60353#取消全部在线状态事件订阅
// 取消指定事件的全部订阅关系
func (api *Api) SubscribeDeleteAll(param *SubscribeDeleteAllParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"subscribe/batchdel.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
