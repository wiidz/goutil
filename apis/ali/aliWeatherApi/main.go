package aliWeatherApi

import (
	"errors"
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"github.com/wiidz/goutil/structs/networkStruct"
)

const Domain = "https://ali-weather.showapi.com"
const QueryByCodeURL = "/phone-post-code-weeather"

// WeatherApi 万维易源-全国天气预报查询接口
// https://market.aliyun.com/products/57096001/cmapi010812.html?spm=5176.730005.result.2.766b3524cFRHS9&innerSource=search_%E5%A4%A9%E6%B0%94#sku=yuncode481200005
// 50元=10万次
type WeatherApi struct {
	Config *configStruct.AliApiConfig
}

// NewWeatherApi  天气接口
func NewWeatherApi(config *configStruct.AliApiConfig) *WeatherApi {
	return &WeatherApi{
		Config: config,
	}
}

// QueryByCode 根据城市编码获取天气预报
func (api *WeatherApi) QueryByCode(code string) (resStruct *QueryByCodeRes, err error) {

	temp, _, statusCode, err := networkHelper.RequestJsonWithStruct(networkStruct.Get, QueryByCodeURL, map[string]interface{}{
		"post_code": code, // 邮编，比如上海200000
	}, map[string]string{
		"Authorization": "APPCODE " + api.Config.AppCode,
	}, &QueryByCodeRes{})

	if err != nil {
		return nil, err
	}

	if statusCode == 200 {
		resStruct = temp.(*QueryByCodeRes)

		if resStruct.ShowapiResError != "" {
			err = errors.New(resStruct.ShowapiResError)
			return
		}
	} else if statusCode == 555 {
		// 不扣费，不扣调用次数
		err = errors.New("未知错误")
		return
	}

	return temp.(*QueryByCodeRes), nil
}
