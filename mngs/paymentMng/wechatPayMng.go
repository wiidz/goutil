package paymentMng

import (
	"errors"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/util"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/go-pay/gopay/wechat"
	"github.com/wiidz/goutil/mngs/configMng"
	"log"
	"strconv"
	"time"
)

type WechatPayMng struct {
	Config *configMng.WechatPayConfig
	Client *wechat.Client
}

// getWechatPayInstance 获取微信支付实例
func getWechatPayInstance(config *configMng.WechatPayConfig) *WechatPayMng {

	var wechatPayMng = WechatPayMng{
		Config: config,
		Client: wechat.NewClient(config.AppID, config.MchID, config.PayKey, config.IsProd),
	}

	// 打开Debug开关，输出请求日志，默认关闭
	wechatPayMng.Client.DebugSwitch = gopay.DebugOn

	// 设置国家：不设置默认 中国国内
	//    wechat.China：中国国内
	//    wechat.China2：中国国内备用
	//    wechat.SoutheastAsia：东南亚
	//    wechat.Other：其他国家
	wechatPayMng.Client.SetCountry(wechat.China)

	//certFilePath := "./configs/certs/"
	//keyPath:= ""

	//_ = wechatPayMng.Client.AddCertPemFilePath(Config.CertPath,Config.CertKeyPath)
	_ = wechatPayMng.Client.AddCertPkcs12FileContent([]byte(config.CertFileContent))
	//_ = wechatPayMng.Client.AddCertPkcs12FilePath(config.CertPath)

	// 添加微信pem证书
	return &wechatPayMng
}

// NewWechatPayMngSingle 从本地configs里的配置，生成单例
func NewWechatPayMngSingle() *WechatPayMng {
	config := configMng.GetWechatPay()
	return getWechatPayInstance(config)
}

// NewWechatPayMng 根据传入的config，获取信的微信支付
func NewWechatPayMng(config *configMng.WechatPayConfig) *WechatPayMng {
	return getWechatPayInstance(config)
}

// UnifiedOrder 统一下单获取paysign totalFee 是分为单位
func (mng *WechatPayMng) UnifiedOrder(param *UnifiedOrderParam) (mWebUrl string, err error) {
	//初始化参数Map
	totalFee := param.TotalAmount * 100 // 分为单位

	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", util.GetRandomString(32)).
		Set("body", param.Title).
		Set("out_trade_no", param.OutTradeNo).
		Set("total_fee", totalFee).
		Set("spbill_create_ip", param.IP).
		Set("notify_url", mng.Config.NotifyURL).
		Set("trade_type", wechat.TradeType_H5).
		Set("device_info", "WEB").
		Set("sign_type", wechat.SignType_MD5).
		SetBodyMap("scene_info", func(bm gopay.BodyMap) {
			bm.SetBodyMap("h5_info", func(bm gopay.BodyMap) {
				bm.Set("type", "Wap")
				bm.Set("wap_url", param.ReturnURL)
				bm.Set("wap_name", param.AppName)
			})
		})

	//请求支付下单，成功后得到结果
	var wxRsp *wechat.UnifiedOrderResponse
	wxRsp, err = mng.Client.UnifiedOrder(bm)
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Debug("Response：", wxRsp)
	xlog.Debug("wxRsp.MwebUrl:", wxRsp.MwebUrl)
	xlog.Debug("wxRsp.ResultCode:", wxRsp.ResultCode)
	xlog.Debug("wxRsp.ReturnCode:", wxRsp.ReturnCode)
	xlog.Debug("wxRsp.ReturnMsg:", wxRsp.ReturnMsg)
	xlog.Debug("wxRsp.MwebUrl:", wxRsp.ErrCode)
	xlog.Debug("wxRsp.MwebUrl:", wxRsp.ErrCodeDes)

	if len(wxRsp.ErrCodeDes) != 0 {
		err = errors.New(wxRsp.ErrCodeDes)
		return
	}

	mWebUrl = wxRsp.MwebUrl
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	log.Println("timeStamp", timeStamp)
	return
}

// UnifiedOrderJs 统一下单获取paysign totalFee 是分为单位
func (mng *WechatPayMng) UnifiedOrderJs(param *UnifiedOrderParam,openID string) (data map[string]interface{}, err error) {

	totalFee := param.TotalAmount * 100 // 分为单位

	//初始化参数Map
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", util.GetRandomString(32)).
		Set("body", param.Title).
		Set("out_trade_no", param.OutTradeNo).
		Set("total_fee", totalFee).
		Set("spbill_create_ip", param.IP).
		Set("notify_url", mng.Config.NotifyURL).
		Set("trade_type", wechat.TradeType_JsApi).
		Set("device_info", "WEB").
		Set("sign_type", wechat.SignType_MD5).
		SetBodyMap("scene_info", func(bm gopay.BodyMap) {
			bm.SetBodyMap("h5_info", func(bm gopay.BodyMap) {
				bm.Set("type", "Wap")
				bm.Set("wap_url", param.ReturnURL)
				bm.Set("wap_name", param.AppName)
			})
		}).Set("openid", openID) //js支付必填

	//请求支付下单，成功后得到结果
	wxRsp, err := mng.Client.UnifiedOrder(bm)
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Debug("Response：", wxRsp)
	xlog.Debug("wxRsp.MwebUrl:", wxRsp.MwebUrl)

	//获取Jsapi支付需要的paySign
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)

	pac := "prepay_id=" + wxRsp.PrepayId
	paySign := wechat.GetJsapiPaySign(mng.Config.AppID, wxRsp.NonceStr, pac, wechat.SignType_MD5, timeStamp, mng.Config.PayKey)
	xlog.Debug("paySign:", paySign)

	return map[string]interface{}{
		"app_id":    mng.Config.AppID,
		"timestamp": timeStamp,
		"nonce_str": wxRsp.NonceStr,
		"package":   "prepay_id=" + wxRsp.PrepayId,
		"sign_type": wechat.SignType_MD5,
		"pay_sign":  paySign,
	}, nil
}

// Refund 退款
func (mng *WechatPayMng) Refund(param *RefundParam) (err error) {

	xlog.Debug("out_refund_no:", param.OutTradeNo)

	totalFee :=  param.TotalAmount * 100 // 分为单位
	refundFee :=  param.RefundAmount * 100 // 分为单位

	// 初始化参数结构体
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", param.OutTradeNo).
		Set("nonce_str", util.GetRandomString(32)).
		Set("sign_type", wechat.SignType_MD5).
		Set("out_refund_no", param.OrderRefundNo).
		Set("total_fee", totalFee).
		Set("refund_fee", refundFee).
		Set("notify_url", mng.Config.RefundNotifyURL)

	//请求申请退款（沙箱环境下，证书路径参数可传空）
	//    body：参数Body
	wxRsp, resBm, err := mng.Client.Refund(bm)
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Debug("wxRsp：", wxRsp)
	xlog.Debug("resBm:", resBm)
	return
}
