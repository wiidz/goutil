package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type UpdateTeamNickParam struct {
	Tid    string `json:"tid" validate:"required"`   // 群唯一标识，创建群时网易云信服务器产生并返回
	Owner  string `json:"owner" validate:"required"` // 群主 accid
	Accid  string `json:"accid" validate:"required"` // 要修改群昵称的群成员 accid
	Nick   string `json:"nick,omitempty"`            // accid 对应的群昵称，最大长度32字符
	Custom string `json:"custom,omitempty"`          // 自定义扩展字段，最大长度1024字节
}

// UpdateTeamNickParam 修改群昵称（群名片）
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#修改群昵称
// 修改指定账号在群内的昵称
func (api *Api) UpdateTeamNickParam(param *UpdateTeamNickParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"updateTeamNick.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
