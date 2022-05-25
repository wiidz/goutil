package aliOcrApi

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	ocr20191230 "github.com/alibabacloud-go/ocr-20191230/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wiidz/goutil/structs/configStruct"
)

// AliOcrApi 文字识别Api
// 1,000点 = ¥0.01
// 10,000点 = ¥550.00
// 100,000点 = ¥3300.00
// 车牌识别:1点/次
// 驾驶证识别:1点/次
// 行驶证识别:1点/次
// 通用文字识别:1点/次
// 身份证识别:1点/次
// 银行卡识别:1点/次
// 营业执照识别:1点/次
// 二维码识别:0.2点/次
// 增值税发票识别:2.18点/次
// VIN码识别:0.4点/次
// PDF识别:2.18点/次
// 定额发票识别:1点/次
// 增值税发票卷票识别:1点/次
type AliOcrApi struct {
	Client *ocr20191230.Client
	Config *configStruct.AliRamConfig
}

// NewAliOcrApi 返回ocrAPI管理器
func NewAliOcrApi(config *configStruct.AliRamConfig) (api *AliOcrApi, err error) {

	var client *ocr20191230.Client
	client, err = ocr20191230.NewClient(&openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: &(config.AccessKeyID),
		// 您的AccessKey Secret
		AccessKeySecret: &(config.AccessKeySecret),
		Endpoint:        tea.String("ocr.cn-shanghai.aliyuncs.com"),
	})
	if err != nil {
		return
	}

	api = &AliOcrApi{
		Client: client,
		Config: config,
	}
	return
}

// CheckIDCard 身份证识别（1点/次）
// https://help.aliyun.com/document_detail/151899.html
//【功能描述】
// 身份证识别可以识别二代身份证关键字段内容，关键字段包括：姓名、性别、民族、身份证号、出生日期、地址信息、有效起始时间、签发机关，同时可输出身份证区域位置和人脸位置信息。
//【应用场景】
// 远程注册：识别用户提交的身份证内容，自动完成用户身份信息填充。
//【特色优势】
// 识别完整：支持识别身份证内各项内容。
//【输入限制】
// 图像格式：JPEG、JPG、PNG、BMP、GIF。
// 图像大小：不超过3MB。
// 图像分辨率：大于15×15像素，小于4096×4096像素。
// URL地址中不能包含中文字符。
func (api *AliOcrApi) CheckIDCard(imgURL string, side CardSide) (res *ocr20191230.RecognizeIdentityCardResponse, err error) {

	temp := string(side)
	params := &ocr20191230.RecognizeIdentityCardRequest{
		ImageURL: &imgURL, // 图像URL地址。当前仅支持上海地域的OSS链接，如何生成URL请参见生成URL。
		Side:     &temp,   // 身份证正反面类型。 face：人像面。 back：国徽面。
	}

	runtime := &util.RuntimeOptions{} // 一些配置，暂时用不到
	res, err = api.Client.RecognizeIdentityCardWithOptions(params, runtime)
	return
}

// CheckDrivingLicense 行驶证识别(1点/次)
//【功能描述】
// 行驶证识别能力可以识别行驶证首页和副页关键字段内容，输出品牌型号、车辆类型、车牌号码、检验记录、核定载质量、核定载人数等21个关键字段信息。
//【前提条件】
//请确保您已开通文字识别服务，若未开通服务请立即开通。
//【输入限制】
// 图像格式：JPEG、JPG、PNG、BMP、GIF。
// 图像大小：不超过3 MB。
// 图像分辨率：不限制图片分辨率，但图片分辨率太高可能会导致API识别超时，超时时间为5秒。
// URL地址中不能包含中文字符。
func (api *AliOcrApi) CheckDrivingLicense(imgURL string, side CardSide) (res *ocr20191230.RecognizeDrivingLicenseResponse, err error) {
	temp := string(side)
	params := &ocr20191230.RecognizeDrivingLicenseRequest{
		ImageURL: &imgURL, // 图像URL地址。当前仅支持上海地域的OSS链接，如何生成URL请参见生成URL。
		Side:     &temp,   // 身份证正反面类型。 face：人像面。 back：国徽面。
	}

	//runtime := &util.RuntimeOptions{} // 一些配置，暂时用不到
	//res, err = api.Client.RecognizeDrivingLicenseWithOptions(params, runtime)
	res, err = api.Client.RecognizeDrivingLicense(params)
	return
}

// CheckDriverLicense 驾驶证识别(1点/次)
//【功能描述】
// 驾驶证识别能力可以识别驾驶证首页和副页关键字段内容，包括：档案编号、姓名、有效期时长、性别、发证日期、驾驶证号、驾驶证准驾车型、有效期开始时间、地址，共9个关键字段信息。
//【前提条件】
// 请确保您已开通文字识别服务，若未开通服务请立即开通。
//【输入限制】
// 图像格式：JPEG、JPG、PNG、BMP、GIF。
// 图像大小：不超过4 MB。
// 图像分辨率：大于15×15像素，小于4096×4096像素。
// URL地址中不能包含中文字符。
func (api *AliOcrApi) CheckDriverLicense(imgURL string, side CardSide) (res *ocr20191230.RecognizeDriverLicenseResponse, err error) {
	temp := string(side)
	params := &ocr20191230.RecognizeDriverLicenseRequest{
		ImageURL: &imgURL, // 图像URL地址。当前仅支持上海地域的OSS链接，如何生成URL请参见生成URL。
		Side:     &temp,   // 身份证正反面类型。 face：人像面。 back：国徽面。
	}

	//runtime := &util.RuntimeOptions{} // 一些配置，暂时用不到
	//res, err = api.Client.RecognizeDriverLicenseWithOptions(params, runtime)
	res, err = api.Client.RecognizeDriverLicense(params)
	return
}

// CheckLicensePlate 车牌识别
//【功能描述】
// 车牌识别能力可以准确识别出图像中车牌位置，输出车牌位置坐标、车牌类型、车牌号码、车牌号码置信度、车牌置信度，共5个关键字段信息。
//【前提条件】
// 请确保您已开通文字识别服务，若未开通服务请立即开通。
//【输入限制】
// 图像格式：JPEG、JPG、PNG、BMP、GIF。
// 图像大小：不超过4 MB。
// 图像分辨率：大于15×15像素，小于4096×4096像素。
// URL地址中不能包含中文字符。
func (api *AliOcrApi) CheckLicensePlate(imgURL string) (res *ocr20191230.RecognizeLicensePlateResponse, err error) {
	params := &ocr20191230.RecognizeLicensePlateRequest{
		ImageURL: &imgURL, // 图像URL地址。当前仅支持上海地域的OSS链接，如何生成URL请参见生成URL。
	}

	//runtime := &util.RuntimeOptions{} // 一些配置，暂时用不到
	//res, err = api.Client.RecognizeLicensePlateWithOptions(params, runtime)
	res, err = api.Client.RecognizeLicensePlate(params)
	return
}
