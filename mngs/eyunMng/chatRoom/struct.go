package chatRoom

import (
	"github.com/wiidz/goutil/mngs/eyunMng/base"
	"github.com/wiidz/goutil/mngs/redisMng"
)

type Api struct {
	Config   *base.Config
	RedisMng *redisMng.RedisMng
}

type GetChatRoomMemberResp struct {
	Message string        `json:"message"`
	Code    string        `json:"code"`
	Data    []*MemberData `json:"data"`
}

type MemberData struct {
	ChatRoomId      string      `json:"chatRoomId"`
	UserName        string      `json:"userName"` // 群成员微信号 （假如需要手机上显示的微信号或更详细的信息，则需要再调用获取群成员详情接口获取）
	NickName        string      `json:"nickName"` // 群成员默认昵称
	ChatRoomOwner   interface{} `json:"chatRoomOwner"`
	BigHeadImgUrl   string      `json:"bigHeadImgUrl"`   // 大头像
	SmallHeadImgUrl string      `json:"smallHeadImgUrl"` // 小头像
	V1              interface{} `json:"v1"`
	MemberCount     int         `json:"memberCount"`
	DisplayName     string      `json:"displayName"` // 群成员修改后的昵称
	ChatRoomMembers interface{} `json:"chatRoomMembers"`
	InviterUserName string      `json:"inviterUserName"` // 邀请人微信号（仅有群主和管理可以看到）
}

type GetChatRoomMemberInfoResp struct {
	Message string            `json:"message"`
	Code    string            `json:"code"`
	Data    []*MemberInfoData `json:"data"`
}

type MemberInfoData struct {
	UserName  string      `json:"userName"` // 微信id
	NickName  string      `json:"nickName"` // 昵称
	Remark    string      `json:"remark"`
	Signature string      `json:"signature"` // 签名
	Sex       int         `json:"sex"`       // 性别
	AliasName string      `json:"aliasName"` // 微信号
	Country   interface{} `json:"country"`
	BigHead   string      `json:"bigHead"`   // 大头像
	SmallHead string      `json:"smallHead"` // 小头像
	LabelList interface{} `json:"labelList"`
	V1        string      `json:"v1"`
	V2        string      `json:"v2"`

	PlatformID uint64 `json:"platform_id"` // 自己加的，在平台里的ID
}

// GetChatRoomInfoResp 获取群信息
type GetChatRoomInfoResp struct {
	Message string              `json:"message"`
	Code    string              `json:"code"`
	Data    []*ChatRoomInfoData `json:"data"`
}

type ChatRoomInfoData struct {
	ChatRoomId      string      `json:"chatRoomId"` // 群号
	UserName        interface{} `json:"userName"`
	NickName        string      `json:"nickName"`        // 群名称
	ChatRoomOwner   string      `json:"chatRoomOwner"`   // 群主
	BigHeadImgUrl   interface{} `json:"bigHeadImgUrl"`   // 大头像
	SmallHeadImgUrl string      `json:"smallHeadImgUrl"` // 小头像
	V1              string      `json:"v1"`
	MemberCount     int         `json:"memberCount"` // 群成员数
	ChatRoomMembers []struct {
		UserName        string `json:"userName"`        // 群成员微信号
		NikeName        string `json:"nikeName"`        // 群成员昵称
		InviterUserName string `json:"inviterUserName"` // 邀请人微信号（仅有群主和管理可以看到）
	} `json:"chatRoomMembers"`
}
