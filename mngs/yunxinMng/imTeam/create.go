package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type CreateParam struct {
	TeamName           string `json:"tname" validate:"required"`    // 群名称，最大长度64字符
	Owner              string `json:"owner" validate:"required"`    // 群主用户帐号，最大长度32字符
	Members            string `json:"members" validate:"required"`  // 邀请的群成员列表。["aaa","bbb"](JSONArray对应的accid，如果解析出错会报414)，members与owner总和上限为200。members中无需再加owner自己的账号。
	Announcement       string `json:"announcement,omitempty"`       // 群公告，最大长度1024字符
	Intro              string `json:"intro,omitempty"`              // 群描述，最大长度512字符
	Msg                string `json:"msg" validate:"required"`      // 邀请发送的文字，最大长度150字符
	ManagerAgree       int    `json:"magree" validate:"required"`   // 管理后台建群时，0不需要被邀请人同意加入群，1需要被邀请人同意才可以加入群。其它会返回414
	JoinMode           int    `json:"joinmode" validate:"required"` // 群建好后，sdk操作时，0不用验证，1需要验证,2不允许任何人加入。其它返回414
	Custom             string `json:"custom,omitempty"`             // 自定义高级群扩展属性，第三方可以跟据此属性自定义扩展自己的群属性。（建议为json）,最大长度1024字符
	Icon               string `json:"icon,omitempty"`               // 群头像，最大长度1024字符
	BeInviteMode       int    `json:"beinvitemode,omitempty"`       // 被邀请人同意方式，0-需要同意(默认),1-不需要同意。其它返回414
	InviteMode         int    `json:"invitemode,omitempty"`         // 谁可以邀请他人入群，0-管理员(默认),1-所有人。其它返回414
	UpdateTeamInfoMode int    `json:"uptinfomode,omitempty"`        // 谁可以修改群资料，0-管理员(默认),1-所有人。其它返回414
	UpdateCustomMode   int    `json:"upcustommode,omitempty"`       // 谁可以更新群自定义属性，0-管理员(默认),1-所有人。其它返回414
	TeamMemberLimit    int    `json:"teamMemberLimit,omitempty"`    // 该群最大人数(包含群主)，范围：2至应用定义的最大群人数(默认:200)。其它返回414
	Attach             string `json:"attach,omitempty"`             // 	自定义扩展字段，最大长度512
	Bid                string `json:"bid,omitempty"`                // 反垃圾业务ID，JSON字符串，{"textbid":"","picbid":""}，若不填则使用原来的反垃圾配置
}

type CreateResp struct {
	*imClient.CommonResp
	TeamID string `json:"tid"`
	FAccid struct {
		Accid []string `json:"accid"` // ["a","b","c"]
		Msg   string   `json:"msg"`   // team count exceed
	} `json:"faccid"` // 如果创建时邀请的成员中存在加群数量超过限制的情况，会返回faccid
}

// Create 创建群
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#创建群
// 1.创建高级群，以邀请的方式发送给用户；
// 2.custom 字段是给第三方的扩展字段，第三方可以基于此字段扩展高级群的功能，构建自己需要的群；
// 3.建群成功会返回tid，需要保存，以便于加人与踢人等后续操作；
// 4.每个用户可创建的群数量有限制，限制值由 IM 套餐的群组配置决定，可登录管理后台查看。
func (api *Api) Create(param *CreateParam) (*CreateResp, error) {
	res, err := api.Client.Post(SubDomain+"create.action", param, &CreateResp{})
	return res.(*CreateResp), err
}
