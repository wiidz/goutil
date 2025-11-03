package eyunMng

import (
	"errors"

	"github.com/wiidz/goutil/mngs/eyunMng/account"
	"github.com/wiidz/goutil/mngs/eyunMng/base"
	"github.com/wiidz/goutil/mngs/eyunMng/chatRoom"
	"github.com/wiidz/goutil/mngs/eyunMng/login"
	"github.com/wiidz/goutil/mngs/eyunMng/msgReceive"
	"github.com/wiidz/goutil/mngs/eyunMng/msgSend"
	"github.com/wiidz/goutil/mngs/redisMng"
)

// NewEYunMng 返回一个e云管家实例
func NewEYunMng(config *base.Config, redisM *redisMng.RedisMng) (mng *EYunMng, err error) {

	if redisM == nil {
		err = errors.New("e云管家生成失败，请传入redis")
		return
	}

	mng = &EYunMng{
		Config:   config,
		RedisMng: redisM,
	}

	if mng.Config.Authorization == "" {
		var authorization string
		authorization, err = mng.GetLogin().LoginEYun()
		if err != nil {
			return
		}
		mng.Config.Authorization = authorization
	}
	return
}

// GetLogin 获取登录
func (mng *EYunMng) GetLogin() *login.Api {
	if mng.loginApi == nil {
		mng.loginApi = login.NewApi(mng.Config)
	}
	return mng.loginApi
}

// GetMsgSend 获取消息发送
func (mng *EYunMng) GetMsgSend() *msgSend.Api {

	if mng.msgSendApi == nil {
		mng.msgSendApi = msgSend.NewApi(mng.Config)
	}
	return mng.msgSendApi
}

// GetMsgReceive 获取消息接收
func (mng *EYunMng) GetMsgReceive() *msgReceive.Api {

	if mng.msgReceiveApi == nil {
		mng.msgReceiveApi = msgReceive.NewApi(mng.Config)
	}
	return mng.msgReceiveApi
}

// GetChatRoom 获取群聊
func (mng *EYunMng) GetChatRoom() *chatRoom.Api {

	if mng.chatRoom == nil {
		mng.chatRoom = chatRoom.NewApi(mng.Config, mng.RedisMng)
	}
	return mng.chatRoom
}

// GetAccount 获取账户
func (mng *EYunMng) GetAccount() *account.Api {

	if mng.account == nil {
		mng.account = account.NewApi(mng.Config)
	}
	return mng.account
}
