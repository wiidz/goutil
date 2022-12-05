package paymentMng

import (
	"context"
	"errors"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/util"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/go-pay/gopay/wechat"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
	"strconv"
	"time"
)

type WechatPayMngV2 struct {
	Config *configStruct.WechatPayConfigV2
	Client *wechat.Client
}

// getWechatPayInstance 获取微信支付实例
func getWechatPayInstance(config *configStruct.WechatPayConfigV2) (wechatPayMng *WechatPayMngV2, err error) {

	wechatPayMng = &WechatPayMngV2{
		Config: config,
		Client: wechat.NewClient(config.AppID, config.MchID, config.ApiKey, !config.Debug),
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

	//_ = wechatPayMng.Client.AddCertPemFilePath(Config.CertPath, Config.CertKeyPath)
	//err = wechatPayMng.Client.AddCertPkcs12FileContent([]byte(config.PEMCertContent))
	_ = wechatPayMng.Client.AddCertPkcs12FilePath(config.P12CertFilePath)

	// 添加微信pem证书
	return
}

// NewWechatPayMngV2 根据传入的config，获取信的微信支付
func NewWechatPayMngV2(config *configStruct.WechatPayConfigV2) (*WechatPayMngV2, error) {
	return getWechatPayInstance(config)
}

// H5 H5场景 totalFee 是分为单位
func (mng *WechatPayMngV2) H5(ctx context.Context, param *UnifiedOrderParam) (mWebUrl string, err error) {
	//初始化参数Map
	totalFee := param.TotalAmount * 100 // 分为单位

	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", util.RandomString(32)).
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
	wxRsp, err = mng.Client.UnifiedOrder(ctx, bm)
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

	if wxRsp.ReturnCode == "FAIL" {
		err = errors.New(wxRsp.ReturnMsg)
		return
	} else if len(wxRsp.ErrCodeDes) != 0 {
		err = errors.New(wxRsp.ErrCodeDes)
		return
	}

	mWebUrl = wxRsp.MwebUrl
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	log.Println("timeStamp", timeStamp)
	return
}

// Js js场景 统一下单获取 totalFee 是分为单位
func (mng *WechatPayMngV2) Js(ctx context.Context, param *UnifiedOrderParam, openID string) (data map[string]interface{}, err error) {

	totalFee := param.TotalAmount * 100 // 分为单位

	//初始化参数Map
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", util.RandomString(32)).
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
	wxRsp, err := mng.Client.UnifiedOrder(ctx, bm)
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Debug("Response：", wxRsp)
	xlog.Debug("wxRsp.MwebUrl:", wxRsp.MwebUrl)

	//获取Jsapi支付需要的paySign
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)

	pac := "prepay_id=" + wxRsp.PrepayId
	paySign := wechat.GetJsapiPaySign(mng.Config.AppID, wxRsp.NonceStr, pac, wechat.SignType_MD5, timeStamp, mng.Config.ApiKey)
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

// Mini 小程序场景下单
func (mng *WechatPayMngV2) Mini(ctx context.Context, param *UnifiedOrderParam, openID string) (timestampStr, packageStr, nonceStr, paySign string, err error) {
	//初始化参数Map
	totalFee := param.TotalAmount * 100 // 分为单位

	nonceStr = util.RandomString(32)

	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", nonceStr).
		Set("body", param.Title).
		Set("out_trade_no", param.OutTradeNo).
		Set("total_fee", totalFee).
		Set("spbill_create_ip", param.IP).
		Set("notify_url", mng.Config.NotifyURL).
		Set("trade_type", wechat.TradeType_Mini).
		Set("device_info", "WEB").
		Set("sign_type", wechat.SignType_MD5).
		Set("openid", openID). //js支付必填
		SetBodyMap("scene_info", func(bm gopay.BodyMap) {
			bm.SetBodyMap("h5_info", func(bm gopay.BodyMap) {
				bm.Set("type", "Wap")
				bm.Set("wap_url", param.ReturnURL)
				bm.Set("wap_name", param.AppName)
			})
		})

	//请求支付下单，成功后得到结果
	var wxRsp *wechat.UnifiedOrderResponse
	wxRsp, err = mng.Client.UnifiedOrder(ctx, bm)
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

	if wxRsp.ReturnCode == "FAIL" {
		err = errors.New(wxRsp.ReturnMsg)
		return
	} else if len(wxRsp.ErrCodeDes) != 0 {
		err = errors.New(wxRsp.ErrCodeDes)
		return
	}

	//pac := "prepay_id=" + wxRsp.PrepayId
	//paySign := wechat.GetMiniPaySign("wxdaa2ab9ef87b5497", wxRsp.NonceStr, pac, wechat.SignType_MD5, timeStamp, "GFDS8j98rewnmgl45wHTt980jg543abc")
	//xlog.Debug("paySign:", paySign)

	timestampStr = strconv.FormatInt(time.Now().Unix(), 10)
	packageStr = "prepay_id=" + wxRsp.PrepayId
	paySign = wechat.GetMiniPaySign(mng.Config.AppID, wxRsp.NonceStr, packageStr, wechat.SignType_MD5, timestampStr, mng.Config.ApiKey)
	xlog.Debug("paySign:", paySign)
	return
}

// Refund 退款
func (mng *WechatPayMngV2) Refund(ctx context.Context, param *RefundParam) (err error) {

	xlog.Debug("out_refund_no:", param.OutTradeNo)

	totalFee := param.TotalAmount * 100   // 分为单位
	refundFee := param.RefundAmount * 100 // 分为单位

	// 初始化参数结构体
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", param.OutTradeNo).
		Set("nonce_str", util.RandomString(32)).
		Set("sign_type", wechat.SignType_MD5).
		Set("out_refund_no", param.OrderRefundNo).
		Set("total_fee", totalFee).
		Set("refund_fee", refundFee).
		Set("notify_url", mng.Config.RefundNotifyURL)

	//请求申请退款（沙箱环境下，证书路径参数可传空）
	//    body：参数Body
	wxRsp, resBm, err := mng.Client.Refund(ctx, bm)
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Debug("wxRsp：", wxRsp)
	xlog.Debug("resBm:", resBm)
	return
}

// ScanPay 扫用户付款码收款
func (mng *WechatPayMngV2) ScanPay(ctx context.Context, param *ScanPayParam) (wxRsp *wechat.MicropayResponse, err error) {

	xlog.Debug("out_refund_no:", param.OutTradeNo)

	//初始化参数Map
	totalFee := param.TotalAmount * 100 // 分为单位

	nonceStr := util.RandomString(32)

	bm := make(gopay.BodyMap)
	bm.Set("appid", mng.Config.AppID).
		Set("mch_id", mng.Config.MchID).
		Set("device_info", param.DeviceNo).      //【是-String(32)】终端设备号(商户自定义，如门店编号)
		Set("nonce_str", nonceStr).              //【是-String(32)】随机字符串，不长于32位。推荐随机数生成算法
		Set("sign", "").                         //【是-String(32)】签名，详见签名生成算法 ，gopay会自动填充
		Set("sign_type", wechat.SignType_MD5).   //【是-String(32)】签名类型，目前支持HMAC-SHA256和MD5，默认为MD5
		Set("body", param.Title).                //【是-String(127)】商品简单描述，该字段须严格按照规范传递，具体请见参数规定
		Set("detail", nonceStr).                 //【否-String(6000)】单品优惠功能字段，需要接入详见单品优惠详细说明
		Set("attach", param.Attach).             //【否-String(127)】附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
		Set("out_trade_no", param.OutTradeNo).   //【是-String(32)】商户系统内部订单号，要求32个字符内（最少6个字符），只能是数字、大小写字母_-|*且在同一个商户号下唯一。详见商户订单号
		Set("total_fee", totalFee).              //【是-int】订单总金额，单位为分，只能为整数，详见支付金额
		Set("fee_type", "CNY").                  //【否-String(16)】符合ISO4217标准的三位字母代码，默认人民币：CNY，详见货币类型
		Set("spbill_create_ip", param.DeviceIP). //【否-String(64)】支持IPV4和IPV6两种格式的IP地址。调用微信支付API的机器IP
		Set("goods_tag", "").                    //【否-String(32)】订单优惠标记，代金券或立减优惠功能的参数，详见代金券或立减优惠
		Set("time_start", "").                   //【否-String(14)】订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
		Set("time_expire", "").                  //【否-String(14)】订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。
		Set("receipt", "").                      //【否-String(14)】电子发票入口开放标识,Y，传入Y时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效
		Set("auth_code", param.AuthCode).        //【是-String(128)】扫码支付付款码，设备读取用户微信中的条码或者二维码信息 （用户付款码规则：18位纯数字，前缀以10、11、12、13、14、15开头）
		Set("profit_sharing", "").               //【是-String(16)】Y-是，需要分账 N-否，不分账 字母要求大写，不传默认不分账
		SetBodyMap("scene_info", func(bm gopay.BodyMap) {
			bm.SetBodyMap("store_info", func(bm gopay.BodyMap) {
				bm.Set("id", param.DeviceNo)  // 门店ID
				bm.Set("name", param.AppName) // 名称
				bm.Set("area_code", "")       // 编码
				bm.Set("address", "")         // 地址
			})
		}) // 【否-String(256)】该字段用于上报场景信息，目前支持上报实际门店信息。该字段为JSON对象数据，对象格式为{"store_info":{"id": "门店ID","name": "名称","area_code": "编码","address": "地址" }} ，字段详细说明请点击行前的+展开

	wxRsp, err = mng.Client.Micropay(ctx, bm)
	if err != nil {
		xlog.Error(err)
		return
	}

	xlog.Debug("Response：", wxRsp)
	xlog.Debug("wxRsp.CashFee:", wxRsp.CashFee)
	xlog.Debug("wxRsp.CashFeeType:", wxRsp.CashFeeType)
	xlog.Debug("wxRsp.ResultCode:", wxRsp.ResultCode)
	xlog.Debug("wxRsp.ReturnCode:", wxRsp.ReturnCode)
	xlog.Debug("wxRsp.ReturnMsg:", wxRsp.ReturnMsg)
	xlog.Debug("wxRsp.MwebUrl:", wxRsp.ErrCode)
	xlog.Debug("wxRsp.MwebUrl:", wxRsp.ErrCodeDes)
	xlog.Debug("wxRsp.Sign:", wxRsp.Sign)
	xlog.Debug("wxRsp.Appid:", wxRsp.Appid)
	xlog.Debug("wxRsp.Attach:", wxRsp.Attach)
	xlog.Debug("wxRsp.CouponFee:", wxRsp.CouponFee)
	xlog.Debug("wxRsp.DeviceInfo:", wxRsp.DeviceInfo)
	xlog.Debug("wxRsp.MchId:", wxRsp.MchId)
	xlog.Debug("wxRsp.TransactionId:", wxRsp.TransactionId)
	xlog.Debug("wxRsp.IsSubscribe:", wxRsp.IsSubscribe)
	xlog.Debug("wxRsp.FeeType:", wxRsp.FeeType)
	xlog.Debug("wxRsp.FeeType:", wxRsp.FeeType)

	if wxRsp.ReturnCode == "FAIL" {
		err = errors.New(wxRsp.ReturnMsg)
		return
	} else if len(wxRsp.ErrCodeDes) != 0 {
		err = errors.New(wxRsp.ErrCodeDes)
		return
	}

	return
}
