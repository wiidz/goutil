package imMsg

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type UploadMultipartParams struct {
	Content   string `json:"content"`   // 最大15M的字符流
	Type      string `json:"type"`      // 上传文件类型
	IsHttps   bool   `json:"ishttps"`   // 返回的url是否需要为https的url，true或false，默认false
	ExpireSec int    `json:"expireSrc"` // 文件过期时长，单位：秒，必须大于等于86400
	Tag       string `json:"tag"`       // 文件的应用场景，不超过32个字符
}

type UploadMultipartResp struct {
	*imClient.CommonResp
	URL string `json:"url"` // 文件地址
}

// UploadMultipart 文件上传
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DEwMTE3NzQ?platformId=60353#文件上传
// 文件上传，字符流需要base64编码，最大15M。
func (api *Api) UploadMultipart(param *UploadMultipartParams) (*UploadMultipartResp, error) {
	res, err := api.Client.Post(SubDomain+"fileUpload.action", param, &UploadMultipartResp{})
	return res.(*UploadMultipartResp), err
}
