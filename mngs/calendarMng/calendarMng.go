package calendarMng

import (
	"fmt"

	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"github.com/wiidz/goutil/structs/networkStruct"
)

const URL = "https://jmhlysjjr.market.alicloudapi.com"

// 以下方法均为聚美智数在阿里云市场中提供的服务
// https://market.aliyun.com/detail/cmapi00066017?spm=5176.730005.result.2.325e414avf3OMU&innerSource=search_%E6%97%A5%E5%8E%86#sku=yuncode6001700003
type CalendarMng struct {
	Debug  bool
	Config *configStruct.AliApiConfig
}

// NewCalendarMng : 返回日历管理器
func NewCalendarMng(config *configStruct.AliApiConfig, Debug bool) *CalendarMng {
	return &CalendarMng{
		Debug:  Debug,
		Config: config,
	}
}

// request 统一的请求处理函数
// 检查错误和响应码，返回 Data 字段
func request[T any](
	mng *CalendarMng,
	method networkStruct.Method,
	url string,
	params map[string]interface{},
	respStruct *BaseResp[T],
) (*T, error) {
	resStr, _, _, err := networkHelper.RequestJsonWithStruct(
		method,
		url,
		params,
		map[string]string{
			"Authorization": "APPCODE " + mng.Config.AppCode,
		},
		respStruct,
		mng.Debug,
	)

	if err != nil {
		return nil, err
	}

	resp, ok := resStr.(*BaseResp[T])
	if !ok {
		return nil, fmt.Errorf("响应类型转换失败")
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("API返回错误: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return resp.Data, nil
}

// GetHolidayList 获取节假日列表
// date string 需要查询的日期，格式：YYYYMMDD（选填）
// needDesc string 是否需要返回当日公众日、国际日和我国传统节日的简介，1-返回，默认不返回（选填）
func (mng *CalendarMng) GetHolidayList(date string, needDesc string) (data *HolidayData, err error) {
	params := make(map[string]interface{})
	if date != "" {
		params["date"] = date
	}
	if needDesc != "" {
		params["needDesc"] = needDesc
	}

	return request(mng, networkStruct.Post, URL+"/holiday/list", params, &HolidayResponse{})
}

// GetHolidayDetail 获取节假日详情
// year string 需要查询的年份，格式：YYYY
func (mng *CalendarMng) GetHolidayDetail(year string) (data *HolidayDetailData, err error) {
	return request(mng, networkStruct.Post, URL+"/holiday/detail", map[string]interface{}{
		"year": year,
	}, &HolidayDetailResp{})
}

// GetAuspiciousTime 获取时辰吉凶数据
// date string 需要查询的日期，格式：YYYYMMDD
func (mng *CalendarMng) GetAuspiciousTime(date string) (data *AuspiciousTimeData, err error) {
	return request(mng, networkStruct.Post, URL+"/luck-tendency/auspicious-time", map[string]interface{}{
		"date": date,
	}, &AuspiciousTimeResp{})
}

// GetAuspiciousDemon 获取宜忌神煞数据
// date string 需要查询的日期，格式：YYYYMMDD
func (mng *CalendarMng) GetAuspiciousDemon(date string) (data *AuspiciousDemonData, err error) {
	return request(mng, networkStruct.Post, URL+"/luck-tendency/auspicious-demon", map[string]interface{}{
		"date": date,
	}, &AuspiciousDemonResp{})
}

// GetAlmanac 获取黄历数据
// date string 需要查询的日期，格式：YYYYMMDD
func (mng *CalendarMng) GetAlmanac(date string) (data *AlmanacData, err error) {
	return request(mng, networkStruct.Post, URL+"/luck-tendency/almanac", map[string]interface{}{
		"date": date,
	}, &AlmanacResp{})
}
