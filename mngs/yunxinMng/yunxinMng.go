package yunxinMng

import (
	"github.com/wiidz/goutil/mngs/yunxinMng/imClient"
	"github.com/wiidz/goutil/mngs/yunxinMng/imEvent"
	"github.com/wiidz/goutil/mngs/yunxinMng/imFriend"
	"github.com/wiidz/goutil/mngs/yunxinMng/imHistory"
	"github.com/wiidz/goutil/mngs/yunxinMng/imMsg"
	"github.com/wiidz/goutil/mngs/yunxinMng/imTeam"
	"github.com/wiidz/goutil/mngs/yunxinMng/imUser"
	"github.com/wiidz/goutil/structs/configStruct"
)

type YunxinMng struct {
	Config *configStruct.YunxinConfig
	Client *imClient.Client
}

// NewYunxinMng 获取云信管理器
func NewYunxinMng(config *configStruct.YunxinConfig) *YunxinMng {
	return &YunxinMng{
		Config: config,
		Client: imClient.NewClient(config.AppKey, config.AppSecret),
	}
}

// GetUser 获取用户Api
func (mng *YunxinMng) GetUser() *imUser.Api {
	return &imUser.Api{
		Client: mng.Client,
	}
}

// GetFriend 获取好友关系Api
func (mng *YunxinMng) GetFriend() *imFriend.Api {
	return &imFriend.Api{
		Client: mng.Client,
	}
}

// GetHistory 获取历史记录Api
func (mng *YunxinMng) GetHistory() *imHistory.Api {
	return &imHistory.Api{
		Client: mng.Client,
	}
}

// GetTeam 获取群组Api
func (mng *YunxinMng) GetTeam() *imTeam.Api {
	return &imTeam.Api{
		Client: mng.Client,
	}
}

// GetEvent 获取事件Api
func (mng *YunxinMng) GetEvent() *imEvent.Api {
	return &imEvent.Api{
		Client: mng.Client,
	}
}

// GetMsg 获取消息Api
func (mng *YunxinMng) GetMsg() *imMsg.Api {
	return &imMsg.Api{
		Client: mng.Client,
	}
}
