package msgReceive

import (
	"errors"

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

// SetHttpCallbackUrl 设置消息接收地址
// 用户需提供接收微信消息的公网接口URL，并将此url在此接口 配置（PS：简单理解就是腾讯服务器会将消息请求到你们编写的接口服务）
// 公网接口需流畅，微信消息是Http Post Json请求，默认最高6秒内建立连接并发送数据，通讯时长超过6秒，不发送回调消息。
// 配置成功后，会接收一条test的回调。
// 开发者若未配置此接口，消息默认推送至后台系统-消息推送模块。
// 注意：机器人微信自己通过接口发送的消息不会有回调，因为回调是接收消息，发送不属于接收，但是手机微信发送的消息也会有，因为这属于消息同步（同步其他客户端的消息至本客户端，IM原理）。
func (api *Api) SetHttpCallbackUrl(httpURL string) (err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/setHttpCallbackUrl"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"httpUrl": httpURL, // 开发者接口回调地址
		"type":    2,
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

// CancelHttpCallbackUrl  取消消息接收
func (api *Api) CancelHttpCallbackUrl() (err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/cancelHttpCallbackUrl"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, nil, map[string]string{
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

// ParsedFromMsg 解析发来的消息数据
func (api *Api) ParsedFromMsg(jsonStr string) (data *FromMsgResp, err error) {

	err = typeHelper.JsonDecodeWithStruct(jsonStr, &data)
	return
}

// ParsedTextParam 解析发来的消息数据
func (api *Api) ParsedTextParam(jsonStr string) (data *TextParam, err error) {
	err = typeHelper.JsonDecodeWithStruct(jsonStr, &data)
	return
}

// ParsedAudioParam 解析发来的消息数据
func (api *Api) ParsedAudioParam(jsonStr string) (data *AudioParam, err error) {
	err = typeHelper.JsonDecodeWithStruct(jsonStr, &data)
	return
}

// ParsedGroupUpdateParam 解析发来的群聊消息修改数据
func (api *Api) ParsedGroupUpdateParam(jsonStr string) (data *GroupUpdateParam, err error) {
	err = typeHelper.JsonDecodeWithStruct(jsonStr, &data)
	return
}

// ParsedGroupExitParam 解析发来的退出群聊（本人退出）
func (api *Api) ParsedGroupExitParam(jsonStr string) (data *GroupExitParam, err error) {
	err = typeHelper.JsonDecodeWithStruct(jsonStr, &data)
	return
}

// ParsedGroupInviteParam xxx邀请xxx加入群聊
func (api *Api) ParsedGroupInviteParam(jsonStr string) (data *GroupInviteParam, err error) {
	err = typeHelper.JsonDecodeWithStruct(jsonStr, &data)
	return
}

// GetAudioDownloadURL 获取语音文件下载地址
func (api *Api) GetAudioDownloadURL(data *AudioData) (url string, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/getMsgVoice"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":      data.WId, // 开发者接口回调地址
		"msgId":    data.MsgId,
		"length":   data.Length,
		"bufId":    data.BufId,
		"fromUser": data.FromUser,
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &DownloadAudioResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*DownloadAudioResp)
	if resp.Code == string(base.Success) {
		url = resp.Data.URL
	} else {
		err = errors.New(resp.Message)
	}

	return
}
