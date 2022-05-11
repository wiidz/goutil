package yunxinMng

type UserCreateParam struct {
	// 网易云信IM账号，最大长度32字符，必须保证一个
	// APP内唯一（只允许字母、数字、半角下划线_、
	// @、半角点以及半角-组成，不区分大小写，
	// 会统一小写处理，请注意以此接口返回结果中的accid为准）。
	// 是否必须: true
	Accid string `json:"accid"`

	// 网易云信IM账号昵称，最大长度64字符。
	// 是否必须: false
	Name string `json:"name,omitempty"`

	// json属性，开发者可选填，最大长度1024字符。该参数已不建议使用。
	// 是否必须: false
	Props string `json:"props,omitempty"`

	// 网易云信IM账号头像URL，开发者可选填，最大长度1024
	// 是否必须: false
	Icon string `json:"icon,omitempty"`

	// 网易云信IM账号可以指定登录IM token值，最大长度128字符，
	// 并更新，如果未指定，会自动生成token，并在
	// 创建成功后返回
	// 是否必须: false
	Token string `json:"token,omitempty"`

	// 用户签名，最大长度256字符
	// 是否必须: false
	Sign string `json:"sign,omitempty"`

	// 用户email，最大长度64字符
	// 是否必须: false
	Email string `json:"email,omitempty"`

	// 用户生日，最大长度16字符
	// 是否必须: false
	Birth string `json:"birth,omitempty"`

	// 用户mobile，最大长度32字符，非中国大陆手机号码需要填写国家代码(如美国：+1-xxxxxxxxxx)或地区代码(如香港：+852-xxxxxxxx)
	// 是否必须: false
	Mobile string `json:"mobile,omitempty"`

	// 用户性别，0表示未知，1表示男，2女表示女，其它会报参数错误
	// 是否必须: false
	Gender int `json:"gender,omitempty"`

	// 用户名片扩展字段，最大长度1024字符，用户可自行扩展，建议封装成JSON字符串
	// 是否必须: false
	Ex string `json:"ex,omitempty"`
}

type UserCreateResp struct {
	*CommonResp
	Info struct {
		Token string `json:"token"`
		Accid string `json:"accid"`
		Name  string `json:"name"`
	} `json:"info"`
}

// UserCreate 创建网易云信IM账号
// doc: https://dev.yunxin.163.com/docs/product/IM%E5%8D%B3%E6%97%B6%E9%80%9A%E8%AE%AF/%E6%9C%8D%E5%8A%A1%E7%AB%AFAPI%E6%96%87%E6%A1%A3/%E7%BD%91%E6%98%93%E4%BA%91%E9%80%9A%E4%BF%A1ID
func (mng *YunxinMng) UserCreate(param *UserCreateParam) (*UserCreateResp, error) {
	res, err := mng.Post(UserCreateURL, param, &UserCreateResp{})
	return res.(*UserCreateResp), err
}
