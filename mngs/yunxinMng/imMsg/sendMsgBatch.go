package imMsg

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type SendMsgBatchParam struct {
	FromAccid         string  `json:"fromAccid" validate:"required"` // 发送者accid，用户帐号，最大32字符
	ToAccids          string  `json:"toAccids" validate:"required"`  // ["aaa","bbb"]（JSONArray对应的accid，如果解析出错，会报414错误），限500人
	Type              MsgType `json:"type" validate:"required"`      // 0 表示文本消息,  1 表示图片，  2 表示语音， 3 表示视频，  4 表示地理位置信息， 6 表示文件， 10 表示提示消息， 100 自定义消息类型（特别注意，对于未对接易盾反垃圾功能的应用，该类型的消息不会提交反垃圾系统检测）
	Body              string  `json:"body" validate:"required"`      // 最大长度5000字符，JSON格式。
	MsgDesc           string  `json:"msgdesc,omitempty"`             // 消息描述文本，针对非Text、Tip类型的消息有效，最大长度500字符。该描述信息可用于云端历史消息关键词检索。
	Option            Option  `json:"option"`                        // 发消息时特殊指定的行为选项,JSON格式，可用于指定消息的漫游，存云端历史，发送方多端同步，推送，消息抄送等特殊行为;option中字段不填时表示默认值 ，option示例: {"push":false,"roam":true,"history":false,"sendersync":true,"route":false,"badge":false,"needPushNick":true}
	PushContent       string  `json:"pushcontent,omitempty"`         // 推送文案,最长500个字符。具体参见 推送配置参数详解。
	Payload           string  `json:"payload,omitempty"`             // 必须是JSON,不能超过2k字符。该参数与APNs推送的payload含义不同。具体参见 推送配置参数详解。
	Ext               string  `json:"ext,omitempty"`                 //  开发者扩展字段，长度限制1024字符
	Bid               string  `json:"bid,omitempty"`                 // 可选，反垃圾业务ID，实现“单条消息配置对应反垃圾”，若不填则使用原来的反垃圾配置
	UseYiDun          int     `json:"useYidun,omitempty"`            // 可选，单条消息是否使用易盾反垃圾，可选值为0。0：（在开通易盾的情况下）不使用易盾反垃圾而是使用通用反垃圾，包括自定义消息。 若不填此字段，即在默认情况下，若应用开通了易盾反垃圾功能，则使用易盾反垃圾来进行垃圾消息的判断
	YiDunAntiCheating string  `json:"yidunAntiCheating,omitempty"`   // 可选，易盾反垃圾增强反作弊专属字段，限制json，长度限制1024字符（详见易盾反垃圾接口文档反垃圾防刷版专属字段）
	YiDunAntiSpamExt  string  `json:"yidunAntiSpamExt,omitempty"`    // 可选，透传给易盾的反垃圾增强版的检测参数，格式为json，长度限制 1024 字符。（具体请参见易盾的反垃圾增强版用户可扩展字段https://support.dun.163.com/documents/588434200783982592?docId=476559002902757376#/%E7%94%A8%E6%88%B7%E6%89%A9%E5%B1%95%E5%8F%82%E6%95%B0）。
	MarkRead          int     `json:"markRead,omitempty"`            // 可选，群消息是否需要已读业务（仅对群消息有效），0:不需要，1:需要
	CheckFriend       bool    `json:"checkFriend,omitempty"`         // 是否为好友关系才发送消息，默认否 使用该参数需要先开通功能服务
	ReturnMsgID       bool    `json:"returnMsgid,omitempty"`         // 是否需要返回消息ID false：不返回消息ID（默认值） true：返回消息ID（toAccids包含的账号数量不可以超过100个）
	Env               string  `json:"env,omitempty"`                 // 所属环境，根据env可以配置不同的抄送地址
}

type SendMsgBatchResp struct {
	*imClient.CommonResp
	MsgIDs     map[string]int64 `json:"msgids"`     // //消息接受者对应的消息ID，returnMsgId参数为true时才返回
	TimeTag    int64            `json:"timetag"`    // 消息发送的时间戳
	Unregister []string         `json:"unregister"` // 未注册的帐号
}

// SendBatch 批量发送点对点普通消息
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DEwMTE3NzQ?platformId=60353#批量发送点对点普通消息
// 1.给用户发送点对点普通消息，包括文本，图片，语音，视频，地理位置和自定义消息。
// 2.最大限500人，只能针对个人,如果批量提供的帐号中有未注册的帐号，会提示并返回给用户。
// 3.此接口受频率控制，一个应用一分钟最多调用120次，超过会返回416状态码，并且被屏蔽一段时间；
func (api *Api) SendBatch(param *SendMsgBatchParam) (*SendMsgBatchResp, error) {
	res, err := api.Client.Post(SubDomain+"sendBatchMsg.action", param, &SendMsgBatchResp{})
	return res.(*SendMsgBatchResp), err
}
