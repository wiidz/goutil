package aliOssApi

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/wiidz/goutil/apis/ali/aliSTSApi"
	"github.com/wiidz/goutil/helpers/timeHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"hash"
	"io"
	"os"
	"strings"
	"time"
)

const base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"

var coder = base64.NewEncoding(base64Table)

// NewOssApi 创建ossApi
func NewOssApi(config *configStruct.AliOssConfig) (*OssApi, error) {
	ossApi := OssApi{
		Config: config,
	}

	err := ossApi.refreshClient()
	return &ossApi, err
}

// refreshClient 在sts没过期的时候直接返回，否则重启
func (ossApi *OssApi) refreshClient() (err error) {

	//【1】获取配置
	config := ossApi.Config
	flag, err := ossApi.refreshSTSData()
	if err != nil {
		return
	}
	stsData := ossApi.STSData

	//【2】获取客户端链接，获取STS临时凭证后，您可以通过其中的安全令牌（SecurityToken）和临时访问密钥（AccessKeyId和AccessKeySecret）生成OSSClient。
	if flag || ossApi.Client == nil {
		//ossApi.Client, err = oss.New(config.EndPoint, config.AccessKeyID, config.AccessKeySecret, oss.SecurityToken(ossApi.STSData.SecurityToken))
		// 我们现在统一使用临时身份去实例化client，而不是系统账户了
		ossApi.Client, err = oss.New(config.EndPoint, stsData.AccessKeyId, stsData.AccessKeySecret, oss.SecurityToken(stsData.SecurityToken))
		if err != nil {
			return
		}

		//【3】实例化bucket
		ossApi.Bucket, err = ossApi.Client.Bucket(ossApi.Config.BucketName)
	}

	return
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

	deByte := base64.StdEncoding.EncodeToString([]byte(result))
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(ossApi.Config.AccessKeySecret))
	_, _ = io.WriteString(h, deByte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	//【5】填充属性
	policyToken = PolicyToken{
		BucketName:  ossApi.Config.BucketName,
		AccessKeyId: ossApi.Config.AccessKeyID,
		Host:        ossApi.Config.Host,
		Expire:      expireEnd,
		Signature:   signedStr,
		Dir:         remotePath,
		Policy:      deByte,
	}
	return
}

// getSTSData 获取临时身份数据
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

// refreshSTSData 刷新临时身份数据
func (ossApi *OssApi) refreshSTSData() (isNewStsData bool, err error) {

	if ossApi.STSData == nil {
		isNewStsData = true
		// 为空
		err = ossApi.getSTSData()
		if err != nil {
			return
		}
	} else if typeHelper.Str2Int64(ossApi.STSData.Expiration) < time.Now().Unix() {
		isNewStsData = true
		// 过期
		err = ossApi.getSTSData()
		if err != nil {
			return
		}
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

// Upload 上传本地文件
func (ossApi *OssApi) Upload(filePath, objectName string) (remoteURL string, err error) {

	//【1】重启一下服务器（可能token过期了）
	err = ossApi.refreshClient()
	if err != nil {
		return
	}

	// 读取本地文件。
	fd, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer fd.Close()

	// 上传文件流。
	err = ossApi.Bucket.PutObject(objectName, fd)
	if err != nil {
		return
	}

	remoteURL = ossApi.GetHost() + "/" + objectName
	return
}

// GetRemotePath 组合远程文件夹路径（目录+时间+用户名+随机数）
func (ossApi *OssApi) GetRemotePath(object string) (remotePath string) {
	now := time.Now().Unix()
	dateStamp := time.Unix(now, 0).Format("20060102")
	remotePath = object + "/" + dateStamp + "/"
	return
}

// GetPrivateObjectURL 获取私密文件的url
func (ossApi *OssApi) GetPrivateObjectURL(object string, expireSec int64) (signedURL string, err error) {

	//【1】重启一下服务器（可能token过期了）
	err = ossApi.refreshClient()
	if err != nil {
		return
	}

	//【2】组合url
	signedURL, err = ossApi.Bucket.SignURL(object, oss.HTTPGet, expireSec) // 使用签名URL将OSS文件下载到流。
	//url = ossApi.GetHost() + "/" + object + "?Expires=" + ossApi.STSData.Expiration + "&OSSAccessKeyId=" + ossApi.STSData.AccessKeyId + "&Signature=" + ossApi.STSData.SecurityToken
	return
}

// SignPrivateURL 给不带签名的私密链接，加上签名，使其可以访问
// 例如url = http://21bcn-private.oss-cn-shanghai.aliyuncs.com/private/20220531/1653982464704.jpg
// domainSuffix = com 或者 cn 这样
// expireSec = 60
func (ossApi *OssApi) SignPrivateURL(url, domainSuffix string, expireSec int64) (signedURL string, err error) {
	//【1】重启一下服务器（可能token过期了）
	err = ossApi.refreshClient()
	if err != nil {
		return
	}
	//【2】组合url
	temp := typeHelper.ExplodeStr(url, "."+domainSuffix+"/")
	if len(temp) < 2 {
		err = errors.New("解析url时失败，无法定位object")
		return
	}
	object := temp[1]
	signedURL, err = ossApi.Bucket.SignURL(object, oss.HTTPGet, expireSec) // 使用签名URL将OSS文件下载到流。
	return
}

// DropSign 将链接中的签名数据全都去除
func DropSign(signedURL string) (rawURL string) {
	// 这里采取最暴力的方式，就是找到第一个?，然后把问号到结尾所有的字符串全部抛弃
	index := strings.Index(signedURL, "?")
	if index == -1 {
		return signedURL
	}

	return signedURL[0 : index-1]
}
