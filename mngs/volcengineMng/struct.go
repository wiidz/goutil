package volcengineMng

import (
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/wiidz/goutil/structs/configStruct"
)

type Config struct {
	ApiKey string
}

type VolcengineMng struct {
	Config *configStruct.VolcengineConfig
	Client *arkruntime.Client
}

type AIModel string

const (
	Doubao AIModel = "doubao-1-5-pro-32k-250115"
)

type Role string

const (
	User   Role = "user"
	System Role = "system"
)

type ChatParam struct {
	Role Role
	Text string
}

type ThinkingType string

const (
	Disabled ThinkingType = "disabled"
	Enabled  ThinkingType = "enabled"
	Auto     ThinkingType = "auto"
)

func (m ThinkingType) GetThinkingType() model.ThinkingType {
	if m == Disabled {
		return model.ThinkingTypeDisabled
	} else if m == Enabled {
		return model.ThinkingTypeEnabled
	} else if m == Auto {
		return model.ThinkingTypeAuto
	} else {
		return model.ThinkingTypeDisabled
	}
}
