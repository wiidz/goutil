package imUser

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type UnBlockParam struct {
	Accid string `validate:"required" json:"accid"` // Accid 用户帐号
}

// UnBlock 解禁网易云信IM账号
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/Dc2NTM1NzI?platformId=60353#%E5%B0%81%E7%A6%81%E7%BD%91%E6%98%93%E4%BA%91%E4%BF%A1IM%E8%B4%A6%E5%8F%B7
// 1.封禁网易云信IM账号后，此ID将不能再次登录。若封禁时，该id处于登录状态，则当前登录不受影响，仍然可以收发消息。封禁效果会在下次登录时生效。因此建议，将needkick设置为true，让该账号同时被踢出登录。
// 2.出于安全目的，账号创建后只能封禁，不能删除；封禁后账号仍计入应用内账号总数。
func (api *Api) UnBlock(param *UnBlockParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post("/user/unblock.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
