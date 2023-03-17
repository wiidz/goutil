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
	Country       string        `json:"country"`  // 国家：中国
	Province      string        `json:"province"` // 省份：北京市
	City          string        `json:"city"`     // 城市：
	CityCode      string        `json:"citycode"` // 城市编码：010
	District      string        `json:"district"` // 街道：朝阳区
	AdCode        string        `json:"adcode"`   // 邮政编码：110105
	Township      string        `json:"township"` // 乡镇：望京街道
	TownCode      string        `json:"towncode"` // 乡镇编号：110105026000
	Neighborhood  *Neighborhood `json:"neighborhood"`
	Building      *Neighborhood `json:"building"`
	StreetNumber  *Street       `json:"streetNumber"`
	BusinessAreas []*MetaPlace  `json:"businessAreas"`
}

type MetaPlace struct {
	ID       string `json:"id"`
	Name     string `json:"name"`     // 名称：方恒国际中心
	Location string `json:"location"` // 经纬度：（经度，纬度）
}

type ReGeoResWithoutBusinessAreas struct {
	*BaseRes
	ReGeoCode *ReGeoDataWithoutBusinessAreas
}

// ReGeoDataWithoutBusinessAreas 上面那个可能报错，
type ReGeoDataWithoutBusinessAreas struct {
	FormattedAddress string                                `json:"formatted_address"` // 格式化的地址：北京市朝阳区望京街道方恒国际中心B座方恒国际中心
	AddressComponent *AddressComponentWithoutBusinessAreas `json:"addressComponent"`
}

// AddressComponentWithoutBusinessAreas 上面那个可能报错，因为BusinessAreas的值可能是[[]]这样的
type AddressComponentWithoutBusinessAreas struct {
	Country      string        `json:"country"`  // 国家：中国
	Province     string        `json:"province"` // 省份：北京市
	City         string        `json:"city"`     // 城市：
	CityCode     string        `json:"citycode"` // 城市编码：010
	District     string        `json:"district"` // 街道：朝阳区
	AdCode       string        `json:"adcode"`   // 邮政编码：110105
	Township     string        `json:"township"` // 乡镇：望京街道
	TownCode     string        `json:"towncode"` // 乡镇编号：110105026000
	Neighborhood *Neighborhood `json:"neighborhood"`
	Building     *Neighborhood `json:"building"`
	StreetNumber *Street       `json:"streetNumber"`
	//BusinessAreas []*MetaPlace  `json:"businessAreas"`
}

type Neighborhood struct {
	//ID       string `json:"id"`
	Name []string `json:"name"` // 名称：方恒国际中心
	Type []string `json:"type"` // 类型：商务住宅;楼宇;商务写字楼
	//Location string `json:"location"` // 经纬度：（经度，纬度）
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
type RouteShowFields struct {
	Cost *Cost `json:"cost"` // 设置后可返回方案所需时间及费用成本
	Tmcs *Tmcs `json:"tmcs"` // 设置后可返回分段路况详情
	//Navi     `json:"navi"`     // 设置后可返回详细导航动作指令
	//Cities   `json:"cities"`   // 设置后可返回分段途径城市信息
	//polyline `json:"Polyline"` // 设置后可返回分路段坐标点串，两点间用“,”分隔
}

type Cost struct {
	Duration      string `json:"duration"`       // 线路耗时，分段step中的耗时
	Tolls         string `json:"tolls"`          // 线路耗时，分段step中的耗时
	TollDistance  string `json:"toll_distance"`  // 线路耗时，分段step中的耗时
	TollRoad      string `json:"toll_road"`      // 线路耗时，分段step中的耗时
	TrafficLights string `json:"traffic_lights"` // 线路耗时，分段step中的耗时
}

type Tmcs struct {
	TmcStatus    string `json:"tmc_status"`   // 路况信息，包括：未知、畅通、缓行、拥堵、严重拥堵
	Tolls        string `json:"tmc_distance"` // 从当前坐标点开始step中路况相同的距离
	TollDistance string `json:"tmc_polyline"` // 此段路况涉及的道路坐标点串，点间用","分隔
}

type Route struct {
	Origin      string  `json:"origin"`      // 起点经纬度
	Destination string  `json:"destination"` // 终点经纬度
	TaxiCost    string  `json:"taxi_cost"`   // 预计出租车费用，单位：元
	Paths       []*Path `json:"paths"`       // 算路方案详情
}

type Path struct {
	// Basic
	Distance string `json:"distance"` // 方案距离，单位：米 27876
	//Duration      string      `json:"duration"`       // 4197
	//Strategy      string      `json:"strategy"`       // 速度最快
	//Tolls         interface{} `json:"tolls"`          // 0
	//TollDistance  interface{} `json:"toll_distance"`  // 0
	Restriction string `json:"restriction"` // 0 代表限行已规避或未限行，即该路线没有限行路段  1 代表限行无法规避，即该线路有限行路段
	//TrafficLights string  `json:"traffic_lights"` // 23
	Steps []*Step `json:"steps"` // 路线分段

	// RouteShowFields
	Cost *Cost `json:"cost"` // 设置后可返回方案所需时间及费用成本
	Tmcs *Tmcs `json:"tmcs"` // 设置后可返回分段路况详情
	//Navi     `json:"navi"`     // 设置后可返回详细导航动作指令
	//Cities   `json:"cities"`   // 设置后可返回分段途径城市信息
	//polyline `json:"Polyline"` // 设置后可返回分路段坐标点串，两点间用“,”分隔
}

type Step struct {
	Instruction  string `json:"instruction"`   // 行驶指示 向西南行驶44米右转进入主路
	Orientation  string `json:"orientation"`   // 进入道路方向 西南
	RoadName     string `json:"road_name"`     // 分段道路名称
	StepDistance string `json:"step_distance"` // 分段距离信息 44

	//Tolls           interface{} `json:"tolls"`            // 0
	//TollDistance    interface{} `json:"toll_distance"`    // 0
	//TollRoad        interface{} `json:"toll_road"`        // []
	//Duration        string      `json:"duration"`         // 5
	//Polyline        string      `json:"polyline"`         // 116.481216,39.989532;116.48101,39.989311;116.480957,39.989262;116.480904,39.989216
	//Action          interface{} `json:"action"`           // 右转 或者 []
	//AssistantAction interface{} `json:"assistant_action"` // 进入主路 或者 []
}
