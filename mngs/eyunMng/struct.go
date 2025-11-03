package eyunMng

import (
	"github.com/wiidz/goutil/mngs/eyunMng/account"
	"github.com/wiidz/goutil/mngs/eyunMng/base"
	"github.com/wiidz/goutil/mngs/eyunMng/chatRoom"
	"github.com/wiidz/goutil/mngs/eyunMng/login"
	"github.com/wiidz/goutil/mngs/eyunMng/msgReceive"
	"github.com/wiidz/goutil/mngs/eyunMng/msgSend"
	"github.com/wiidz/goutil/mngs/redisMng"
)

type EYunMng struct {
	Config        *base.Config
	RedisMng      *redisMng.RedisMng
	loginApi      *login.Api      // 登录
	msgSendApi    *msgSend.Api    // 发送消息
	msgReceiveApi *msgReceive.Api // 接受消息
	chatRoom      *chatRoom.Api   // 聊天群
	account       *account.Api    // 账号
}

type TestReturn struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Data    struct {
		CallbackUrl   interface{} `json:"callbackUrl"`
		Status        int         `json:"status"`
		Authorization string      `json:"Authorization"`
	} `json:"data"`
}

type QrcodeReturn struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Data    struct {
		WId       string `json:"wId"`
		QrCodeUrl string `json:"qrCodeUrl"`
	} `json:"data"`
}

type AddressListReturn struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Chatrooms []string `json:"chatrooms"`
		Friends   []string `json:"friends"`
		Ghs       []string `json:"ghs"`
		Others    []string `json:"others"`
	} `json:"data"`
}

type SendMsgReturn struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Type       int    `json:"type"`
		MsgId      int64  `json:"msgId"`
		NewMsgId   int64  `json:"newMsgId"`
		CreateTime int    `json:"createTime"`
		WcId       string `json:"wcId"`
	} `json:"data"`
}

type LoginReturn struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Country         string `json:"country"`
		WAccount        string `json:"wAccount"`
		DeviceType      string `json:"deviceType"`
		City            string `json:"city"`
		Signature       string `json:"signature"`
		NickName        string `json:"nickName"`
		Sex             int    `json:"sex"`
		HeadUrl         string `json:"headUrl"`
		Type            int    `json:"type"`
		SmallHeadImgUrl string `json:"smallHeadImgUrl"`
		WcId            string `json:"wcId"`
		WId             string `json:"wId"`
		MobilePhone     string `json:"mobilePhone"`
		Uin             int    `json:"uin"`
		Status          int    `json:"status"`
		Username        string `json:"username"`
	} `json:"data"`
}
