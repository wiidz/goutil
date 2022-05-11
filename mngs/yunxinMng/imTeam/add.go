package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type AddParam struct {
	TeamID       string `json:"tid" validate:"required"`     // 网易云信服务器产生，群唯一标识，创建群时会返回，最大长度128字符
	Owner        string `json:"owner" validate:"required"`   // 用户帐号，最大长度32字符，按照群属性invitemode传入
	Members      string `json:"members" validate:"required"` // ["aaa","bbb"](JSONArray对应的accid，如果解析出错会报414)，一次最多拉200个成员
	ManagerAgree int    `json:"magree" validate:"required"`  // 管理后台建群时，0不需要被邀请人同意加入群，1需要被邀请人同意才可以加入群。其它会返回414
	Msg          string `json:"msg" validate:"required"`     // 邀请发送的文字，最大长度150字符
	Attach       string `json:"attach,omitempty"`            // 自定义扩展字段，最大长度512
}

type AddResp struct {
	*imClient.CommonResp
	TeamID string `json:"tid"`
	FAccid struct {
		Accid []string `json:"accid"` // ["a","b","c"]
		Msg   string   `json:"msg"`   // team count exceed
	} `json:"faccid"` // 如果创建时邀请的成员中存在加群数量超过限制的情况，会返回faccid
}

// Add 拉人入群
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#拉人入群
// 1.可以批量邀请，邀请时需指定群主；
// 2.当群成员达到上限时，再邀请某人入群返回失败；
// 3.当群成员达到上限时，被邀请人“接受邀请"的操作也将返回失败。
func (api *Api) Add(param *AddParam) (*AddResp, error) {
	res, err := api.Client.Post(SubDomain+"add.action", param, &AddResp{})
	return res.(*AddResp), err
}
