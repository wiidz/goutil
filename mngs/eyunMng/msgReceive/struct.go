package msgReceive

import "github.com/wiidz/goutil/mngs/eyunMng/base"

type Api struct {
	Config *base.Config
}

// FromMsgResp 基础
type FromMsgResp struct {
	Account     string      `json:"account"`     // 账号
	Data        interface{} `json:"data"`        // 消息体
	MessageType string      `json:"messageType"` // 消息类型
	WcId        string      `json:"wcId"`        // 微信id
}

// TextParam 群聊文本信息
type TextParam struct {
	Account     string    `json:"account"`     // 账号
	Data        *TextData `json:"data"`        // 消息体
	MessageType string    `json:"messageType"` // 消息类型
	WcId        string    `json:"wcId"`        // 微信id
}

type TextData struct {
	WId string `json:"wId"` //

	Content   string `json:"content"`   // 消息内容
	FromUser  string `json:"fromUser"`  // 发送微信id（二选一）
	FromGroup string `json:"fromGroup"` // 发送群号（二选一）
	Img       string `json:"img"`
	MsgId     int    `json:"msgId"`    // 消息msgId
	NewMsgId  int64  `json:"newMsgId"` // 消息newMsgId
	Self      bool   `json:"self"`     //是否是自己发送的消息
	Timestamp int    `json:"timestamp"`
	ToUser    string `json:"toUser"` // 接收微信id
}

type MessageType int

const GroupTextMsg MessageType = 80001  // 群组文本消息
const GroupUpdate MessageType = 85001   // 群聊信息变更通知
const GroupExit MessageType = 85015     // 退出群聊
const GroupInviteIn MessageType = 85008 // xxx邀请xxx加入群聊
const GroupQrcodeIn MessageType = 85009 // xxx通过扫描xxx的二维码加入群聊

// GroupUpdateParam 群组信息修改
type GroupUpdateParam struct {
	Account string `json:"account"` // 账号
	Data    struct {
		AddContactScene int    `json:"addContactScene"`
		BitMask         int64  `json:"bitMask"`
		ChatRoomNotify  int    `json:"chatRoomNotify"`
		ChatRoomOwner   string `json:"chatRoomOwner"`
		ChatRoomStatus  int    `json:"chatRoomStatus"`
		Description     string `json:"description"`
		NickName        string `json:"nickName"`
		SmallHeadImgUrl string `json:"smallHeadImgUrl"`
		UserName        string `json:"userName"`
		WId             string `json:"wId"`
	} `json:"data"` // 消息体
	MessageType string `json:"messageType"` // 消息类型
	WcId        string `json:"wcId"`        // 微信id
}

// GroupExitParam 群组信息修改
type GroupExitParam struct {
	Account string `json:"account"` // 账号
	Data    struct {
		UserName string `json:"userName"` // 群聊的ID
		WId      string `json:"wId"`
	} `json:"data"` // 消息体
	MessageType string `json:"messageType"` // 消息类型
	WcId        string `json:"wcId"`        // 微信id
}

// GroupInviteParam 邀请进入群聊
type GroupInviteParam struct {
	Account string `json:"account"`
	Data    struct {
		Content   string `json:"content"`
		FromGroup string `json:"fromGroup"`
		FromUser  string `json:"fromUser"`
		MsgId     int    `json:"msgId"`
		MsgType   int    `json:"msgType"`
		NewMsgId  int64  `json:"newMsgId"`
		Self      bool   `json:"self"`
		Timestamp int    `json:"timestamp"`
		ToUser    string `json:"toUser"`
		WId       string `json:"wId"`
	} `json:"data"`
	MessageType string `json:"messageType"`
	WcId        string `json:"wcId"`
}
