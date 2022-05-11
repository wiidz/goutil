package imClient

const IMDomain = "https://api.netease.im/nimserver"

type CommonResp struct {
	Code int         `json:"code"`
	Desc string      `json:"desc"`
	Info interface{} `json:"info"`
}

func (resp *CommonResp) GetCode() int {
	return resp.Code
}
func (resp *CommonResp) GetDesc() string {
	return resp.Desc
}
func (resp *CommonResp) GetInfo() interface{} {
	return resp.Info
}

type RespInterface interface {
	GetCode() int
	GetDesc() string
	GetInfo() interface{}
}
