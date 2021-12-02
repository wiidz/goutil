package ossMng

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/wiidz/goutil/structs/configStruct"
)

// easyjson:skip
type OssMng struct {
	Config *configStruct.OssConfig
	Client *oss.Client //一个db连接
}

// easyjson:json
type PolicyConfig struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

// easyjson:json
type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Dir         string `json:"dir"`
}
