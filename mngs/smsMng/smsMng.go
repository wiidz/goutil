package smsMng

import (
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
)

const Domain = "dysmsapi.aliyuncs.com"

type SmsMng struct {
	Client *dysmsapi.Client
}

func NewSmsMng(params *configStruct.AliSmsConfig) (smsM *SmsMng, err error) {

	log.Println("AccessKeyID", params.AccessKeyID)
	log.Println("AccessKeySecret", params.AccessKeySecret)
	config := &openapi.Config{
		AccessKeyId:     &params.AccessKeyID,     // 您的AccessKey ID
		AccessKeySecret: &params.AccessKeySecret, // 您的AccessKey Secret
	}

	// 访问的域名
	var client *dysmsapi.Client
	config.Endpoint = tea.String(Domain)
	client, err = dysmsapi.NewClient(config)

	smsM = &SmsMng{
		Client: client,
	}
	return
}

// SendOne 发送一条短信
func (mng *SmsMng) SendOne(signName, templateCode, phoneNumber, paramsJson string) ( err error) {

	sendSmsRequest := &dysmsapi.SendSmsRequest{
		SignName:      &signName,
		TemplateCode:  &templateCode,
		PhoneNumbers:  &phoneNumber,
		TemplateParam: &paramsJson,
	}

	// 复制代码运行请自行打印 API 的返回值
	var res *dysmsapi.SendSmsResponse
	res,err = mng.Client.SendSms(sendSmsRequest)

	if err != nil {
		return
	}
	if *(res.Body.Code) != "ok" {
		err = errors.New(*(res.Body.Message))
	}

	return
}

// SendMany 发送多条短信
func (mng *SmsMng) SendMany(signName, templateCode string, phoneNumbers []string, paramsJson string) (res *dysmsapi.SendBatchSmsResponse, err error) {

	numbers, _ := typeHelper.JsonEncode(phoneNumbers)
	sendSmsRequest := &dysmsapi.SendBatchSmsRequest{
		SignNameJson:      &signName,
		TemplateCode:      &templateCode,
		PhoneNumberJson:   &numbers,
		TemplateParamJson: &paramsJson,
	}

	// 复制代码运行请自行打印 API 的返回值
	return mng.Client.SendBatchSms(sendSmsRequest)
}
