package imFriend

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

// 更新好友相关信息
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DQ0MTY1NzI?platformId=60353
// 更新好友相关信息，如加备注名，必须是好友才可以

type UpdateParam struct {
	Accid    string `json:"accid" validate:"required"`  // Accid 加好友发起者accid
	FAccid   string `json:"faccid" validate:"required"` // FAccid 加好友接收者accid
	Alias    int    `json:"alias"`                      // 给好友增加备注名，限制长度128，可设置为空字符串
	Ex       string `json:"ex"`                         // 修改ex字段，限制长度256，可设置为空字符串
	ServerEX string `json:"serverex"`                   // 服务器端扩展字段，限制长度256 此字段client端只读，server端读写
}

func (api *Api) Update(param *UpdateParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"update.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
