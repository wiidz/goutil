package wechatMng

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/encryptor"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
	"github.com/wiidz/goutil/structs/configStruct"
)

// MiniMng 微信小程序管理器
type MiniMng struct {
	AppID     string
	AppSecret string
	//AccessToken string
	Client *miniprogram.MiniProgram
}

// NewMiniMng 获取小程序管理器
func NewMiniMng(redisC *configStruct.RedisConfig, miniC *configStruct.WechatMiniConfig) *MiniMng {

	//【1】使用redis缓存accessToken
	// memory := cache.NewMemory() // accessToken存在内存中
	redisOpts := &cache.RedisOpts{
		Host:     redisC.Host + ":" + redisC.Port,
		Password: redisC.Password,
	}
	redisCache := cache.NewRedis(redisOpts)

	//【2】创建mini实例
	cfg := &miniConfig.Config{
		AppID:     miniC.AppID,
		AppSecret: miniC.AppSecret,
		Cache:     redisCache,
	}
	wc := wechat.NewWechat()
	mini := wc.GetMiniProgram(cfg)

	//【3】返回
	var wechatMng = MiniMng{
		AppID:     miniC.AppID,
		AppSecret: miniC.AppSecret,
		Client:    mini,
	}
	return &wechatMng
}

// Login 微信小程序登陆
func (mng *MiniMng) Login(code string) (*auth.ResCode2Session, error) {
	authClient := mng.Client.GetAuth()
	res, err := authClient.Code2Session(code)
	return &res, err
}

// GetUserInfo 获取微信资料
func (mng *MiniMng) GetUserInfo(sessionKey, encryptedData, iv string) (*encryptor.PlainData, error) {
	encryptorClient := mng.Client.GetEncryptor()
	return encryptorClient.Decrypt(sessionKey, encryptedData, iv)
}

// GetPhone 获取微信手机号
func (mng *MiniMng) GetPhone(sessionKey, encryptedData, iv string) (*encryptor.PlainData, error) {
	encryptorClient := mng.Client.GetEncryptor()
	res, err := encryptorClient.Decrypt(sessionKey, encryptedData, iv)
	return res, err
}

// GetQRCode 获取二维码
func (mng *MiniMng) GetQRCode(qrCoder qrcode.QRCoder) ([]byte, error) {
	qrCodeApi := mng.Client.GetQRCode()
	res, err := qrCodeApi.GetWXACodeUnlimit(qrCoder)
	return res, err
}

// TextCheck 文字检测
func (mng *MiniMng) TextCheck(content string) (err error) {
	securityApi := mng.Client.GetContentSecurity()
	err = securityApi.CheckText(content)
	return
}

// ImgCheck 网络图片检测
func (mng *MiniMng) ImgCheck(imgURL string) (err error) {
	securityApi := mng.Client.GetContentSecurity()
	err = securityApi.CheckImage(imgURL)
	return
}

// ImgsCheck 网络图片检测
func (mng *MiniMng) ImgsCheck(imgURLs []string) (err error) {
	securityApi := mng.Client.GetContentSecurity()

	for k := range imgURLs {
		err = securityApi.CheckImage(imgURLs[k])
		if err != nil {
			break
		}
	}
	return
}

//func ImgCheck(url string)error{
//	// 网络图片检测
//	// @url 要检测的图片网络路径
//	// @token 接口调用凭证(access_token)
//	access_token := getWxApiAccessToken()
//	if len(url)==0{
//		return errors.New("图片地址为空")
//	}
//	res, err := weapp.IMGSecCheckFromNet(url,access_token)
//	if err != nil {
//		return err
//	}
//
//	if res.Errcode!=0{
//		return errors.New("图片可能存在敏感内容")
//	}
//	return nil
//}
//
//func TextCheck(content string)error{
//	// 文本检测
//	// @content 要检测的文本内容，长度不超过 500KB，编码格式为utf-8
//	// @token 接口调用凭证(access_token)
//	access_token := getWxApiAccessToken()
//	if len(content)==0{
//		return nil
//	}
//	res, err := weapp.MSGSecCheck(content,access_token)
//	if err != nil {
//		return err
//	}
//
//	if res.Errcode!=0{
//		return errors.New("文字可能存在敏感内容")
//	}
//	return nil
//}
