package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type KickParam struct {
	TeamID  string `json:"tid" validate:"required"`   // 网易云信服务器产生，群唯一标识，创建群时会返回，最大长度128字符
	Owner   string `json:"owner" validate:"required"` // 用户帐号，最大长度32字符，按照群属性invitemode传入
	Member  string `json:"member"`                    // 被移除人的accid，用户账号，最大长度32字符;注：member或members任意提供一个，优先使用member参数
	Members string `json:"members"`                   // ["aaa","bbb"]（JSONArray对应的accid，如果解析出错，会报414）一次最多操作200个accid; 注：member或members任意提供一个，优先使用member参数
	Attach  string `json:"attach,omitempty"`          // 自定义扩展字段，最大长度512
}

// Kick 踢人出群
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#踢人出群
// 高级群踢人出群，需要提供群主accid以及要踢除人的accid。
func (api *Api) Kick(param *KickParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"kick.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
