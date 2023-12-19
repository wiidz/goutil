package ipMng

type Resp struct {
	Ret   int       `json:"ret"` // 200
	Msg   string    `json:"msg"` // success
	Data  *RespData `json:"data"`
	LogId string    `json:"log_id"` // 9122a1dcc5a2424a9ef2d8723481adba
}

type RespData struct {
	IP         string `json:"ip"`          // 218.18.228.178
	LongIP     string `json:"long_ip"`     // 3658671282
	ISP        string `json:"isp"`         // 电信
	Area       string `json:"area"`        // 华南
	RegionID   string `json:"region_id"`   // 440000
	Region     string `json:"region"`      // 广东
	CityID     string `json:"city_id"`     // 440300
	City       string `json:"city"`        // 深圳
	District   string `json:"district"`    // 南山区
	DistrictID string `json:"district_id"` // 440305
	CountryID  string `json:"country_id"`  // CN
	Country    string `json:"country"`     // 中国
	Lat        string `json:"lat"`         // 22.528499
	Lng        string `json:"lng"`         // 113.923552
}
