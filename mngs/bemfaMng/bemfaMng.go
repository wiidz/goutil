package bemfaMng

import (
	"errors"
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/networkStruct"
)

const Domain = "https://apis.bemfa.com"

// NewBemfaMng  获取es管理器
func NewBemfaMng(UID, topicID string) (es *BemfaMng) {
	return &BemfaMng{
		UID:     UID,
		TopicID: topicID,
	}
}

// SwitchOn 打开1路开关
func (mng *BemfaMng) SwitchOn(weMsg string) (data *ReturnBase, err error) {
	return mng.SendMsg("a1", weMsg)
}

// SwitchOff 关闭1路开关
func (mng *BemfaMng) SwitchOff(weMsg string) (data *ReturnBase, err error) {
	return mng.SendMsg("b1", weMsg)
}

// GetSwitchStatus 获取开关的状态
func (mng *BemfaMng) GetSwitchStatus(weMsg string) (isOn bool, err error) {
	data, err := mng.SendMsg("q1", weMsg)
	if err != nil {
		return
	}

	if data.Data == "n1" {
		isOn = true
	}

	return isOn, err
}

// SendMsg 推送消息
// 向主题推送消息，支持POST协议
func (mng *BemfaMng) SendMsg(msg, weMsg string) (data *ReturnBase, err error) {
	var url = Domain + "/va/postJsonMsg"
	var sendMap = map[string]interface{}{
		"uid":   mng.UID,     // 必填，用户私钥，巴法云控制台获取
		"topic": mng.TopicID, // 必填，主题名，可在控制台创建
		"type":  3,           // 必填，主题类型，当type=1时是MQTT协议，3是TCP协议
		"msg":   msg,         // 必填，消息体，要推送的消息，自定义即可，比如on，或off等等
		//"wemsg": weMsg,       // 选填，发送到微信的消息，自定义即可。如果携带此字段，会将消息发送到微信
	}
	if weMsg != "" {
		sendMap["wemsg"] = weMsg
	}

	//res, _, _, err := networkHelper.RequestJsonWithStructTest(networkStruct.Post, url, sendMap, nil, &ReturnBase{})
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, url, sendMap, nil, &ReturnBase{})

	if err != nil {
		return
	}

	data = res.(*ReturnBase)
	err = ReturnCode(data.Code).GetError()
	return
}

// GetMsg 获取消息
// 获取主题消息，支持GET协议
func (mng *BemfaMng) GetMsg(msgAmount int) (data *GetMsgResult, err error) {
	var url = Domain + "/va/getmsg"
	res, _, _, err := networkHelper.RequestJsonWithStruct(networkStruct.Get, url, map[string]interface{}{
		"uid":   mng.UID,     // 必填，用户私钥，巴法云控制台获取
		"topic": mng.TopicID, // 必填，主题名，可在控制台创建
		"type":  3,           // 必填，主题类型，当type=1时是MQTT协议，3是TCP协议
		"num":   msgAmount,   // 选填，获取的历史数据条数，不填默认默认是1，最大5000
	}, map[string]string{}, &GetMsgResult{})
	return res.(*GetMsgResult), err
}

// GetAllTopic 获取全部主题
func (mng *BemfaMng) GetAllTopic() (data *AllTopicResult, err error) {
	var url = Domain + "/va/alltopic"
	res, _, _, err := networkHelper.RequestJsonWithStruct(networkStruct.Get, url, map[string]interface{}{
		"uid":  mng.UID, // 必填，用户私钥，巴法云控制台获取
		"type": 3,       // 必填，主题类型，当type=1时是MQTT协议，3是TCP协议
	}, map[string]string{}, &GetMsgResult{})
	return res.(*AllTopicResult), err
}

// IsOnline 判断设备是否在线
// 获取主题消息，支持GET协议
func (mng *BemfaMng) IsOnline() (isOnline bool, err error) {
	var url = Domain + "/va/online"
	res, _, _, err := networkHelper.RequestJsonWithStruct(networkStruct.Get, url, map[string]interface{}{
		"uid":   mng.UID,     // 必填，用户私钥，巴法云控制台获取
		"topic": mng.TopicID, // 必填，主题名，可在控制台创建
		"type":  3,           // 必填，主题类型，当type=1时是MQTT协议，3是TCP协议
	}, map[string]string{}, &IsOnlineResult{})

	if err != nil {
		return
	}

	isOnline = res.(*IsOnlineResult).Data
	return
}

// SetTimer 设置定时操作
func (mng *BemfaMng) SetTimer(msg string, hour, min, second int) (ok bool, err error) {
	var url = Domain + "/cloud/settime/v1/"
	var timeStr = typeHelper.Int2Str(hour) + ":" + typeHelper.Int2Str(min) + ":" + typeHelper.Int2Str(second)
	res, _, _, err := networkHelper.RequestWithStruct(networkStruct.Post, networkStruct.BodyForm, url, map[string]interface{}{
		"uid":    mng.UID,     // 必填，用户私钥，巴法云控制台获取
		"topic":  mng.TopicID, // 必填，主题名，可在控制台创建
		"type":   3,           // 必填，主题类型，当type=1时是MQTT协议，3是TCP协议
		"msg":    msg,         // 必填，消息体，即定时发送的消息,比如等于on或者off
		"time":   timeStr,     // 必填，时间，格式为 小时:分钟:秒，中间":"为英文格式符号，小时0-23，分钟0-59，秒0-59
		"action": "add",       // 必填，动作，action=add时是添加定时，action=del是删除定时
	}, map[string]string{}, &ReturnStatus{})

	if err != nil {
		return
	}

	data := res.(*ReturnStatus)
	if data.Code == int(AddSuccess) {
		return true, nil
	}

	return false, errors.New(data.Status)

}

// DeleteTimer 删除定时操作
func (mng *BemfaMng) DeleteTimer() (ok bool, err error) {
	var url = Domain + "/cloud/settime/v1/"
	res, _, _, err := networkHelper.RequestWithStruct(networkStruct.Delete, networkStruct.BodyForm, url, map[string]interface{}{
		"uid":   mng.UID,     // 必填，用户私钥，巴法云控制台获取
		"topic": mng.TopicID, // 必填，主题名，可在控制台创建
		"type":  3,           // 必填，主题类型，当type=1时是MQTT协议，3是TCP协议
		//"msg":    msg,         // 必填，消息体，即定时发送的消息,比如等于on或者off
		//"time":   timeStr,     // 必填，时间，格式为 小时:分钟:秒，中间":"为英文格式符号，小时0-23，分钟0-59，秒0-59
		"action": "del", // 必填，动作，action=add时是添加定时，action=del是删除定时
	}, map[string]string{}, &ReturnStatus{})

	if err != nil {
		return
	}

	data := res.(*ReturnStatus)
	if data.Code == int(DeleteSuccess) {
		return true, nil
	}

	return false, errors.New(data.Status)

}

// GetTimer 获取定时操作
func (mng *BemfaMng) GetTimer() (data *ReturnBase, err error) {
	var url = Domain + "/cloud/settime/v1/"
	res, _, _, err := networkHelper.RequestJsonWithStruct(networkStruct.Get, url, map[string]interface{}{
		"uid":   mng.UID,     // 必填，用户私钥，巴法云控制台获取
		"topic": mng.TopicID, // 必填，主题名，可在控制台创建
		"type":  3,           // 必填，设备类型，当type=3时为创客云设备，当等于2时为设备云设备，当等于1时为MQTT设备
	}, map[string]string{}, &ReturnBase{})
	return res.(*ReturnBase), err
}
