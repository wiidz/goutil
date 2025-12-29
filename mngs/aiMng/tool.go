package aiMng

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type productPriceToolArgs struct {
	ID     uint64 `json:"id"`
	NameZh string `json:"name_zh"`
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
}

// searchText 是 productPriceToolArgs 的辅助方法，生成用于搜索的关键词。
func (a productPriceToolArgs) searchText() string {
	if s := strings.TrimSpace(a.NameZh); s != "" {
		return s
	}
	if s := strings.TrimSpace(a.Query); s != "" {
		return s
	}
	return ""
}

type wikiSearchToolArgs struct {
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit"`
}

type wikiNormalizeNameArgs struct {
	Name  string `json:"name"`
	Limit int    `json:"limit"`
}

// productFocusToolArgs 用于 identify_product_focus 工具携带模型识别到的产品候选。
type productFocusToolArgs struct {
	Candidates []string `json:"candidates"`
}

// buildToolResponses 是工具调用处理函数，负责执行模型指定的函数调用并构造回复消息。
func (m *Manager) buildToolResponses(ctx context.Context, sessionID string, toolCalls []*model.ToolCall) ([]*model.ChatCompletionMessage, error) {
	wikiM := m.wikiMng
	messages := make([]*model.ChatCompletionMessage, 0, len(toolCalls))
	var firstErr error

	for _, call := range toolCalls {
		payload := map[string]interface{}{"success": false}

		switch call.Function.Name {
		case "query_product_price":
			var args productPriceToolArgs
			if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
				payload["message"] = fmt.Sprintf("参数解析失败: %v", err)
				if firstErr == nil {
					firstErr = err
				}
				break
			}

			searchText := args.searchText()
			if m.productPriceQuery == nil {
				return nil, errors.New("产品价格查询接口未初始化")
			}
			records, total, debugInfo, err := m.productPriceQuery.QueryProductPrices(
				ctx,
				ProductPriceQueryParams{
					ID:         args.ID,
					SearchText: searchText,
					Limit:      args.Limit,
				},
			)
			if debugInfo != nil {
				payload["debug"] = debugInfo
			}
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
				payload["message"] = err.Error()
				break
			}
			if total == 0 {
				payload["message"] = "未查询到匹配的产品价格"
				break
			}
			data, message := buildProductRecordsPayload(records, total)
			payload["success"] = true
			payload["data"] = data
			if message != "" {
				payload["message"] = message
			}
		case "search_wiki":
			if wikiM == nil {
				err := errors.New("知识库工具未配置")
				payload["message"] = err.Error()
				if firstErr == nil {
					firstErr = err
				}
				break
			}
			var args wikiSearchToolArgs
			if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
				payload["message"] = fmt.Sprintf("参数解析失败: %v", err)
				if firstErr == nil {
					firstErr = err
				}
				break
			}
			articles, totalCount, matchedCount, rewriteQuery, usage, err := m.queryWikiArticles(ctx, wikiM, args.Keyword, args.Limit)
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
				payload["message"] = err.Error()
				break
			}
			if matchedCount == 0 {
				keyword := strings.TrimSpace(args.Keyword)
				payload["message"] = fmt.Sprintf("未在知识库中找到与\"%s\"相关的内容", keyword)
				if keyword != "" && m.productPriceQuery != nil {
					productRecords, productTotal, debugInfo, prodErr := m.productPriceQuery.QueryProductPrices(
						ctx,
						ProductPriceQueryParams{
							SearchText: keyword,
							Limit:      m.config.DirectSkuDisplayLimit,
						},
					)
					if debugInfo != nil {
						payload["product_query_debug"] = debugInfo
					}
					if prodErr == nil && productTotal > 0 && len(productRecords) > 0 {
						data, message := buildProductRecordsPayload(productRecords, productTotal)
						if data != nil {
							data["source"] = "product_db"
						}
						payload["success"] = true
						payload["data"] = data
						if message != "" {
							payload["message"] = fmt.Sprintf("知识库未命中，但产品库%s", message)
						} else {
							payload["message"] = "知识库未命中，但产品库中找到相关产品"
						}
						payload["fallback_source"] = "product_db"
						break
					}
				}
				break
			}
			payload["success"] = true
			items := wikiArticlesToPayloadItems(articles)
			data := map[string]interface{}{
				"total":           totalCount,
				"matched":         matchedCount,
				"returned":        len(items),
				"items":           items,
				"query_rewritten": rewriteQuery,
			}
			if usageMap := wikiTokenUsageToMap(usage); usageMap != nil {
				data["token_usage"] = usageMap
			}
			payload["data"] = data
			if len(items) < matchedCount {
				payload["message"] = fmt.Sprintf("知识库命中%d/%d条，已展示前%d条", matchedCount, totalCount, len(items))
			} else {
				payload["message"] = fmt.Sprintf("知识库命中%d/%d条", matchedCount, totalCount)
			}
		case "normalize_product_name":
			if wikiM == nil {
				err := errors.New("知识库工具未配置")
				payload["message"] = err.Error()
				if firstErr == nil {
					firstErr = err
				}
				break
			}
			var args wikiNormalizeNameArgs
			if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
				payload["message"] = fmt.Sprintf("参数解析失败: %v", err)
				if firstErr == nil {
					firstErr = err
				}
				break
			}
			results, rewriteQuery, usage, err := m.queryWikiNameNormalizations(ctx, wikiM, args.Name, args.Limit)
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
				payload["message"] = err.Error()
				break
			}
			if len(results) == 0 {
				payload["message"] = fmt.Sprintf("未找到与“%s”匹配的标准化名称", strings.TrimSpace(args.Name))
				break
			}
			payload["success"] = true
			items := make([]map[string]interface{}, 0, len(results))
			for _, result := range results {
				item := map[string]interface{}{
					"standard_name": result.StandardName,
				}
				if len(result.Aliases) > 0 {
					item["aliases"] = result.Aliases
				}
				if len(result.MatchedValues) > 0 {
					item["matched_values"] = result.MatchedValues
				}
				if result.Score > 0 {
					item["score"] = result.Score
				}
				if result.DocID != "" {
					item["doc_id"] = result.DocID
				}
				if result.DocTitle != "" {
					item["doc_title"] = result.DocTitle
				}
				if result.Source != "" {
					item["source"] = result.Source
				}
				if result.Link != "" {
					item["link"] = result.Link
				}
				items = append(items, item)
			}
			data := map[string]interface{}{
				"total":           len(results),
				"returned":        len(items),
				"items":           items,
				"query_rewritten": rewriteQuery,
			}
			if usageMap := wikiTokenUsageToMap(usage); usageMap != nil {
				data["token_usage"] = usageMap
			}
			payload["data"] = data
			payload["message"] = fmt.Sprintf("找到%d条标准化名称候选", len(items))
		case "identify_product_focus":
			if sessionID == "" {
				err := errors.New("缺少会话标识，无法更新产品焦点")
				payload["message"] = err.Error()
				if firstErr == nil {
					firstErr = err
				}
				break
			}
			var args productFocusToolArgs
			if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
				payload["message"] = fmt.Sprintf("参数解析失败: %v", err)
				if firstErr == nil {
					firstErr = err
				}
				break
			}
			seen := make(map[string]struct{})
			focusEntries := make([]productFocusLock, 0, len(args.Candidates))
			items := make([]map[string]interface{}, 0, len(args.Candidates))

			for _, raw := range args.Candidates {
				query := strings.TrimSpace(raw)
				if query == "" {
					continue
				}
				lower := strings.ToLower(query)
				if _, ok := seen[lower]; ok {
					continue
				}
				seen[lower] = struct{}{}

				result, err := m.lookupSkuRecords(ctx, query, m.config.DirectSkuDisplayLimit)
				if err != nil {
					if firstErr == nil {
						firstErr = err
					}
					if _, ok := payload["message"]; !ok {
						payload["message"] = fmt.Sprintf("查询产品“%s”失败: %v", query, err)
					}
					continue
				}
				if result == nil || len(result.Records) == 0 {
					continue
				}

				focusEntries = append(focusEntries, productFocusLock{
					Query:  query,
					Result: result,
				})

				recordData, note := buildProductRecordsPayload(result.Records, result.Total)
				item := map[string]interface{}{
					"query":         query,
					"standard_name": strings.TrimSpace(result.StandardName),
					"search_used":   result.SearchUsed,
					"total":         result.Total,
					"products":      recordData,
				}
				if note != "" {
					item["note"] = note
				}
				if len(result.Images) > 0 {
					item["images"] = result.Images
				}
				items = append(items, item)
			}

			if len(focusEntries) == 0 {
				if _, ok := payload["message"]; !ok {
					payload["message"] = "未识别到可锁定的产品焦点"
				}
				break
			}

			if err := m.applyProductFocusFromLookup(ctx, sessionID, focusEntries); err != nil {
				if firstErr == nil {
					firstErr = err
				}
				payload["message"] = fmt.Sprintf("更新产品上下文失败: %v", err)
				break
			}

			payload["success"] = true
			payload["data"] = map[string]interface{}{
				"matched": len(items),
				"items":   items,
			}
			payload["message"] = fmt.Sprintf("已锁定%d个产品焦点", len(items))
		default:
			err := fmt.Errorf("未实现的工具：%s", call.Function.Name)
			payload["message"] = err.Error()
			if firstErr == nil {
				firstErr = err
			}
		}

		if _, ok := payload["message"]; !ok {
			if success, _ := payload["success"].(bool); success {
				payload["message"] = "ok"
			} else {
				payload["message"] = "执行失败"
			}
		}

		raw, err := json.Marshal(payload)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			raw = []byte(`{"success":false,"message":"响应编码失败"}`)
		}

		messages = append(messages, &model.ChatCompletionMessage{
			Role:       model.ChatMessageRoleTool,
			ToolCallID: call.ID,
			Content: &model.ChatCompletionMessageContent{
				StringValue: volcengine.String(string(raw)),
			},
		})
	}

	return messages, firstErr
}

// buildProductRecordsPayload 是数据整形函数，将产品价格结果封装为模型可读的 JSON。
func buildProductRecordsPayload(records []ProductPriceRecord, total int64) (map[string]interface{}, string) {
	items := make([]map[string]interface{}, 0, len(records))
	for _, record := range records {
		item := map[string]interface{}{
			"id":           record.ID,
			"sku_code":     record.SKUCode,
			"display_name": record.DisplayName,
			"attributes":   record.Attributes,
		}
		if record.PriceIn != nil {
			item["price_in"] = *record.PriceIn
		}
		if record.PricePromot != nil {
			item["price_promot"] = *record.PricePromot
		}
		if record.Currency != "" {
			item["currency"] = record.Currency
		}
		if record.Thumbnail != "" {
			item["thumbnail_img"] = record.Thumbnail
		}
		items = append(items, item)
	}
	returned := len(items)
	data := map[string]interface{}{
		"total":    total,
		"returned": returned,
		"items":    items,
	}

	message := ""
	if int64(returned) < total {
		message = fmt.Sprintf("共找到%d条结果，已展示前%d条", total, returned)
	} else {
		message = fmt.Sprintf("共找到%d条结果", total)
	}
	return data, message
}

type productToolPayload struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Items []struct {
			Thumbnail string `json:"thumbnail_img"`
		} `json:"items"`
	} `json:"data"`
}

// collectProductThumbnails 是工具输出解析函数，用于提取产品缩略图链接。
func collectProductThumbnails(toolMessages []*model.ChatCompletionMessage) []string {
	if len(toolMessages) == 0 {
		return nil
	}
	result := make([]string, 0)
	seen := make(map[string]struct{})

	for _, msg := range toolMessages {
		if msg == nil || msg.Content == nil || msg.Content.StringValue == nil {
			continue
		}

		raw := strings.TrimSpace(*msg.Content.StringValue)
		if raw == "" || !strings.HasPrefix(raw, "{") {
			continue
		}

		var payload productToolPayload
		if err := json.Unmarshal([]byte(raw), &payload); err != nil {
			continue
		}
		if len(payload.Data.Items) == 0 {
			continue
		}
		for _, item := range payload.Data.Items {
			url := strings.TrimSpace(item.Thumbnail)
			if url == "" {
				continue
			}
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			result = append(result, url)
		}
	}

	return result
}
