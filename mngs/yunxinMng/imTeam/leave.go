package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type LeaveParam struct {
	Tid    string `json:"tid" validate:"required"`   // 网易云信服务器产生，群唯一标识，创建群时会返回
	Accid  string `json:"accid" validate:"required"` // 退群的accid
	Attach string `json:"attach"`                    // 自定义扩展字段，最大长度512
}

// Leave 主动退群
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#主动退群
// 高级群主动退群
func (api *Api) Leave(param *LeaveParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"leave.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
