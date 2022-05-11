package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type ChangeOwnerParam struct {
	Tid      string `json:"tid" validate:"required"`      // 网易云信服务器产生，群唯一标识，创建群时会返回，最大长度128字符
	Owner    string `json:"owner" validate:"required"`    // 群主用户帐号，最大长度32字符
	NewOwner string `json:"newowner" validate:"required"` // 新群主帐号，最大长度32字符
	Leave    int    `json:"leave" validate:"required"`    // 1:群主解除群主后离开群，2：群主解除群主后成为普通成员。其它414
}

// ChangeOwner 移交群主
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#移交群主
// 1.转换群主身份；
// 2.群主可以选择离开此群，还是留下来成为普通成员。
func (api *Api) ChangeOwner(param *ChangeOwnerParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"changeOwner.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
