package imMsg

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type SendAttachMsgParam struct {
	From             string       `json:"from" validate:"required" `   // 发送者accid，用户帐号，最大32字符，APP内唯一
	Ope              Ope          `json:"msgtype" validate:"required"` // 0：点对点自定义通知，1：群消息自定义通知，其他返回414，这里特意把属性改成了Ope不是MsgType，防止混淆
	To               string       `json:"to" validate:"required"`      // msgtype==0是表示accid即用户id，msgtype==1表示tid即群id
	Attach           string       `json:"attach" validate:"required"`  // 自定义通知内容，第三方组装的字符串，建议是JSON串，最大长度4096字符
	PushContent      string       `json:"pushcontent,omitempty"`       // 推送文案，最长500个字符 https://faq.yunxin.163.com/kb/main/#/item/KB0291
	Payload          string       `json:"payload,omitempty"`           // 必须是JSON,不能超过2k字符。该参数与APNs推送的payload含义不同 https://faq.yunxin.163.com/kb/main/#/item/KB0291
	Sound            string       `json:"sound,omitempty"`             // 如果有指定推送，此属性指定为客户端本地的声音文件名，长度不要超过30个字符，如果不指定，会使用默认声音
	Save             string       `json:"save,omitempty"`              // 1表示只发在线，2表示会存离线，其他会报414错误。默认会存离线
	Option           NoticeOption `json:"option,omitempty"`            // 发消息时特殊指定的行为选项,Json格式，可用于指定消息计数等特殊行为;option中字段不填时表示默认值。 option示例： {"badge":false,"needPushNick":false,"route":false}
	IsForcePush      string       `json:"isForcePush,omitempty"`       // 发自定义通知时，是否强制推送
	ForcePushContent string       `json:"forcePushContent,omitempty"`  // 发自定义通知时，强制推送文案，最长500个字符
	ForcePushAll     string       `json:"forcePushAll,omitempty"`      // 发群自定义通知时，强推列表是否为群里除发送者外的所有有效成员
	ForcePushList    string       `json:"forcePushList,omitempty"`     // 发群自定义通知时，强推列表，格式为JSONArray，如"accid1","accid2"
	Env              string       `json:"env,omitempty"`               // 所属环境，根据env可以配置不同的抄送地址
}

// SendAttachMsg 发送自定义系统通知
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DEwMTE3NzQ?platformId=60353#发送自定义系统通知
// 1.自定义系统通知区别于普通消息，方便开发者进行业务逻辑的通知；
// 2.目前支持两种类型：点对点类型和群类型（仅限高级群），根据msgType有所区别。
// 应用场景：如某个用户给另一个用户发送好友请求信息等，具体attach为请求消息体，第三方可以自行扩展，建议是json格式
func (api *Api) SendAttachMsg(param *SendAttachMsgParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"sendAttachMsg.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
