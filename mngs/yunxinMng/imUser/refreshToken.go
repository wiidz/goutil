package imUser

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type RefreshTokenParam struct {
	Accid string `validate:"required" json:"accid"` // Accid 用户帐号
}
type RefreshTokenResp struct {
	*imClient.CommonResp
	Info struct {
		Token string `json:"token"`
		Accid string `json:"accid"`
	} `json:"info"`
}

// ResetToken 重置网易云信IM token
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/Dc2NTM1NzI?platformId=60353#重置网易云信IM token
// 1.由云信webserver随机重置网易云信IM账号的token，同时将新的token返回，更新后请开发者务必做好本地的维护。
// 2.此接口与更新网易云信IM token 接口最大的区别在于：前者的token是由云信服务器指定，后者的token是由开发者自己指定。
func (api *Api) ResetToken(param *RefreshTokenParam) (*RefreshTokenResp, error) {
	res, err := api.Client.Post("/user/refreshToken.action", param, &RefreshTokenResp{})
	return res.(*RefreshTokenResp), err
}
