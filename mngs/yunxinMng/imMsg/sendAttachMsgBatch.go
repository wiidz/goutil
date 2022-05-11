package imMsg

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type SendAttachMsgBatchParam struct {
	FromAccid        string       `json:"fromAccid" validate:"required"` // 发送者accid，用户帐号，最大32字符，APP内唯一
	ToAccids         string       `json:"toAccids" validate:"required"`  // ["aaa","bbb"]（JSONArray对应的accid，如果解析出错，会报414错误），最大限500人
	Attach           string       `json:"attach" validate:"required"`    // 自定义通知内容，第三方组装的字符串，建议是JSON串，最大长度4096字符
	PushContent      string       `json:"pushcontent,omitempty"`         // 推送文案，最长500个字符 https://faq.yunxin.163.com/kb/main/#/item/KB0291
	Payload          string       `json:"payload,omitempty"`             // 必须是JSON,不能超过2k字符。该参数与APNs推送的payload含义不同 https://faq.yunxin.163.com/kb/main/#/item/KB0291
	Sound            string       `json:"sound,omitempty"`               // 如果有指定推送，此属性指定为客户端本地的声音文件名，长度不要超过30个字符，如果不指定，会使用默认声音
	Save             string       `json:"save,omitempty"`                // 1表示只发在线，2表示会存离线，其他会报414错误。默认会存离线
	Option           NoticeOption `json:"option,omitempty"`              // 发消息时特殊指定的行为选项,Json格式，可用于指定消息计数等特殊行为;option中字段不填时表示默认值。 option示例： {"badge":false,"needPushNick":false,"route":false}
	IsForcePush      string       `json:"isForcePush,omitempty"`         // 发自定义通知时，是否强制推送
	ForcePushContent string       `json:"forcePushContent,omitempty"`    // 发自定义通知时，强制推送文案，最长500个字符
	Env              string       `json:"env,omitempty"`                 // 所属环境，根据env可以配置不同的抄送地址
}

type SendAttachMsgBatchResp struct {
	*imClient.CommonResp
	Unregister []string `json:"unregister"` // 未注册的帐号
}

// SendAttachMsgBatch 批量发送点对点自定义系统通知
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DEwMTE3NzQ?platformId=60353#批量发送点对点自定义系统通知
// 1.系统通知区别于普通消息，应用接收到直接交给上层处理，客户端可不做展示；
// 2.目前支持类型：点对点类型；
// 3.最大限500人，只能针对个人,如果批量提供的帐号中有未注册的帐号，会提示并返回给用户；
// 4.此接口受频率控制，一个应用一分钟最多调用120次，超过会返回416状态码，并且被屏蔽一段时间；
// 应用场景：如某个用户给另一个用户发送好友请求信息等，具体attach为请求消息体，第三方可以自行扩展，建议是json格式
func (api *Api) SendAttachMsgBatch(param *SendAttachMsgBatchParam) (*SendAttachMsgBatchResp, error) {
	res, err := api.Client.Post(SubDomain+"sendBatchAttachMsg.action", param, &SendAttachMsgBatchResp{})
	return res.(*SendAttachMsgBatchResp), err
}
