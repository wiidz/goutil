package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type AddManagerParam struct {
	Tid     string `json:"tid" validate:"required"`     // 网易云信服务器产生，群唯一标识，创建群时会返回，最大长度128字符
	Owner   string `json:"owner" validate:"required"`   // 群主用户帐号，最大长度32字符
	Members string `json:"members" validate:"required"` // ["aaa","bbb"](JSONArray对应的accid，如果解析出错会报414)，长度最大1024字符（一次添加最多10个管理员）
	Attach  string `json:"attach,omitempty"`            // 自定义扩展字段，最大长度512
}

// AddManager 任命管理员
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#任命管理员
// 提升普通成员为群管理员，可以批量，但是一次添加最多不超过10个人。
func (api *Api) AddManager(param *AddManagerParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"addManager.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
