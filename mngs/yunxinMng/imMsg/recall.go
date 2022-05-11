package imMsg

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type RecallParam struct {
	DeleteMsgID string `json:"deleteMsgid"`              // 要撤回消息的msgid
	TimeTag     int64  `json:"timetag"`                  // 要撤回消息的创建时间
	Type        int    `json:"type"`                     // 7:表示点对点消息撤回，8:表示群消息撤回，其它为参数错误
	From        string `json:"from" validate:"required"` // 发消息的accid
	To          string `json:"to" validate:"required"`   // 如果点对点消息，为接收消息的accid,如果群消息，为对应群的tid
	Msg         string `json:"msg"`                      // 可以带上对应的描述
	IgnoreTime  string `json:"ignoreTime"`               // 1表示绕过撤回时间检测，其它为非法参数，最多撤回近30天内的消息。如果需要撤回时间检测，不填即可。
	PushContent string `json:"pushcontent,omitempty"`    // 推送文案，android以此为推送显示文案；ios若未填写payload，显示文案以pushcontent为准。超过500字符后，会对文本进行截断。
	Payload     string `json:"payload,omitempty"`        // 推送对应的payload,必须是JSON,不超过2K字符
	Env         string `json:"env,omitempty"`            // 所属环境，根据env可以配置不同的抄送地
	Attach      string `json:"attach"`                   // 扩展字段，最大5000字符
}

// Recall 消息撤回
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DEwMTE3NzQ?platformId=60353#消息撤回
// 消息撤回接口，可以撤回一定时间内的点对点与群消息
func (api *Api) Recall(param *RecallParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"recall.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
