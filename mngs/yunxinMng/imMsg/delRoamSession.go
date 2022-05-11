package imMsg

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type DelRoamSessionParam struct {
	Type int    `json:"type" validate:"required"` // 会话类型，1-p2p会话，2-群会话，其他返回414
	From string `json:"from" validate:"required"` // 发送者accid, 用户帐号，最大长度32字节
	To   string `json:"to" validate:"required"`   // type=1表示对端accid，type=2表示tid
}

// DelRoamSession 删除会话漫游
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DEwMTE3NzQ?platformId=60353#删除会话漫游
// 按会话删除漫游消息，可以删除p2p/群会话
func (api *Api) DelRoamSession(param *DelRoamSessionParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"delRoamSession.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
