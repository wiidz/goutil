package audioToTextMng

import "github.com/wiidz/goutil/structs/configStruct"

type AudioToTextMng struct {
	Config *configStruct.AliApiConfig
}

type Resp struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Data    struct {
		OrderNo string   `json:"orderNo"`
		Result  []string `json:"result"`
	} `json:"data"`
}
