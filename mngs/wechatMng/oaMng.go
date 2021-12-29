package wechatMng

import (
	"fmt"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
	"github.com/wiidz/goutil/structs/configStruct"
	"net/http"
)

// WechatOaMng 微信公众号管理器
type WechatOaMng struct {
	Config *configStruct.WechatOaConfig
	//AccessToken string
	Client *officialaccount.OfficialAccount
}

// NewWechatOaMng 获取公众号管理器
func NewWechatOaMng(config *configStruct.WechatOaConfig) *WechatOaMng {

	//【1】使用redis缓存accessToken
	memory := cache.NewMemory() // accessToken存在内存中

	//【2】创建mini实例
	cfg := &offConfig.Config{
		AppID:          config.AppID,
		AppSecret:      config.AppSecret,
		Token:          config.Token,
		EncodingAESKey: config.EncodingAESKey,
		Cache:          memory,
	}
	wc := wechat.NewWechat()
	off := wc.GetOfficialAccount(cfg)

	//【3】返回
	var wechatMng = WechatOaMng{
		Config: config,
		Client: off,
	}
	return &wechatMng
}

// Login 微信公众号登陆
func (mng *WechatOaMng) Login(code string) (*oauth.ResAccessToken, error) {
	u := mng.Client.GetOauth()
	res, err := u.GetUserAccessToken(code)
	return &res, err
}

// Notify 微信公众号登陆
func (mng *WechatOaMng) Notify(rw http.ResponseWriter, req *http.Request) {
	oaServer := mng.Client.GetServer(req, rw)
	//设置接收消息的处理方法
	oaServer.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		//TODO
		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err := oaServer.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	_ = oaServer.Send()
}
