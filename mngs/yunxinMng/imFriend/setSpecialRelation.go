package imFriend

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type SetSpecialRelationParam struct {
	Accid        string `json:"accid" validate:"required"`     // 发起者accid
	TargetAccid  string `json:"targetAcc" validate:"required"` // 被加黑或加静音的帐号
	RelationType int    `json:"relationType"`                  // 本次操作的关系类型,1:黑名单操作，2:静音列表操作
	Value        int    `json:"value"`                         // 操作值，0:取消黑名单或静音，1:加入黑名单或静音
}

// SetSpecialRelation 设置黑名单/静音
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DQ0MTY1NzI?platformId=60353
// 拉黑/取消拉黑；设置静音/取消静音
func (api *Api) SetSpecialRelation(param *SetSpecialRelationParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"setSpecialRelation.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
