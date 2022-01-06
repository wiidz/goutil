package aliIotMng

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	iot "github.com/alibabacloud-go/iot-20180120/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
)

type AliIotMng struct {
	Client *iot.Client
	Config *configStruct.AliIotConfig
}

// NewAliIotMng 返回一个api实例
func NewAliIotMng(config *configStruct.AliIotConfig) (mng *AliIotMng, err error) {
	mng = &AliIotMng{
		Config: config,
	}
	mng.Client, err = mng.createClient(&config.AccessKeyID, &config.AccessKeySecret, &config.EndPoint)
	return
}

// CreateClient 创建客户端
func (mng *AliIotMng)  createClient(accessKeyId, accessKeySecret, endPoint *string) (result *iot.Client, err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,     // 您的AccessKey ID
		AccessKeySecret: accessKeySecret, // 您的AccessKey Secret
		Endpoint:        endPoint,        // 访问的域名
	}
	result = &iot.Client{}
	result, err = iot.NewClient(config)
	return result, err
}

// GetAttributes 获取设备属性
func (mng *AliIotMng) GetAttributes(iotInstanceID,productKey, deviceName string) (res *iot.QueryDevicePropertyStatusResponse,err error){
	queryDevicePropertyStatusRequest := &iot.QueryDevicePropertyStatusRequest{
		IotInstanceId: tea.String(iotInstanceID),
		ProductKey: tea.String(productKey),
		DeviceName: tea.String(deviceName),
	}
	// 复制代码运行请自行打印 API 的返回值
	res, err = mng.Client.QueryDevicePropertyStatus(queryDevicePropertyStatusRequest)
	return
}

// SetAttributes 设置设备属性
func (mng *AliIotMng) SetAttributes(iotInstanceID,productKey, deviceName string, items map[string]interface{}) (res *iot.SetDevicePropertyResponse,err error){

	temp,_  := typeHelper.JsonEncode(items)
	setDevicePropertyRequest := &iot.SetDevicePropertyRequest{
		IotInstanceId: tea.String(iotInstanceID),
		ProductKey:    tea.String(productKey),
		DeviceName:    tea.String(deviceName),
		IotId:         nil,
		Items:          tea.String(temp),
	}

	// 复制代码运行请自行打印 API 的返回值
	res, err = mng.Client.SetDeviceProperty(setDevicePropertyRequest)
	return
}
