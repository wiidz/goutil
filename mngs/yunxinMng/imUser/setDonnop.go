package imUser

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type SetDonnopParam struct {
	Accid      string `json:"accid" validate:"required"`       // 用户帐号
	DonnopOpen string `json:"donnop_open" validate:"required"` // 桌面端在线时，移动端是否不推送： true:移动端不需要推送，false:移动端需要推送
}

// SetDonnop 用户设置
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/TU5Mjc5MTg?platformId=60353
// 设置桌面端在线时，移动端是否需要推送
func (api *Api) SetDonnop(param *SetDonnopParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post("/user/setDonnop.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
