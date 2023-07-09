package bemfaMng

import (
	"errors"
	"strconv"
)

// BemfaMng 巴法云
type BemfaMng struct {
	UID     string `json:"uid"`
	TopicID string `json:"topic_id"`
}

type ReturnBase struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ReturnStatus struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

type GetMsgResult struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    []HistoryData `json:"data"`
}

type HistoryData struct {
	Msg  string `json:"msg"`  // 获取的主题消息
	Time string `json:"time"` // 消息发送的时间，时区UTC/GMT+08:00
	Unix int    `json:"unix"` // 消息发送的时间戳
}

// AllTopicResult 获取全部主题
type AllTopicResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    []TopicData `json:"data"`
}

type TopicData struct {
	Topic  string `json:"topic"`  // 主题值
	Msg    string `json:"msg"`    // 消息体
	Name   string `json:"name"`   // 主题名字
	Online bool   `json:"online"` // 是否在线
	Tid    string `json:"tid"`    // 设备类型
	Sid    string `json:"sid"`    // 如果是分享设备，此字段是分享者密钥
	Time   string `json:"time"`   // 消息发送的时间，时区UTC/GMT+08:00
	Unix   int    `json:"unix"`   // 消息发送的时间戳
}

type TID string

type IsOnlineResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    bool   `json:"data"`
}

const Outlet TID = "outlet"             // 插座
const Light TID = "light"               // 灯
const Fan TID = "fan"                   // 风扇
const Sensor TID = "sensor"             // 传感器
const AirCondition TID = "aircondition" // 空调
const Switch TID = "switch"             // 开关
const Curtain TID = "curtain"           // 窗帘

type ReturnCode int

const Success ReturnCode = 0 // 成功

const ParamError ReturnCode = 10002 // 请求参数有误

const UnknownError ReturnCode = 40000 // 未知错误
const KeyError ReturnCode = 40004     // 私钥或主题错误

func (code ReturnCode) GetError() (err error) {

	var errorStr = ""
	switch code {
	case Success:
	case AddSuccess:
	case DeleteSuccess:
	case Getok:

	case ParamError:
		errorStr = "请求参数有误"
	case UnknownError:
		errorStr = "未知错误"
	case KeyError:
		errorStr = "私钥或主题错误"

	case NoUIDError:
		errorStr = "缺少uid字段"
	case WrongUIDError:
		errorStr = "uid值为空或不正确"
	case NoTopicError:
		errorStr = "缺少topic字段"
	case WrongTopicError:
		errorStr = "topic值为空或不正确"
	case NoTypeError:
		errorStr = "缺少type字段"
	case WrongTypeError:
		errorStr = "type值为空或不正确"
	case NoTimeError:
		errorStr = "缺少time字段"
	case WrongTimeError:
		errorStr = "time值为空或不正确"

	case NoMsgError:
		errorStr = "缺少msg字段"
	case WrongMsgError:
		errorStr = "time值为空或不正确"
	case NoActionError:
		errorStr = "缺少action字段"
	case WrongActionError:
		errorStr = "action值为空或不正确"

	case Existed:
		errorStr = "定时已存在"

	case NoUIDErrorOld:
		errorStr = "缺少uid字段"
	case WrongUIDErrorOld:
		errorStr = "uid值为空或不正确"
	case NoTopicErrorOld:
		errorStr = "缺少topic字段"
	case WrongTopicErrorOld:
		errorStr = "topic值为空或不正确"

	default:
		errorStr = strconv.Itoa(int(code))
	}

	if errorStr == "" {
		return nil
	}
	return errors.New(errorStr)
}

// 定时相关
const NoUIDError ReturnCode = 4003001       // 缺少uid字段
const WrongUIDError ReturnCode = 4003002    // uid值为空或不正确
const NoTopicError ReturnCode = 4003003     // 缺少topic字段
const WrongTopicError ReturnCode = 4003004  // topic值为空或不正确
const NoTypeError ReturnCode = 4003005      // 缺少type字段
const WrongTypeError ReturnCode = 4003006   // type值为空或不正确
const NoTimeError ReturnCode = 4003007      // 缺少time字段
const WrongTimeError ReturnCode = 4003008   // time值为空或不正确
const NoMsgError ReturnCode = 4003009       // 缺少msg字段
const WrongMsgError ReturnCode = 4003010    // msg值为空或不正确
const NoActionError ReturnCode = 4003011    // 缺少action字段
const WrongActionError ReturnCode = 4003012 // action值为空或不正确
const AddSuccess ReturnCode = 4003013       // 添加成功
const DeleteSuccess ReturnCode = 4003014    // 删除
const Existed ReturnCode = 4003015          // 定时已存在

// 似乎是老版本的数据
const Getok ReturnCode = 40020              // 获取设备状态的时候 getok
const WrongUIDErrorOld ReturnCode = 40012   // uid值为空或不正确 no uid
const NoUIDErrorOld ReturnCode = 41020      // 缺少uid字段 no uid
const NoTopicErrorOld ReturnCode = 41021    // 缺少topic字段 no uid
const WrongTopicErrorOld ReturnCode = 41022 // topic值为空或不正确 no uid

type GetSwitchStatusResult struct {
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Time   string `json:"time"`
	Status string `json:"status"`
}
