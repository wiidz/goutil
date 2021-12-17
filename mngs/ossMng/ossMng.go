package ossMng

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/wiidz/goutil/helpers/timeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"hash"
	"io"
	"log"
	"os"
	"time"
)

const base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"

var coder = base64.NewEncoding(base64Table)

// NewOssMngSingle 创建单例ossMng
func NewOssMngSingle(ossC *configStruct.OssConfig) (*OssMng, error) {

	ossMng := OssMng{}
	client, err := oss.New(ossC.EndPoint, ossC.AccessKeyID, ossC.AccessKeySecret)

	if err != nil {
		return &ossMng, err
	}

	ossMng.Client = client
	return &ossMng, nil
}

// NewOssMng 创建ossMng
func NewOssMng(config *configStruct.OssConfig) (*OssMng, error) {
	ossMng := OssMng{
		Config: config,
	}
	client, err := oss.New(config.EndPoint, config.AccessKeyID, config.AccessKeySecret)

	if err != nil {
		return &ossMng, err
	}

	ossMng.Client = client
	return &ossMng, nil
}

// getPolicyToken 获取token
func (ossMng *OssMng) getPolicyToken(remotePath string) (PolicyToken, error) {
	policyToken := PolicyToken{}

	//【1】拼接目标目录名称
	//tm := time.Unix(now, 0)
	//tm1 := tm.Format("20060102")
	//upload_dir := object + "/" + tm1 + "/"
	//fmt.Println(upload_dir)

	//【2】token过期时间
	expireEnd := time.Now().Unix() + ossMng.Config.ExpireTime
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
	result, err := policyConfig.MarshalJSON()
	if err != nil {
		return policyToken, err
	}

	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(ossMng.Config.AccessKeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	//【5】填充属性
	policyToken.AccessKeyId = ossMng.Config.AccessKeyID
	policyToken.Host = ossMng.Config.Host
	policyToken.Expire = expireEnd
	policyToken.Signature = string(signedStr)
	policyToken.Dir = remotePath
	policyToken.Policy = string(debyte)

	return policyToken, nil
}

// GetHost 获取域名
func (ossMng *OssMng) GetHost() string {
	return ossMng.Config.Host
}

// GetSign 获取签名
func (ossMng *OssMng) GetSign(object string) (PolicyToken, error) {
	response, err := ossMng.getPolicyToken(object)
	return response, err
}

// Upload 上传
func (ossMng *OssMng) Upload(filePath, objectName string) (string, error) {
	//response, err := ossMng.getPolicyToken(object)

	// 获取存储空间。
	bucket, err := ossMng.Client.Bucket(ossMng.Config.BucketName)
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
	ossPath := ossMng.GetHost() + "/" + objectName
	return ossPath, nil
}

// GetBucketInfo 获取Bucket信息
func (ossMng *OssMng) GetBucketInfo() {
	//response, err := ossMng.getPolicyToken(object)

	// 获取存储空间。
	bucket, _ := ossMng.Client.Bucket(ossMng.Config.BucketName)
	log.Println("bucket", bucket)
	//log.Println("bucket",bucket.Get)

}

// SimpleGetOssSign 简单获取签名
func SimpleGetOssSign(ossMng *OssMng, object string) (msg string, data interface{}, statusCode int) {

	remotePath := GetRemotePath(object)

	res, err := ossMng.GetSign(remotePath)

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
