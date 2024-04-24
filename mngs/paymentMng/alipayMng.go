package paymentMng

import (
	"context"
	"errors"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/xlog"
	"github.com/wiidz/goutil/helpers/strHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
)

type AliPayMng struct {
	Client *alipay.Client
	Config *configStruct.AliPayConfig
}

// getAliPayInstance 获取实例
func getAliPayInstance(config *configStruct.AliPayConfig) (*AliPayMng, error) {

	client, err := alipay.NewClient(config.AppID, config.PrivateKey, !config.Debug)

	if err != nil {
		return nil, err
	}

	var alipayMng = &AliPayMng{
		Config: config,
		Client: client,
	}

	//配置公共参数
	alipayMng.Client.SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetNotifyUrl(config.NotifyURL).
		SetCertSnByContent([]byte(config.AppCertPublicKey), []byte(config.RootCert), []byte(config.CertPublicKey))

	// 打开Debug开关，输出日志，默认关闭
	alipayMng.Client.DebugSwitch = gopay.DebugOn

	// 设置支付宝请求 公共参数
	//    注意：具体设置哪些参数，根据不同的方法而不同，此处列举出所有设置参数
	//client.SetLocation().                       // 设置时区，不设置或出错均为默认服务器时间
	//	SetPrivateKeyType().                    // 设置 支付宝 私钥类型，alipay.PKCS1 或 alipay.PKCS8，默认 PKCS1
	//	SetAliPayRootCertSN().                  // 设置支付宝根证书SN，通过 alipay.GetRootCertSN() 获取
	//	SetAppCertSN().                         // 设置应用公钥证书SN，通过 alipay.GetCertSN() 获取
	//	SetAliPayPublicCertSN().                // 设置支付宝公钥证书SN，通过 alipay.GetCertSN() 获取
	//	SetCharset("utf-8").                    // 设置字符编码，不设置默认 utf-8
	//	SetSignType(alipay.RSA2).               // 设置签名类型，不设置默认 RSA2
	//	SetReturnUrl("https://www.fmm.ink").    // 设置返回URL
	//	SetNotifyUrl("https://www.fmm.ink").    // 设置异步通知URL
	//	SetAppAuthToken()                       // 设置第三方应用授权

	// 自动同步验签（只支持证书模式）
	// 传入 alipayCertPublicKey_RSA2.crt 内容
	//client.AutoVerifySign("alipayCertPublicKey_RSA2 bytes")

	// 证书路径
	//err := alipay.Client.SetCertSnByPath("appCertPublicKey.crt", "alipayRootCert.crt", "alipayCertPublicKey_RSA2.crt")
	// 证书内容
	//err := client.SetCertSnByContent("appCertPublicKey bytes", "alipayRootCert bytes", "alipayCertPublicKey_RSA2 bytes")
	return alipayMng, nil
}

// NewAliPayMng 根据传入的config生成管理器
func NewAliPayMng(config *configStruct.AliPayConfig) (*AliPayMng, error) {
	return getAliPayInstance(config)
}

// WapPay H5支付
func (aliPayMng *AliPayMng) WapPay(ctx context.Context, params *UnifiedOrderParam) (payURL string, err error) {

	aliPayMng.Client.SetReturnUrl(params.ReturnURL)

	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", params.Title)
	body.Set("out_trade_no", params.OutTradeNo)
	body.Set("quit_url", params.ReturnURL)
	body.Set("total_amount", typeHelper.Float64ToStr(params.TotalAmount)) // 元为单位
	body.Set("product_code", "QUICK_WAP_WAY")

	//手机网站支付请求
	payURL, err = aliPayMng.Client.TradeWapPay(ctx, body)
	if err != nil {
		xlog.Error("err:", err)
		temp := typeHelper.JsonDecodeMap(err.Error())
		err = errors.New(temp["sub_msg"].(string))
		return
	}
	xlog.Debug("payUrl:", payURL)
	return
}

// Refund 退款
func (aliPayMng *AliPayMng) Refund(ctx context.Context, param *RefundParam) (err error) {

	//请求参数
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", param.OutTradeNo)
	body.Set("refund_amount", param.RefundAmount)
	body.Set("out_request_no", param.OrderRefundNo) // 退款单号
	body.Set("refund_reason", param.Reason)

	//发起退款请求
	aliRsp, err := aliPayMng.Client.TradeRefund(ctx, body)
	//aliPayMng.Client.FundTransRefund()
	if err != nil {
		xlog.Error("err:", err)
		temp := typeHelper.JsonDecodeMap(err.Error())
		err = errors.New(temp["sub_msg"].(string))
		return
	}
	xlog.Debug("aliRsp:", *aliRsp)

	return err
}

// GetRefund 查询退款情况
func (aliPayMng *AliPayMng) GetRefund(ctx context.Context, outTradeNo, orderRefundNo string) (resp *alipay.TradeFastpayRefundQueryResponse, err error) {
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", outTradeNo)
	body.Set("out_request_no", orderRefundNo)
	//发起退款查询请求
	aliRsp, err := aliPayMng.Client.TradeFastPayRefundQuery(ctx, body)
	if err != nil {
		xlog.Error("err:", err)
		return
	}
	xlog.Debug("aliRsp:", *aliRsp)
	return aliRsp, err
}

// ScanPay 扫码支付
func (aliPayMng *AliPayMng) ScanPay(ctx context.Context, params *ScanPayParam) (resp *alipay.TradePayResponse, err error) {

	var scene = "bar_code"
	if strHelper.Exist(params.AuthCode, "fp") {
		scene = "security_code"
	}

	body := make(gopay.BodyMap)
	body.Set("out_trade_no", params.OutTradeNo)  //【是-String(64)】商户订单号。 由商家自定义，64个字符以内，仅支持字母、数字、下划线且需保证在商户端不重复。
	body.Set("total_amount", params.TotalAmount) //【是-Price(11)】订单总金额。 单位为元，精确到小数点后两位，取值范围：[0.01,100000000] 。
	body.Set("subject", params.Title)            //【是-String(256)】订单标题。 注意：不可使用特殊字符，如 /，=，& 等
	body.Set("auth_code", params.AuthCode)       //【是-String(64)】  支付授权码。 当面付场景传买家的付款码（25~30开头的长度为16~24位的数字，实际字符串长度以开发者获取的付款码长度为准）或者刷脸标识串（fp开头的35位字符串）。
	body.Set("scene", scene)                     //【是-String(64)】支付场景。 枚举值： bar_code：当面付条码支付场景； security_code：当面付刷脸支付场景，对应的auth_code为fp开头的刷脸标识串； 默认值为bar_code。

	//body.Set("product_code", "") //【否-String(64)】产品码。 商家和支付宝签约的产品码。当面付场景下，如果签约的是当面付快捷版，则传 OFFLINE_PAYMENT;其它支付宝当面付产品传 FACE_TO_FACE_PAYMENT；不传则默认使用FACE_TO_FACE_PAYMENT。
	//body.Set("seller_id", "")    //【否-String(28)】卖家支付宝用户ID。 当需要指定收款账号时，通过该参数传入，如果该值为空，则默认为商户签约账号对应的支付宝用户ID。 收款账号优先级规则：门店绑定的收款账户>请求传入的seller_id>商户签约账号对应的支付宝用户ID； 注：直付通和机构间联场景下seller_id无需传入或者保持跟pid一致；如果传入的seller_id与pid不一致，需要联系支付宝小二配置收款关系；
	//body.Set("goods_detail", "")                 //【否-String】订单包含的商品列表信息，json格式。
	//body.Set("extend_params", params.Attach) //【否-String】业务扩展参数
	//body.Set("promo_params", params.Attach)  //【否-String】优惠明细参数，通过此属性补充营销参数。 注：仅与支付宝协商后可用。
	//body.Set("store_id", params.DeviceNo) //【否-String(32)】商户门店编号。 指商户创建门店时输入的门店编号。
	//body.Set("operator_id", "")           //【否-String(28)】商户操作员编号。
	//body.Set("query_options", "")        //【否-String(1024)】返回参数选项。商户通过传递该参数来定制同步需要额外返回的信息字段，数组格式。如：["fund_bill_list","voucher_detail_list","discount_goods_detail"]
	resp, err = aliPayMng.Client.TradePay(ctx, body)

	if err != nil {
		xlog.Error("err:", err)
		return
	}
	xlog.Debug("aliRsp:", *resp)
	return
}

// TradeQuery 查询订单详情
func (aliPayMng *AliPayMng) TradeQuery(ctx context.Context, tradeNo, outTradeNo string) (resp *alipay.TradeQueryResponse, err error) {
	body := make(gopay.BodyMap)
	//body.Set("app_id", aliPayMng.Config.AppID) //【是-String(32)】支付宝分配给开发者的应用ID
	//body.Set("charset", alipay.UTF8)           //【是-String(10)】请求使用的编码格式，如utf-8,gbk,gb2312等
	body.Set("trade_no", tradeNo)        //【二者取一-String(64)】支付宝交易号，和商户订单号不能同时为空
	body.Set("out_trade_no", outTradeNo) //【二者取一-String(64)】订单支付时传入的商户订单号,和支付宝交易号不能同时为空。 trade_no,out_trade_no如果同时存在优先取trade_no

	//body.Set("query_options", "") //【否-String(1024)】查询选项，商户通过上送该参数来定制同步需要额外返回的信息字段，数组格式。支持枚举如下：

	// fund_bill_list：交易支付使用的资金渠道；
	// voucher_detail_list：交易支付时使用的所有优惠券信息；
	// discount_goods_detail：交易支付所使用的单品券优惠的商品优惠信息；
	// mdiscount_amount：商家优惠金额；

	resp, err = aliPayMng.Client.TradeQuery(ctx, body)
	if err != nil {
		xlog.Error("err:", err)
		return
	}
	xlog.Debug("aliRsp:", *resp)
	return
}
