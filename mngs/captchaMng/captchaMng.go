package captchaMng

import (
	cp "github.com/mojocn/base64Captcha"
	"image/color"
	"strings"
)

// 目前是单机部署版本，后续改redis

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
func GenerateCaptcha(width,height,noiseCount,length int) (id, b64s string, err error) {
	driver := cp.NewDriverString(height, width, noiseCount, cp.OptionShowHollowLine,
		length, cp.TxtSimpleCharaters, &color.RGBA{254, 254, 254, 254}, []string{"Flim-Flam.ttf"})
	captcha := cp.NewCaptcha(driver, cp.DefaultMemStore)

	return captcha.Generate()
}
