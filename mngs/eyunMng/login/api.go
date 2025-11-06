package login

import (
	"errors"
	"time"

	"github.com/wiidz/goutil/helpers/mathHelper"
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/eyunMng/base"
	"github.com/wiidz/goutil/structs/networkStruct"
)

func NewApi(config *base.Config) (api *Api) {
	return &Api{
		Config: config,
	}
}

// LoginEYun 登录E云平台（第一步）
func (api *Api) LoginEYun() (authorization string, err error) {

	//【1】组合URL
	var URL = api.Config.Host + "/member/login"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"account":  api.Config.Account,
		"password": api.Config.Password,
	}, nil, &EYunLoginResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*EYunLoginResp)
	if resp.Code == string(base.Success) {
		if resp.Data.Status == 0 {
			authorization = resp.Data.Authorization
		} else if resp.Data.Status == 1 {
			err = errors.New("登录失败，账户已冻结")
		} else if resp.Data.Status == 2 {
			err = errors.New("登录失败，账户已到期")
		}
	} else {
		err = errors.New(resp.Message)
	}

	return
}

// LocalIPadLogin 获取二维码（第二步-方式1）
// 本方式需要用户下载app/exe获取值并上传到接口,方可取码（虽然繁琐，但是建议使用本接口取码）
// 开发者将本接口返回的二维码让用户去扫码，手机扫码结束后，需要调用第三步接口才会登录成功，且手机顶部显示ipad已登录
// 此方法登录ip显示在苏州
func (api *Api) LocalIPadLogin(wcID, ttuID string) (wID, qrcodeURL string, err error) {

	//【1】组合URL
	var URL = api.Config.Host + "/localIPadLogin"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wcId":  wcID,  // 微信原始id （首次登录平台的号传""，掉线后必须传值，否则会频繁掉线！！！） 第三步会返回此字段，记得入库保存
		"ttuid": ttuID, // 用户需安装app/pc，且上传app/pc中的字段，若是开发者公司有app/pc也可直接集成sdk至app/pc中，可以做到无需用户上传，且无需下载我司提供的软件
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &WechatQrcodeResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*WechatQrcodeResp)
	if resp.Code == string(base.Success) {
		wID = resp.Data.WId
		qrcodeURL = resp.Data.QrCodeUrl
	} else {
		err = errors.New(resp.Message)
	}

	return
}

// IPadLogin 获取二维码（第二步-方式2）
// 本方式需要用户下载app/exe获取值并上传到接口,方可取码（虽然繁琐，但是建议使用本接口取码）
// 开发者将本接口返回的二维码让用户去扫码，手机扫码结束后，需要调用第三步接口才会登录成功，且手机顶部显示ipad已登录
// 此方法登录ip显示在苏州
func (api *Api) IPadLogin(wcID string) (wID, qrcodeURL string, err error) {

	//【1】组合URL
	var URL = api.Config.Host + "/iPadLogin"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wcId":          wcID,                            // 微信原始id （首次登录平台的号传空，掉线重登必须传值，否则会频繁掉线！！！） 第三步会返回此字段，记得入库保存
		"proxy":         8,                               // 测试长效代理 8=杭州，会被覆盖不用管
		"proxyIp":       api.Config.ProxyConfig.Host,     // 自定义长效代理IP+端口
		"proxyUser":     api.Config.ProxyConfig.Username, // 自定义长效代理IP平台账号
		"proxyPassword": api.Config.ProxyConfig.Password, // 自定义长效代理IP平台密码
		//"ttuid": TTUid,
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &WechatQrcodeResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*WechatQrcodeResp)
	if resp.Code == string(base.Success) {
		wID = resp.Data.WId
		qrcodeURL = resp.Data.QrCodeUrl
		api.Config.WechatData.WID = wID
	} else {
		err = errors.New(resp.Message)
	}

	return
}

// WechatLogin 执行微信登录（第三步）
// 此接口为检测耗时接口，最长250S返回请求，用户VX扫码了会返回结果，且扫码成功后手机上会显示ipad登录成功，才可以收发消息及调用其它接口！
// 首次登录平台，24小时内会掉线1次，且72小时内不能发送朋友圈，掉线后必须传wcid调用获取二维码接口再次扫码登录即可实现3月内不掉线哦， 详细规范点击这里(第1大类1小节) PS：若出现登录60S内无故掉线也看这里哦!
func (api *Api) WechatLogin(verifyCode string) (loginData *base.WechatData, err error) {

	//【1】组合URL
	var URL = api.Config.Host + "/getIPadLoginInfo"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":        api.Config.WechatData.WID, // 登录实例标识
		"verifyCode": verifyCode,                // 验证码 默认不传，若扫码结束后，本接口返回提示"请在ipad上输入验证码" ，则再调用1次本接口且传验证码即可（PS：极少情况下会出现此情况,可忽略此字段）
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &WechatLoginResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*WechatLoginResp)
	if resp.Code == string(base.Success) {
		loginData = resp.Data
	} else {
		err = errors.New(resp.Message)
	}

	return
}

// InitAddressList 初始化通讯录列表（第四步）
func (api *Api) InitAddressList() (err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/initAddressList"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId": api.Config.WechatData.WID, // 登录实例标识
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &base.BaseResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*base.BaseResp)
	if resp.Code == string(base.Success) {

	} else {
		err = errors.New(resp.Message)
	}

	return
}

// GetAddressList 获取通讯录列表
// 获取通讯录列表之前，必须调用初始化通讯录列表接口。
// 此接口不会返回好友/群的详细信息，如需获取详细信息，请调用获取联系人详情接口
// 本接口的返回群聊的是保存到通讯录的群聊详细规范点击这里(第5大类3小节)
func (api *Api) GetAddressList() (data *AddressListData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/getAddressList"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId": api.Config.WechatData.WID, // 登录实例标识
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &AddressListResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*AddressListResp)
	if resp.Code == string(base.Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}

// GetContact 获取联系人信息
func (api *Api) GetContact(wcIDs []string) (err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/getContact"

	var amount = len(wcIDs)
	if amount == 0 {
		return errors.New("操作失败，wcID数据为空")
	} else if amount > 20 {
		return errors.New("操作失败，wcID数据不能超过20个")
	}

	wcIDStr := typeHelper.Implode(wcIDs, ",")

	//【2】请求数据
	// 随机间隔300ms-1500ms，频繁调用容易导致掉线
	random := time.Duration(mathHelper.GetRandomInt(300, 1500)) * time.Microsecond
	time.Sleep(random)

	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":  api.Config.WechatData.WID, // 登录实例标识
		"wcId": wcIDStr,                   // 好友微信id/群id,多个好友/群 以","分隔每次最多支持20个微信/群号,记得本接口随机间隔300ms-1500ms，频繁调用容易导致掉线
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &AddressListResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*AddressListResp)
	if resp.Code == string(base.Success) {

	} else {
		err = errors.New(resp.Message)
	}

	return
}

// Update 登录E云平台（第一步）
func (api *Api) Update() {
	api.Config.WechatData.WID = "updated"
}
