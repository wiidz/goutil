package paymentMng

type PaymentWay int8 // 支付方式

const UnknownWay PaymentWay = 0           // 未知
const Cash PaymentWay = 1                 // 现金
const WechatPay PaymentWay = 2            // 微信支付
const AliPay PaymentWay = 3               // 支付宝支付
const OfficalBankTransfer PaymentWay = 4  // 对公账户转账
const PersonalBankTransfer PaymentWay = 5 // 私人账户转账
const Check PaymentWay = 6                // 支票

type PaymentKind int8 // 支付方式

const UnknownKind PaymentKind = 0 // 未知
const Pay PaymentKind = 1         // 支付
const Refund PaymentKind = 2      // 退款（订单退款、充值退款等）
const Withdraw PaymentKind = 3    // 提现（钱包余额）

// AliNotifyData 支付宝回调参数
type AliNotifyData struct {
	NotifyTime        string  `json:"notify_time"`         // 通知的发送时间。格式为 yyyy-MM-dd HH:mm:ss。示例值：2015-14-27 15:45:58
	NotifyType        string  `json:"notify_type"`         // 通知的类型。 示例值：trade_status_sync
	NotifyID          string  `json:"notify_id"`           // 通知校验 ID。示例值：ac05099524730693a8b330c5ecf72da9786
	AppID             string  `json:"app_id"`              // 支付宝分配给开发者的应用 ID。示例值：2014072300007148
	Charset           string  `json:"charset"`             // 编码格式，如 utf-8、gbk、gb2312 等。 示例值：utf-8
	Version           string  `json:"version"`             // 调用的接口版本，固定为：1.0。示例值：1.0
	SignType          string  `json:"sign_type"`           // 商户生成签名字符串所使用的签名算法类型，目前支持 RSA2 和 RSA，推荐使用 RSA2。示例值：RSA2
	Sign              string  `json:"sign"`                // 签名。详见下文 异步返回结果的验签。示例值：601510b7970e52cc63db0f44997cf70e
	TradeNo           string  `json:"trade_no"`            // 支付宝交易凭证号。示例值：2013112011001004330000121536
	OutTradeNo        string  `json:"out_trade_no"`        // 	原支付请求的商户订单号。示例值：6823789339978248
	OutBizNo          string  `json:"out_biz_no"`          // 商户业务 ID，主要是退款通知中返回退款申请的流水号。示例值：HZRF001
	BuyerID           string  `json:"buyer_id"`            // 	买家支付宝账号对应的支付宝唯一用户号。以 2088 开头的纯 16 位数字。示例值：2088102122524333
	BuyerLogonID      string  `json:"buyer_logon_id"`      // 买家支付宝账号。示例值：159﹡﹡﹡﹡﹡﹡20
	SellerID          string  `json:"seller_id"`           // 卖家支付宝用户号。示例值：2088101106499364
	SellerEmail       string  `json:"seller_email"`        // 卖家支付宝账号。示例值：zhu﹡﹡﹡@alitest.com
	TradeStatus       string  `json:"trade_status"`        // 	交易目前所处的状态。详见下文 交易状态说明。	示例值：TRADE_CLOSED
	TotalAmount       float64 `json:"total_amount"`        // 本次交易支付的订单金额，单位为人民币（元）。示例值：20
	ReceiptAmount     float64 `json:"receipt_amount"`      // 商家在交易中实际收到的款项，单位为人民币（元）。示例值：15
	InvoiceAmount     float64 `json:"invoice_amount"`      // 用户在交易中支付的可开发票的金额。示例值：10.00
	BuyerPayAmount    float64 `json:"buyer_pay_amount"`    // 用户在交易中支付的金额。示例值：13.88
	PointAmount       float64 `json:"point_amount"`        // 使用集分宝支付的金额。示例值：12.00
	RefundFee         float64 `json:"refund_fee"`          // 退款通知中，返回总退款金额，单位为人民币（元），支持两位小数。示例值：2.58
	Subject           string  `json:"subject"`             // 商品的标题/交易标题/订单标题/订单关键字等，是请求时对应的参数，原样通知回来。示例值：当面付交易
	Body              string  `json:"body"`                // 该订单的备注、描述、明细等。对应请求时的 body 参数，原样通知回来。示例值：当面付交易内容
	GmtCreate         string  `json:"gmt_create"`          // 该笔交易创建的时间。格式为yyyy-MM-dd HH:mm:ss示例值：2015-04-27 15:45:57
	GmtPayment        string  `json:"gmt_payment"`         // 该笔交易的买家付款时间。格式为yyyy-MM-dd HH:mm:ss 示例值：2015-04-27 15:45:57
	GmtRefund         string  `json:"gmt_refund"`          // 该笔交易的退款时间。格式为yyyy-MM-dd HH:mm:ss.S 示例值：2015-04-28 15:45:57.320
	GmtClose          string  `json:"gmt_close"`           // 该笔交易结束时间。格式为yyyy-MM-dd HH:mm:ss 示例值：2015-04-27 15:45:57
	FundBillList      string  `json:"fund_bill_list"`      // 支付成功的各个渠道金额信息。详见下文 资金明细信息说明。 示例值：[{"amount":"15.00","fundChannel":"ALIPAYACCOUNT"}]
	PassbackParams    string  `json:"passback_params"`     // 公共回传参数，如果请求时传递了该参数，则返回给商户时会在异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝。示例值：merchantBizType%3d3C%26merchantBizNo%3d2016010101111
	VoucherDetailList string  `json:"voucher_detail_list"` // 本交易支付时所使用的所有优惠券信息，详见下文 优惠券信息说明。示例值：[{"amount":"0.20","merchantContribute":"0.00","name":"一键创建券模板的券名称","otherContribute":"0.20","type":"ALIPAY_BIZ_VOUCHER","memo":"学生卡8折优惠"]
}

// AliUnifiedOrderParam 支付宝统一下单参数
type AliUnifiedOrderParam struct {
	Title       string  // 订单标题
	OutTradeNo  string  // 外部订单号
	TotalAmount float64 // 金额（元为单位）
	ReturnURL   string  // 支付后返回的页面URL
	IP          string  // 下单人的IP
}

// RefundParam 退款参数
type RefundParam struct {
	TransactionID string  // 原支付交易对应的微信订单号（二选一）
	OutTradeNo    string  // 原支付交易对应的商户订单号（二选一）
	OrderRefundNo string  // 商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	TotalAmount   float64 // 订单总金额
	RefundAmount  float64 // 退款金额
	Reason        string  // 退款原因，若商户传入，会在下发给用户的退款消息中体现退款原因
}

// UnifiedOrderParam 微信统一下单参数
type UnifiedOrderParam struct {
	Title       string  // 订单标题
	OutTradeNo  string  // 外部订单号
	TotalAmount float64 // 总金额（元为单位）
	ReturnURL   string  // 支付后返回的页面URL
	IP          string  // 下单人的IP
	AppName     string  // 我们的项目名称
}
