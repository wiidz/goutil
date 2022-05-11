package imMsg

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type BroadcastMsgParam struct {
	Body      string `json:"body" validate:"required"` // 广播消息内容，最大4096字符
	From      string `json:"from,omitempty"`           // 发送者accid, 用户帐号，最大长度32字符，必须保证一个APP内唯一
	IsOffline string `json:"isOffline,omitempty"`      // 是否存离线，true或false，默认false
	Ttl       int    `json:"ttl,omitempty"`            // 存离线状态下的有效期，单位小时，默认7天
	TargetOs  string `json:"targetOs,omitempty"`       // 目标客户端，默认所有客户端，jsonArray，格式：["ios","aos","pc","web","mac"]
}

type BroadcastMsgResp struct {
	*imClient.CommonResp
	Msg struct {
		ExpireTime  int64    `json:"expireTime"`
		Body        string   `json:"body"`
		CreateTime  int64    `json:"createTime"`
		IsOffline   bool     `json:"isOffline"`
		BroadcastId int64    `json:"broadcastId"`
		TargetOs    []string `json:"targetOs"`
	} `json:"msg"`
}

// BroadcastMsg 发送广播消息
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DEwMTE3NzQ?platformId=60353#发送广播消息
// 1、使用广播消息前，请务必阅读注意事项，详见关于广播消息。
// 2、广播消息，可以对应用内的所有用户发送广播消息，广播消息目前暂不支持第三方推送（APNS、小米、华为等）；
// 3、广播消息支持离线存储，并可以自定义设置离线存储的有效期，最多保留最近100条离线广播消息；
// 4、此接口受频率控制，一个应用一分钟最多调用10次，一天最多调用1000次，超过会返回416状态码；
// 5、该功能目前需申请开通，详情可咨询您的客户经理。
func (api *Api) BroadcastMsg(param *BroadcastMsgParam) (*BroadcastMsgResp, error) {
	res, err := api.Client.Post(SubDomain+"broadcastMsg.action", param, &BroadcastMsgResp{})
	return res.(*BroadcastMsgResp), err
}
