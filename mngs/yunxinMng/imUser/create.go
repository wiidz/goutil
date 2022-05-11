package imUser

import (
	"github.com/wiidz/goutil/mngs/yunxinMng/imClient"
)

type CreateParam struct {
	Accid  string `json:"accid" validate:"required"` // 网易云信IM账号，最大长度32字符，必须保证一个 APP内唯一（只允许字母、数字、半角下划线_、 @、半角点以及半角-组成，不区分大小写， 会统一小写处理，请注意以此接口返回结果中的accid为准）。
	Name   string `json:"name,omitempty"`            // 网易云信IM账号昵称，最大长度64字符。
	Props  string `json:"props,omitempty"`           // json属性，开发者可选填，最大长度1024字符。该参数已不建议使用。
	Icon   string `json:"icon,omitempty"`            // 网易云信IM账号头像URL，开发者可选填，最大长度1024
	Token  string `json:"token,omitempty"`           // 网易云信IM账号可以指定登录IM token值，最大长度128字符，并更新，如果未指定，会自动生成token，并在创建成功后返回
	Sign   string `json:"sign,omitempty"`            // 用户签名，最大长度256字符
	Email  string `json:"email,omitempty"`           //用户email，最大长度64字符
	Birth  string `json:"birth,omitempty"`           // 用户生日，最大长度16字符
	Mobile string `json:"mobile,omitempty"`          // 用户mobile，最大长度32字符，非中国大陆手机号码需要填写国家代码(如美国：+1-xxxxxxxxxx)或地区代码(如香港：+852-xxxxxxxx)
	Gender int    `json:"gender,omitempty"`          //用户性别，0表示未知，1表示男，2女表示女，其它会报参数错误
	Ex     string `json:"ex,omitempty"`              // 用户名片扩展字段，最大长度1024字符，用户可自行扩展，建议封装成JSON字符串
}
type CreateResp struct {
	*imClient.CommonResp
	Info struct {
		Token string `json:"token"`
		Accid string `json:"accid"`
		Name  string `json:"name"`
	} `json:"info"`
}

// UserCreate 创建网易云信IM账号
// doc: https://doc.yunxin.163.com/docs/TM5MzM5Njk/Dc2NTM1NzI?platformId=60353
// 1.第三方帐号导入到网易云信平台。注册成功后务必在自身的应用服务器上维护好accid与token。
// 2.注意 IM accid，name 长度以及考虑管理 IM token。
// 3.云信应用内的accid若涉及字母，请一律为小写，并确保服务端与所有客户端均保持小写。
func (api *Api) UserCreate(param *CreateParam) (*CreateResp, error) {
	res, err := api.Client.Post("/user/create.action", param, &CreateResp{})
	return res.(*CreateResp), err
}
