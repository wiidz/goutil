package volcengineMng

import (
	"context"
	"errors"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
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

// CreateChatCompletionRequest 创建文本对话
func (mng *VolcengineMng) CreateChatCompletionRequest(ctx context.Context, aiModel AIModel, thinking ThinkingType, params []*ChatParam) (response model.ChatCompletionResponse, err error) {

	//【1】处理对话文本
	var msgs = []*model.ChatCompletionMessage{}
	for _, v := range params {
		temp := &model.ChatCompletionMessage{
			Role: string(v.Role),
			Content: &model.ChatCompletionMessageContent{
				StringValue: volcengine.String(v.Text),
			},
		}
		msgs = append(msgs, temp)
	}

	//【2】组合参数
	req := model.CreateChatCompletionRequest{
		Model:    string(aiModel),
		Messages: msgs,
		Thinking: &model.Thinking{
			Type: thinking.GetThinkingType(),
		},
	}

	//【3】发送请求
	response, err = mng.Client.CreateChatCompletion(ctx, req)
	return
}
