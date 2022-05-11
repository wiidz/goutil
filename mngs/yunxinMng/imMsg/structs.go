package imMsg

type Option struct {
	Roam          bool `json:"roam"`          // 该消息是否需要漫游，默认true（需要app开通漫游消息功能）
	History       bool `json:"history"`       // 该消息是否存云端历史，默认true
	SenderSync    bool `json:"sendersync"`    // 该消息是否需要发送方多端同步，默认true
	Push          bool `json:"push"`          // 该消息是否需要APNS推送或安卓系统通知栏推送，默认true
	Route         bool `json:"route"`         // 该消息是否需要抄送第三方；默认true (需要app开通消息抄送功能)
	Badge         bool `json:"badge"`         // 该消息是否需要计入到未读计数中，默认true
	NeedPushNick  bool `json:"NeedPushNick"`  // 推送文案是否需要带上昵称，不设置该参数时默认true
	Persistent    bool `json:"persistent"`    // 是否需要存离线消息，不设置该参数时默认true
	SessionUpdate bool `json:"sessionUpdate"` // 是否将本消息更新到会话列表服务里本会话的lastmsg，默认true
}

// Ope 消息类型
type Ope int

const P2P Ope = 0
const Group Ope = 1

type NoticeOption struct {
	Badge        bool `json:"badge"`        // 该消息是否需要计入到未读计数中，默认true
	NeedPushNick bool `json:"NeedPushNick"` // 推送文案是否需要带上昵称，不设置该参数时默认false(ps:注意与sendMsg.action接口有别);
	Route        bool `json:"route"`        // 该消息是否需要抄送第三方；默认true (需要app开通消息抄送功能)
}
