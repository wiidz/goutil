package captchaMng

import (
	"context"
	"errors"
	"image/color"
	"log"
	"strings"
	"time"

	cp "github.com/mojocn/base64Captcha"
	"github.com/wiidz/goutil/helpers/mathHelper"
	"github.com/wiidz/goutil/helpers/strHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/memoryMng"
	"github.com/wiidz/goutil/mngs/redisMng"
	"github.com/wiidz/goutil/structs/dataSourceStruct"
)

// LoginCaptchaChars 登录验证码字符集：纯数字，排除易混淆的 0/1/2/5/8。
const LoginCaptchaChars = "346789"

type CaptchaMng struct {
	DataSource dataSourceStruct.DataSource
	RedisMng   *redisMng.RedisMng
	MemoryMng  *memoryMng.MemoryMng
}

// GraphCaptchaOption 图形验证码配置。
type GraphCaptchaOption struct {
	Width           int
	Height          int
	Length          int
	NoiseCount      int
	ShowLineOptions int
	Source          string
	Fonts           []string
	BgColor         *color.RGBA
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

// DefaultLoginGraphCaptchaOption 登录页默认可读性优先的配置。
func DefaultLoginGraphCaptchaOption() GraphCaptchaOption {
	return GraphCaptchaOption{
		Width:           200,
		Height:          56,
		Length:          4,
		NoiseCount:      0,
		ShowLineOptions: 0,
		Source:          LoginCaptchaChars,
		Fonts:           []string{"actionj.ttf"},
		BgColor:         &color.RGBA{R: 248, G: 250, B: 252, A: 255},
	}
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

// GenerateLoginGraphCaptcha 生成登录页图形验证码（纯数字、无干扰线、清晰字体）。
func (mng *CaptchaMng) GenerateLoginGraphCaptcha() (id, b64s string, err error) {
	return mng.GenerateGraphCaptchaWithOption(DefaultLoginGraphCaptchaOption())
}

// GenerateGraphCaptchaWithOption 按配置生成图形验证码。
func (mng *CaptchaMng) GenerateGraphCaptchaWithOption(opt GraphCaptchaOption) (id, b64s string, err error) {
	driver := cp.NewDriverString(
		opt.Height,
		opt.Width,
		opt.NoiseCount,
		opt.ShowLineOptions,
		opt.Length,
		opt.Source,
		opt.BgColor,
		opt.Fonts,
	)
	captcha := cp.NewCaptcha(driver, cp.DefaultMemStore)
	return captcha.Generate()
}

// GenerateGraphCaptcha 生成图形验证码 base64。
// width、height 为图片像素宽高，对应前端展示区域尺寸。
func (mng *CaptchaMng) GenerateGraphCaptcha(width, height, noiseCount, length int) (id, b64s string, err error) {
	return mng.GenerateGraphCaptchaWithOption(GraphCaptchaOption{
		Width:           width,
		Height:          height,
		Length:          length,
		NoiseCount:      noiseCount,
		ShowLineOptions: 0,
		Source:          cp.TxtSimpleCharaters,
		Fonts:           []string{"actionj.ttf"},
		BgColor:         &color.RGBA{R: 254, G: 254, B: 254, A: 254},
	})
}

// GenerateNumberGraphCaptcha 生成图形验证码 base64（仅数字）
func (mng *CaptchaMng) GenerateNumberGraphCaptcha(width, height, noiseCount, length int) (id, b64s string, err error) {
	return mng.GenerateGraphCaptchaWithOption(GraphCaptchaOption{
		Width:           width,
		Height:          height,
		Length:          length,
		NoiseCount:      noiseCount,
		ShowLineOptions: 0,
		Source:          cp.TxtNumbers,
		Fonts:           []string{"actionj.ttf"},
		BgColor:         &color.RGBA{R: 254, G: 254, B: 254, A: 254},
	})
}

// GetNumberCaptcha 获取数字验证码
func (mng *CaptchaMng) GetNumberCaptcha(ctx context.Context, identify string) (id, captchaStr string, err error) {

	captcha := mathHelper.GetRandomInt(100000, 999999) // 默认六位
	captchaStr = typeHelper.Int2Str(captcha)
	id = strHelper.GetRandomString(10)

	_ = mng.SetCache(ctx, identify+id, captchaStr, time.Second*300) // 300秒有效
	return
}

// VerifyNumberCaptcha 验证数字验证码
func (mng *CaptchaMng) VerifyNumberCaptcha(ctx context.Context, identifyKey, id, captchaStr string) (err error) {

	keyName := identifyKey + id
	captchaCache, err := mng.GetCache(ctx, keyName)
	if err != nil {
		return
	}
	if captchaCache == "" {
		return errors.New("验证码已失效")
	}

	if captchaCache != captchaStr {
		return errors.New("验证码错误")
	}

	_ = mng.SetCache(ctx, keyName, "0", 0)

	return nil
}

// SetCache 记录缓存
func (mng *CaptchaMng) SetCache(ctx context.Context, keyName, value string, expire time.Duration) (err error) {
	if mng.DataSource == dataSourceStruct.Redis {
		err = mng.RedisMng.Set(ctx, keyName, value, expire)
	} else if mng.DataSource == dataSourceStruct.Memory {
		mng.MemoryMng.Set(keyName, value, expire)
	}

	log.Println("keyName", keyName)
	log.Println("value", value)

	return err
}

// GetCache 读取缓存
func (mng *CaptchaMng) GetCache(ctx context.Context, keyName string) (string, error) {
	log.Println("keyName", keyName)
	if mng.DataSource == dataSourceStruct.Redis {
		return mng.RedisMng.GetString(ctx, keyName)
	} else if mng.DataSource == dataSourceStruct.Memory {
		value, exist := mng.MemoryMng.GetString(keyName)
		if exist == false {
			return "", errors.New("指定的keyName不存在")
		}
		return value, nil
	}
	return "", errors.New("未知数据源")
}
