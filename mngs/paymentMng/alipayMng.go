package paymentMng

import (
	"errors"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/configMng"
)

type AliPayMng struct {
	Client *alipay.Client
	Config *configMng.AliPayConfig
}

// getAliPayInstance 获取实例
func getAliPayInstance(config *configMng.AliPayConfig) *AliPayMng {

	var alipayMng = AliPayMng{
		Config: config,
		Client:alipay.NewClient(config.AppID, config.PrivateKey, config.IsProd),
	}

	//配置公共参数
	alipayMng.Client.SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetPrivateKeyType(alipay.PKCS1).
		SetNotifyUrl(config.NotifyURL)

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
	return &alipayMng
}

// NewAliPayMngSingle 根据configs里的配置文件生成单例
func NewAliPayMngSingle() *AliPayMng {
	config := configMng.GetAliPay()
	return getAliPayInstance(config)
}

// NewAliPayMng 根据传入的config生成管理器
func NewAliPayMng(config *configMng.AliPayConfig) *AliPayMng {
	return getAliPayInstance(config)
}

// WapPay H5支付
func (aliPayMng *AliPayMng) WapPay(params *UnifiedOrderParam) (payURL string, err error) {

	aliPayMng.Client.SetReturnUrl(params.ReturnURL)

	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", params.Title)
	body.Set("out_trade_no", params.OutTradeNo)
	body.Set("quit_url", params.ReturnURL)
	body.Set("total_amount", typeHelper.Float64ToStr(params.TotalAmount)) // 元为单位
	body.Set("product_code", "QUICK_WAP_WAY")

	//手机网站支付请求
	payURL, err = aliPayMng.Client.TradeWapPay(body)
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
func (aliPayMng *AliPayMng) Refund(param *RefundParam) (err error) {

	//请求参数
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", param.OutTradeNo)
	body.Set("refund_amount", param.RefundAmount)
	body.Set("out_request_no", param.OrderRefundNo) // 退款单号
	body.Set("refund_reason", param.Reason)

	//发起退款请求
	aliRsp, err := aliPayMng.Client.TradeRefund(body)
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
func (aliPayMng *AliPayMng) GetRefund(outTradeNo, orderRefundNo string) (resp *alipay.TradeFastpayRefundQueryResponse, err error) {
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", outTradeNo)
	body.Set("out_request_no", orderRefundNo)
	//发起退款查询请求
	aliRsp, err := aliPayMng.Client.TradeFastPayRefundQuery(body)
	if err != nil {
		xlog.Error("err:", err)
		return
	}
	xlog.Debug("aliRsp:", *aliRsp)
	return aliRsp, err
}
