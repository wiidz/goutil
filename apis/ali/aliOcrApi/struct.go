package aliOcrApi

type CardSide string

const Front CardSide = "face" // 首页（人像面等）
const Back CardSide = "back"  // 反面（国徽面等）

type OcrError struct {
	OriginalErr error  `json:"original_err"` // 原始错误信息
	Code        string `json:"Cope"`
	HostId      string `json:"HostId"`
	Message     string `json:"Message"`
	Recommend   string `json:"Recommend"`
	RequestId   string `json:"RequestId"`
}

func (obj *OcrError) Error() string {
	return obj.Message
}
