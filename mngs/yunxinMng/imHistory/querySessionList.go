package imHistory

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type QuerySessionListParam struct {
	Accid        string `json:"accid" validate:"required"` // 账号
	MinTimestamp int64  `json:"minTimestamp,omitempty"`    // 查询的最小时间戳，单位毫秒，默认0
	MaxTimestamp int64  `json:"maxTimestamp,omitempty"`    // 查询的最大时间戳，单位毫秒，默认当前时间戳
	Limit        int    `json:"limit,omitempty"`           // 最大100，默认100
	NeedLastMsg  int    `json:"needLastMsg,omitempty"`     // 是否需要同步返回最后一条消息，0表示不需要，1表示需要，默认不需要
}

type QuerySessionListResp struct {
	*imClient.CommonResp
	Data struct {
		Sessions []struct {
			LastMsgType string `json:"lastMsgType"` // RECALL 表示最后一条消息是撤回，MESSAGE 表示最后一条消息是消息
			UpdateTime  int64  `json:"updateTime"`
			LastMsg     struct {
				DeleteMsgIdClient   string      `json:"deleteMsgIdClient,omitempty"`   // 【RECALL】被撤回消息的消息id（客户端）
				DeleteMsgCreateTime int64       `json:"deleteMsgCreateTime,omitempty"` // 【RECALL】被撤回消息的发送时间
				From                string      `json:"from"`                          // 【RECALL】撤回操作者 【MESSAGE】消息发送者
				Time                interface{} `json:"time"`                          // 【RECALL】撤回时间 【MESSAGE】发送时间
				DeleteMsgIdServer   string      `json:"deleteMsgIdServer,omitempty"`   // 【RECALL】被撤回消息的消息id（服务器）
				MsgIdClient         string      `json:"msgIdClient,omitempty"`         // 【MESSAGE】消息id（客户端）
				FromClientType      string      `json:"fromClientType,omitempty"`      // 【MESSAGE】消息发送者客户端类型
				Body                string      `json:"body,omitempty"`                // 【MESSAGE】消息内容，包括body、attach、ext三个字段，本例子仅包含了body字段
				Type                int         `json:"type,omitempty"`                // 【MESSAGE】消息类型
				MsgIdServer         int         `json:"msgIdServer,omitempty"`         // 【MESSAGE】消息id（服务器）
			} `json:"lastMsg"`
			SessionType int    `json:"sessionType"`
			Accid       string `json:"accid,omitempty"`
			TeamID      int    `json:"tid,omitempty"` // 如果是点对点会话，则包括accid字段；如果是群和超大群，则包括tid字段
		} `json:"sessions"`
		HasMore bool `json:"hasMore"`
	} `json:"data"`
}

// QuerySessionList 查询云端会话列表
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DE0MTk0OTY?platformId=60353#查询云端会话列表
// 查询云端会话列表，需要先开通云端会话列表功能
func (api *Api) QuerySessionList(param *QuerySessionListParam) (*QuerySessionListResp, error) {
	res, err := api.Client.Post(SubDomain+"querySessionList.action", param, &QuerySessionListResp{})
	return res.(*QuerySessionListResp), err
}
