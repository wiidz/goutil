package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type UpdateParam struct {
	TeamID             string `json:"tid"`                       // 网易云信服务器产生，群唯一标识，创建群时会返回
	TeamName           string `json:"tname,omitempty"`           // 群名称，最大长度64字符
	Owner              string `json:"owner"`                     // 群主用户帐号，最大长度32字符
	Announcement       string `json:"announcement,omitempty"`    // 群公告，最大长度1024字符
	Intro              string `json:"intro,omitempty"`           // 群描述，最大长度512字符
	JoinMode           int    `json:"joinmode,omitempty"`        // 群建好后，sdk操作时，0不用验证，1需要验证,2不允许任何人加入。其它返回414
	Custom             string `json:"custom,omitempty"`          // 自定义高级群扩展属性，第三方可以跟据此属性自定义扩展自己的群属性。（建议为json）,最大长度1024字符
	Icon               string `json:"icon,omitempty"`            // 群头像，最大长度1024字符
	BeInviteMode       int    `json:"beinvitemode,omitempty"`    // 被邀请人同意方式，0-需要同意(默认),1-不需要同意。其它返回414
	InviteMode         int    `json:"invitemode,omitempty"`      // 谁可以邀请他人入群，0-管理员(默认),1-所有人。其它返回414
	UpdateTeamInfoMode int    `json:"uptinfomode,omitempty"`     // 谁可以修改群资料，0-管理员(默认),1-所有人。其它返回414
	UpdateCustomMode   int    `json:"upcustommode,omitempty"`    // 谁可以更新群自定义属性，0-管理员(默认),1-所有人。其它返回414
	TeamMemberLimit    int    `json:"teamMemberLimit,omitempty"` // 该群最大人数(包含群主)，范围：2至应用定义的最大群人数(默认:200)。其它返回414
	Attach             string `json:"attach,omitempty"`          // 	自定义扩展字段，最大长度512
	Bid                string `json:"bid,omitempty"`             // 反垃圾业务ID，JSON字符串，{"textbid":"","picbid":""}，若不填则使用原来的反垃圾配置
}

// Update 编辑群资料
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#编辑群资料
// 高级群基本信息修改
func (api *Api) Update(param *UpdateParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post(SubDomain+"update.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
