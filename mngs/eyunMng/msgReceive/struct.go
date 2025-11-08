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
type AudioParam struct {
	Account     string     `json:"account"`     // 账号
	Data        *AudioData `json:"data"`        // 消息体
	MessageType string     `json:"messageType"` // 消息类型
	WcId        string     `json:"wcId"`        // 微信id
}
type AudioData struct {
	WId string `json:"wId"`

	BufId       string `json:"bufId"` // 下载语音时会用到
	Content     string `json:"content"`
	FromUser    string `json:"fromUser"`
	Length      int    `json:"length"` // 语音文件大小(b)
	MsgId       int    `json:"msgId"`
	NewMsgId    int64  `json:"newMsgId"`
	Self        bool   `json:"self"`
	Timestamp   int    `json:"timestamp"`
	ToUser      string `json:"toUser"`
	VoiceLength int    `json:"voiceLength"` // 语音时长(ms)
}

type DownloadAudioResp struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Data    *DownloadAudioData `json:"data"`
}

type DownloadAudioData struct {
	URL string `json:"url"`
}

type MessageType int

const PersonalTextMsg MessageType = 60001   // 私聊文本
const PersonalPic MessageType = 60002       // 私聊图片
const PersonalVideo MessageType = 60003     // 私聊视频
const PersonalAudio MessageType = 60004     // 私聊语音
const PersonalIDCard MessageType = 60005    // 私聊名片
const PersonalEmoji MessageType = 60006     // 私聊emoji
const PersonalLink MessageType = 60007      // 私聊链接
const PersonalFile MessageType = 60008      // 私聊文件
const PersonalFileDone MessageType = 60009  // 私聊文件发送完成消息
const PersonalMini MessageType = 60010      // 私聊小程序
const PersonalChatLog MessageType = 60011   // 私聊聊天记录
const PersonalTel MessageType = 60012       // 私聊语音请求
const PersonalTelCancel MessageType = 60013 // 语音聊天挂断
const PersonalQuote MessageType = 60014     // 引用消息
const PersonalTransfer MessageType = 60015  // 转账
const PersonalRedPack MessageType = 60016   // 红包
const PersonalVideoAcc MessageType = 60017  // 视频号
const PersonalCallback MessageType = 60018  // 撤回
const PersonalPai MessageType = 60019       // 拍一拍
const PersonalLocation MessageType = 60020  // 位置

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
