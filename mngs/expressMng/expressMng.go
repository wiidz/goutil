package expressMng

import (
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"github.com/wiidz/goutil/structs/networkStruct"
)

const DriveRouteURL = "http://wuliu.market.alicloudapi.com/kdi"

// ExpressMng 全国快递物流查询-快递查询接口
// https://market.aliyun.com/products/57126001/cmapi021863.html?spm=5176.2020520132.101.1.2b9f7218CZkxvN#sku=yuncode1586300000
// 1元=210次
type ExpressMng struct {
	Config *configStruct.AliApiConfig
}

// NewExpressMng : 返回快递管理器
func NewExpressMng(config *configStruct.AliApiConfig) *ExpressMng {
	return &ExpressMng{
		Config: config,
	}
}

// GetDetailInfo 全国快递物流查询
// Docs: https://market.aliyun.com/products/57126001/cmapi021863.html?spm=5176.2020520132.101.1.2b9f7218CZkxvN#sku=yuncode1586300000
func (mng *ExpressMng) GetDetailInfo(expressNo, expressType string) (*DetailRes, error) {
	resStr, _, _, err := networkHelper.RequestJsonWithStruct(networkStruct.Get, DriveRouteURL, map[string]interface{}{
		"no":   expressNo,   // 快递单号 【顺丰和丰网请输入单号 : 收件人或寄件人手机号后四位。例如：123456789:1234】
		"type": expressType, // 快递公司字母简写：不知道可不填 95%能自动识别，填写查询速度会更快【见产品详情】
	}, map[string]string{
		"Authorization": "APPCODE " + mng.Config.AppCode,
	}, &DetailRes{})

	if err != nil {
		return nil, err
	}

	return resStr.(*DetailRes), nil
}
