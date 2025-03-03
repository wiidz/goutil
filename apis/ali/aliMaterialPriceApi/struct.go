package aliMaterialPriceApi

// RegionCode 地区代码
type RegionCode string

const (
	Domestic      RegionCode = "domestic"
	International RegionCode = "inter"
)

type KLineType int

const (
	Day    KLineType = 0
	MinOne KLineType = 1

	MinFive   KLineType = 5
	MinThirty KLineType = 30

	HourOne  KLineType = 60
	HourTwo  KLineType = 120
	HourFour KLineType = 240
)

// MaterialCode 金属代码
// 具体查看 https://market.aliyun.com/apimarket/detail/cmapi00068808#sku=yuncode6280800002
// 或者交易市场，本api调用注意国内和国际
type MaterialCode string

const CU MaterialCode = "CU2407" // 沪铜2407
const AL MaterialCode = "AL2406" // 沪铝2406
const SS MaterialCode = "SS2406" // 不锈钢2406

type PriceParam struct {
	Region RegionCode
	Symbol MaterialCode
}

type KLineParam struct {
	Region RegionCode
	Symbol MaterialCode
	Type   KLineType
	Limit  int
}

type PriceResp struct {
	Data struct {
		AskVol     string `json:"ask_vol"`
		Settle     string `json:"settle"`
		Bid1Vol    string `json:"bid1_vol"`
		Hold       string `json:"hold"`
		UpdateTime int    `json:"update_time"`
		High       string `json:"high"`
		Low        string `json:"low"`
		Price      string `json:"price"`
		Ask1Vol    string `json:"ask1_vol"`
		Bid2Vol    string `json:"bid2_vol"`
		ChangeRate string `json:"changeRate"`
		Ask5Vol    string `json:"ask5_vol"`
		Presettle  string `json:"presettle"`
		AvgPx      string `json:"avg_px"`
		Bid4Vol    string `json:"bid4_vol"`
		Change     string `json:"change"`
		Bid3Vol    string `json:"bid3_vol"`
		BidVol     string `json:"bid_vol"`
		Ask2Vol    string `json:"ask2_vol"`
		Volume     string `json:"volume"`
		Bid5Vol    string `json:"bid5_vol"`
		Ask4Vol    string `json:"ask4_vol"`
		Ask5       string `json:"ask5"`
		Ask2       string `json:"ask2"`
		Ask1       string `json:"ask1"`
		Bid5       string `json:"bid5"`
		Ask4       string `json:"ask4"`
		Ask3       string `json:"ask3"`
		Name       string `json:"name"`
		Ask        string `json:"ask"`
		Bid3       string `json:"bid3"`
		Bid4       string `json:"bid4"`
		Bid1       string `json:"bid1"`
		Bid2       string `json:"bid2"`
		Bid        string `json:"bid"`
		Open       string `json:"open"`
		Ask3Vol    string `json:"ask3_vol"`
	} `json:"data"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	TaskNo  string `json:"taskNo"`
}
