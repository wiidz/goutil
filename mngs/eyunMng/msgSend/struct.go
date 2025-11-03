package msgSend

import "github.com/wiidz/goutil/mngs/eyunMng/base"

type ReturnCode string

const Success ReturnCode = "1000"
const Failed ReturnCode = "1001"

type Api struct {
	Config *base.Config
}

type SendMsgResp struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Data    *SendMsgData `json:"data"`
}

type SendMsgData struct {
	Type       int    `json:"type"`       // 类型
	MsgID      int64  `json:"msgId"`      // 消息msgId
	NewMsgID   int64  `json:"newMsgId"`   // 消息newMsgId
	CreateTime int    `json:"createTime"` // 消息发送时间戳
	WcID       string `json:"wcId"`       // 消息接收方id
}

type SendMiniParam struct {
	DisplayName string `json:"displayName"` // 小程序的名称，例如：京东
	IconUrl     string `json:"iconUrl"`     // 小程序卡片图标的url(50KB以内的png/jpg)
	AppId       string `json:"appId"`       // 小程序的appID,例如：wx7c544xxxxxx
	PagePath    string `json:"pagePath"`    // 点击小程序卡片跳转的url
	ThumbUrl    string `json:"thumbUrl"`    // 小程序卡片缩略图的url (50KB以内的png/jpg)
	Title       string `json:"title"`       // 标题
	UserName    string `json:"userName"`    // 小程序所有人的ID,例如：gh_1c0daexxxx@app
}

// TextMsgSendParam 发送文本信息
type TextMsgSendParam struct {
	WcID    string   `json:"wc_id"`   // 发送给谁
	Content string   `json:"content"` // 发送的内容
	AtIDs   []string `json:"at_ids"`  // 被艾特的wcIDs
	AtAll   bool     `json:"at_all"`  // 艾特所有人
}
