package orderMap

type OrderStatus int8

const WaitConfirm OrderStatus = 0 // （可跳过此步）待确认
const WaitPay OrderStatus = 1     // （确认完毕）待支付
const WaitCheck OrderStatus = 2   // （支付完毕）待核验
const Preparing OrderStatus = 3   // （核验完成）准备中/备货中
const Packing OrderStatus = 4     // （备货完毕）打包中
const WaitCollect OrderStatus = 5 // （打包完成）揽件中
const Delivering OrderStatus = 6  // （揽件完成）运输中
const Signed OrderStatus = 7      // （运输完成）已签收（直接已完成）

//const Refunding OrderStatus = 98 // 申请退款中

const Done OrderStatus = 99      // 已完成
const Canceled OrderStatus = 100 // 已取消
const Closed OrderStatus = 101   // 已关闭
//const Refunded OrderStatus = 102 // 已退款

func (status OrderStatus) ToString() string {
	switch status {
	case WaitConfirm:
		return "待确认"
	case WaitPay:
		return "待支付"
	case WaitCheck:
		return "备货中"
	case Preparing:
		return "待核验"
	case Packing:
		return "打包中"
	case WaitCollect:
		return "揽件中"
	case Delivering:
		return "运输中"
	case Signed:
		return "已签收"
	case Done:
		return "已完成"
	case Canceled:
		return "已取消"
	case Closed:
		return "已关闭"
	default:
		return "未知"
	}
}
