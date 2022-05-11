package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type MuteMemberParam struct {
	Tid    string `json:"tid" validate:"required"`   // 网易云信服务器产生，群唯一标识，创建群时会返回
	Owner  string `json:"owner" validate:"required"` // 群主accid
	Accid  string `json:"accid" validate:"required"` // 禁言对象的accid
	Mute   int    `json:"mute" validate:"required"`  // 1-禁言，0-解禁
	Attach string `json:"attach"`                    // 自定义扩展字段，最大长度512
}

// MuteMember 禁言群成员
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#禁言群成员
// 高级群禁言群成员
func (api *Api) MuteMember(param *MuteMemberParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"muteTlist.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
