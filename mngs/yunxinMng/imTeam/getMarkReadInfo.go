package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type GetMarkReadInfoParam struct {
	TeamID    int64  `json:"tid" validate:"required"`       // 群id，群唯一标识，创建群时会返回
	MsgID     int64  `json:"msgid" validate:"required"`     // 发送群已读业务消息时服务器返回的消息ID
	FromAccid string `json:"fromaccid" validate:"required"` // 消息发送者账号
	SnapShot  bool   `json:"snapshot"`                      // 是否返回已读、未读成员的accid列表，默认为false
}

type GetMarkReadInfoResp struct {
	*imClient.CommonResp
	Data struct {
		ReadSize     int      `json:"readSize"`
		UnreadSize   int      `json:"unreadSize"`
		ReadAccids   []string `json:"readAccids"`
		UnreadAccids []string `json:"unreadAccids"`
	} `json:"data"`
}

// GetMarkReadInfo 获取群组已读消息的已读详情信息
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#获取群组已读消息的已读详情信息
// 获取群组已读消息的已读详情信息
func (api *Api) GetMarkReadInfo(param *GetMarkReadInfoParam) (*GetMarkReadInfoResp, error) {
	res, err := api.Client.Post(SubDomain+"getMarkReadInfo.action", param, &GetMarkReadInfoResp{})
	return res.(*GetMarkReadInfoResp), err
}
