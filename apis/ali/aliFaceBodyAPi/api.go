package aliFaceBodyAPi

import (
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	facebody20200910 "github.com/alibabacloud-go/facebody-20200910/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
)

// AliFaceBodyApi 阿里云人脸人体识别
// https://help.aliyun.com/document_detail/146428.html
type AliFaceBodyApi struct {
	Client *facebody20200910.Client
	Config *configStruct.AliRamConfig
}

// NewAliOcrApi 返回ocrAPI管理器
func NewAliOcrApi(config *configStruct.AliRamConfig) (api *AliFaceBodyApi, err error) {

	var client *facebody20200910.Client
	client, err = facebody20200910.NewClient(&openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: &(config.AccessKeyID),
		// 您的AccessKey Secret
		AccessKeySecret: &(config.AccessKeySecret),
		Endpoint:        tea.String("facebody.cn-shanghai.aliyuncs.com"),
	})
	if err != nil {
		return
	}

	api = &AliFaceBodyApi{
		Client: client,
		Config: config,
	}
	return
}

// FinanceLevelIdentifyCheck 金融级身份验证（身份证+人脸照片）
// 人脸核身服务端 https://help.aliyun.com/document_detail/201377.htm?spm=a2c4g.11186623.0.0.2f9b1e0fwC6FUI#doc-api-facebody-ExecuteServerSideVerification
//【应用场景】
// 金融机构网上开户：在金融行业手机App注册开户，通过实人认证，验证开户用户真实身份，降低运营成本，提升风控水平。
// 线上政务平台注册：疫情期间，政府在App中推出线上口罩预约功能，实施实人认证，可有效防止同一人当天重复领取口罩，导致分配不均。
// 修改密码或手机号码：在移动互联网App修改密码，或绑定手机号码时，通过实人认证进行用户真实身份确认。
// 网约车司机认证：网约车当前运营司机身份确认，防止冒用身份驾驶运营车辆。
// 金融风控：在支付、挂失、解冻、转账、取款、信贷、理财等各个环节进行用户身份验证，做好金融风险管控。
//【特色优势】
// 金融级的指标：误识率低于1/100000，准确率高于99%。
// 成熟行业应用：服务超过2亿互联网金融用户，保障超过20亿次交易安全。
// 秒级活体检测：无需复杂交互动作，只需秒级即可完成活体检测，更高效，同时也具备更高级别私密性，更高安全性。
// 通过金融级防攻击测试：抵御各种真实发生的伪造攻击，权威数据源验证。
// 低成本落地方案：纯软件方案，支持普通摄像头，成本极低，适配室内外。
// 国内外权威认证：通过公安部认证、ISO 27001信息安全体系认证，ISO30107-3人脸活体防攻击认证（iBeta PAD Level1）， ISO/TC68。
//【输入限制】
// 图像格式：仅支持JPG格式。
// 图像大小：不超过1 MB。
// 图片分辨率：大于640×480像素，小于2048×2048像素，长宽比小于等于2。
// URL地址中不能包含中文字符。
//func (api *AliFaceBodyApi) FinanceLevelIdentifyCheck(trueName, idNo, imgURL string, imgData []byte) (res *facebody20200910.ExecuteServerSideVerificationResponse, err error) {
func (api *AliFaceBodyApi) FinanceLevelIdentifyCheck(trueName, idNo, imgURL string, imgData []byte) (valid bool, err error) {

	param := &facebody20200910.ExecuteServerSideVerificationRequest{
		CertificateName:   &trueName,            // 真实姓名
		CertificateNumber: &idNo,                // 身份证号
		FacialPictureData: nil,                  // 待比对的图像数据 Base64格式
		FacialPictureUrl:  nil,                  // 待比对的图像URL地址。当前仅支持上海地域的OSS链接，如何生成URL请参见生成URL。
		SceneType:         tea.String("server"), // 场景类型，默认为server。
	}
	if imgURL != "" {
		param.FacialPictureUrl = &imgURL
	} else if imgData != nil {
		param.FacialPictureData = imgData
	} else {
		err = errors.New("图片地址和数据不能同时为空")
		return
	}

	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)

	var res *facebody20200910.ExecuteServerSideVerificationResponse
	res, err = api.Client.ExecuteServerSideVerificationWithOptions(param, headers, runtime)

	log.Println("res", res)
	log.Println("err", err)

	if err != nil {
		return
	}

	if *res.Body.Data.Pass {
		valid = true
		return
	}

	err = errors.New(*res.Body.Message)
	return
}
