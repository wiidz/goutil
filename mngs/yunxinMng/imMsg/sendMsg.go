package imMsg

import (
	"github.com/wiidz/goutil/mngs/yunxinMng/imClient"
)

type SendMsgParam struct {
	From               string `json:"from" validate:"required"`     // 发送者accid，用户帐号，最大32字符，
	Ope                int    `json:"ope" validate:"required"`      // 0：点对点个人消息，1：群消息（高级群），其他返回414
	To                 string `json:"to" validate:"required"`       // ope==0是表示accid即用户id，ope==1表示tid即群id
	Type               int    `json:"type" validate:"required"`     // 0 表示文本消息,  1 表示图片，  2 表示语音， 3 表示视频，  4 表示地理位置信息， 6 表示文件， 10 表示提示消息， 100 自定义消息类型（特别注意，对于未对接易盾反垃圾功能的应用，该类型的消息不会提交反垃圾系统检测）
	Body               string `json:"body" validate:"required"`     // 最大长度5000字符，JSON格式。
	MsgDesc            string `json:"msgdesc,omitempty"`            // 消息描述文本，针对非Text、Tip类型的消息有效，最大长度500字符。该描述信息可用于云端历史消息关键词检索。
	Antispam           string `json:"antispam,omitempty"`           // 对于对接了易盾反垃圾功能的应用，本消息是否需要指定经由易盾检测的内容（antispamCustom）。true或false, 默认false。 只对消息类型为：100 自定义消息类型 的消息生效。
	AntispamCustom     string `json:"antispamCustom,omitempty"`     // 在antispam参数为true时生效。 自定义的反垃圾检测内容, JSON格式，长度限制同body字段，不能超过5000字符，要求antispamCustom格式如下： {"type":1,"data":"custom content"} 字段说明： 1. type: 1：文本，2：图片。 2. data: 文本内容or图片地址。
	Option             int    `json:"option,omitempty"`             // 发消息时特殊指定的行为选项,JSON格式，可用于指定消息的漫游，存云端历史，发送方多端同步，推送，消息抄送等特殊行为;option中字段不填时表示默认值 ，option示例: {"push":false,"roam":true,"history":false,"sendersync":true,"route":false,"badge":false,"needPushNick":true}
	PushContent        string `json:"pushcontent,omitempty"`        // 推送文案,最长500个字符。具体参见 推送配置参数详解。
	Payload            string `json:"payload,omitempty"`            // 必须是JSON,不能超过2k字符。该参数与APNs推送的payload含义不同。具体参见 推送配置参数详解。
	Ext                string `json:"ext,omitempty"`                //  开发者扩展字段，长度限制1024字符
	ForcePushList      string `json:"forcepushlist,omitempty"`      // 发送群消息时的强推用户列表（云信demo中用于承载被@的成员），格式为JSONArray，如["accid1","accid2"]。若forcepushall为true，则forcepushlist为除发送者外的所有有效群成员
	ForcePushContent   string `json:"forcepushcontent,omitempty"`   // 发送群消息时，针对强推列表forcepushlist中的用户，强制推送的内容
	ForcePushAll       string `json:"forcepushall,omitempty"`       // 发送群消息时，强推列表是否为群里除发送者外的所有有效成员，true或false，默认为false
	Bid                string `json:"bid,omitempty"`                // 可选，反垃圾业务ID，实现“单条消息配置对应反垃圾”，若不填则使用原来的反垃圾配置
	UseYiDun           int    `json:"useYidun,omitempty"`           // 可选，单条消息是否使用易盾反垃圾，可选值为0。0：（在开通易盾的情况下）不使用易盾反垃圾而是使用通用反垃圾，包括自定义消息。 若不填此字段，即在默认情况下，若应用开通了易盾反垃圾功能，则使用易盾反垃圾来进行垃圾消息的判断
	YiDunAntiCheating  string `json:"yidunAntiCheating,omitempty"`  // 可选，易盾反垃圾增强反作弊专属字段，限制json，长度限制1024字符（详见易盾反垃圾接口文档反垃圾防刷版专属字段）
	MarkRead           int    `json:"markRead,omitempty"`           // 可选，群消息是否需要已读业务（仅对群消息有效），0:不需要，1:需要
	CheckFriend        bool   `json:"checkFriend,omitempty"`        // 是否为好友关系才发送消息，默认否 使用该参数需要先开通功能服务
	SubType            int    `json:"subType,omitempty"`            // 自定义消息子类型，大于0
	MsgSenderNoSense   int    `json:"msgSenderNoSense,omitempty"`   // 发送方是否无感知。0-有感知，1-无感知。若无感知，则消息发送者无该消息的多端、漫游、历史记录等。
	MsgReceiverNoSense int    `json:"msgReceiverNoSense,omitempty"` // 接受方是否无感知。0-有感知，1-无感知。若无感知，则消息接收者者无该消息的多端、漫游、历史记录等
	Env                string `json:"env,omitempty"`                // 所属环境，根据env可以配置不同的抄送地址
}

type SendMsgResp struct {
	*imClient.CommonResp
	Data struct {
		MsgID    uint64 `json:"msgid"`
		TimeTag  uint64 `json:"timetag"`
		Antispam bool   `json:"antispam"`
	} `json:"info"`
}

// SendMsg 发送普通消息
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DEwMTE3NzQ?platformId=60353#发送普通消息
// 一秒内默认最多调用发送消息接口100次。如需上调上限，请在官网首页通过QQ、在线消息或电话等方式咨询商务人员。
// 给用户或者高级群发送普通消息，包括文本，图片，语音，视频和地理位置，具体消息参考下面描述。
func (api *Api) SendMsg(param *SendMsgParam) (*SendMsgResp, error) {
	res, err := api.Client.Post(SubDomain+"sendMsg.action", param, &SendMsgResp{})
	return res.(*SendMsgResp), err
}
