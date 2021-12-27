package paymentMng

import (
	"errors"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
	"net/http"
	"time"
)

type WechatPayMngV3 struct {
	Config *configStruct.WechatPayConfig
	Client *wechat.ClientV3
}

// getWechatPayV3Instance 获取微信支付V3实例
func getWechatPayV3Instance(config *configStruct.WechatPayConfig) (mng *WechatPayMngV3, err error) {

	var client *wechat.ClientV3
	client, err = wechat.NewClientV3(config.MchID, config.CertSerialNo, config.ApiKeyV3, config.CertContent)
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
func NewWechatPayMngV3(config *configStruct.WechatPayConfig) (*WechatPayMngV3, error) {
	return getWechatPayV3Instance(config)
}

// Mini 小程序场景下单
func (mng *WechatPayMngV3) Mini(params *UnifiedOrderParam, openID string) (timestampStr, packageStr, nonceStr, paySign string, err error) {

	//【1】获取prepayID
	prepayRes, err := mng.jsApiPlaceOrder(params, openID)
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

// Js 公众号支付
func (mng *WechatPayMngV3) Js(params *UnifiedOrderParam, openID string) (appID, timestampStr, nonceStr, packageStr, paySign, signType string, err error) {

	//【1】获取prepayID
	prepayRes, err := mng.jsApiPlaceOrder(params, openID)
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

// H5 网页支付
func (mng *WechatPayMngV3) H5(params *UnifiedOrderParam, openID string) (H5Url string, err error) {

	//【1】构建结构体
	bm := make(gopay.BodyMap)
	bm.Set("appid", mng.Config.AppID).
		Set("mchid", mng.Config.MchID).
		Set("description", params.Title).
		Set("out_trade_no", params.OutTradeNo).
		Set("notify_url", mng.Config.NotifyURL).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", params.TotalAmount).
				Set("currency", "CNY")
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", openID)
		}).
		SetBodyMap("scene_info", func(bm gopay.BodyMap) {
			bm.Set("payer_client_ip", params.IP).
				SetBodyMap("h5_info", func(bm gopay.BodyMap) {
					bm.Set("type", "H5")
				})
		})

	log.Println("bm", bm)
	log.Println("bm", bm)

	//【2】获取签名
	var res *wechat.H5Rsp
	res, err = mng.Client.V3TransactionH5(bm)
	if err != nil {
		return
	}
	if res.Code != 0 {
		// 出错
		return "", errors.New(res.Error)
	}

	return res.Response.H5Url, err
}

// Refund 退款
func (mng *WechatPayMngV3) Refund(param *RefundParam) (wxRsp *wechat.RefundRsp, err error) {
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", param.OutTradeNo).
		Set("out_refund_no", param.OrderRefundNo).
		Set("reason", param.Reason).
		Set("notify_url", mng.Config.RefundNotifyURL).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("refund", param.RefundAmount).
				Set("total", param.TotalAmount).
				Set("currency", "CNY")
		})

	wxRsp, err = mng.Client.V3Refund(bm)
	return
}

// jsApiPlaceOrder JSAPI/小程序下单API totalFee 是分为单位
func (mng *WechatPayMngV3) jsApiPlaceOrder(params *UnifiedOrderParam, openID string) (prepayRsp *wechat.PrepayRsp, err error) {

	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	totalFee := params.TotalAmount * 100 // 分为单位

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

	prepayRsp, err = mng.Client.V3TransactionJsapi(bm)
	if err != nil {
		xlog.Error(err)
		return
	}

	if prepayRsp.Code != 0 {
		err = errors.New(prepayRsp.Error)
		return
	}

	return
}

// h5PlaceOrder H5下单
//func (mng *WechatPayMngV3) h5PlaceOrder(params *UnifiedOrderParam, openID string) (prepayRsp *wechat.H5Rsp, err error) {
//
//	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
//	totalFee := params.TotalAmount * 100 // 分为单位
//
//	// 初始化 BodyMap
//	bm := make(gopay.BodyMap)
//	bm.Set("appid", "appid").
//		Set("mchid", "mchid").
//		Set("sub_mchid", "sub_mchid").
//		Set("description", params.Title).
//		Set("out_trade_no", params.OutTradeNo).
//		Set("time_expire", expire).
//		Set("notify_url", mng.Config.NotifyURL).
//		SetBodyMap("amount", func(bm gopay.BodyMap) {
//			bm.Set("total", totalFee).
//				Set("currency", "CNY")
//		}).
//		SetBodyMap("payer", func(bm gopay.BodyMap) {
//			bm.Set("openid", openID)
//		})
//
//	prepayRsp, err = mng.Client.V3TransactionH5(bm)
//	if err != nil {
//		xlog.Error(err)
//		return
//	}
//
//	if len(prepayRsp.Error) != 0 {
//		err = errors.New(prepayRsp.Error)
//		return
//	}
//
//	return
//}

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
