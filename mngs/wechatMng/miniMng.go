package wechatMng

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/encryptor"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
	"github.com/wiidz/goutil/mngs/configMng"
	"log"
)

// WechatMiniMng 微信小程序管理器
type WechatMiniMng struct {
	AppID     string
	AppSecret string
	//AccessToken string
	Client *miniprogram.MiniProgram
}

// NewWechatMiniMng 获取小程序管理器
func NewWechatMiniMng(appID string, appSecret string) *WechatMiniMng {

	//【1】使用redis缓存accessToken
	// memory := cache.NewMemory() // accessToken存在内存中
	redisConfig := configMng.GetRedis()
	redisOpts := &cache.RedisOpts{
		Host:     redisConfig.IP + ":" + redisConfig.Port,
		Password: redisConfig.Password,
	}
	redisCache := cache.NewRedis(redisOpts)

	//【2】创建mini实例
	cfg := &miniConfig.Config{
		AppID:     appID,
		AppSecret: appSecret,
		Cache:     redisCache,
	}
	wc := wechat.NewWechat()
	mini := wc.GetMiniProgram(cfg)

	//【3】返回
	var wechatMng = WechatMiniMng{
		AppID:     appID,
		AppSecret: appSecret,
		Client:    mini,
	}
	return &wechatMng
}

// Login 微信小程序登陆
func (mng *WechatMiniMng) Login(code string) (*auth.ResCode2Session, error) {
	authClient := mng.Client.GetAuth()
	res, err := authClient.Code2Session(code)
	return &res, err
}

// GetUserInfo 获取微信资料
func (mng *WechatMiniMng) GetUserInfo(sessionKey, encryptedData, iv string) (*encryptor.PlainData, error) {
	encryptorClient := mng.Client.GetEncryptor()
	res, err := encryptorClient.Decrypt(sessionKey, encryptedData, iv)
	return res, err
}

// GetPhone 获取微信手机号
func (mng *WechatMiniMng) GetPhone(sessionKey, encryptedData, iv string) (*encryptor.PlainData, error) {
	encryptorClient := mng.Client.GetEncryptor()
	res, err := encryptorClient.Decrypt(sessionKey, encryptedData, iv)
	return res, err
}

// GetQRCode 获取二维码
func (mng *WechatMiniMng) GetQRCode(qrCoder qrcode.QRCoder) ([]byte, error) {
	qrCodeApi := mng.Client.GetQRCode()
	res, err := qrCodeApi.GetWXACodeUnlimit(qrCoder)
	log.Println("res", res)
	log.Println("err", err)
	return res, err
}
