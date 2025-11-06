package login

import "github.com/wiidz/goutil/mngs/eyunMng/base"

type Api struct {
	Config *base.Config
}

type DeviceType string

const IPad DeviceType = "ipad"
const Mac DeviceType = "mac"

// EYunLoginResp e云管家登录
type EYunLoginResp struct {
	Message string `json:"message"` // 反馈信息
	Code    string `json:"code"`    // 1000成功，1001失败
	Data    struct {
		CallbackUrl   interface{} `json:"callbackUrl"`   // 消息回调地址
		Status        int         `json:"status"`        // 状态（0：正常，1：冻结，2：到期）
		Authorization string      `json:"Authorization"` // 授权密钥，生成后永久有效
	} `json:"data"`
}

// WechatQrcodeResp 获取微信登录二维码
type WechatQrcodeResp struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Data    struct {
		WId       string `json:"wId"`       // 登录实例标识 （本值非固定的，每次重新登录会返回新的，数据库记得实时更新wid）
		QrCodeUrl string `json:"qrCodeUrl"` // 扫码登录地址
	} `json:"data"`
}

// WechatLoginResp 微信登录
type WechatLoginResp struct {
	Code    string           `json:"code"`    // 1000成功，1001失败
	Message string           `json:"message"` // 反馈信息
	Data    *base.WechatData `json:"data"`
}

//type WechatLoginData struct {
//	WId        string `json:"wId"`        // 登录实例标识（"25d50610-1a82-4531-b9db-dd80c5a3c14a"）
//	DeviceType string `json:"deviceType"` // 扫码的设备类型（"android"）
//	Type       int    `json:"type"`
//	Uin        int    `json:"uin"`    // 识别码
//	Status     int    `json:"status"` // 保留字段
//
//	WcId        string `json:"wcId"`        // 微信id (唯一值）
//	WAccount    string `json:"wAccount"`    // 手机上显示的微信号（用户若手机改变微信号，本值会变）
//	MobilePhone string `json:"mobilePhone"` // 绑定手机
//	Username    string `json:"username"`    // 登录用户名（手机号）
//
//	// 个人信息
//	NickName        string `json:"nickName"`        // 昵称
//	HeadUrl         string `json:"headUrl"`         // 头像url
//	SmallHeadImgUrl string `json:"smallHeadImgUrl"` // 头像缩略图
//	Sex             int    `json:"sex"`             // 性别
//	Signature       string `json:"signature"`       // 个性签名
//	Country         string `json:"country"`         // 国家（"CN"）
//	City            string `json:"city"`            // 城市
//}

// AddressListResp 通讯录列表
type AddressListResp struct {
	Code    string           `json:"code"`
	Message string           `json:"message"`
	Data    *AddressListData `json:"data"`
}

type AddressListData struct {
	ChatRooms []string `json:"chatrooms"` // 群组列表
	Friends   []string `json:"friends"`   // 好友列表
	Ghs       []string `json:"ghs"`       // 公众号列表
	Others    []string `json:"others"`    // 微信其他相关
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
