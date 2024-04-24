package paymentMng

import (
	"context"
	"errors"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
	"math"
	"net/http"
	"time"
)

// 欢迎使用微信支付！
// 附件中的三份文件（证书pkcs12格式、证书pem格式、证书密钥pem格式）,为接口中强制要求时需携带的证书文件。
// 证书属于敏感信息，请妥善保管不要泄露和被他人复制。
// 不同开发语言下的证书格式不同，以下为说明指引：

// 证书pkcs12格式（apiclient_cert.p12）
// 包含了私钥信息的证书文件，为p12(pfx)格式，由微信支付签发给您用来标识和界定您的身份
// 部分安全性要求较高的API需要使用该证书来确认您的调用身份
// windows上可以直接双击导入系统，导入过程中会提示输入证书密码，证书密码默认为您的商户号（如：1900006031）

// 证书pem格式（apiclient_cert.pem）
// 从apiclient_cert.p12中导出证书部分的文件，为pem格式，请妥善保管不要泄漏和被他人复制
// 部分开发语言和环境，不能直接使用p12文件，而需要使用pem，所以为了方便您使用，已为您直接提供
// 您也可以使用openssl命令来自己导出：openssl pkcs12 -clcerts -nokeys -in apiclient_cert.p12 -out apiclient_cert.pem

// 证书密钥pem格式（apiclient_key.pem）
// 从apiclient_cert.p12中导出密钥部分的文件，为pem格式
// 部分开发语言和环境，不能直接使用p12文件，而需要使用pem，所以为了方便您使用，已为您直接提供
// 您也可以使用openssl命令来自己导出：openssl pkcs12 -nocerts -in apiclient_cert.p12 -out apiclient_key.pem
// 备注说明：
// 由于绝大部分操作系统已内置了微信支付服务器证书的根CA证书,  2018年3月6日后, 不再提供CA证书文件（rootca.pem）下载

// 注意 原始文件并没有给public
// 我们手动从给的 私钥 apiclient_key.pem 中生成 公钥
// openssl rsa -in apiclient_key.pem -pubout -out apiclient_key_public.pem
// https://blog.csdn.net/u011580177/article/details/106222865

// 以上信息不用看了，不用我们去维护公钥！！！

type WechatPayMngV3 struct {
	Config *configStruct.WechatPayConfigV3
	Client *wechat.ClientV3
}

// getWechatPayV3Instance 获取微信支付V3实例
func getWechatPayV3Instance(config *configStruct.WechatPayConfigV3) (mng *WechatPayMngV3, err error) {

	var client *wechat.ClientV3
	client, err = wechat.NewClientV3(config.MchID, config.CertSerialNo, config.ApiKeyV3, config.PEMPrivateKeyContent)
	if err != nil {
		return
	}
	client.DebugSwitch = gopay.DebugOn // 打开Debug开关，输出请求日志，默认关闭
	client.AutoVerifySign()            // 启用自动同步返回验签，并定时更新微信平台API证书

	mng = &WechatPayMngV3{
		Config: config,
		Client: client,
	}

	return
}

// NewWechatPayMngV3 根据传入的config，获取信的微信支付
func NewWechatPayMngV3(config *configStruct.WechatPayConfigV3) (*WechatPayMngV3, error) {
	return getWechatPayV3Instance(config)
}

// Mini 小程序场景下单，注意params.TotalAmount是元为单位
func (mng *WechatPayMngV3) Mini(ctx context.Context, params *UnifiedOrderParam, openID string) (timestampStr, packageStr, nonceStr, paySign string, err error) {

	//【1】获取prepayID
	prepayRes, err := mng.jsApiPlaceOrder(ctx, params, openID)
	if err != nil {
		return
	}

	//【2】获取签名
	var applet *wechat.AppletParams
	applet, err = mng.Client.PaySignOfApplet(mng.Config.AppID, prepayRes.Response.PrepayId)
	if err != nil {
		return
	}

	timestampStr = applet.TimeStamp
	packageStr = "prepay_id=" + prepayRes.Response.PrepayId
	paySign = applet.PaySign
	nonceStr = applet.NonceStr
	return
}

// Js 公众号支付，注意params.TotalAmount是元为单位
func (mng *WechatPayMngV3) Js(ctx context.Context, params *UnifiedOrderParam, openID string) (appID, timestampStr, nonceStr, packageStr, paySign, signType string, err error) {

	//【1】获取prepayID
	prepayRes, err := mng.jsApiPlaceOrder(ctx, params, openID)
	if err != nil {
		return
	}

	//【2】获取签名
	var jsapi *wechat.JSAPIPayParams
	jsapi, err = mng.Client.PaySignOfJSAPI(mng.Config.AppID, prepayRes.Response.PrepayId)
	if err != nil {
		return
	}

	appID = jsapi.AppId
	timestampStr = jsapi.TimeStamp
	nonceStr = jsapi.NonceStr
	packageStr = jsapi.Package
	paySign = jsapi.PaySign
	signType = jsapi.SignType
	return
}

// H5 网页支付，注意params.TotalAmount是元为单位
func (mng *WechatPayMngV3) H5(ctx context.Context, params *UnifiedOrderParam, openID string) (H5Url string, err error) {

	//【1】构建结构体
	totalFee := int(math.Round(params.TotalAmount * 100)) // 分为单位，注意这里的近似
	bm := gopay.BodyMap{}
	bm.Set("appid", mng.Config.AppID).
		Set("mchid", mng.Config.MchID).
		Set("description", params.Title).
		Set("out_trade_no", params.OutTradeNo).
		Set("notify_url", mng.Config.NotifyURL).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", totalFee).
				Set("currency", "CNY")
		}).
		SetBodyMap("scene_info", func(bm gopay.BodyMap) {
			bm.Set("payer_client_ip", params.IP).
				SetBodyMap("h5_info", func(bm gopay.BodyMap) {
					bm.Set("type", "Wap")
				})
		})

	//【2】获取签名
	var wxRsp *wechat.H5Rsp
	wxRsp, err = mng.Client.V3TransactionH5(ctx, bm)
	if err != nil {
		return
	}
	if wxRsp.Code != 0 {
		wechatErr := WechatError{}
		err = typeHelper.JsonDecodeWithStruct(wxRsp.Error, &wechatErr)
		if err == nil {
			err = errors.New(wechatErr.Message)
		} else {
			err = errors.New(wxRsp.Error)
		}
	}

	return wxRsp.Response.H5Url, err
}

// Refund 退款
func (mng *WechatPayMngV3) Refund(ctx context.Context, param *RefundParam) (wxRsp *wechat.RefundRsp, err error) {
	refundFee := int(param.RefundAmount * 100)
	totalFee := int(param.TotalAmount * 100)
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", param.OutTradeNo).
		Set("out_refund_no", param.OrderRefundNo).
		Set("reason", param.Reason).
		Set("notify_url", mng.Config.RefundNotifyURL).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("refund", refundFee).
				Set("total", totalFee).
				Set("currency", "CNY")
		})

	wxRsp, err = mng.Client.V3Refund(ctx, bm)
	if err != nil {
		return
	}
	if wxRsp.Code != 0 {
		wechatErr := WechatError{}
		err = typeHelper.JsonDecodeWithStruct(wxRsp.Error, &wechatErr)
		if err == nil {
			err = errors.New(wechatErr.Message)
		} else {
			err = errors.New(wxRsp.Error)
		}
	}
	return
}

// jsApiPlaceOrder JSAPI/小程序下单API totalFee 是分为单位
func (mng *WechatPayMngV3) jsApiPlaceOrder(ctx context.Context, params *UnifiedOrderParam, openID string) (wxRsp *wechat.PrepayRsp, err error) {

	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	//totalFee := int(params.TotalAmount * 100) // 分为单位
	totalFee := int(math.Round(params.TotalAmount * 100)) // 分为单位，注意这里的近似

	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("appid", mng.Config.AppID).
		Set("mchid", mng.Config.MchID).
		Set("description", params.Title).
		Set("out_trade_no", params.OutTradeNo).
		Set("time_expire", expire).
		Set("notify_url", mng.Config.NotifyURL).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", totalFee).
				Set("currency", "CNY")
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", openID)
		})

	wxRsp, err = mng.Client.V3TransactionJsapi(ctx, bm)

	if err != nil {
		xlog.Error(err)
		return
	}

	if wxRsp.Code != 0 {
		wechatErr := WechatError{}
		err = typeHelper.JsonDecodeWithStruct(wxRsp.Error, &wechatErr)
		if err == nil {
			err = errors.New(wechatErr.Message)
		} else {
			err = errors.New(wxRsp.Error)
		}
	}

	return
}

// NotifyPayment 一般支付回调
func (mng *WechatPayMngV3) NotifyPayment(req *http.Request) (res *wechat.V3DecryptResult, err error) {

	var notifyReq *wechat.V3NotifyReq
	notifyReq, err = wechat.V3ParseNotify(req)
	if err != nil {
		xlog.Error(err)
		return
	}

	// 普通支付通知解密
	res, err = notifyReq.DecryptCipherText(mng.Config.ApiKeyV3)
	return
}

// NotifyRefund 退款回调
func (mng *WechatPayMngV3) NotifyRefund(req *http.Request) (res *wechat.V3DecryptRefundResult, err error) {

	var notifyReq *wechat.V3NotifyReq
	notifyReq, err = wechat.V3ParseNotify(req)
	if err != nil {
		xlog.Error(err)
		return
	}

	// 普通支付通知解密
	res, err = notifyReq.DecryptRefundCipherText(mng.Config.ApiKeyV3)
	return
}

// BatchPayUser 批量付款给用户（用户的真实姓名要么都填，要么都不填，大于2000必填）
func (mng *WechatPayMngV3) BatchPayUser(ctx context.Context, params *TransferUserParam, transferList []*TransferUserDetailList) (res *wechat.TransferRsp, err error) {

	// 【1】为名称加密
	for k := range transferList {
		if transferList[k].UserName == "" {
			continue
		}
		//transferList[k].UserName, err = wechat.V3EncryptText(transferList[k].UserName, []byte(mng.Config.PEMPublicKeyContent)) // 不用我们去维护公钥！！！
		transferList[k].UserName, err = mng.Client.V3EncryptText(transferList[k].UserName)
		if err != nil {
			return
		}
	}

	// 初始化参数结构体
	bm := make(gopay.BodyMap)
	bm.Set("appid", mng.Config.AppID). // 直连商户的appid，申请商户号的appid或商户号绑定的appid（企业号corpid即为此appid）
						Set("out_batch_no", params.OutBatchNo).
						Set("batch_name", params.BatchName).
						Set("batch_remark", params.BatchRemark).
						Set("total_amount", params.TotalAmount).
						Set("total_num", params.TotalNum).
						Set("transfer_detail_list", transferList).
						Set("notify_url", mng.Config.NotifyURL)

	//bm.Set("nonce_str", util.RandomString(32)).
	//	Set("partner_trade_no", util.RandomString(32)).
	//	Set("openid", "o0Df70H2Q0fY8JXh1aFPIRyOBgu8").
	//	Set("check_name", "FORCE_CHECK"). // NO_CHECK：不校验真实姓名 , FORCE_CHECK：强校验真实姓名
	//	Set("re_user_name", "付明明"). // 收款用户真实姓名。 如果check_name设置为FORCE_CHECK，则必填用户真实姓名
	//	Set("amount", 30). // 企业付款金额，单位为分
	//	Set("desc", "测试转账"). // 企业付款备注，必填。注意：备注中的敏感词会被转成字符*
	//	Set("spbill_create_ip", "127.0.0.1")

	// 企业向微信用户个人付款（不支持沙箱环境）
	//    body：参数Body
	res, err = mng.Client.V3Transfer(ctx, bm)
	xlog.Debug("Response：", res)
	xlog.Debug("err：", err)

	if err != nil {
		return
	}

	if res.Code != 0 {
		wechatErr := WechatError{}
		err = typeHelper.JsonDecodeWithStruct(res.Error, &wechatErr)
		if err == nil {
			err = errors.New(wechatErr.Message)
		} else {
			err = errors.New(res.Error)
		}
	}

	return
}

// ScanPay 扫用户付款码收款（腾讯那边 V3没有完成，暂时用V2）
func (mng *WechatPayMngV3) ScanPay(ctx context.Context) (err error) {
	return
}

// TransactionQueryOrder 查询订单
func (mng *WechatPayMngV3) TransactionQueryOrder(ctx context.Context, transactionID, outTradeNo string) (res *wechat.QueryOrderRsp, err error) {

	if transactionID != "" {
		res, err = mng.Client.V3TransactionQueryOrder(ctx, wechat.TransactionId, transactionID)
	} else if outTradeNo != "" {
		res, err = mng.Client.V3TransactionQueryOrder(ctx, wechat.OutTradeNo, outTradeNo)
	}

	log.Println("res", res)
	log.Println("res.Code", res.Code)
	log.Println("res.Response", res.Response)
	log.Println("res.SignInfo", res.SignInfo)

	if res.Code != 0 {
		wechatErr := WechatError{}
		err = typeHelper.JsonDecodeWithStruct(res.Error, &wechatErr)
		if err == nil {
			err = errors.New(wechatErr.Message)
		} else {
			err = errors.New(res.Error)
		}
	}
	if res.Response.TradeState != "SUCCESS" {
		err = errors.New(res.Response.TradeStateDesc)
	}

	return
}

type DetailStatus string

const (
	All      DetailStatus = "ALL"
	WAIT_PAY DetailStatus = "WAIT_PAY"
	SUCCESS  DetailStatus = "SUCCESS"
	FAIL     DetailStatus = "FAIL"
)

func (status DetailStatus) String() string {
	switch status {
	case All:
		return "ALL"
	case SUCCESS:
		return "SUCCESS"
	case FAIL:
		return "FAIL"
	case WAIT_PAY:
		return "WAIT_PAY"
	default:
		return "ALL"
	}
}

type TransferMerchantError struct {
	Code   string `json:"code"`
	Detail struct {
		Location string `json:"location"`
		Value    int    `json:"value"`
	} `json:"detail"`
	Message string `json:"message"`
}

// TransferMerchantQuery 商家转账到零钱 查询转账批次
func (mng *WechatPayMngV3) TransferMerchantQuery(ctx context.Context, outBatchNo string, offset, limit int, detailStatus DetailStatus) (res *wechat.TransferMerchantQueryRsp, err error) {

	//【need_query_detail】：boolean,query枚举值： true：是；false：否，默认否。商户可选择是否查询指定状态的转账明细单，当转账批次单状态为“FINISHED”（已完成）时，才会返回满足条件的转账明细单。示例值：true
	//【offset】：int,请求资源起始位置，该次请求资源（转账明细单）的起始位置，从0开始，默认值为0，示例值：1
	//【limit】：int,该次请求可返回的最大资源（转账明细单）条数，最小20条，最大100条，不传则默认20条。不足20条按实际条数返回,示例值：20
	//【detail_status】：string[1,32]，查询指定状态的转账明细单，当need_query_detail为true时，该字段必填，ALL：全部。需要同时查询转账成功和转账失败的明细单 ，SUCCESS：转账成功。只查询转账成功的明细单 ，FAIL：转账失败。需要通过查询明细单接口确认明细失败原因后，再决定是否重新发起对该笔明细单的转账（并非整个转账批次单）

	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("need_query_detail", true).
		Set("offset", offset).
		Set("limit", limit).
		Set("detail_status", detailStatus.String())

	res, err = mng.Client.V3TransferMerchantQuery(ctx, outBatchNo, bm)
	if err != nil {
		return
	}
	err = mng.handleError(res.Error)

	return
}

// handleError 处理错误信息
func (mng *WechatPayMngV3) handleError(errStr string) (err error) {
	if errStr != "" {
		// {"code":"PARAM_ERROR","detail":{"location":"query","value":1},"message":"输入源“/query/limit”映射到数值字段“最大资源条数”规则校验失败，值低于最小值 20"}
		var errObj TransferMerchantError
		parseErr := typeHelper.JsonDecodeWithStruct(errStr, &errObj)
		if parseErr != nil {
			// 解析失败
			err = errors.New(errStr)
			return
		}

		err = errors.New(errObj.Message)
	}
	return
}

// V3TransferMerchantDetail 商家明细单号查询明细单API
func (mng *WechatPayMngV3) V3TransferMerchantDetail(ctx context.Context, outBatchNo, outDetailNo string) (res *wechat.TransferMerchantDetailRsp, err error) {
	res, err = mng.Client.V3TransferMerchantDetail(ctx, outBatchNo, outDetailNo)
	if err != nil {
		return
	}
	err = mng.handleError(res.Error)

	return
}
