package imHistory

import (
	"github.com/wiidz/goutil/mngs/yunxinMng/imClient"
	"github.com/wiidz/goutil/mngs/yunxinMng/imStruct"
)

type QuerySessionMsgParam struct {
	From      string `json:"from" validate:"required"`      // 发送者accid
	To        string `json:"to" validate:"required"`        // 接收者accid
	BeginTime string `json:"begintime" validate:"required"` // 开始时间，毫秒级
	EndTime   string `json:"endtime" validate:"required"`   // 截止时间，毫秒级
	Limit     int    `json:"limit" validate:"required"`     // 本次查询的消息条数上限(最多100条),小于等于0，或者大于100，会提示参数错误
	Reverse   int    `json:"reverse,omitempty"`             // 1按时间正序排列，2按时间降序排列。其它返回参数414错误.默认是按降序排列，即时间戳最晚的消息排在最前面。
	Type      string `json:"type,omitempty"`                // 查询指定的多个消息类型，类型之间用","分割，不设置该参数则查询全部类型消息格式示例： 0,1,2,3，类型支持： 1:图片，2:语音，3:视频，4:地理位置，5:通知，6:文件，10:提示，11:Robot，100:自定义
}

type QuerySessionMsgResp struct {
	*imClient.CommonResp
	Size int           `json:"size"` // 总共消息条数
	Msgs []*HistoryMsg `json:"msgs"` // 消息
}

type HistoryMsg struct {
	From           string                `json:"from"`
	MsgID          int                   `json:"msgid"`
	SendTime       int64                 `json:"sendtime"`       //发送时间ms
	Type           int                   `json:"type"`           // 消息类型，对应去看yunxinMng.MsgType
	FromClientType int                   `json:"fromclienttype"` // //1：android、2:iOS、4：PC、16:WEB、32:REST、64:MAC
	MsgIDClient    string                `json:"msgidclient"`
	Body           imStruct.MsgInterface `json:"boy"`
}

// QuerySessionMsg 单聊云端历史消息查询
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DE0MTk0OTY?platformId=60353#单聊云端历史消息查询
// 查询存储在网易云信服务器中的单人聊天历史消息，只能查询在保存时间范围内的消息
// 1. 根据时间段查询点对点消息，每次最多返回100条；
// 2. 不提供分页支持，第三方需要根据时间段来查询。
// 3. begintime需要早于endtime，否则会返回{"desc": "bad time", "code": 414}。
func (api *Api) QuerySessionMsg(param *QuerySessionMsgParam) (*QuerySessionMsgResp, error) {
	res, err := api.Client.Post(SubDomain+"querySessionMsg.action", param, &QuerySessionMsgResp{})
	return res.(*QuerySessionMsgResp), err
}
