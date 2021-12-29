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
	FormattedAddress string            `json:"formatted_address"` // 格式化的地址：北京市朝阳区望京街道方恒国际中心B座方恒国际中心
	AddressComponent *AddressComponent `json:"addressComponent"`
}

type AddressComponent struct {
	Country       string       `json:"country"`  // 国家：中国
	Province      string       `json:"province"` // 省份：北京市
	City          string       `json:"city"`     // 城市：
	CityCode      string       `json:"citycode"` // 城市编码：010
	District      string       `json:"district"` // 街道：朝阳区
	AdCode        string       `json:"adcode"`   // 邮政编码：110105
	Township      string       `json:"township"` // 乡镇：望京街道
	TownCode      string       `json:"towncode"` // 乡镇编号：110105026000
	Neighborhood  *MetaPlace   `json:"neighborhood"`
	Building      *MetaPlace   `json:"building"`
	StreetNumber  *Street      `json:"streetNumber"`
	BusinessAreas []*MetaPlace `json:"businessAreas"`
}

type MetaPlace struct {
	ID       string `json:"id"`
	Name     []string `json:"name"`     // 名称：方恒国际中心
	Type     []string `json:"type"`     // 类型：商务住宅;楼宇;商务写字楼
	Location string `json:"location"` // 经纬度：（经度，纬度）
}

type Street struct {
	Street    string `json:"street"`    // 街道名：阜通东大街
	Number    string `json:"number"`    // ？？：6-2号楼
	Location  string `json:"location"`  // 经纬度：（经度，纬度）
	Direction string `json:"Direction"` // 方向：西南
	Distance  string `json:"distance"`  // 距离：25.9205
}

// RouteRes 驾车路线
type RouteRes struct {
	Count string `json:"count"`
	Route *Route `json:"route"`
	*BaseRes
}

type Route struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Paths       []*Path `json:"paths"`
}

type Path struct {
	Distance      string      `json:"distance"`       // 27876
	Duration      string      `json:"duration"`       // 4197
	Strategy      string      `json:"strategy"`       // 速度最快
	Tolls         interface{} `json:"tolls"`          // 0
	TollDistance  interface{} `json:"toll_distance"`  // 0
	Restriction   string      `json:"restriction"`    // 0
	TrafficLights string      `json:"traffic_lights"` // 23
	Steps         []*Step     `json:"steps"`
}

type Step struct {
	Instruction     string      `json:"instruction"`      // 向西南行驶44米右转进入主路
	Orientation     string      `json:"orientation"`      // 西南
	Distance        string      `json:"distance"`         // 44
	Tolls           interface{} `json:"tolls"`            // 0
	TollDistance    interface{} `json:"toll_distance"`    // 0
	TollRoad        interface{} `json:"toll_road"`        // []
	Duration        string      `json:"duration"`         // 5
	Polyline        string      `json:"polyline"`         // 116.481216,39.989532;116.48101,39.989311;116.480957,39.989262;116.480904,39.989216
	Action          interface{} `json:"action"`           // 右转 或者 []
	AssistantAction interface{} `json:"assistant_action"` // 进入主路 或者 []
}
