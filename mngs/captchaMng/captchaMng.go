package captchaMng

import (
	cp "github.com/mojocn/base64Captcha"
	"image/color"
	"strings"
)

// 目前是单机部署版本，后续改redis

var captcha *cp.Captcha

const (
	height = 43
	width  = 200
	length = 4
)

//创建字符串验证码实例
func init() {
	driver := cp.NewDriverString(height, width, 0, cp.OptionShowHollowLine,
		length, cp.TxtSimpleCharaters, &color.RGBA{254, 254, 254, 254}, []string{"Flim-Flam.ttf"})
	cape := cp.NewCaptcha(driver, cp.DefaultMemStore)
	captcha = cape
}

//验证是否有效
func VerifyCaptcha(id, answer string) bool {
	get := cp.DefaultMemStore.Get(id, false)
	if get == "" {
		return false
	}
	if strings.ToLower(strings.TrimSpace(answer)) != strings.ToLower(get) {
		return false
	}
	return true
}

//生成base64
func GenerateCaptcha() (id, b64s string, err error) {
	return captcha.Generate()
}
