package kookMng

import "github.com/wiidz/goutil/structs/configStruct"

var suffixs = struct {
	ChanelMessageCreate string
	ChanelMessageDelete string
	DirectMsgCreate     string
	UserChatCreate      string
	EditUserNickname    string
	AddUserGroup        string
	DelUserGroup        string
	GetUserGroupList    string
	AddUserBlackList    string
	DelUserBlackList    string
}{
	ChanelMessageCreate: "/message/create",        // 发送频道消息
	ChanelMessageDelete: "/message/delete",        // 删除频道消息
	UserChatCreate:      "/user-chat/create",      // 发送私聊消息
	DirectMsgCreate:     "/direct-message/create", // 创建私聊会话
	EditUserNickname:    "/guild/nickname",        // 修改服务器中用户的昵称
	AddUserGroup:        "/guild-role/grant",      // 赋予用户角色
	DelUserGroup:        "/guild-role/revoke",     // 删除用户角色
	GetUserGroupList:    "/guild-role/list",       // 获取服务器角色列表
	AddUserBlackList:    "/blacklist/create",      //  加入黑名单
	DelUserBlackList:    "/blacklist/delete",      //  移除黑名单
}

type KookMng struct {
	Config *configStruct.KookConfig
}

type ReceiveData struct {
	S int      `json:"s"`
	D *MsgData `json:"d"`
}

type MsgData struct {
	ChannelType string `json:"channel_type"` // 消息通道类型, GROUP 为组播消息, PERSON 为单播消息, BROADCAST 为广播消息
	Type        int    `json:"type"`         // 1:文字消息, 2:图片消息，3:视频消息，4:文件消息， 8:音频消息，9:KMarkdown，10:card 消息，255:系统消息, 其它的暂未开放

	// 消息独有
	TargetID     string `json:"target_id"`     // 发送目的, 频道消息类时, 代表的是频道 channel_id，如果 channel_type 为 GROUP 组播且 type 为 255 系统消息时，则代表服务器 guild_id
	AuthorID     string `json:"author_id"`     // 发送者 id, 1 代表系统
	Content      string `json:"content"`       // 消息内容, 文件，图片，视频时，content 为 url
	MsgID        string `json:"msg_id"`        // 消息的 id
	MsgTimeStamp int    `json:"msg_timestamp"` // 消息发送时间的毫秒时间戳
	Nonce        string `json:"nonce"`         // 随机串，与用户消息发送 api 中传的 nonce 保持一致

	// 当 type 非系统消息(255)时
	Extra *Extra `json:"extra"` // 不同的消息类型，结构不一致

	// challenge独有
	VerifyToken string `json:"verify_token"`
	Challenge   string `json:"challenge"`
}

// Extra 表示事件消息的数据结构
type Extra struct {
	//Type string `json:"type"` // 事件类型，与ReceiveData的type相同
	//Type int `json:"type"` // 事件类型，与ReceiveData的type相同

	// 系统事件消息
	Body map[string]interface{} `json:"body"`

	// 消息(非255)独有
	GuildID      string   `json:"guild_id"`      // 服务器 ID
	ChannelName  string   `json:"channel_name"`  // 频道名
	Mention      []string `json:"mention"`       // 提及到的用户 ID 列表
	MentionAll   bool     `json:"mention_all"`   // 是否提及所有用户
	MentionRoles []string `json:"mention_roles"` // 提及用户角色的 ID 数组
	MentionHere  bool     `json:"mention_here"`  // 是否提及在线用户
	Author       User     `json:"author"`        // 用户信息
}

// User 表示用户对象结构体，补充一个示例
type User struct {
	ID       string `json:"id"`       // 用户 ID
	Username string `json:"username"` // 用户名
	Avatar   string `json:"avatar"`   // 用户头像地址
}

// ApiResponse 接口返回的消息
type ApiResponse struct {
	Code    int       `json:"code"`    // integer, 错误码，0代表成功，非0代表失败，具体的错误码参见错误码一览
	Message string    `json:"message"` // string, 错误消息，具体的返回消息会根据Accept-Language来返回。
	Data    *RespData `json:"data"`    // mixed, 具体的数据。
}

type RespData struct {
	MsgID        string `json:"msg_id"` // string
	MsgTimestamp int    `json:"msg_timestamp"`
	Nonce        string `json:"nonce"`
}

// CreateChanelMsgParam 表示发送频道消息时请求的结构体
type CreateChanelMsgParam struct {
	Type         int    `json:"type,omitempty"`           // 消息类型, 默认为 9 (kmarkdown), 10 代表卡片消息
	TargetID     string `json:"target_id"`                // 目标频道 id
	Content      string `json:"content"`                  // 消息内容
	Quote        string `json:"quote,omitempty"`          // 回复某条消息的 msgId
	Nonce        string `json:"nonce,omitempty"`          // nonce, 服务端不处理, 原样返回
	TempTargetID string `json:"temp_target_id,omitempty"` // 用户 id, 代表临时消息, 只会推送给该用户，不存数据库
	TemplateID   string `json:"template_id,omitempty"`    // 模板消息id, 使用时content为模板消息input
}

// DeleteChanelMsgParam 删除频道消息
type DeleteChanelMsgParam struct {
	MsgID string `json:"msg_id"` // 消息ID
}

// CreateUserChatMsgParam 表示发送私聊消息时请求的结构体
type CreateUserChatMsgParam struct {
}

// CreateDirectMessageParam 表示私信消息发送的请求结构体
type CreateDirectMessageParam struct {
	Type       int    `json:"type,omitempty"`        // 消息类型, 默认为 1 (文本), 9 表示 kmarkdown, 10 表示卡片消息
	TargetID   string `json:"target_id,omitempty"`   // 目标用户 id，后端自动创建会话，与 chat_code 二选一必填
	ChatCode   string `json:"chat_code,omitempty"`   // 目标会话 code，与 target_id 二选一必填
	Content    string `json:"content"`               // 消息内容
	Quote      string `json:"quote,omitempty"`       // 回复某条消息的 msgId
	Nonce      string `json:"nonce,omitempty"`       // nonce，服务端不做处理，原样返回
	TemplateID string `json:"template_id,omitempty"` // 模板消息id，使用时 content 作为模板消息 input
}

// CreateUserChatParam 创建私聊会话
type CreateUserChatParam struct {
	TargetID string `json:"target_id,omitempty"` // 目标用户 id，后端自动创建会话，与 chat_code 二选一必填
}

type EditUserNicknameParam struct {
	GuildID  string `json:"guild_id"` // 服务器的 ID
	Nickname string `json:"nickname"` // 昵称，2 - 64 长度，不传则清空昵称
	UserID   string `json:"user_id"`  // 要修改昵称的目标用户 ID，不传则修改当前登陆用户的昵称
}

// GetUserGroupListParam 获取服务器角色列表
type GetUserGroupListParam struct {
	GuildID  string `json:"guild_id"`  // 服务器 id
	Page     int    `json:"page"`      // 目标页数
	PageSize int    `json:"page_size"` // 每页数据数量
}

// UpdateUserGroupParam 赋予用户角色
type UpdateUserGroupParam struct {
	GuildID string `json:"guild_id"` // 服务器 id
	UserID  string `json:"user_id"`  // 要修改昵称的目标用户 ID，不传则修改当前登陆用户的昵称
	RoleID  uint64 `json:"role_id"`  // 服务器角色 id
}

// AddUserBlackListParam 加入黑名单
type AddUserBlackListParam struct {
	GuildID    string `json:"guild_id"`     // 服务器 id(必填)
	TargetID   string `json:"target_id"`    // 目标用户 id(必填)
	Remark     string `json:"remark"`       // 加入黑名单的原因
	DelMsgDays int    `json:"del_msg_days"` // 删除最近几天的消息，最大 7 天, 默认 0
}

// DelUserBlackListParam 移除黑名单
type DelUserBlackListParam struct {
	GuildID  string `json:"guild_id"`  // 服务器 id(必填)
	TargetID string `json:"target_id"` // 目标用户 id(必填)
}
