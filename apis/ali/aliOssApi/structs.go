package aliOssApi

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/wiidz/goutil/structs/configStruct"
)

// easyjson:skip
type OssApi struct {
	Config  *configStruct.AliOssConfig
	Client  *oss.Client //一个db连接
	STSData *sts.Credentials
}

// easyjson:json
type PolicyConfig struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

// easyjson:json
type PolicyToken struct {
	BucketName  string `json:"bucket_name"`
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Dir         string `json:"dir"`
}
