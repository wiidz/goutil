package imUser

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type UpdateTokenParam struct {
	Accid string `json:"accid" validate:"required"` // 网易云信IM账号，最大长度32字符，必须保证一个APP内唯一
	Props string `json:"props"`                     // 该参数已不建议使用
	Token string `json:"token" validate:"required"` // 网易云信IM账号可以指定登录token值，最大长度128字符
}

// UpdateToken 更新网易云信IM token
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/Dc2NTM1NzI?platformId=60353#更新网易云信IM token
// 1.更新网易云信IM token。通过该接口，可以对accid更新到指定的IM token，更新后请开发者务必做好本地的维护。更新后，需要确保客户端SDK再次登录时携带的token保持最新。
func (api *Api) UpdateToken(param *UpdateTokenParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post("/user/update.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
