package genderMap

type Gender int8

const (
	Unknown Gender = 0 // 未知性别
	Male    Gender = 1 // 男
	FeMale  Gender = 2 // 女
)

func (gender Gender) ToString() string {
	switch gender {
	case Male:
		return "男"
	case FeMale:
		return "女"
	default:
		return "未知"
	}
}

// Euphemism 获取雅称
func (gender Gender) Euphemism() string {
	switch gender {
	case Male:
		return "先生"
	case FeMale:
		return "女士"
	default:
		return "用户"
	}
}

// Emoji 获取性别图标
func (gender Gender) Emoji() string {
	switch gender {
	case Male:
		return "🚹"
	case FeMale:
		return "🚺"
	default:
		return "❓"
	}
}
