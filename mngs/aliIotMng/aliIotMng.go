package aliIotMng

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	iot "github.com/alibabacloud-go/iot-20180120/v2/client"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
)

// AliIotMng 这是一个实例的管理器
type AliIotMng struct {
	IotInstanceID string // 实例ID
	Client        *iot.Client
	ApiConfig     *configStruct.AliIotConfig
}

// NewAliIotMng 返回一个物联网实例管理器
func NewAliIotMng(config *configStruct.AliIotConfig, iotInstanceID string) (mng *AliIotMng, err error) {
	mng = &AliIotMng{
		IotInstanceID: iotInstanceID,
		ApiConfig:     config,
	}
	mng.Client, err = mng.createClient(&config.AccessKeyID, &config.AccessKeySecret, &config.EndPoint)
	return
}

// CreateClient 创建客户端
func (mng *AliIotMng) createClient(accessKeyId, accessKeySecret, endPoint *string) (result *iot.Client, err error) {

	config := &openapi.Config{
		AccessKeyId:     accessKeyId,     // 您的AccessKey ID
		AccessKeySecret: accessKeySecret, // 您的AccessKey Secret
		Endpoint:        endPoint,        // 访问的域名
	}

	return iot.NewClient(config)
}

// GetAttributes 获取设备属性
func (mng *AliIotMng) GetAttributes(productKey, deviceName string) (res *iot.QueryDevicePropertyStatusResponse, err error) {

	queryDevicePropertyStatusRequest := &iot.QueryDevicePropertyStatusRequest{
		IotInstanceId: &mng.IotInstanceID,
		ProductKey:    &productKey,
		DeviceName:    &deviceName,
	}

	res, err = mng.Client.QueryDevicePropertyStatus(queryDevicePropertyStatusRequest)
	return
}

// SetAttributes 设置设备属性
func (mng *AliIotMng) SetAttributes(productKey, deviceName string, items map[string]interface{}) (res *iot.SetDevicePropertyResponse, err error) {

	temp, _ := typeHelper.JsonEncode(items)
	setDevicePropertyRequest := &iot.SetDevicePropertyRequest{
		IotInstanceId: &mng.IotInstanceID,
		ProductKey:    &productKey,
		DeviceName:    &deviceName,
		Items:         &temp,
	}

	res, err = mng.Client.SetDeviceProperty(setDevicePropertyRequest)
	return
}
