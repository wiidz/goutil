package aliExpressApi

// DetailRes 物流信息
type DetailRes struct {
	Status string `json:"status"` // status 0:正常查询 201:快递单号错误 203:快递公司不存在 204:快递公司识别失败 205:没有信息 207:该单号被限制，错误单号 - 【"0"】
	Msg    string `json:"msg"`    // - 【"ok"】
	Result struct {
		Number string `json:"number"` // 快递号 - 【"780098068058"】
		Type   string `json:"type"`   // 快递公司代号 - 【"zto"】
		List   []struct {
			Time   string `json:"time"`   // 时间 - 【"2018-03-09 11:59:26"】
			Status string `json:"status"` // 状态文字 - 【"【石家庄市】快件已在【长安三部】 签收,签收人: 本人,感谢使用中通快递,期待再次为您服务!"】
		} `json:"list"`
		DeliveryStatus string `json:"deliverystatus"` // 状态 0：快递收件(揽件)1.在途中 2.正在派件 3.已签收 4.派送失败 5.疑难件 6.退件签收  - 【"3"】
		IsSign         string `json:"issign"`         // 是否签收：0=否，1=是签收
		ExpName        string `json:"expName"`        // 快递公司名称 - 【"中通快递"】
		ExpSite        string `json:"expSite"`        // 快递公司官网 - 【"www.zto.com"】
		ExpPhone       string `json:"expPhone"`       // 快递公司电话 - 【"95311"】
		Courier        string `json:"courier"`        // 快递员 或 快递站(没有则为空）- 【"容晓光"】
		CourierPhone   string `json:"courierPhone"`   // 快递员电话 (没有则为空) - 【"13081105270"】
		UpdateTime     string `json:"updateTime"`     // 快递轨迹信息最新时间 - 【"2019-08-27 13:56:19"】
		TakeTime       string `json:"takeTime"`       // 发货到收货消耗时长 (截止最新轨迹) - 【"2天20小时14分"】
		Logo           string `json:"logo"`           // 快递公司LOGO- 【"https://img3.fegine.com/express/zto.jpg"】
	} `json:"result"`
}
