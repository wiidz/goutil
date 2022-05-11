package imEvent

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type SubscribeAddParam struct {
	Accid           string `json:"accid" validate:"required"`           // 事件订阅人账号
	EventType       int    `json:"eventType" validate:"required"`       // 事件类型，固定设置为1，即 eventType=1
	PublisherAccids string `json:"publisherAccids" validate:"required"` // 被订阅人的账号列表，最多100个账号，JSONArray格式。示例：["pub_user1","pub_user2"]
	Ttl             int64  `json:"ttl" validate:"required"`             // 有效期，单位：秒。取值范围：60～2592000（即60秒到30天）
}
type SubscribeAddResp struct {
	*imClient.CommonResp
	FailedAccid []string `json:"failed_accid"` // 订阅失败的账号数组
}

// SubscribeAdd 订阅在线状态事件
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc5NDQwODk?platformId=60353#订阅在线状态事件
// 订阅指定人员的在线状态事件，每个账号最大有效订阅账号不超过3000个
func (api *Api) SubscribeAdd(param *SubscribeAddParam) (*SubscribeAddResp, error) {
	res, err := api.Client.Post(SubDomain+"subscribe/add.action", param, &SubscribeAddResp{})
	return res.(*SubscribeAddResp), err
}
