package aliMaterialPriceApi

import (
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"github.com/wiidz/goutil/structs/networkStruct"
)

const Domain = "https://jmgjs.market.alicloudapi.com/precious-metal"
const PriceURL = "/price"
const KLineURL = "/kline"

// AliMaterialPriceApi 阿里云贵金属价格查询接口
// https://market.aliyun.com/apimarket/detail/cmapi00068808#sku=yuncode6280800002
type AliMaterialPriceApi struct {
	Config *configStruct.AliApiConfig
}

// NewAliMaterialPriceApi 贵金属价格
func NewAliMaterialPriceApi(config *configStruct.AliApiConfig) *AliMaterialPriceApi {
	return &AliMaterialPriceApi{
		Config: config,
	}
}

// GetPrice 获取价格
func (api *AliMaterialPriceApi) GetPrice(param *PriceParam) (resp *PriceResp, err error) {
	resStr, _, _, err := networkHelper.RequestJsonWithStruct(networkStruct.Post, Domain+string(param.Region)+PriceURL, map[string]interface{}{
		"symbol": param.Symbol,
	}, map[string]string{
		"Authorization": "APPCODE " + api.Config.AppCode,
	}, &PriceResp{})

	if err != nil {
		return
	}

	return resStr.(*PriceResp), nil
}

// GetKLine 获取贵金属K线
// amount默认是10条
func (api *AliMaterialPriceApi) GetKLine(param *KLineParam) (resp *PriceResp, err error) {
	resStr, _, _, err := networkHelper.RequestJsonWithStruct(networkStruct.Post, Domain+string(param.Region)+KLineURL, map[string]interface{}{
		"symbol": param.Symbol,    // 国际贵金属品种，详见国际贵金属现货，详见国际贵金属期货
		"type":   int(param.Type), // k线类型 0：日k，1：1分钟，5：五分钟，30：30分钟，60：60分钟，120：120分钟，240：240分钟
		"limit":  param.Limit,     // 返回条数 默认10
	}, map[string]string{
		"Authorization": "APPCODE " + api.Config.AppCode,
	}, &PriceResp{})

	if err != nil {
		return
	}

	return resStr.(*PriceResp), nil
}

// 国际贵金属期货合约
// 国际贵金属报价
// 国内贵金属期货合约
// 国内贵金属K线
