package aiMng

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrServiceBusy AI服务正在执行中
	ErrServiceBusy = errors.New("AI服务正在执行中，请稍后再试")
	// ErrTimeout AI服务执行超时
	ErrTimeout = errors.New("AI服务执行超时")
)

// 名称标准化字符替换器
var nameNormalizeReplacer = strings.NewReplacer(
	"-", "",
	"_", "",
	"/", "",
	"（", "",
	"）", "",
	"(", "",
	")", "",
	" ", "",
	"·", "",
	".", "",
	",", "",
	"，", "",
	"×", "x",
)

var (
	skuIntentKeywords = []string{
		"规格", "型号", "选型", "多少种", "几种", "多少", "参数", "价格", "多少钱", "尺寸", "外径",
	}
	skuIntentKeywordsLower = []string{
		"spec", "model", "size", "variant", "sku",
	}
	skuCodePattern        = regexp.MustCompile(`(?i)\b[a-z]{1,6}\s*[-/]?\s*\d+(?:\.\d+)?[a-z0-9]*\b`)
	selectionIndexPattern = regexp.MustCompile(`(?i)(?:第)?([1-9]|10)(?:款|个|条|项|种|号)?`)
	skuFollowUpStopWords  = map[string]struct{}{
		"款": {}, "个": {}, "种": {}, "要": {}, "再": {}, "和": {}, "还": {}, "有": {}, "的": {},
	}
)
