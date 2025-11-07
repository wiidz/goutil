package msgSend

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

// SendText 发送文本消息
func (api *Api) SendText(wcID, content string) (data *SendMsgData, err error) {

	if content == "" {
		err = errors.New("发送失败，消息体为空")
		return
	}

	//【1】组合URL
	var URL = api.Config.Host + "/sendText"

	if wcID == "" {
		err = errors.New("发送失败，对象主体不能为空")
	} else if content == "" {
		err = errors.New("发送失败，内容不能为空")
	}

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStruct(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":     api.Config.WechatData.WID,
		"wcId":    wcID,
		"content": content,
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &SendMsgResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*SendMsgResp)
	if resp.Code == string(Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}

// SendTextAt 发送文本消息带@
func (api *Api) SendTextAt(wcID, content string, atIDs []string, atAll bool) (data *SendMsgData, err error) {

	if content == "" {
		err = errors.New("发送失败，消息体为空")
		return
	}

	//【1】组合URL
	var URL = api.Config.Host + "/sendText"

	//【2】组合atIDs
	var atIDStr string
	if atAll {
		atIDStr = "notify@all"
	} else {
		atIDStr = typeHelper.ImplodeStr(atIDs, ",")
	}

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStruct(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":     api.Config.WechatData.WID,
		"wcId":    wcID,
		"content": content,
		"at":      atIDStr,
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &SendMsgResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*SendMsgResp)
	if resp.Code == string(Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}

// SendImage 发送图片消息
func (api *Api) SendImage(wcID, url string) (data *SendMsgData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/sendImage"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStruct(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":     api.Config.WechatData.WID,
		"wcId":    wcID,
		"content": url,
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &SendMsgResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*SendMsgResp)
	if resp.Code == string(Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}

// SendVideo 发送视频消息
func (api *Api) SendVideo(wcID, url string) (data *SendMsgData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/sendVideo"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStruct(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":       api.Config.WechatData.WID, // 登录实例标识
		"wcId":      wcID,                      // 接收人微信id/群id
		"path":      url,                       // 视频url链接
		"thumbPath": "",                        // 视频封面url链接，可自定义（也可自己服务器获取视频首帧）
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &SendMsgResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*SendMsgResp)
	if resp.Code == string(Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}

// SendURL 发送链接消息
func (api *Api) SendURL(wcID, url, title, desc, thumbURL string) (data *SendMsgData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/sendVideo"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStruct(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":         api.Config.WechatData.WID, // 登录实例标识
		"wcId":        wcID,                      // 接收人微信id/群id
		"url":         url,                       // 链接
		"title":       title,                     // 标题
		"description": desc,                      // 描述
		"thumbUrl":    thumbURL,                  // 图标url
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &SendMsgResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*SendMsgResp)
	if resp.Code == string(Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}

// SendNameCard 发送名片消息
func (api *Api) SendNameCard(wcID, nameCardID string) (data *SendMsgData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/sendNameCard"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStruct(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":        api.Config.WechatData.WID, // 登录实例标识
		"wcId":       wcID,                      // 接收人微信id/群id
		"nameCardId": nameCardID,                // 要发送的名片微信id
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &SendMsgResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*SendMsgResp)
	if resp.Code == string(Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}

// SendEmoji 发送发送emoji表情
func (api *Api) SendEmoji(wcID, imageMd5, imgSize string) (data *SendMsgData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/sendEmoji"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStruct(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":      api.Config.WechatData.WID, // 登录实例标识
		"wcId":     wcID,                      // 接收人微信id/群id
		"imageMd5": imageMd5,                  // 取回调中xml中md5字段值
		"imgSize":  imgSize,                   // 取回调中xml中len字段值
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &SendMsgResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*SendMsgResp)
	if resp.Code == string(Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}

// SendApplet 发送APP类消息
func (api *Api) SendApplet(wcID, content string) (data *SendMsgData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/sendApplet"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStruct(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":     api.Config.WechatData.WID, // 登录实例标识
		"wcId":    wcID,                      // 接收人微信id/群id
		"content": content,                   // 消息xml回调内容, (此回调的XML需要去掉部分，截取appmsg开头的，具体请看请求参数示例）
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &SendMsgResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*SendMsgResp)
	if resp.Code == string(Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}

// SendApplets 发送小程序
func (api *Api) SendApplets(wcID, param *SendMiniParam) (data *SendMsgData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/sendApplet"

	//【2】请求数据

	res, _, _, err := networkHelper.RequestWithStruct(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":         api.Config.WechatData.WID, // 登录实例标识
		"wcId":        wcID,                      // 接收人微信id/群id
		"displayName": param.DisplayName,
		"iconUrl":     param.IconUrl,
		"appId":       param.AppId,
		"pagePath":    param.PagePath,
		"thumbUrl":    param.ThumbUrl,
		"title":       param.Title,
		"userName":    param.UserName,
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &SendMsgResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*SendMsgResp)
	if resp.Code == string(Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}

// RevokeMsg 撤回消息
func (api *Api) RevokeMsg(wcID, msgID, newMsgID, createTime string) (data *SendMsgData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/sendApplet"

	//【2】请求数据

	res, _, _, err := networkHelper.RequestWithStruct(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":        api.Config.WechatData.WID, // 登录实例标识
		"wcId":       wcID,                      // 接收人微信id/群id
		"msgID":      msgID,                     // 消息msgId(发送类接口返回的msgId)
		"newMsgId":   newMsgID,                  // 消息newMsgId(发送类接口返回的msgId)
		"createTime": createTime,                // 发送时间（选填）
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &SendMsgResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*SendMsgResp)
	if resp.Code == string(Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}
