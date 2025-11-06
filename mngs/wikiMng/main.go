package wikiMng

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/volcengine/volc-sdk-golang/base"
	"github.com/wiidz/goutil/mngs/volcengineMng"
	"github.com/wiidz/goutil/structs/configStruct"
)

// 配置常量 - 请根据实际情况修改
const (
	// Domain 知识库服务域名
	Domain = "api-knowledgebase.mlp.cn-beijing.volces.com"
	// SearchKnowledgePath  知识库检索接口，建议您首次接入时使用该检索接口，其他检索接口后续不再进行维护
	SearchKnowledgePath = "/api/knowledge/collection/search_knowledge"
	// ChatCompletionPath 大模型对话接口，可以和检索接口接合串联RAG流程，也可以单独使用进行生成
	ChatCompletionPath = "/api/knowledge/chat/completions"
)

func NewWikiMng(config *configStruct.WikiConfig) *WikiMng {
	return &WikiMng{
		Config: config,
	}
}

// SearchKnowledge 从知识库搜索
func (mng *WikiMng) SearchKnowledge(ctx context.Context, searchKnowledgeReqParams *CollectionSearchKnowledgeRequest) (*CollectionSearchKnowledgeResponse, error) {
	searchKnowledgeReqParamsBytes, err := mng.serializeToJsonBytesUseNumber(searchKnowledgeReqParams)
	if err != nil {
		return nil, err
	}
	req := mng.prepareRequest("POST", SearchKnowledgePath, searchKnowledgeReqParamsBytes)
	client := &http.Client{Timeout: mng.Config.SimpleTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var searchKnowledgeResp *CollectionSearchKnowledgeResponse
	err = ParseJsonUseNumber(body, &searchKnowledgeResp)
	if err != nil {
		return nil, err
	}
	return searchKnowledgeResp, nil
}

// ChatCompletion 非流式调用
func (mng *WikiMng) ChatCompletion(ctx context.Context, chatCompletionReqParams *CollectionChatCompletionRequest) (*CollectionChatCompletionResponse, error) {
	chatCompletionReqParams.Stream = false
	chatCompletionReqParamsBytes, err := mng.serializeToJsonBytesUseNumber(chatCompletionReqParams)
	if err != nil {
		return nil, err
	}

	request := mng.prepareRequest("POST", ChatCompletionPath, chatCompletionReqParamsBytes)
	client := &http.Client{
		Timeout: mng.Config.StreamTimeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var chatCompletionResp *CollectionChatCompletionResponse
	err = ParseJsonUseNumber(body, &chatCompletionResp)
	if err != nil {
		return nil, err
	}
	return chatCompletionResp, nil
}

// ChatCompletionStream 流式调用
func (mng *WikiMng) ChatCompletionStream(ctx context.Context, chatCompletionReqParams *CollectionChatCompletionRequest) (answer string, usage *ModelTokenUsage, err error) {
	chatCompletionReqParams.Stream = true
	chatCompletionReqParamsBytes, err := mng.serializeToJsonBytesUseNumber(chatCompletionReqParams)
	if err != nil {
		return "", nil, err
	}

	request := mng.prepareRequest("POST", ChatCompletionPath, chatCompletionReqParamsBytes)
	client := &http.Client{
		Timeout: mng.Config.StreamTimeout,
	}
	request.Header.Set("Accept", "text/event-stream")
	resp, err := client.Do(request)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	// 读取流式返回
	scanner := bufio.NewScanner(resp.Body)
	// 指定分隔符函数
	scanner.Split(scanDoubleCRLF)

	var answerBuilder strings.Builder
	var modelTokenUsage ModelTokenUsage

	buf := make([]byte, 0, 150*1024)
	scanner.Buffer(buf, 150*1024) // 可以按需调整scanner的大小

	// 读取数据
	for scanner.Scan() {
		streamLine := scanner.Text()
		if len(streamLine) < 5 {
			continue
		}
		streamJson := streamLine[5:]
		var chatCompletionResponse CollectionChatCompletionResponse
		err := ParseJsonUseNumber([]byte(streamJson), &chatCompletionResponse)
		if err != nil {
			return "", nil, err
		}
		// 获取流式返回的内容
		fmt.Println(chatCompletionResponse.Data.GenerateAnswer)

		answerBuilder.WriteString(chatCompletionResponse.Data.GenerateAnswer)

		// 最后一条流式返回中，携带本次请求token使用信息
		if chatCompletionResponse.Data.Usage != "" {
			err := ParseJsonUseNumber([]byte(chatCompletionResponse.Data.Usage), &modelTokenUsage)
			if err != nil {
				return "", nil, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", nil, err
	}
	return answerBuilder.String(), &modelTokenUsage, nil
}

// getContentForPrompt 生成内容提示
func (mng *WikiMng) getContentForPrompt(item *CollectionSearchResponseItem, imageNum int) string {
	content := item.Content

	if item.OriginalQuestion != "" {
		return fmt.Sprintf("当询问到相似问题时，请参考对应答案进行回答：问题：“%s”。答案：“%s”",
			item.OriginalQuestion, content)
	}

	if imageNum > 0 && len(item.ChunkAttachmentList) > 0 && item.ChunkAttachmentList[0].Link != "" {
		placeholder := fmt.Sprintf("<img>图片%d</img>", imageNum)
		return content + placeholder
	}

	return content
}

// GeneratePrompt 将知识库检索到的多个文档片段，按照用户配置的字段要求，拼接成一个完整的提示词。
// basePrompt 基础提示词
// resp 知识库搜索结果
func (mng *WikiMng) GeneratePrompt(model volcengineMng.AIModel, basePrompt string, fieldSetting *PromptExtraContext, resp *CollectionSearchKnowledgeResponse) (finalPrompt string, imageURLs []string, err error) {

	//【1】判断查询结果
	if resp == nil {
		err = fmt.Errorf("response is nil")
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf(resp.Message)
		return
	}
	if fieldSetting == nil {
		err = fmt.Errorf("fieldSetting is nil")
		return
	}

	var promptBuilder strings.Builder
	imageCnt := 0

	for _, point := range resp.Data.ResultList {
		// 对vision模型需要额外处理图片链接
		//if model.IsVisionModel() && len(point.ChunkAttachmentList) > 0 {
		if false && len(point.ChunkAttachmentList) > 0 {
			link := point.ChunkAttachmentList[0].Link
			if link != "" {
				imageURLs = append(imageURLs, link)
				imageCnt++
			}
		}

		// 处理系统字段
		docInfo := point.DocInfo
		// 拼接用户指定的系统字段
		for _, sysField := range fieldSetting.SystemFields {
			switch sysField {
			case SysFieldDocName:
				promptBuilder.WriteString(fmt.Sprintf("%s: %s\n", sysField, docInfo.DocName))
			case SysFieldTitle:
				promptBuilder.WriteString(fmt.Sprintf("%s: %s\n", sysField, docInfo.Title))
			case SysFieldChunkTitle:
				promptBuilder.WriteString(fmt.Sprintf("%s: %s\n", sysField, point.ChunkTitle))
			case SysFieldContent:
				promptBuilder.WriteString(fmt.Sprintf("%s: %s\n", sysField, mng.getContentForPrompt(point, imageCnt)))
			}
		}

		// 结构化数据- 拼接用户指定的自定义字段
		for _, selfField := range fieldSetting.SelfDefineFields {
			for _, tableChunkField := range point.TableChunkFields {
				if tableChunkField.FieldName == selfField {
					promptBuilder.WriteString(fmt.Sprintf("%s: %v\n", tableChunkField.FieldName, tableChunkField.FieldValue))
				}
			}
		}
		promptBuilder.WriteString("---\n")
	}

	// 基础提示词模板替换
	finalPrompt = strings.Replace(basePrompt, "{prompt}", promptBuilder.String(), -1)
	return finalPrompt, imageURLs, nil
}

// RAG 检索增强生成流程串联
func (mng *WikiMng) RAG(ctx context.Context, model volcengineMng.AIModel, basePrompt string, fieldSetting *PromptExtraContext, chatParams *CollectionChatCompletionRequest, params *CollectionSearchKnowledgeRequest) error {

	// 知识库检索
	searchResp, err := mng.SearchKnowledge(ctx, params)
	if err != nil {
		return err
	}

	// 生成提示词
	prompt, images, err := mng.GeneratePrompt(model, basePrompt, fieldSetting, searchResp)
	if err != nil {
		return err
	}
	fmt.Printf("提示词：%s\n", prompt)

	// 生成Chat的Message结构体，拼接message对话, 问题对应role为user，系统对应role为system, 答案对应role为assistant, 内容对应content
	var messages []*MessageParam
	if len(images) > 0 {
		// 对于Vision模型，需要将图片链接拼接到Message中
		var multiModalMessage []*ChatCompletionMessageContentPart
		multiModalMessage = append(multiModalMessage, &ChatCompletionMessageContentPart{
			Type: ChatCompletionMessageContentPartTypeText,
			Text: params.Query,
		})
		for _, imageURL := range images {
			multiModalMessage = append(multiModalMessage, &ChatCompletionMessageContentPart{
				Type:     ChatCompletionMessageContentPartTypeImageURL,
				ImageURL: &ChatMessageImageURL{URL: imageURL},
			})
		}

		messages = []*MessageParam{
			{
				Role:    "system",
				Content: prompt,
			},
			{
				Role:    "user",
				Content: multiModalMessage,
			},
		}
	} else {
		// 如果使用的是普通的文本LLM模型，使用该分支拼接生成message
		messages = []*MessageParam{
			{
				Role:    "system",
				Content: prompt,
			},
			{
				Role:    "user",
				Content: params.Query,
			},
		}
	}
	chatParams.Messages = messages

	if chatParams.Stream {
		// 流式调用
		answer, usage, err := mng.ChatCompletionStream(ctx, chatParams)
		if err != nil {
			return err
		}
		fmt.Printf("大模型流式调用返回结果：%s\n", answer)
		fmt.Printf("大模型流式调用返回token使用情况：%+v\n", usage)
	} else {
		// 非流式调用
		ChatCompletionResponse, err := mng.ChatCompletion(ctx, chatParams)
		if err != nil {
			return err
		}
		if ChatCompletionResponse.Code != 0 {
			return fmt.Errorf(ChatCompletionResponse.Message)
		}

		answer := ChatCompletionResponse.Data.GenerateAnswer
		fmt.Printf("非流式大模型返回结果：%s\n", answer)

		var modelTokenUsage ModelTokenUsage
		err = ParseJsonUseNumber([]byte(ChatCompletionResponse.Data.Usage), &modelTokenUsage)
		if err != nil {
			return err
		}
		fmt.Printf("非流式大模型返回token使用情况：%+v\n", modelTokenUsage)
	}
	return nil
}

// prepareRequest 准备发送请求
func (mng *WikiMng) prepareRequest(method string, path string, body []byte) *http.Request {
	u := url.URL{
		Scheme: "https",
		Host:   Domain,
		Path:   path,
	}
	req, _ := http.NewRequest(strings.ToUpper(method), u.String(), bytes.NewReader(body))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", Domain)
	credential := base.Credentials{
		AccessKeyID:     mng.Config.AccessKeyID,
		SecretAccessKey: mng.Config.SecretKey,
		Service:         "air",
		Region:          "cn-north-1",
	}
	req = credential.Sign(req)
	return req
}
