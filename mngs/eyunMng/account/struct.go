package account

import "github.com/wiidz/goutil/mngs/eyunMng/base"

type Api struct {
	Config *base.Config
}

type QueryLoginWxResp struct {
	Message string         `json:"message"`
	Code    string         `json:"code"`
	Data    []*LoginWxData `json:"data"`
}

type LoginWxData struct {
	WcId string `json:"wcId"` // 微信id
	WId  string `json:"wId"`  // 登录实例标识
}

type IsOnlineResp struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Data    struct {
		IsOnline bool `json:"isOnline"`
	} `json:"data"`
}
