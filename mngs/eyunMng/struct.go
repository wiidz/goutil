package eyunMng

type TestReturn struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Data    struct {
		CallbackUrl   interface{} `json:"callbackUrl"`
		Status        int         `json:"status"`
		Authorization string      `json:"Authorization"`
	} `json:"data"`
}

type QrcodeReturn struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Data    struct {
		WId       string `json:"wId"`
		QrCodeUrl string `json:"qrCodeUrl"`
	} `json:"data"`
}

type AddressListReturn struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Chatrooms []string `json:"chatrooms"`
		Friends   []string `json:"friends"`
		Ghs       []string `json:"ghs"`
		Others    []string `json:"others"`
	} `json:"data"`
}

type SendMsgReturn struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Type       int    `json:"type"`
		MsgId      int64  `json:"msgId"`
		NewMsgId   int64  `json:"newMsgId"`
		CreateTime int    `json:"createTime"`
		WcId       string `json:"wcId"`
	} `json:"data"`
}

type LoginReturn struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Country         string `json:"country"`
		WAccount        string `json:"wAccount"`
		DeviceType      string `json:"deviceType"`
		City            string `json:"city"`
		Signature       string `json:"signature"`
		NickName        string `json:"nickName"`
		Sex             int    `json:"sex"`
		HeadUrl         string `json:"headUrl"`
		Type            int    `json:"type"`
		SmallHeadImgUrl string `json:"smallHeadImgUrl"`
		WcId            string `json:"wcId"`
		WId             string `json:"wId"`
		MobilePhone     string `json:"mobilePhone"`
		Uin             int    `json:"uin"`
		Status          int    `json:"status"`
		Username        string `json:"username"`
	} `json:"data"`
}
