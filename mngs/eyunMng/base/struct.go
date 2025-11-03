package base

type ReturnCode string

const Success ReturnCode = "1000"
const Failed ReturnCode = "1001"

type Config struct {
	Account       string       // e云管家账号
	Password      string       // e云管家密码
	Authorization string       // e云管家token
	Host          string       // e云管家API开通信息里的IP加PORT(http协议的)
	ProxyConfig   *ProxyConfig // e云管家代理设置

	//WechatConfig *WechatConfig          // 微信设置
	WechatData *WechatData // 微信账户数据
	//WechatAccount string // 微信账号
	//WID           string // 登录实例标识（每次重新登录都会不一样）
	//WcID          string // 微信原始ID
}

// ProxyConfig 代理服务器设置
type ProxyConfig struct {
	Host     string // 自定义长效代理IP+端口
	Username string // 自定义长效代理IP平台账号
	Password string // 自定义长效代理IP平台密码
}

// BaseResp 基础
type BaseResp struct {
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Data    interface{} `json:"data"`
}

// WechatConfig 微信账号信息
type WechatConfig struct {
	WechatAccount string // 微信账号
	WID           string // 登录实例标识（每次重新登录都会不一样）
	WcID          string // 微信原始ID
}

type WechatData struct {
	WID        string `json:"wId"`        // 登录实例标识（"25d50610-1a82-4531-b9db-dd80c5a3c14a"）
	DeviceType string `json:"deviceType"` // 扫码的设备类型（"android"）
	Type       int    `json:"type"`
	Uin        int    `json:"uin"`    // 识别码
	Status     int    `json:"status"` // 保留字段

	WcID        string `json:"wcId"`        // 微信id (唯一值）
	WAccount    string `json:"wAccount"`    // 手机上显示的微信号（用户若手机改变微信号，本值会变）
	MobilePhone string `json:"mobilePhone"` // 绑定手机
	Username    string `json:"username"`    // 登录用户名（手机号）

	// 个人信息
	NickName        string `json:"nickName"`        // 昵称
	HeadUrl         string `json:"headUrl"`         // 头像url
	SmallHeadImgUrl string `json:"smallHeadImgUrl"` // 头像缩略图
	Sex             int    `json:"sex"`             // 性别
	Signature       string `json:"signature"`       // 个性签名
	Country         string `json:"country"`         // 国家（"CN"）
	City            string `json:"city"`            // 城市

	//IsOnline bool   `json:"is_online"` // 是否在线（远程不返回的，我们是自己添加的字段）
	LoginAt string `json:"login_at"` // 登录时间
}
