package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type RemoveParam struct {
	TeamID string `json:"tid" validate:"required"`   // 网易云信服务器产生，群唯一标识，创建群时会返回，最大长度128字符
	Owner  string `json:"owner" validate:"required"` // 群主用户帐号，最大长度32字符
	Attach string `json:"attach,omitempty"`          // 自定义扩展字段，最大长度512
}

// Remove 解散群
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#解散群
// 删除整个群，会解散该群，需要提供群主accid，谨慎操作！
func (api *Api) Remove(param *RemoveParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"remove.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
