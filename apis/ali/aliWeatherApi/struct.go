package aliWeatherApi

// QueryByCodeRes 根据区号查询响应体
type QueryByCodeRes struct {
	ShowapiResError string          `json:"showapi_res_error"`
	ShowapiResId    string          `json:"showapi_res_id"`
	ShowapiResCode  int             `json:"showapi_res_code"`
	ShowapiFeeNum   int             `json:"showapi_fee_num"`
	ShowapiResBody  QueryByCodeBody `json:"showapi_res_body"`
}

// QueryByCodeBody 根据区号查询数据体
type QueryByCodeBody struct {
	Time    string `json:"time"`
	Remark  string `json:"remark"`
	RetCode int    `json:"ret_code"`

	CityInfo struct {
		Longitude float64 `json:"longitude"` //经度（100.222）
		Latitude  float64 `json:"latitude"`  // 纬度（26.903）

		C0  string `json:"c0"`
		C1  string `json:"c1"`  // 区域ID（101291401）
		C2  string `json:"c2"`  // 城市英文名（lijiang）
		C3  string `json:"c3"`  // 城市中文名（丽江）
		C4  string `json:"c4"`  // 城市所在市英文名（lijiang）
		C5  string `json:"c5"`  // 城市所在市中文名（丽江）
		C6  string `json:"c6"`  // 城市所在省英文名（yunnan）
		C7  string `json:"c7"`  // 城市所在省中文名（云南）
		C8  string `json:"c8"`  // 城市所在国家英文名（china）
		C9  string `json:"c9"`  // 城市所在国家中文名
		C10 string `json:"c10"` // 城市级别（2）
		C11 string `json:"c11"` // 城市区号（0888）
		C12 string `json:"c12"` // 邮编（674100）
		C15 string `json:"c15"` // 海拔（2394）
		C16 string `json:"c16"` // 雷达站号（AZ9888）
		C17 string `json:"c17"` // （+8）
	} `json:"cityInfo"` // 查询的地区基本资料

	Now *NowWeather `json:"now"` // 现在实时的天气情况

	F1 *DailyWeather `json:"f1"` // 今天的天气预报 ,f2是明天，f3是后天，直到f7
	F2 *DailyWeather `json:"f2"`
	F3 *DailyWeather `json:"f3"`
}

// DailyWeather 天气预测
type DailyWeather struct {
	Day     string `json:"day"`     // 日期（20151023）
	Weekday int    `json:"weekday"` // 星期几（5）

	Ziwaixian   string `json:"ziwaixian"`     // 紫外线强度（中等）
	SunBeginEnd string `json:"sun_begin_end"` // 日出|日落时间（06:35|17:23）
	AirPress    string `json:"air_press"`     // 大气压（1008 hPa）
	Jiangshui   string `json:"jiangshui"`     // 降水概率（3%）

	DayWeather        string `json:"day_weather"` // 白天天气标识（晴）
	DayWeatherCode    string `json:"day_weather_code"`
	DayWeatherPic     string `json:"day_weather_pic"`     // 白天天气图标（http://app1.showapi.com/weather/icon/day/00.png）
	DayAirTemperature string `json:"day_air_temperature"` // 白天天气温度摄氏度（18）
	DayWindPower      string `json:"day_wind_power"`      // 白天风力（微风<10m/h）
	DayWindDirection  string `json:"day_wind_direction"`  // 白天风向（无持续风向）

	NightWeather        string `json:"night_weather"` // 晚上天气标识（晴）
	NightWeatherCode    string `json:"night_weather_code"`
	NightWeatherPic     string `json:"night_weather_pic"`     // 晚上天气图标（http://app1.showapi.com/weather/icon/night/00.png）
	NightAirTemperature string `json:"night_air_temperature"` //晚上气温摄氏度（9）
	NightWindPower      string `json:"night_wind_power"`      // 晚上风力（微风<10m/h）
	NightWindDirection  string `json:"night_wind_direction"`  // 晚上风向（无持续风向）

}

// NowWeather 当前天气
type NowWeather struct {
	Weather     string `json:"weather"`     // 天气文字标识（晴）
	WeatherPic  string `json:"weather_pic"` // 天气图标（http://app1.showapi.com/weather/icon/day/00.png）
	WeatherCode string `json:"weather_code"`

	WindPower     string `json:"wind_power"`     // 风力（1级）
	WindDirection string `json:"wind_direction"` // 风向名称（北风）

	Aqi       string     `json:"aqi"`       // 空气质量指数，越小越好（71）
	AqiDetail *AqiDetail `json:"aqiDetail"` // aqi明细数据

	Temperature     string `json:"temperature"`      // 当前气温摄氏度（15）
	TemperatureTime string `json:"temperature_time"` // 采集时间（18：30）

	Rain string `json:"rain"` //
	Sd   string `json:"sd"`   // 空气湿度（70%）
}

type AqiDetail struct {
	Area     string `json:"area"`      // 地区中文（北京）
	AreaCode string `json:"area_code"` // 地区拼音（beijing）

	Aqi string `json:"aqi"` // 空气质量指数，越小越好（71）

	Quality string `json:"quality"` // 空气质量指数类别，有“优质、良好、轻度污染、中度污染、重度污染、严重污染”6类 （良）

	PrimaryPollutant string `json:"primary_pollutant"` //（颗粒物(PM2.5)）
	Pm25             string `json:"pm2_5"`             //颗粒物（粒径小于等于2.5μm）1小时平均（51）
	Pm10             string `json:"pm10"`              //颗粒物（粒径小于等于10μm）1小时平均（56）

	Co   string `json:"co"`    // 一氧化碳1小时平均（0.817）
	So2  string `json:"so2"`   // 二氧化硫1小时平均（3）
	No2  string `json:"no2"`   // 二氧化氮1小时平均（52）
	O3   string `json:"o3"`    // 臭氧1小时平均（33）
	O38H string `json:"o3_8h"` // 臭氧8小时平均（9）

	Num string `json:"num"` //

}
