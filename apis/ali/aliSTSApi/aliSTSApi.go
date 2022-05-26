package aliSTSApi

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/wiidz/goutil/structs/configStruct"
)

const RegionID = "cn-hangzhou" // 写死了

type AliSTSApi struct {
	Client *sts.Client
	Config *configStruct.AliRamConfig
}

func NewAliSTSApi(config *configStruct.AliRamConfig) (api *AliSTSApi, err error) {

	//构建一个阿里云客户端, 用于发起请求。
	//设置调用者（RAM用户或RAM角色）的AccessKey ID和AccessKey Secret。
	var client *sts.Client
	client, err = sts.NewClientWithAccessKey(RegionID, config.AccessKeyID, config.AccessKeySecret)
	api = &AliSTSApi{
		Client: client,
		Config: config,
	}
	return
}

// AssumeRole 调用AssumeRole获取一个扮演RAM角色的临时身份凭证（STS Token）。
// https://help.aliyun.com/document_detail/28763.htm?spm=a2c4g.11186623.0.0.38ab3767aV8tng#reference-clc-3sv-xdb
//【前提条件】
// 请确保已为调用者（RAM用户或RAM角色）授予STS的管理权限（AliyunSTSAssumeRoleAccess）。
// 否则，会报如下错误：
// You are not authorized to do this action. You should be authorized by RAM.
// 问题原因和解决方法如下：
// 该调用者缺少允许STS扮演角色的权限策略：请为该调用者添加系统策略（AliyunSTSAssumeRoleAccess）或自定义策略。具体操作，请参见能否指定RAM用户具体可以扮演哪个RAM角色、为RAM用户授权。
// RAM角色的信任策略不包含调用者，即RAM角色不允许该调用者扮演：请为RAM角色添加允许该调用者扮演的信任策略。具体操作，请参见修改RAM角色的信任策略。
//【最佳实践】
// STS Token自颁发后将在一段时间内有效，建议您设置合理的Token有效期，并在有效期内重复使用，以避免业务请求速率上升后，STS Token颁发的速率限制影响到业务。具体速率限制，请参见STS服务调用次数是否有上限。您可以通过请求参数DurationSeconds设置Token有效期。
// 在移动端上传或下载OSS文件等场景下，其访问量较大，即使重复使用STS Token也可能无法满足限流要求。为避免STS的限流成为OSS访问量的瓶颈，您可以尝试OSS的在URL中包含签名的方案。更多信息，请参见在URL中包含签名和服务端签名后直传。
func (api *AliSTSApi) AssumeRole(roleArn, roleSessionName string) (res *sts.AssumeRoleResponse, err error) {
	//构建请求对象。
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	request.RoleArn = roleArn // 要扮演的RAM角色ARN。
	// 该角色是可信实体为阿里云账号类型的RAM角色。更多信息，请参见创建可信实体为阿里云账号的RAM角色或CreateRole。
	// 格式：acs:ram::<account_id>:role/<role_name> 。
	// 您可以通过RAM控制台或API查看角色ARN。具体如下：
	// RAM控制台：请参见查看RAM角色的ARN。
	// API：请参见ListRoles或GetRole。

	request.RoleSessionName = roleSessionName // 角色会话名称。
	// 该参数为用户自定义参数。通常设置为调用该API的用户身份，例如：用户名。在操作审计日志中，即使是同一个RAM角色执行的操作，也可以根据不同的RoleSessionName来区分实际操作者，以实现用户级别的访问审计。
	// 长度为2~64个字符，可包含英文字母、数字、半角句号（.）、at（@）、短划线（-）和下划线（_）。

	//发起请求，并得到响应。
	res, err = api.Client.AssumeRole(request)
	return
}
