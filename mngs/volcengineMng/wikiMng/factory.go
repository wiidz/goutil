package wikiMng

// GetSearchKnowledgeReqParams 工厂模式返回搜索知识库的参数
// 知识库检索请求参数生成：以下详细展示了部分参数的传递规则，其余参数请参考官方接口文档，如果您想快速接入进行测试，可以只传入参数集合的最小集，其余参数可以使用默认值。
// 必传参数如下：
// 1. resourceId (也可使用resource_id 或 name + project, 二选一)
// 2. query：用户问题
func (mng *WikiMng) GetSearchKnowledgeReqParams(collectionName, project, query string) *CollectionSearchKnowledgeRequest {
	return &CollectionSearchKnowledgeRequest{
		Name:    collectionName, // 知识库名称
		Project: project,        // 知识库项目名称
		//ResourceId:  ResourceID,     // 知识库resource_id (二选一，查询时，可使用resource_id 或 name + project)
		Query:       query, // 用户问题
		Limit:       10,    // 返回数量, 不传递默认返回10条
		DenseWeight: 0.5,   //混合搜索的权重
		Preprocessing: &PreProcessing{
			NeedInstruction:  true,
			ReturnTokenUsage: true,
			Rewrite:          false, // 问题改写开关，默认不开启
			Messages: []*MessageParam{ // 仅在使用改写或意图识别时需要传且必传Messages
				{
					Role:    "system",
					Content: ChatCompletionMessageContent{},
				},
				{
					Role: "user",
					Content: ChatCompletionMessageContent{
						StringValue: &query,
					},
				},
			},
		},
		Postprocessing: &PostProcessing{
			RerankSwitch:        false, // 重排开关，默认不开启
			RetrieveCount:       25,    //进入重排的切片数量，重排打开时生效，需要大于limit,当limit=10，默认值为25
			GetAttachmentLink:   true,  // 是否返回原始图片，仅当创建知识库开启 OCR 时生效，否则自动跳过图片
			ChunkGroup:          true,  //是否对召回切片按照文档进行聚合
			ChunkDiffusionCount: 0,     //切片扩散数量-是否召回切片的临近切片，如 1 代表额外召回当前切片的上下各一个切片
		},
	}
}

func (mng *WikiMng) GenerateChatCompletionReqParams(modelSetting *AIModelSetting, stream bool, messages []*MessageParam) *CollectionChatCompletionRequest {
	return &CollectionChatCompletionRequest{
		Model:            string(modelSetting.Model), // 如果使用私有ep，此处替换为私有ep即可，格式 ep-xxx-xxx
		ModelVersion:     modelSetting.ModelVersion,  // 模型版本，使用公有接入点时，可以选择指定模型版本，不指定则服务会自动指定默认版本
		Stream:           stream,                     // 模型结果是否流式返回
		ReturnTokenUsage: true,                       // 是否返回token使用情况
		MaxTokens:        modelSetting.MaxTokens,     // 最大token数
		Temperature:      modelSetting.Temperature,   // 模型温度,取值范围0~1，值越大随机性越大
		APIKey:           mng.Config.ApiKey,          // 使用私有ep时，必须传递此参数才能生效
		Messages:         messages,                   // 模型对话信息
	}
}
