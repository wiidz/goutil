package aiMng

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"github.com/wiidz/goutil/mngs/eyunMng/msgSend"
	"github.com/wiidz/goutil/mngs/volcengineMng"
)

// aiSearch 调用模型执行一次简单会话。
func (m *Manager) aiSearch(ctx *gin.Context, prompt string, content string, returnTokenUsage bool) (replyContent string, err error) {
	if err = m.acquireLock(true); err != nil {
		return "", err
	}
	defer m.releaseLock()

	resp, err := m.aiMng.CreateChatCompletionRequestSimple(ctx, volcengineMng.Doubao, volcengineMng.Disabled, []*volcengineMng.ChatParam{
		{
			Role: volcengineMng.System,
			Text: prompt,
		},
		{
			Role: volcengineMng.User,
			Text: content,
		},
	})
	if err != nil {
		fmt.Printf("standard chat error: %v\n", err)
		return "", err
	}

	tokenUsedStr := ""
	if returnTokenUsage {
		tokenUsedStr = fmt.Sprintf("\n本次消耗：%dtokens", resp.Usage.TotalTokens)
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == nil || resp.Choices[0].Message.Content.StringValue == nil {
		return "", fmt.Errorf("AI 未返回有效内容")
	}

	return *resp.Choices[0].Message.Content.StringValue + tokenUsedStr, nil
}

// searchWithProductPrice 支持通过工具函数调用数据库和知识库。
func (m *Manager) searchWithProductPrice(ctx context.Context, sessionID string, prompt string, content string, messageTime *time.Time, returnTokenUsage bool) (messages []msgSend.MessagePayload, err error) {

	//【1】锁
	if err = m.acquireLock(true); err != nil {
		return nil, err
	}
	defer m.releaseLock()

	//【2】读取会话历史
	conversation, historyErr := m.loadSessionConversation(ctx, sessionID)
	if historyErr != nil {
		fmt.Printf("load session conversation error: %v\n", historyErr)
	}

	//【3】是否重置会话
	resetRequested := strings.TrimSpace(content) == m.config.ResetSignal
	if resetRequested {
		conversation = nil
		if delErr := m.redis.Set(ctx, m.sessionRedisKey(sessionID), "", 0); delErr != nil {
			fmt.Printf("reset session error: %v\n", delErr)
		}
		if clearErr := m.clearDirectSkuContext(ctx, sessionID); clearErr != nil {
			fmt.Printf("clear direct sku context error: %v\n", clearErr)
		}
		if clearErr := m.clearSessionProductContext(ctx, sessionID); clearErr != nil {
			fmt.Printf("clear product context error: %v\n", clearErr)
		}
		clearSessionProductSnapshot(sessionID)
		m.recordSessionMessageTimes(sessionID, nil)
		m.debugPrintSession(sessionID, "conversation reset", conversation)
	} else {
		m.debugPrintSession(sessionID, "loaded history", conversation)
	}

	if len(conversation) == 0 && prompt != "" {
		conversation = append(conversation, m.newChatMessage(sessionID, model.ChatMessageRoleSystem, prompt, nil))
	}

	if resetRequested {
		text := "已开启新的对话，请继续提问"
		saveErr := m.saveSessionConversation(ctx, sessionID, m.trimConversationHistory(conversation))
		if saveErr != nil {
			fmt.Printf("save session conversation error: %v\n", saveErr)
		}
		m.debugPrintSession(sessionID, "after reset reply", conversation)
		return []msgSend.MessagePayload{
			&msgSend.TextMessage{
				WcID:    sessionID,
				Content: text,
			},
		}, nil
	}

	conversation = append(conversation, m.newChatMessage(sessionID, model.ChatMessageRoleUser, content, messageTime))
	conversation = m.appendProductContextMessage(ctx, sessionID, conversation)
	m.debugPrintSession(sessionID, "after user message", conversation)

	if directMessages, handled, handleErr := m.maybeHandleDirectSkuQuery(ctx, sessionID, conversation, content); handled {
		if handleErr != nil {
			return nil, handleErr
		}
		return directMessages, nil
	}

	tools := []*model.Tool{
		{
			Type: model.ToolTypeFunction,
			Function: &model.FunctionDefinition{
				Name:        "query_product_price",
				Description: "根据产品编号或名称查询最新的产品价格（数据来源于PostgreSQL的sku表）。",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "number",
							"description": "产品编号，例如SKU或内部物料编号；如果已知请优先使用。",
						},
						"name_zh": map[string]interface{}{
							"type":        "string",
							"description": "产品名称，当没有编号时使用，支持模糊匹配。",
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"minimum":     1,
							"maximum":     20,
							"description": "返回的最大结果数量，默认5，最大20。",
						},
					},
					"required":             []string{},
					"additionalProperties": false,
				},
			},
		},
	}

	tools = append(tools, &model.Tool{
		Type: model.ToolTypeFunction,
		Function: &model.FunctionDefinition{
			Name:        "search_wiki",
			Description: "检索内部知识库以获取产品或技术问答摘要，并返回来源信息。",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"keyword": map[string]interface{}{
						"type":        "string",
						"description": "需要检索的关键词或问题。",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"minimum":     1,
						"maximum":     m.config.MaxWikiSearchLimit,
						"description": fmt.Sprintf("返回的最大结果数量，默认%d，最大%d。", m.config.DefaultWikiSearchLimit, m.config.MaxWikiSearchLimit),
					},
				},
				"required":             []string{"keyword"},
				"additionalProperties": false,
			},
		},
	})
	tools = append(tools, &model.Tool{
		Type: model.ToolTypeFunction,
		Function: &model.FunctionDefinition{
			Name:        "normalize_product_name",
			Description: "根据给定的产品别名或缩写，查询知识库返回标准化的产品名称以及匹配到的别名。",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "需要标准化的名称、别名或缩写，例如\"bh0.66\"。",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"minimum":     1,
						"maximum":     m.config.MaxWikiSearchLimit,
						"description": fmt.Sprintf("返回的最大候选数量，默认%d，最大%d。", m.config.DefaultWikiSearchLimit, m.config.MaxWikiSearchLimit),
					},
				},
				"required":             []string{"name"},
				"additionalProperties": false,
			},
		},
	})
	tools = append(tools, &model.Tool{
		Type: model.ToolTypeFunction,
		Function: &model.FunctionDefinition{
			Name:        "identify_product_focus",
			Description: "识别用户对话中提到的电气产品型号，并返回数据库中匹配的产品信息。",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"candidates": map[string]interface{}{
						"type":        "array",
						"description": "用户可能关注的产品名称、型号或别名列表。",
						"items":       map[string]interface{}{"type": "string"},
					},
				},
				"required":             []string{"candidates"},
				"additionalProperties": false,
			},
		},
	})

	totalTokens := 0
	const maxToolIterations = 3
	imageSet := make(map[string]struct{})
	imageList := make([]string, 0)

	for round := 0; round < maxToolIterations; round++ {
		conversation = m.appendProductContextMessage(ctx, sessionID, conversation)
		req := model.CreateChatCompletionRequest{
			Model:      string(volcengineMng.Doubao),
			Messages:   conversation,
			Tools:      tools,
			ToolChoice: model.ToolChoiceStringTypeAuto,
			Thinking: &model.Thinking{
				Type: volcengineMng.Disabled.GetThinkingType(),
			},
		}

		resp, reqErr := m.aiMng.Client.CreateChatCompletion(ctx, req)
		if reqErr != nil {
			fmt.Printf("chat with product tool error: %v\n", reqErr)
			return nil, reqErr
		}

		totalTokens += resp.Usage.TotalTokens
		if len(resp.Choices) == 0 {
			return nil, errors.New("AI未返回有效结果")
		}

		choice := resp.Choices[0]
		assistantMsg := choice.Message

		if len(choice.Message.ToolCalls) == 0 {
			var text string
			if assistantMsg.Content != nil && assistantMsg.Content.StringValue != nil {
				text = *assistantMsg.Content.StringValue
			}
			if returnTokenUsage {
				text += fmt.Sprintf("\n本次消耗：%dtokens", totalTokens)
			}

			trimmedText := strings.TrimSpace(text)

			assistantCopy := assistantMsg
			conversation = append(conversation, &assistantCopy)
			if trimmedText != "" {
				m.updateProductContextFromAssistant(ctx, sessionID, trimmedText)
			}
			m.debugPrintSession(sessionID, "before final save", conversation)

			saveErr := m.saveSessionConversation(ctx, sessionID, m.trimConversationHistory(conversation))
			if saveErr != nil {
				fmt.Printf("save session conversation error: %v\n", saveErr)
			}
			m.debugPrintSession(sessionID, "after final save", conversation)

			text = trimmedText
			if text != "" {
				messages = append(messages, &msgSend.TextMessage{
					WcID:    sessionID,
					Content: text,
				})
			}
			for _, url := range imageList {
				messages = append(messages, &msgSend.ImageMessage{
					WcID: sessionID,
					URL:  url,
				})
			}
			if len(messages) == 0 {
				return nil, errors.New("AI未返回有效结果")
			}
			return messages, nil
		}

		assistantCopy := assistantMsg
		conversation = append(conversation, &assistantCopy)

		toolMessages, toolErr := m.buildToolResponses(ctx, sessionID, choice.Message.ToolCalls)
		if toolErr != nil {
			fmt.Printf("tool execution error: %v\n", toolErr)
		}
		if thumbs := collectProductThumbnails(toolMessages); len(thumbs) > 0 {
			for _, url := range thumbs {
				if _, ok := imageSet[url]; ok {
					continue
				}
				imageSet[url] = struct{}{}
				imageList = append(imageList, url)
			}
		}
		conversation = append(conversation, toolMessages...)
		m.debugPrintSession(sessionID, fmt.Sprintf("after tool round %d", round+1), conversation)
	}

	saveErr := m.saveSessionConversation(ctx, sessionID, m.trimConversationHistory(conversation))
	if saveErr != nil {
		fmt.Printf("save session conversation error: %v\n", saveErr)
	}
	m.debugPrintSession(sessionID, "after timeout save", conversation)

	return nil, errors.New("工具调用超过最大次数限制")
}

// newChatMessage 构造一条模型对话消息并记录时间信息。
func (m *Manager) newChatMessage(sessionID string, role string, text string, messageTime *time.Time) *model.ChatCompletionMessage {
	msg := &model.ChatCompletionMessage{
		Role: role,
		Content: &model.ChatCompletionMessageContent{
			StringValue: volcengine.String(text),
		},
	}

	if m == nil || sessionID == "" {
		return msg
	}

	ts := time.Now()
	if messageTime != nil && !messageTime.IsZero() {
		ts = *messageTime
	}

	times := m.sessionMessageTimesFor(sessionID)
	times = append(times, messageTimestamp{
		Value: ts.Format(time.DateTime),
		Unix:  ts.Unix(),
	})
	m.recordSessionMessageTimes(sessionID, times)
	return msg
}

