package lbsMng

type BaseRes struct {
	Status   string
	Info     string
	InfoCode string
}

type ReGeoRes struct {
	*BaseRes
	ReGeoCode *ReGeoData
}

type ReGeoData struct {
	FormattedAddress string           `json:"formatted_address"` // 格式化的地址：北京市朝阳区望京街道方恒国际中心B座方恒国际中心
	AddressComponent AddressComponent `json:"addressComponent"`
}

type AddressComponent struct {
	Country       string      `json:"country"`  // 国家：中国
	Province      string      `json:"province"` // 省份：北京市
	City          string      `json:"city"`     // 城市：
	CityCode      string      `json:"citycode"` // 城市编码：010
	District      string      `json:"district"` // 街道：朝阳区
	AdCode        string      `json:"adcode"`   // 邮政编码：110105
	Township      string      `json:"township"` // 乡镇：望京街道
	TownCode      string      `json:"towncode"` // 乡镇编号：110105026000
	Neighborhood  MetaPlace   `json:"neighborhood"`
	Building      MetaPlace   `json:"building"`
	StreetNumber  Street      `json:"streetNumber"`
	BusinessAreas []MetaPlace `json:"businessAreas"`
}

type MetaPlace struct {
	ID       string `json:"id"`
	Name     string `json:"name"`     // 名称：方恒国际中心
	Type     string `json:"type"`     // 类型：商务住宅;楼宇;商务写字楼
	Location string `json:"location"` // 经纬度：（经度，纬度）
}

type Street struct {
	Street    string `json:"street"`    // 街道名：阜通东大街
	Number    string `json:"number"`    // ？？：6-2号楼
	Location  string `json:"location"`  // 经纬度：（经度，纬度）
	Direction string `json:"Direction"` // 方向：西南
	Distance  string `json:"distance"`  // 距离：25.9205
}
