package volcengineMng

import (
	"errors"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/wiidz/goutil/structs/configStruct"
)

// NewVolcengineMng 创建一个Volcengine管理器
func NewVolcengineMng(config *configStruct.VolcengineConfig) (mng *VolcengineMng, err error) {
	if config.ApiKey == "" {
		err = errors.New("idr配置参数有误")
		return
	}
	mng = &VolcengineMng{
		Config: config,
	}

	mng.initClient()

	return
}

func (mng *VolcengineMng) initClient() {
	mng.Client = arkruntime.NewClientWithApiKey(mng.Config.ApiKey)
	return
}
