package paymentMng

import (
	"context"
	"errors"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"math"
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
	client, err = wechat.NewClientV3(config.MchID, config.CertSerialNo, config.ApiKeyV3, config.PEMKeyContent)
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
		if err != nil {
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
		if err != nil {
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
		if err != nil {
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
		transferList[k].UserName, err = wechat.V3EncryptText(transferList[k].UserName, []byte(mng.Config.PEMKeyContent))
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
						Set("transfer_detail_list", transferList)

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
		xlog.Error(err)
		return
	}

	return
}
