package aliOssApi

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/wiidz/goutil/apis/ali/aliSTSApi"
	"github.com/wiidz/goutil/helpers/timeHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"hash"
	"io"
	"log"
	"os"
	"time"
)

const base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"

var coder = base64.NewEncoding(base64Table)

// NewOssApiSingle 创建单例ossApi
func NewOssApiSingle(ossC *configStruct.AliOssConfig) (*OssApi, error) {

	ossApi := OssApi{}

	// 获取STS临时凭证后，您可以通过其中的安全令牌（SecurityToken）和临时访问密钥（AccessKeyId和AccessKeySecret）生成OSSClient。
	client, err := oss.New(ossC.EndPoint, ossC.AccessKeyID, ossC.AccessKeySecret, oss.SecurityToken(ossC.SecurityToken))

	if err != nil {
		return &ossApi, err
	}

	ossApi.Client = client
	return &ossApi, nil
}

// NewOssApi 创建ossApi
func NewOssApi(config *configStruct.AliOssConfig) (*OssApi, error) {
	ossApi := OssApi{
		Config: config,
	}
	client, err := oss.New(config.EndPoint, config.AccessKeyID, config.AccessKeySecret)

	if err != nil {
		return &ossApi, err
	}

	ossApi.Client = client
	return &ossApi, nil
}

// getPolicyToken 获取token
func (ossApi *OssApi) getPolicyToken(remotePath string) (policyToken PolicyToken, err error) {

	//【1】拼接目标目录名称
	//tm := time.Unix(now, 0)
	//tm1 := tm.Format("20060102")
	//upload_dir := object + "/" + tm1 + "/"
	//fmt.Println(upload_dir)

	//【2】token过期时间
	expireEnd := time.Now().Unix() + ossApi.Config.ExpireTime
	tokenExpire := timeHelper.GetISO8601(expireEnd)

	//【3】构建上传策略json
	var policyConfig PolicyConfig
	policyConfig.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, remotePath)
	policyConfig.Conditions = append(policyConfig.Conditions, condition)

	//【4】计算签名
	result, err := typeHelper.JsonEncode(policyConfig)
	if err != nil {
		return
	}

	debyte := base64.StdEncoding.EncodeToString([]byte(result))
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(ossApi.Config.AccessKeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	//【5】填充属性
	policyToken = PolicyToken{
		BucketName:  ossApi.Config.BucketName,
		AccessKeyId: ossApi.Config.AccessKeyID,
		Host:        ossApi.Config.Host,
		Expire:      expireEnd,
		Signature:   signedStr,
		Dir:         remotePath,
		Policy:      debyte,
	}
	return
}

// GetHost 获取域名
func (ossApi *OssApi) GetHost() string {
	return ossApi.Config.Host
}

// GetSign 获取签名
func (ossApi *OssApi) GetSign(object string) (PolicyToken, error) {
	return ossApi.getPolicyToken(object)
}

// Upload 上传
func (ossApi *OssApi) Upload(filePath, objectName string) (string, error) {
	//response, err := ossApi.getPolicyToken(object)

	// 获取存储空间。
	bucket, err := ossApi.Client.Bucket(ossApi.Config.BucketName)
	if err != nil {
		return "", err
	}

	// 读取本地文件。
	fd, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	// 上传文件流。
	err = bucket.PutObject(objectName, fd)
	if err != nil {
		return "", err
	}
	ossPath := ossApi.GetHost() + "/" + objectName
	return ossPath, nil
}

// GetBucketInfo 获取Bucket信息
func (ossApi *OssApi) GetBucketInfo() {
	//response, err := ossApi.getPolicyToken(object)

	// 获取存储空间。
	bucket, _ := ossApi.Client.Bucket(ossApi.Config.BucketName)
	log.Println("bucket", bucket)
	//log.Println("bucket",bucket.Get)

}

// SimpleGetOssSign 简单获取签名
func SimpleGetOssSign(ossApi *OssApi, object string) (msg string, data interface{}, statusCode int) {

	remotePath := GetRemotePath(object)

	res, err := ossApi.GetSign(remotePath)

	if err != nil {
		return "获取签名失败", "", 400
	}

	return "ok", res, 200
}

// GetRemotePath 组合远程文件夹路径（目录+时间+用户名+随机数）
func GetRemotePath(object string) (remotePath string) {
	now := time.Now().Unix()
	dateStamp := time.Unix(now, 0).Format("20060102")
	remotePath = object + "/" + dateStamp + "/"
	return
}

// GetPrivateObjectURL 获取私密文件的url
func (ossApi *OssApi) GetPrivateObjectURL(object string) (url string, err error) {

	//【1】拉取访问权限
	if ossApi.STSData == nil {
		err = ossApi.getSTSData()
		if err != nil {
			return
		}
	}

	//【2】组合url
	url = ossApi.GetHost() + "/" + object + "?Expires=" + ossApi.STSData.Expiration + "&OSSAccessKeyId=" + ossApi.STSData.AccessKeyId + "&Signature=" + ossApi.STSData.SecurityToken
	return
}

func (ossApi *OssApi) getSTSData() (err error) {

	stsApi, err := aliSTSApi.NewAliSTSApi(&configStruct.AliRamConfig{
		AccessKeyID:     ossApi.Config.AccessKeyID,
		AccessKeySecret: ossApi.Config.AccessKeySecret,
	})
	if err != nil {
		return
	}

	res, err := stsApi.AssumeRole(ossApi.Config.ARN, "oss_role")
	if err != nil {
		return
	}

	ossApi.STSData = &res.Credentials
	return
}
