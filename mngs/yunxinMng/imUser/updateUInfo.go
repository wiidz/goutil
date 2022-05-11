package imUser

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type UpdateUInfoParam struct {
	Accid  string `json:"accid" validate:"required"` // 用户帐号，最大长度32字符，必须保证一个APP内唯一
	Name   string `json:"name,omitempty"`            // 用户昵称，最大长度64字符，可设置为空字符串
	Icon   string `json:"icon,omitempty"`            // 用户头像，最大长度1024字节，可设置为空字符串
	Sign   string `json:"sign,omitempty"`            // 用户签名，最大长度256字符，可设置为空字符串
	Email  string `json:"email,omitempty"`           // 用户email，最大长度64字符，可设置为空字符串
	Birth  string `json:"birth,omitempty"`           // 用户生日，最大长度16字符，可设置为空字符串
	Mobile string `json:"mobile,omitempty"`          // 用户mobile，最大长度32字符，非中国大陆手机号码需要填写国家代码(如美国：+1-xxxxxxxxxx)或地区代码(如香港：+852-xxxxxxxx)，可设置为空字符串
	Gender int    `json:"gender,omitempty"`          // 用户性别，0表示未知，1表示男，2女表示女，其它会报参数错误
	Ex     string `json:"ex,omitempty"`              // 用户名片扩展字段，最大长度1024字符，用户可自行扩展，建议封装成JSON字符串，也可以设置为空字符串
}

// UpdateUInfo 更新用户名片
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/zI0NzYyMDQ?platformId=60353
// 更新用户名片。用户名片中包含的用户信息，在群组、聊天室等场景下，会暴露给群组、聊天室内的其他用户。
// 这些字段里mobile，email，birth，gender等字段属于非必填、可能涉及隐私的信息，如果您的业务下，
// 这些信息为敏感信息，建议在通过扩展字段ex填写相关资料并事先加密。
func (api *Api) UpdateUInfo(param *UpdateUInfoParam) (*imClient.CommonResp, error) {
	res, err := api.Client.Post("/user/updateUinfo.action", param, &imClient.CommonResp{})
	return res.(*imClient.CommonResp), err
}
