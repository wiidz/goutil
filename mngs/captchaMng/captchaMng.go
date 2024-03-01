package captchaMng

import (
	"errors"
	cp "github.com/mojocn/base64Captcha"
	"github.com/wiidz/goutil/helpers/mathHelper"
	"github.com/wiidz/goutil/helpers/strHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/memoryMng"
	"github.com/wiidz/goutil/mngs/redisMng"
	"github.com/wiidz/goutil/structs/dataSourceStruct"
	"image/color"
	"log"
	"strings"
	"time"
)

type CaptchaMng struct {
	DataSource dataSourceStruct.DataSource
	RedisMng   *redisMng.RedisMng
	MemoryMng  *memoryMng.MemoryMng
}

func NewCaptchaMngRedis(redisM *redisMng.RedisMng) (*CaptchaMng, error) {
	return &CaptchaMng{
		DataSource: dataSourceStruct.Redis,
		RedisMng:   redisM,
	}, nil
}

func NewCaptchaMngMemory(memoryM *memoryMng.MemoryMng) (*CaptchaMng, error) {
	return &CaptchaMng{
		DataSource: dataSourceStruct.Memory,
		MemoryMng:  memoryM,
	}, nil
}

// VerifyGraphCaptcha 验证图形验证码是否有效
func (mng *CaptchaMng) VerifyGraphCaptcha(id, answer string) bool {
	get := cp.DefaultMemStore.Get(id, false)
	if get == "" {
		return false
	}
	if strings.ToLower(strings.TrimSpace(answer)) != strings.ToLower(get) {
		return false
	}
	
	// 手动删除
	cp.DefaultMemStore.Get(id, true)
	return true
}

// GenerateGraphCaptcha 生成图形验证码base64
func (mng *CaptchaMng) GenerateGraphCaptcha(width, height, noiseCount, length int) (id, b64s string, err error) {
	driver := cp.NewDriverString(height, width, noiseCount, cp.OptionShowHollowLine,
		length, cp.TxtSimpleCharaters, &color.RGBA{254, 254, 254, 254}, []string{"Flim-Flam.ttf"})
	captcha := cp.NewCaptcha(driver, cp.DefaultMemStore)
	return captcha.Generate()
}

// GenerateNumberGraphCaptcha 生成图形验证码base64（仅数字）
func (mng *CaptchaMng) GenerateNumberGraphCaptcha(width, height, noiseCount, length int) (id, b64s string, err error) {
	driver := cp.NewDriverString(height, width, noiseCount, cp.OptionShowHollowLine,
		length, cp.TxtNumbers, &color.RGBA{254, 254, 254, 254}, []string{"Flim-Flam.ttf"})
	captcha := cp.NewCaptcha(driver, cp.DefaultMemStore)
	return captcha.Generate()
}

// GetNumberCaptcha 获取数字验证码
func (mng *CaptchaMng) GetNumberCaptcha(identify string) (id, captchaStr string, err error) {

	captcha := mathHelper.GetRandomInt(100000, 999999) // 默认六位
	captchaStr = typeHelper.Int2Str(captcha)
	id = strHelper.GetRandomString(10)

	_ = mng.SetCache(identify+id, captchaStr, time.Second*300) // 300秒有效
	return
}

// VerifyNumberCaptcha 验证数字验证码
func (mng *CaptchaMng) VerifyNumberCaptcha(identifyKey, id, captchaStr string) (err error) {

	keyName := identifyKey + id
	captchaCache, err := mng.GetCache(keyName)
	if err != nil {
		return
	}
	if captchaCache == "" {
		return errors.New("验证码已失效")
	}

	if captchaCache != captchaStr {
		return errors.New("验证码错误")
	}

	_ = mng.SetCache(keyName, "0", 0)

	return nil
}

// SetCache 记录缓存
func (mng *CaptchaMng) SetCache(keyName, value string, expire time.Duration) (err error) {
	if mng.DataSource == dataSourceStruct.Redis {
		err = mng.RedisMng.Set(keyName, value, expire)
	} else if mng.DataSource == dataSourceStruct.Memory {
		mng.MemoryMng.Set(keyName, value, expire)
	}

	log.Println("keyName", keyName)
	log.Println("value", value)

	return err
}

// GetCache 读取缓存
func (mng *CaptchaMng) GetCache(keyName string) (string, error) {
	log.Println("keyName", keyName)
	if mng.DataSource == dataSourceStruct.Redis {
		return mng.RedisMng.GetString(keyName)
	} else if mng.DataSource == dataSourceStruct.Memory {
		value, exist := mng.MemoryMng.GetString(keyName)
		if exist == false {
			return "", errors.New("指定的keyName不存在")
		}
		return value, nil
	}
	return "", errors.New("未知数据源")
}
