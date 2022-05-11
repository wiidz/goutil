package imEvent

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type SubscribeQueryParam struct {
	Accid           string `json:"accid" validate:"required"`           // 事件订阅人账号
	EventType       int    `json:"eventType" validate:"required"`       // 事件类型，固定设置为1，即 eventType=1
	PublisherAccids string `json:"publisherAccids" validate:"required"` // 被订阅人的账号列表，最多100个账号，JSONArray格式。示例：["pub_user1","pub_user2"]
}
type SubscribeQueryResp struct {
	*imClient.CommonResp
	Subscribes []struct {
		Accid         string `json:"accid"`         // 被订阅人账号
		EventType     int    `json:"eventType"`     // 事件类型
		ExpireTime    int64  `json:"expireTime"`    // 过期时间
		SubscribeTime int64  `json:"subscribeTime"` // 订阅时间
	}
}

// SubscribeQuery 查询在线状态事件订阅关系
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc5NDQwODk?platformId=60353#查询在线状态事件订阅关系
// 查询指定人员的有效在线状态事件订阅关系
func (api *Api) SubscribeQuery(param *SubscribeQueryParam) (*SubscribeQueryResp, error) {
	res, err := api.Client.Post(SubDomain+"subscribe/query.action", param, &SubscribeQueryResp{})
	return res.(*SubscribeQueryResp), err
}
