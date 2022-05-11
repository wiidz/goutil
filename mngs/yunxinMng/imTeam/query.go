package imTeam

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type QueryParam struct {
	TeamIDs       string `json:"tids"  validate:"required"` // 群id列表，如["3083","3084"]
	Ope           int    `json:"ope"  validate:"required"`  // 1表示带上群成员列表，0表示不带群成员列表，只返回群信息
	IgnoreInvalid bool   `json:"ignoreInvalid"`             // 是否忽略无效的tid，默认为false。设置为true时将忽略无效tid，并在响应结果中返回无效的tid
	Size          int    `json:"size" validate:"required"`  // 群内当前人数
}

type QueryResp struct {
	*imClient.CommonResp
	TeamInfos      []*TeamInfo `json:"tinfos"`
	InvalidTeamIDs []int       `json:"invalidTids"` //参数ignoreInvalid=true时才有该字段
}

type TeamInfo struct {
	TeamName     string   `json:"tname"`
	Announcement string   `json:"announcement"`
	Owner        string   `json:"owner"`
	MaxUsers     int      `json:"maxusers"`
	JoinMode     int      `json:"joinmode"`
	Tid          int      `json:"tid"`
	Intro        string   `json:"intro"`
	Size         int      `json:"size"`
	Custom       string   `json:"custom"`
	ClientCustom string   `json:"clientCustom"`
	Mute         bool     `json:"mute"`
	CreateTime   int64    `json:"createtime"`
	UpdateTime   int64    `json:"updatetime"`
	Admins       []string `json:"admins"`  // 查询带群成员的群列表信息
	Members      []string `json:"members"` // members字段中的元素包含管理员，但不包含创建者
}

// Query 群信息与成员列表查询
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/jc2NDgzMTg?platformId=60353#群信息与成员列表查询
// 1.高级群信息与成员列表查询，一次最多查询30个群相关的信息，跟据ope参数来控制是否带上群成员列表；
// 2.查询群成员会稍微慢一些，所以如果不需要群成员列表可以只查群信息；
// 3.此接口受频率控制，某个应用一分钟最多查询30次，超过会返回416，并且被屏蔽一段时间；
// 4.群成员的群列表信息中增加管理员成员admins的返回。
func (api *Api) Query(param *QueryParam) (*QueryResp, error) {
	res, err := api.Client.Post(SubDomain+"query.action", param, &QueryResp{})
	return res.(*QueryResp), err
}
