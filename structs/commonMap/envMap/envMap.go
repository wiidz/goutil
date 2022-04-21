package envMap

type EnvKind int8

const (
	UnknownKind       EnvKind = 0 // 未知环境
	Web               EnvKind = 1 // PC网站
	Wap               EnvKind = 2 // 手机web
	WechatMiniProgram EnvKind = 3 // 微信小程序
	BaiduMiniProgram  EnvKind = 4 // 百度小程序
	TikTokProgram     EnvKind = 5 // 抖音小程序
	AliMiniProgram    EnvKind = 6 // 支付宝小程序
	WechatWap         EnvKind = 7 // 微信浏览器网站
)

// GetStruct 根据值获取对象
func GetStruct(number int8) EnvKind {
	switch number {
	case int8(1):
		return Web
	case int8(2):
		return Wap
	case int8(3):
		return WechatMiniProgram
	case int8(4):
		return BaiduMiniProgram
	case int8(5):
		return TikTokProgram
	case int8(6):
		return AliMiniProgram
	case int8(7):
		return WechatWap
	default:
		return UnknownKind
	}
}

// String 获取中文解释
func (o EnvKind) String() string {
	switch o {
	case Web:
		return "PC网站"
	case Wap:
		return "移动端网站"
	case WechatMiniProgram:
		return "微信小程序"
	case BaiduMiniProgram:
		return "百度小程序"
	case TikTokProgram:
		return "抖音小程序"
	case AliMiniProgram:
		return "支付宝小程序"
	case WechatWap:
		return "微信内浏览器"
	default:
		return "未知环境"
	}
}

// Value 获取对应值
func (o EnvKind) Value() int8 {
	switch o {
	case Web:
		return int8(1)
	case Wap:
		return int8(2)
	case WechatMiniProgram:
		return int8(3)
	case BaiduMiniProgram:
		return int8(4)
	case TikTokProgram:
		return int8(5)
	case AliMiniProgram:
		return int8(6)
	case WechatWap:
		return int8(7)
	default:
		return int8(0)
	}
}
