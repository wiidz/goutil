package yunxinMng

import (
	"errors"
	"github.com/wiidz/goutil/helpers/cryptorHelper"
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/helpers/strHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"github.com/wiidz/goutil/structs/networkStruct"
	"time"
)

type YunxinMng struct {
	Config *configStruct.YunxinConfig
}

type CommonResp struct {
	Code int         `json:"code"`
	Desc string      `json:"desc"`
	Info interface{} `json:"info"`
}

func (resp *CommonResp) GetCode() int {
	return resp.Code
}
func (resp *CommonResp) GetDesc() string {
	return resp.Desc
}
func (resp *CommonResp) GetInfo() interface{} {
	return resp.Info
}

type RespInterface interface {
	GetCode() int
	GetDesc() string
	GetInfo() interface{}
}

func NewYunxinMng(config *configStruct.YunxinConfig) *YunxinMng {
	return &YunxinMng{
		Config: config,
	}
}

// getCheckSum 获取令牌
func (mng *YunxinMng) getCheckSum(nonce, timestampStr string) string {
	return cryptorHelper.SHA1Encrypt(mng.Config.AppSecret + nonce + timestampStr)
}

// Post 发送post请求
func (mng *YunxinMng) Post(url string, params interface{}, iStruct RespInterface) (data interface{}, err error) {

	//【1】构建参数
	paramStr, _ := typeHelper.JsonEncode(params)
	paramMap := typeHelper.JsonDecodeMap(paramStr)

	//【2】获取令牌
	nowTimestamp := typeHelper.Int64ToStr(time.Now().Unix())
	nonce := strHelper.GetRandomString(12)
	checkSum := mng.getCheckSum(nonce, nowTimestamp)

	//【3】发送请求
	var statusCode int
	data, _, statusCode, err = networkHelper.RequestJsonWithStructTest(networkStruct.Post, url, paramMap, map[string]string{
		"AppKey":    mng.Config.AppKey, // 开发者平台分配的 appkey（具体获取方式请参考登录鉴权）
		"AppSecret": mng.Config.AppSecret,
		"Nonce":     nonce,        // 随机数（最大长度128个字符）
		"CurTime":   nowTimestamp, // 当前UTC时间戳，从1970年1月1日0点0 分0 秒开始到现在的秒数(String)
		"CheckSum":  checkSum,     // SHA1(AppSecret + Nonce + CurTime)，三个参数拼接的字符串，
		// 进行SHA1哈希计算，转化成16进制字符(String，小写),出于安全性考虑，每个checkSum的有效期为5分钟(用CurTime计算)，建议每次请求都生成新的checkSum，同时请确认发起请求的服务器是与标准时间同步的，比如有NTP服务。
		// CheckSum检验失败时会返回414错误码，具体参看code状态表。
		// 本文档中提供的所有接口均面向开发者服务器端调用，用于计算CheckSum的AppSecret开发者应妥善保管,可在应用的服务器端存储和使用，但不应存储或传递到客户端，也不应在网页等前端代码中嵌入。
	}, iStruct)

	//【4】判断结果
	if err != nil {
		return
	} else if statusCode != 200 {
		err = errors.New("请求失败")
	} else if iStruct.GetCode() != 200 {
		err = errors.New(iStruct.GetDesc())
	}

	//【5】返回
	return
}
