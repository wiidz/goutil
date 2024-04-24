package paymentMng

import "time"

type TransferMerchantErrorObj struct {
	Code   string `json:"code"`
	Detail struct {
		Location string `json:"location"`
		Value    int    `json:"value"`
	} `json:"detail"`
	Message string `json:"message"`
}

type TransferMerchantError string

// Cn 返回错误信息的中文解释
func (tme TransferMerchantError) Cn() string {
	switch tme {
	case "ACCOUNT_FROZEN":
		return "该用户账户被冻结"
	case "REAL_NAME_CHECK_FAIL":
		return "收款人未实名认证，需要用户完成微信实名认证"
	case "NAME_NOT_CORRECT":
		return "收款人姓名校验不通过，请核实信息"
	case "OPENID_INVALID":
		return "Openid格式错误或者不属于商家公众账号"
	case "TRANSFER_QUOTA_EXCEED":
		return "超过用户单笔收款额度，核实产品设置是否准确"
	case "DAY_RECEIVED_QUOTA_EXCEED":
		return "超过用户单日收款额度，核实产品设置是否准确"
	case "MONTH_RECEIVED_QUOTA_EXCEED":
		return "超过用户单月收款额度，核实产品设置是否准确"
	case "DAY_RECEIVED_COUNT_EXCEED":
		return "超过用户单日收款次数，核实产品设置是否准确"
	case "PRODUCT_AUTH_CHECK_FAIL":
		return "未开通该权限或权限被冻结，请核实产品权限状态"
	case "OVERDUE_CLOSE":
		return "超过系统重试期，系统自动关闭"
	case "ID_CARD_NOT_CORRECT":
		return "收款人身份证校验不通过，请核实信息"
	case "ACCOUNT_NOT_EXIST":
		return "该用户账户不存在"
	case "TRANSFER_RISK":
		return "该笔转账可能存在风险，已被微信拦截"
	case "OTHER_FAIL_REASON_TYPE":
		return "其它失败原因"
	case "REALNAME_ACCOUNT_RECEIVED_QUOTA_EXCEED":
		return "用户账户收款受限，请引导用户在微信支付查看详情"
	case "RECEIVE_ACCOUNT_NOT_PERMMIT":
		return "未配置该用户为转账收款人，请在产品设置中调整，添加该用户为收款人"
	case "PAYEE_ACCOUNT_ABNORMAL":
		return "用户账户收款异常，请联系用户完善其在微信支付的身份信息以继续收款"
	case "PAYER_ACCOUNT_ABNORMAL":
		return "商户账户付款受限，可前往商户平台获取解除功能限制指引"
	case "TRANSFER_SCENE_UNAVAILABLE":
		return "该转账场景暂不可用，请确认转账场景ID是否正确"
	case "TRANSFER_SCENE_INVALID":
		return "你尚未获取该转账场景，请确认转账场景ID是否正确"
	case "TRANSFER_REMARK_SET_FAIL":
		return "转账备注设置失败，请调整后重新再试"
	case "RECEIVE_ACCOUNT_NOT_CONFIGURE":
		return "请前往商户平台-商家转账到零钱-前往功能-转账场景中添加"
	case "BLOCK_B2C_USERLIMITAMOUNT_BSRULE_MONTH":
		return "超出用户单月转账收款20w限额，本月不支持继续向该用户付款"
	case "BLOCK_B2C_USERLIMITAMOUNT_MONTH":
		return "用户账户存在风险收款受限，本月不支持继续向该用户付款"
	case "MERCHANT_REJECT":
		return "商户员工（转账验密人）已驳回转账"
	case "MERCHANT_NOT_CONFIRM":
		return "商户员工（转账验密人）超时未验密"
	default:
		return "未知错误"
	}
}

// V3DecryptTransferMerchantResult 商家转账解密后的数据（gopay没提供，我们自己写）
// 当 event_type为MCHTRANSFER.BATCH.FINISHED时，数据密文ciphertext解密之后的内容
type V3DecryptTransferMerchantResult struct {
	OutBatchNo    string    `json:"out_batch_no"`
	BatchId       string    `json:"batch_id"`
	BatchStatus   string    `json:"batch_status"`
	TotalNum      int       `json:"total_num"`
	TotalAmount   int       `json:"total_amount"`
	SuccessAmount int       `json:"success_amount"`
	SuccessNum    int       `json:"success_num"`
	FailAmount    int       `json:"fail_amount"`
	FailNum       int       `json:"fail_num"`
	MchID         string    `json:"mchid"`
	UpdateTime    time.Time `json:"update_time"`

	CloseReason string `json:"close_reason"` // 当 event_type 为MCHTRANSFER.BATCH.CLOSED时，数据密文ciphertext解密之后的内容
}
