package imUser

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type GetUInfoParam struct {
	Accids string `validate:"required" json:"accids"` // 用户帐号（例如：JSONArray对应的accid串，如：\["zhangsan"\]，如果解析出错，会报414）（一次查询最多为200）
}

type GetUInfoResp struct {
	*imClient.CommonResp
	UInfos []*UInfos `json:"uinfos"`
}

type UInfos struct {
	Accid  string `json:"accid"`            // 用户帐号，最大长度32字符，必须保证一个APP内唯一
	Name   string `json:"name,omitempty"`   // 用户昵称，最大长度64字符，可设置为空字符串
	Icon   string `json:"icon,omitempty"`   // 用户头像，最大长度1024字节，可设置为空字符串
	Sign   string `json:"sign,omitempty"`   // 用户签名，最大长度256字符，可设置为空字符串
	Email  string `json:"email,omitempty"`  // 用户email，最大长度64字符，可设置为空字符串
	Birth  string `json:"birth,omitempty"`  // 用户生日，最大长度16字符，可设置为空字符串
	Mobile string `json:"mobile,omitempty"` // 用户mobile，最大长度32字符，非中国大陆手机号码需要填写国家代码(如美国：+1-xxxxxxxxxx)或地区代码(如香港：+852-xxxxxxxx)，可设置为空字符串
	Gender int    `json:"gender,omitempty"` // 用户性别，0表示未知，1表示男，2女表示女，其它会报参数错误
	Ex     string `json:"ex,omitempty"`     // 用户名片扩展字段，最大长度1024字符，用户可自行扩展，建议封装成JSON字符串，也可以设置为空字符串
	Valid  bool   `json:"valid"`            // 账号是否有效
	Mute   bool   `json:"mute"`             // 账号是否被禁言
}

// GetUInfo 获取用户名片
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/zI0NzYyMDQ?platformId=60353
// 获取用户名片，可批量

func (api *Api) GetUInfo(param *GetUInfoParam) (*GetUInfoResp, error) {
	res, err := api.Client.Post("/user/updateUinfo.action", param, &GetUInfoResp{})
	return res.(*GetUInfoResp), err
}
