package calendarMng

import "encoding/json"

// BaseResp 通用API响应结构
type BaseResp[T any] struct {
	Code   int    `json:"code"`   // 返回码，详见返回码说明
	Msg    string `json:"msg"`    // 返回消息
	TaskNo string `json:"taskNo"` // 请求号
	Data   *T     `json:"data"`   // 响应数据
}

// HolidayResponse 节假日API响应结构
type HolidayResponse = BaseResp[HolidayData]

// Data 节假日数据
type HolidayData struct {
	Count int           `json:"count"` // 一年的节假日数量
	Items []HolidayItem `json:"items"` // 一年的节假日列表
}

// HolidayItem 节假日项
type HolidayItem struct {
	Begin         string   `json:"begin"`          // 节假日的开始时间，格式：YYYYMMDD
	End           string   `json:"end"`            // 节假日的结束时间，格式：YYYYMMDD
	Holiday       string   `json:"holiday"`        // 节假日名称
	HolidayRemark string   `json:"holiday_remark"` // 节假日描述
	InverseDays   []string `json:"inverse_days"`   // 调修日列表，只有2021年之后的数据有值，格式：YYYYMMDD
}

// AlmanacResp 黄历API响应结构
type AlmanacResp = BaseResp[AlmanacData]

// AlmanacData 黄历数据
type AlmanacData struct {
	Gongli    string `json:"gongli"`    // 公历，格式：公元 YYYY年MM月DD日 星期X
	Nongli    string `json:"nongli"`    // 农历
	Jieri     string `json:"jieri"`     // 节日
	Zhiri     string `json:"zhiri"`     // 值日
	Zhishen   string `json:"zhishen"`   // 值神
	Yi        string `json:"yi"`        // 宜
	Ji        string `json:"ji"`        // 忌
	Qixiang   string `json:"qixiang"`   // 气象
	Jieqi24   string `json:"jieqi24"`   // 当前月包含的24节气
	Shengxiao string `json:"shengxiao"` // 生肖
	Xingzuo   string `json:"xingzuo"`   // 星座
	Rulueli   string `json:"rulueli"`   // 儒略历
	Jsyq      string `json:"jsyq"`      // 吉神宜趋
	Xsyj      string `json:"xsyj"`      // 凶神宜忌
	Pzbj      string `json:"pzbj"`      // 彭祖百忌
	Tszf      string `json:"tszf"`      // 胎神占方
	Chongsha  string `json:"chongsha"`  // 冲煞
	Nayin     string `json:"nayin"`     // 纳音
	Dizhi     string `json:"dizhi"`     // 地支
	Ganzhi    string `json:"ganzhi"`    // 干支
}

// AuspiciousDemonResp 宜忌神煞API响应结构
type AuspiciousDemonResp = BaseResp[AuspiciousDemonData]

// AuspiciousDemonData 宜忌神煞数据
type AuspiciousDemonData struct {
	Niansansha   string `json:"niansansha"`   // 年三煞
	Nianqisha    string `json:"nianqisha"`    // 年七煞
	Niankongwang string `json:"niankongwang"` // 年空亡
	Yuezhi       string `json:"yuezhi"`       // 月支
	Yueling      string `json:"yueling"`      // 月令
	Yuexiang     string `json:"yuexiang"`     // 月相
	Yuesansha    string `json:"yuesansha"`    // 月三煞
	Yueqisha     string `json:"yueqisha"`     // 月七煞
	Yuekongwang  string `json:"yuekongwang"`  // 月空亡
	Risansha     string `json:"risansha"`     // 日三煞
	Riqisha      string `json:"riqisha"`      // 日七煞
	Rikongwang   string `json:"rikongwang"`   // 日空亡
	Tjjs         string `json:"tjjs"`         // 推荐吉时
	Taisuiwei    string `json:"taisuiwei"`    // 太岁位
	Fantaisui    string `json:"fantaisui"`    // 犯太岁
	Esbx         string `json:"esbx"`         // 二十八宿
	Jiuxing      string `json:"jiuxing"`      // 九星
	Rilu         string `json:"rilu"`         // 日禄
	Zhongdong    string `json:"zhongdong"`    // 仲冬
	Suipowei     string `json:"suipowei"`     // 岁破位
	Niantaisui   string `json:"niantaisui"`   // 年太岁
	Caishen      string `json:"caishen"`      // 财神
	Xishen       string `json:"xishen"`       // 喜神
	Yangguishen  string `json:"yangguishen"`  // 阳贵神
	Yinguishen   string `json:"yinguishen"`   // 阴贵神
	Fushen       string `json:"fushen"`       // 福神
	Yjgx         string `json:"yjgx"`         // 易经卦象
	Wuhou        string `json:"wuhou"`        // 物候
	Zhishen12    string `json:"zhishen12"`    // 十二值神
	Zhiri12      string `json:"zhiri12"`      // 十二值日
	Liuyao       string `json:"liuyao"`       // 六曜
}

// AuspiciousTimeItem 时辰吉凶数据项
type AuspiciousTimeItem struct {
	Shijian   string `json:"shijian"`   // 时间，格式：HH:MM:SS-HH:MM:SS
	Jixiong   string `json:"jixiong"`   // 吉凶
	Jishen    string `json:"jishen"`    // 吉神
	Xiongshen string `json:"xiongshen"` // 凶神
	Shichong  string `json:"shichong"`  // 时冲
	Shizhu    string `json:"shizhu"`    // 时柱
}

// AuspiciousTimeResp 时辰吉凶API响应结构
type AuspiciousTimeResp = BaseResp[AuspiciousTimeData]

// AuspiciousTimeData 时辰吉凶数据
// key 为时辰名称：zi(子), chou(丑), yin(寅), mao(卯), cheng(辰), si(巳), wu(午), wei(未), shen(申), you(酉), xu(戌), hai(亥)
type AuspiciousTimeData struct {
	UT    string                         `json:"ut,omitempty"`
	Times map[string]*AuspiciousTimeItem `json:"-"`
}

// UnmarshalJSON 兼容包含 ut 字段的响应，将其它键解析为时辰项
func (a *AuspiciousTimeData) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	times := make(map[string]*AuspiciousTimeItem)
	for key, value := range raw {
		if key == "ut" {
			_ = json.Unmarshal(value, &a.UT)
			continue
		}

		var item AuspiciousTimeItem
		if err := json.Unmarshal(value, &item); err != nil {
			// 如果单个时辰解析失败，跳过而不影响其它数据
			continue
		}
		times[key] = &item
	}

	a.Times = times
	return nil
}

// MarshalJSON 将内部结构重新展开为与 API 响应一致的扁平结构
func (a AuspiciousTimeData) MarshalJSON() ([]byte, error) {
	out := make(map[string]any, len(a.Times)+1)
	if a.UT != "" {
		out["ut"] = a.UT
	}
	for key, value := range a.Times {
		out[key] = value
	}
	return json.Marshal(out)
}

// HolidayDetailResp 节假日详情API响应结构
type HolidayDetailResp = BaseResp[HolidayDetailData]

// HolidayDetailItem 节日详情项
type HolidayDetailItem struct {
	Name    string `json:"name"`    // 节日名称
	Genus   string `json:"genus"`   // 节日种类，public表示公众日或国际日，traditional表示传统节日
	Day     string `json:"day"`     // 节日公历日期，格式：MMDD
	LunaDay string `json:"lunaDay"` // 节日的农历日期
	Info    string `json:"info"`    // 节日的简介
	Origin  string `json:"origin"`  // 节日的起源
}

// HolidayDetailData 节假日详情数据
type HolidayDetailData struct {
	Day           string              `json:"day"`            // 查询的日期，格式：YYYYMMDD
	Holiday       string              `json:"holiday"`        // 节日名称，工作日时显示"无"，周末时显示"周末"，节日时显示节日名称
	Type          string              `json:"type"`           // 类型：1为工作日，2为周末，3为节假日
	Begin         string              `json:"begin"`          // 节日或周末开始时间，格式：YYYYMMDD，如果是工作日，此字段为空串
	End           string              `json:"end"`            // 节日或周末结束时间，格式：YYYYMMDD，如果是工作日，此字段为空串
	HolidayRemark string              `json:"holiday_remark"` // 节日备注
	WeekDay       int                 `json:"weekDay"`        // 星期几的数字，0-6（0为周日）
	Cn            string              `json:"cn"`             // 星期几的中文名
	En            string              `json:"en"`             // 星期几的英文名
	H             []HolidayDetailItem `json:"h"`              // 节日列表，如果当日没有节日则返回空List，该字段需要请求参数中needDesc上传1才返回
}
