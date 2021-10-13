package captchaMng

import (
	"errors"
	cp "github.com/mojocn/base64Captcha"
	"github.com/wiidz/goutil/helpers/mathHelper"
	"github.com/wiidz/goutil/helpers/strHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/redisMng"
	"image/color"
	"strings"
)


var redis = redisMng.NewRedisMng()

// 目前是单机部署版本，后续改redis

// 验证是否有效
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

// 生成base64
func GenerateCaptcha(width,height,noiseCount,length int) (id, b64s string, err error) {
	driver := cp.NewDriverString(height, width, noiseCount, cp.OptionShowHollowLine,
		length, cp.TxtSimpleCharaters, &color.RGBA{254, 254, 254, 254}, []string{"Flim-Flam.ttf"})
	captcha := cp.NewCaptcha(driver, cp.DefaultMemStore)

	return captcha.Generate()
}

// GetNumberCaptcha 获取数字验证码
func GetNumberCaptcha()(id,captchaStr string,err error){

	captcha := mathHelper.GetRandomInt(100000,999999) // 默认六位
	captchaStr = typeHelper.Int2Str(captcha)
	id = strHelper.GetRandomString(10)

	_ = redis.Set("captcha-"+id,captchaStr,300) // 300秒有效
	return
}

// VerifyNumberCaptcha 验证数字验证码
func VerifyNumberCaptcha(id,captchaStr string)(err error){
	var captchMem string
	keyName := "captcha-"+id
	captchMem,err = redis.GetString(keyName)
	if err != nil {
		return
	}

	if captchMem != captchaStr {
		return errors.New("验证码错误")
	}

	_ = redis.Set(keyName,0,0)

	return nil
}