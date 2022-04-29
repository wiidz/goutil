package strHelper

import (
	"bytes"
	"encoding/base64"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const (
	base64Str = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

/**
 * @func: Base64_encode  base64编码
 * @author Wiidz
 * @date   2019-11-16
 */
func Base64Encode(src []byte) (str string) {
	coder := base64.NewEncoding(base64Str)
	str = coder.EncodeToString(src)
	return
}

/**
 * @func: Base64_decode 解码
 * @author Wiidz
 * @date   2019-11-16
 */
func Base64Decode(str string) (data []byte) {
	coder := base64.NewEncoding(base64Str)
	var err error
	data, err = coder.DecodeString(str)
	if err != nil {
		return
	}
	return
}

// GetRandomString 获取指定位数的随机字符串
func GetRandomString(l int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tempBytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, tempBytes[r.Intn(len(tempBytes))])
	}
	return string(result)
}

// GetRandomNumbers 获取指定位数的随机字符串
func GetRandomNumbers(l int) string {
	str := "0123456789"
	tempBytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, tempBytes[r.Intn(len(tempBytes))])
	}
	return string(result)
}

/**
 * @func: GetWordLength  文字=1，标点符号=0.5，获取长度
 * @author Wiidz
 * @date   2019-11-16
 */
func GetWordsLength(str string) float64 {
	var total float64
	reg := regexp.MustCompile("/·|，|。|《|》|‘|’|”|“|；|：|【|】|？|（|）|、/")
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || reg.Match([]byte(string(r))) {
			total = total + 1
		} else {
			total = total + 0.5
		}
	}
	return math.Ceil(total)
}

/**
 * @func: StripTags  去除文本中的html标签
 * @author Wiidz
 * @date   2019-11-16
 */
func StripTags(body string) string {

	src := string(body)

	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "，")

	return strings.TrimSpace(src)

}

/**
 * @func: Trim  过滤空字符串（不可靠）
 * @author Wiidz
 * @date   2019-11-16
 */
func Trim(str string) (data string) {
	tmp := bytes.Trim([]byte(str), " ")
	data = string(tmp)
	return
}

/**
 * @func: Substr  截取字符串
 * @author Wiidz
 * @date   2019-11-16
 */
func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

/**
 * @func：正则验证是否是手机号
 * @author Hank
 * @date   2019-02-24
 */
func ValidatePhone(phoneNum string) bool {
	regular := `(?:^1[3456789]|^9[28])\d{9}$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(phoneNum)
}

// Exist 判断目标字符串中是否存在需要的字符
func Exist(targetStr, needleStr string) bool {
	if strings.Index(targetStr, needleStr) == -1 {
		return false
	} else {
		return true
	}
}
