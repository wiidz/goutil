package wechatMng

import (
	"context"
	"errors"
	"unicode/utf8"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/encryptor"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
	"github.com/silenceper/wechat/v2/miniprogram/security"
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/helpers/osHelper"
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
func NewMiniMng(ctx context.Context, redisC *configStruct.RedisConfig, miniC *configStruct.WechatMiniConfig) *MiniMng {

	//【1】使用redis缓存accessToken
	// memory := cache.NewMemory() // accessToken存在内存中
	redisOpts := &cache.RedisOpts{
		Host:     redisC.Host + ":" + redisC.Port,
		Password: redisC.Password,
	}
	redisCache := cache.NewRedis(ctx, redisOpts)

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
func (mng *MiniMng) TextCheck(userOpenID string, content string, scene security.MsgScene) (res security.MsgCheckResponse, err error) {

	totalLen := utf8.RuneCountInString(content)
	if totalLen > 2500 {
		err = errors.New("校验失败，文字不能超过2500字")
		return
	}

	securityApi := mng.Client.GetSecurity()

	res, err = securityApi.MsgCheck(&security.MsgCheckRequest{
		OpenID:    userOpenID, // 用户的openid（用户需在近两小时访问过小程序）
		Scene:     scene,      // 场景枚举值（1 资料；2 评论；3 论坛；4 社交日志）
		Content:   content,    // 需检测的文本内容，文本字数的上限为 2500 字，需使用 UTF-8 编码
		Nickname:  "",         // （非必填）用户昵称，需使用UTF-8编码
		Title:     "",         // （非必填）文本标题，需使用UTF-8编码
		Signature: "",         // （非必填）个性签名，该参数仅在资料类场景有效(scene=1)，需使用UTF-8编码
	})
	return
}

// ImgCheck 本地图片检测
func (mng *MiniMng) ImgCheck(userOpenID string, imgURL string, scene security.MsgScene) (traceID string, err error) {
	securityApi := mng.Client.GetSecurity()
	traceID, err = securityApi.MediaCheckAsync(&security.MediaCheckAsyncRequest{
		MediaURL:  imgURL,       // 要检测的图片或音频的url，支持图片格式包括jpg, jepg, png, bmp, gif（取首帧），支持的音频格式包括mp3, aac, ac3, wma, flac, vorbis, opus, wav
		MediaType: uint8(Image), // 1:音频;2:图片
		OpenID:    userOpenID,   // 用户的openid（用户需在近两小时访问过小程序）
		Scene:     uint8(scene), // 场景枚举值（1 资料；2 评论；3 论坛；4 社交日志）
	})
	return
}

// ImgsCheck 本地图片检测
func (mng *MiniMng) ImgsCheck(userOpenID string, imgURLs []string, scene security.MsgScene) (traceIDs []string, err error) {
	securityApi := mng.Client.GetSecurity()

	traceIDs = []string{}
	for k := range imgURLs {
		var traceID string
		traceID, err = securityApi.MediaCheckAsync(&security.MediaCheckAsyncRequest{
			MediaURL:  imgURLs[k],   // 要检测的图片或音频的url，支持图片格式包括jpg, jepg, png, bmp, gif（取首帧），支持的音频格式包括mp3, aac, ac3, wma, flac, vorbis, opus, wav
			MediaType: uint8(Image), // 1:音频;2:图片
			OpenID:    userOpenID,   // 用户的openid（用户需在近两小时访问过小程序）
			Scene:     uint8(scene), // 场景枚举值（1 资料；2 评论；3 论坛；4 社交日志）
		})
		if err != nil {
			return
		}
		traceIDs = append(traceIDs, traceID)
	}

	return
}

// NetworkImgCheck 网络图片检测
func (mng *MiniMng) NetworkImgCheck(userOpenID string, imgURL string, scene security.MsgScene) (traceID string, err error) {
	securityApi := mng.Client.GetSecurity()
	traceID, err = securityApi.MediaCheckAsync(&security.MediaCheckAsyncRequest{
		MediaURL:  imgURL,       // 要检测的图片或音频的url，支持图片格式包括jpg, jepg, png, bmp, gif（取首帧），支持的音频格式包括mp3, aac, ac3, wma, flac, vorbis, opus, wav
		MediaType: uint8(Image), // 1:音频;2:图片
		OpenID:    userOpenID,   // 用户的openid（用户需在近两小时访问过小程序）
		Scene:     uint8(scene), // 场景枚举值（1 资料；2 评论；3 论坛；4 社交日志）
	})
	return
}

// NetworkImgsCheck 网络图片检测
func (mng *MiniMng) NetworkImgsCheck(userOpenID string, imgURLs []string, scene security.MsgScene) (traceIDs []string, err error) {

	securityApi := mng.Client.GetSecurity()
	localPaths := []string{}

	go func() {
		osHelper.DeleteFiles(localPaths)
	}()

	traceIDs = []string{}
	for k := range imgURLs {
		//【2】下载文件到本地
		var tempPath string
		_, tempPath, err = networkHelper.DownloadFile(imgURLs[k], "")
		if err != nil {
			break
		}
		localPaths = append(localPaths, tempPath)
		var traceID string
		traceID, err = securityApi.MediaCheckAsync(&security.MediaCheckAsyncRequest{
			MediaURL:  imgURLs[k],   // 要检测的图片或音频的url，支持图片格式包括jpg, jepg, png, bmp, gif（取首帧），支持的音频格式包括mp3, aac, ac3, wma, flac, vorbis, opus, wav
			MediaType: 2,            // 1:音频;2:图片
			OpenID:    userOpenID,   // 用户的openid（用户需在近两小时访问过小程序）
			Scene:     uint8(scene), // 场景枚举值（1 资料；2 评论；3 论坛；4 社交日志）
		})
		if err != nil {
			break
		}
		traceIDs = append(traceIDs, traceID)
	}
	return
}

func (mng *MiniMng) AudioCheck(userOpenID string, imgURL string, scene security.MsgScene) (traceID string, err error) {
	securityApi := mng.Client.GetSecurity()
	traceID, err = securityApi.MediaCheckAsync(&security.MediaCheckAsyncRequest{
		MediaURL:  imgURL,       // 要检测的图片或音频的url，支持图片格式包括jpg, jepg, png, bmp, gif（取首帧），支持的音频格式包括mp3, aac, ac3, wma, flac, vorbis, opus, wav
		MediaType: uint8(Audio), // 1:音频;2:图片
		OpenID:    userOpenID,   // 用户的openid（用户需在近两小时访问过小程序）
		Scene:     uint8(scene), // 场景枚举值（1 资料；2 评论；3 论坛；4 社交日志）
	})
	return
}

// AudiosCheck 多个本地音频文件检测
func (mng *MiniMng) AudiosCheck(userOpenID string, audioURLs []string, scene security.MsgScene) (traceIDs []string, err error) {
	securityApi := mng.Client.GetSecurity()
	traceIDs = []string{}
	for k := range audioURLs {
		var traceID string
		traceID, err = securityApi.MediaCheckAsync(&security.MediaCheckAsyncRequest{
			MediaURL:  audioURLs[k], // 要检测的图片或音频的url，支持图片格式包括jpg, jepg, png, bmp, gif（取首帧），支持的音频格式包括mp3, aac, ac3, wma, flac, vorbis, opus, wav
			MediaType: uint8(Audio), // 1:音频;2:图片
			OpenID:    userOpenID,   // 用户的openid（用户需在近两小时访问过小程序）
			Scene:     uint8(scene), // 场景枚举值（1 资料；2 评论；3 论坛；4 社交日志）
		})
		traceIDs = append(traceIDs, traceID)
	}
	return
}
