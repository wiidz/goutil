package aiMng

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/wiidz/goutil/mngs/volcengineMng/wikiMng"
)

type wikiNameNormalization struct {
	StandardName  string   `json:"standard_name"`
	Aliases       []string `json:"aliases,omitempty"`
	MatchedValues []string `json:"matched_values,omitempty"`
	Score         float64  `json:"score,omitempty"`
	DocID         string   `json:"doc_id,omitempty"`
	DocTitle      string   `json:"doc_title,omitempty"`
	Source        string   `json:"source,omitempty"`
	Link          string   `json:"link,omitempty"`
}

type wikiArticle struct {
	Question   string  `json:"question"`
	Answer     string  `json:"answer"`
	Content    string  `json:"content"`
	Score      float64 `json:"score"`
	DocID      string  `json:"doc_id"`
	DocTitle   string  `json:"doc_title"`
	Source     string  `json:"source"`
	ChunkTitle string  `json:"chunk_title"`
	Link       string  `json:"link"`
	RecallRank int     `json:"recall_rank"`
	RerankRank int     `json:"rerank_rank"`
	ChunkID    int     `json:"chunk_id"`
}

type wikiTokenUsage struct {
	EmbeddingTokens int64 `json:"embedding_tokens"`
	RerankTokens    int64 `json:"rerank_tokens"`
	LLMTokens       int64 `json:"llm_tokens"`
	RewriteTokens   int64 `json:"rewrite_tokens"`
	TotalTokens     int64 `json:"total_tokens"`
}

// queryWikiNameNormalizations 是知识库检索函数，用于查询产品名称的标准化结果。
func (m *Manager) queryWikiNameNormalizations(ctx context.Context, wikiM *wikiMng.WikiMng, name string, limit int) ([]wikiNameNormalization, string, *wikiTokenUsage, error) {
	if wikiM == nil {
		return nil, "", nil, errors.New("知识库未初始化")
	}
	query := strings.TrimSpace(name)
	if query == "" {
		return nil, "", nil, errors.New("待标准化名称不能为空")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if limit <= 0 {
		limit = m.config.DefaultNameNormalizationLimit
	}
	if limit > m.config.MaxWikiSearchLimit {
		limit = m.config.MaxWikiSearchLimit
	}

	params := wikiM.GetSearchKnowledgeReqParams(m.config.WikiCollectionID, query, 0.5, 1)
	if params == nil {
		return nil, "", nil, errors.New("构建知识库检索参数失败")
	}
	searchLimit := limit * 2
	if searchLimit < m.config.DefaultNameNormalizationLimit {
		searchLimit = m.config.DefaultNameNormalizationLimit
	}
	if searchLimit > m.config.MaxWikiSearchLimit {
		searchLimit = m.config.MaxWikiSearchLimit
	}
	params.Limit = int32(searchLimit)

	resp, err := wikiM.SearchKnowledge(ctx, params)
	if err != nil {
		if isTimeoutError(err) {
			return nil, "", nil, fmt.Errorf("知识库检索超时，请稍后重试（%s）", err.Error())
		}
		return nil, "", nil, err
	}
	if resp == nil || resp.Data == nil {
		return []wikiNameNormalization{}, "", nil, nil
	}

	normalizedQuery := normalizeNameForMatch(query)
	resultMap := make(map[string]*wikiNameNormalization)

	for _, item := range resp.Data.ResultList {
		if item == nil || item.Score < m.config.WikiNormalizationMinScore {
			continue
		}

		parsed := extractNameNormalizationFromItem(item)
		if parsed.StandardName == "" && len(parsed.Aliases) == 0 {
			continue
		}

		parsed.Score = item.Score
		parsed.DocID = strings.TrimSpace(item.DocInfo.Docid)
		docTitle := strings.TrimSpace(item.DocInfo.DocName)
		if docTitle == "" {
			docTitle = strings.TrimSpace(item.DocInfo.Title)
		}
		parsed.DocTitle = docTitle
		parsed.Source = strings.TrimSpace(item.DocInfo.Source)
		if len(item.ChunkAttachmentList) > 0 {
			parsed.Link = strings.TrimSpace(item.ChunkAttachmentList[0].Link)
		}

		standardNormalized := normalizeNameForMatch(parsed.StandardName)
		matchedValues := make([]string, 0, len(parsed.Aliases)+1)
		if normalizedQuery != "" {
			if standardNormalized == normalizedQuery {
				matchedValues = appendUniqueFold(matchedValues, parsed.StandardName)
			}
			for _, alias := range parsed.Aliases {
				if normalizeNameForMatch(alias) == normalizedQuery {
					matchedValues = appendUniqueFold(matchedValues, alias)
				}
			}
			if len(matchedValues) == 0 {
				combined := strings.ToLower(parsed.StandardName + " " + strings.Join(parsed.Aliases, " "))
				if strings.Contains(combined, strings.ToLower(query)) {
					matchedValues = appendUniqueFold(matchedValues, query)
				}
			}
		}
		parsed.MatchedValues = dedupeStringsCaseInsensitive(matchedValues)

		parsed.Aliases = dedupeStringsCaseInsensitive(parsed.Aliases)
		preferred := choosePreferredStandardName(parsed.StandardName, parsed.Aliases, query, normalizedQuery)
		if !strings.EqualFold(preferred, parsed.StandardName) && strings.TrimSpace(parsed.StandardName) != "" {
			parsed.Aliases = appendUniqueFold(parsed.Aliases, parsed.StandardName)
			parsed.StandardName = preferred
		}
		parsed.Aliases = removeCaseInsensitive(parsed.Aliases, parsed.StandardName)
		parsed.MatchedValues = dedupeStringsCaseInsensitive(parsed.MatchedValues)

		key := normalizeNameForMatch(parsed.StandardName)
		if key == "" {
			key = strings.ToLower(strings.TrimSpace(parsed.StandardName))
		}
		if key == "" && parsed.DocID != "" {
			key = strings.ToLower(parsed.DocID)
		}
		if key == "" {
			continue
		}

		if existing, ok := resultMap[key]; ok {
			for _, alias := range parsed.Aliases {
				existing.Aliases = appendUniqueFold(existing.Aliases, alias)
			}
			for _, mv := range parsed.MatchedValues {
				existing.MatchedValues = appendUniqueFold(existing.MatchedValues, mv)
			}
			if parsed.Score > existing.Score {
				existing.Score = parsed.Score
				if parsed.DocID != "" {
					existing.DocID = parsed.DocID
				}
				if parsed.DocTitle != "" {
					existing.DocTitle = parsed.DocTitle
				}
				if parsed.Source != "" {
					existing.Source = parsed.Source
				}
				if parsed.Link != "" {
					existing.Link = parsed.Link
				}
			}
			continue
		}

		resultCopy := parsed
		resultMap[key] = &resultCopy
	}

	results := make([]wikiNameNormalization, 0, len(resultMap))
	for _, item := range resultMap {
		item.Aliases = dedupeStringsCaseInsensitive(item.Aliases)
		preferred := choosePreferredStandardName(item.StandardName, item.Aliases, query, normalizedQuery)
		if !strings.EqualFold(preferred, item.StandardName) && strings.TrimSpace(item.StandardName) != "" {
			item.Aliases = appendUniqueFold(item.Aliases, item.StandardName)
			item.StandardName = preferred
		}
		if item.StandardName != "" {
			item.Aliases = removeCaseInsensitive(item.Aliases, item.StandardName)
		}
		item.MatchedValues = dedupeStringsCaseInsensitive(item.MatchedValues)
		results = append(results, *item)
	}

	sort.Slice(results, func(i, j int) bool {
		if len(results[i].MatchedValues) != len(results[j].MatchedValues) {
			return len(results[i].MatchedValues) > len(results[j].MatchedValues)
		}
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return strings.ToLower(results[i].StandardName) < strings.ToLower(results[j].StandardName)
	})

	if len(results) > limit {
		results = results[:limit]
	}

	rewriteQuery := strings.TrimSpace(resp.Data.RewriteQuery)
	usage := convertWikiTokenUsage(resp.Data.TokenUsage)

	return results, rewriteQuery, usage, nil
}

// extractNameNormalizationFromItem 是解析函数，将知识库检索结果转为标准化名称结构。
func extractNameNormalizationFromItem(item *wikiMng.CollectionSearchResponseItem) wikiNameNormalization {
	result := wikiNameNormalization{}
	if item == nil {
		return result
	}

	if len(item.TableChunkFields) > 0 {
		aliases := make([]string, 0, len(item.TableChunkFields))
		for _, field := range item.TableChunkFields {
			fieldName := strings.TrimSpace(field.FieldName)
			fieldValue := strings.TrimSpace(interfaceToString(field.FieldValue))
			if fieldName == "" || fieldValue == "" {
				continue
			}
			upperName := strings.ToUpper(fieldName)
			switch {
			case upperName == "标准化名称" || upperName == "标准名称" || upperName == "STANDARD_NAME" || upperName == "STANDARD":
				if result.StandardName == "" {
					result.StandardName = fieldValue
				}
			case upperName == "名称" || upperName == "NAME":
				if result.StandardName == "" {
					result.StandardName = fieldValue
				}
			case strings.Contains(upperName, "别名") || strings.Contains(upperName, "ALIAS") || strings.Contains(upperName, "SYNONYM"):
				for _, alias := range splitAliasValues(fieldValue) {
					aliases = appendUniqueFold(aliases, alias)
				}
			}
		}
		result.Aliases = dedupeStringsCaseInsensitive(aliases)
	}

	if result.StandardName == "" {
		result.StandardName = strings.TrimSpace(item.ChunkTitle)
	}
	if result.StandardName == "" {
		result.StandardName = strings.TrimSpace(item.DocInfo.DocName)
		if result.StandardName == "" {
			result.StandardName = strings.TrimSpace(item.DocInfo.Title)
		}
	}
	if result.StandardName == "" && len(result.Aliases) > 0 {
		result.StandardName = result.Aliases[0]
		result.Aliases = result.Aliases[1:]
	}
	if result.StandardName != "" {
		result.Aliases = removeCaseInsensitive(result.Aliases, result.StandardName)
	}
	return result
}

// splitAliasValues 是文本处理函数，将别名字符串拆分成切片。
func splitAliasValues(raw string) []string {
	if raw == "" {
		return nil
	}
	clean := strings.ReplaceAll(raw, "\r\n", "\n")
	parts := strings.FieldsFunc(clean, func(r rune) bool {
		switch r {
		case ',', '，', '、', ';', '；', '|', '/', '\n', '\r', '\t':
			return true
		default:
			return false
		}
	})
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		trimmed := strings.TrimSpace(raw)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// dedupeStringsCaseInsensitive 是工具函数，按不区分大小写的方式去重字符串列表。
func dedupeStringsCaseInsensitive(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		v := strings.TrimSpace(value)
		if v == "" {
			continue
		}
		key := strings.ToLower(v)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, v)
	}
	return result
}

// appendUniqueFold 是追加函数，为字符串列表追加不重复的元素。
func appendUniqueFold(values []string, value string) []string {
	v := strings.TrimSpace(value)
	if v == "" {
		return values
	}
	for _, existing := range values {
		if strings.EqualFold(existing, v) {
			return values
		}
	}
	return append(values, v)
}

// removeCaseInsensitive 是清理函数，按不区分大小写的方式移除指定字符串。
func removeCaseInsensitive(values []string, target string) []string {
	if len(values) == 0 {
		return []string{}
	}
	targetLower := strings.ToLower(strings.TrimSpace(target))
	if targetLower == "" {
		return values
	}
	result := make([]string, 0, len(values))
	for _, v := range values {
		if strings.ToLower(strings.TrimSpace(v)) == targetLower {
			continue
		}
		result = append(result, v)
	}
	return result
}

// normalizeNameForMatch 是标准化函数，将文本归一化便于匹配。
func normalizeNameForMatch(text string) string {
	if text == "" {
		return ""
	}
	lower := strings.ToLower(strings.TrimSpace(text))
	if lower == "" {
		return ""
	}
	return nameNormalizeReplacer.Replace(lower)
}

// choosePreferredStandardName 是优选函数，从候选中挑选最佳标准名称。
func choosePreferredStandardName(current string, aliases []string, originalQuery string, normalizedQuery string) string {
	currentTrimmed := strings.TrimSpace(current)
	if currentTrimmed == "" {
		for _, alias := range aliases {
			aliasTrimmed := strings.TrimSpace(alias)
			if aliasTrimmed != "" {
				return aliasTrimmed
			}
		}
		return currentTrimmed
	}

	currentNormalized := normalizeNameForMatch(currentTrimmed)
	if strings.Contains(currentTrimmed, "-") || currentNormalized == "" {
		return currentTrimmed
	}

	for _, alias := range aliases {
		if normalizeNameForMatch(alias) != currentNormalized {
			continue
		}
		if strings.Contains(alias, "-") {
			return alias
		}
	}

	if normalizedQuery != "" && currentNormalized == normalizedQuery && strings.Contains(originalQuery, "-") {
		return originalQuery
	}

	return currentTrimmed
}

// queryWikiArticles 是知识库检索函数，用于查询问答类的文章内容。
func (m *Manager) queryWikiArticles(ctx context.Context, wikiM *wikiMng.WikiMng, keyword string, limit int) ([]wikiArticle, int, int, string, *wikiTokenUsage, error) {
	if wikiM == nil {
		return nil, 0, 0, "", nil, errors.New("知识库未初始化")
	}
	kw := strings.TrimSpace(keyword)
	if kw == "" {
		return nil, 0, 0, "", nil, errors.New("知识库检索关键词不能为空")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if limit <= 0 {
		limit = m.config.DefaultWikiSearchLimit
	}
	if limit > m.config.MaxWikiSearchLimit {
		limit = m.config.MaxWikiSearchLimit
	}

	params := wikiM.GetSearchKnowledgeReqParams(m.config.WikiCollectionID, kw, 0.5, 1)
	if params == nil {
		return nil, 0, 0, "", nil, errors.New("构建知识库检索参数失败")
	}
	params.Limit = int32(limit)

	resp, err := wikiM.SearchKnowledge(ctx, params)
	if err != nil {
		if isTimeoutError(err) {
			return nil, 0, 0, "", nil, fmt.Errorf("知识库检索超时，请稍后重试（%s）", err.Error())
		}
		return nil, 0, 0, "", nil, err
	}
	if resp == nil || resp.Data == nil {
		return []wikiArticle{}, 0, 0, "", nil, nil
	}

	totalCount := int(resp.Data.Count)
	rewriteQuery := strings.TrimSpace(resp.Data.RewriteQuery)
	articles := make([]wikiArticle, 0, limit)
	matched := 0

	for _, item := range resp.Data.ResultList {
		if item == nil {
			continue
		}
		if item.Score < m.config.WikiMinScore {
			continue
		}
		matched++
		if len(articles) >= limit {
			continue
		}
		question, answer, contentStr := extractWikiQA(item)
		docTitle := strings.TrimSpace(item.DocInfo.DocName)
		if docTitle == "" {
			docTitle = strings.TrimSpace(item.DocInfo.Title)
		}
		link := ""
		if len(item.ChunkAttachmentList) > 0 {
			link = strings.TrimSpace(item.ChunkAttachmentList[0].Link)
		}
		articles = append(articles, wikiArticle{
			Question:   question,
			Answer:     answer,
			Content:    contentStr,
			Score:      item.Score,
			DocID:      strings.TrimSpace(item.DocInfo.Docid),
			DocTitle:   docTitle,
			Source:     strings.TrimSpace(item.DocInfo.Source),
			ChunkTitle: strings.TrimSpace(item.ChunkTitle),
			Link:       link,
			RecallRank: int(item.RecallPosition),
			RerankRank: int(item.RerankPosition),
			ChunkID:    item.ChunkId,
		})
	}

	usage := convertWikiTokenUsage(resp.Data.TokenUsage)

	return articles, totalCount, matched, rewriteQuery, usage, nil
}

// convertWikiTokenUsage 是转换函数，将底层 Token 使用数据转为内部结构。
func convertWikiTokenUsage(raw *wikiMng.TotalTokenUsage) *wikiTokenUsage {
	if raw == nil {
		return nil
	}
	usage := &wikiTokenUsage{}
	if raw.EmbeddingUsage != nil {
		usage.EmbeddingTokens = raw.EmbeddingUsage.TotalTokens
	}
	if raw.RerankUsage != nil {
		usage.RerankTokens = *raw.RerankUsage
	}
	if raw.LLMUsage != nil {
		usage.LLMTokens = raw.LLMUsage.TotalTokens
	}
	if raw.RewriteUsage != nil {
		usage.RewriteTokens = raw.RewriteUsage.TotalTokens
	}
	usage.TotalTokens = usage.EmbeddingTokens + usage.RerankTokens + usage.LLMTokens + usage.RewriteTokens
	return usage
}

// wikiArticlesToPayloadItems 是整形函数，将文章结果转换为通用结构。
func wikiArticlesToPayloadItems(articles []wikiArticle) []map[string]interface{} {
	items := make([]map[string]interface{}, 0, len(articles))
	for _, art := range articles {
		item := map[string]interface{}{
			"question": art.Question,
			"answer":   art.Answer,
			"score":    art.Score,
		}
		if art.Content != "" {
			item["content"] = art.Content
		}
		if art.DocID != "" {
			item["doc_id"] = art.DocID
		}
		if art.DocTitle != "" {
			item["doc_title"] = art.DocTitle
		}
		if art.Source != "" {
			item["source"] = art.Source
		}
		if art.ChunkTitle != "" {
			item["chunk_title"] = art.ChunkTitle
		}
		if art.Link != "" {
			item["link"] = art.Link
		}
		if art.RecallRank != 0 {
			item["recall_rank"] = art.RecallRank
		}
		if art.RerankRank != 0 {
			item["rerank_rank"] = art.RerankRank
		}
		if art.ChunkID != 0 {
			item["chunk_id"] = art.ChunkID
		}
		items = append(items, item)
	}
	return items
}

// wikiTokenUsageToMap 是转换函数，将 Token 使用信息转换为 map。
func wikiTokenUsageToMap(usage *wikiTokenUsage) map[string]interface{} {
	if usage == nil {
		return nil
	}
	return map[string]interface{}{
		"embedding_tokens": usage.EmbeddingTokens,
		"rerank_tokens":    usage.RerankTokens,
		"llm_tokens":       usage.LLMTokens,
		"rewrite_tokens":   usage.RewriteTokens,
		"total_tokens":     usage.TotalTokens,
	}
}

// isTimeoutError 是判断函数，检测错误是否由超时引起。
func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return strings.Contains(err.Error(), "Client.Timeout")
}

// extractWikiQA 是解析函数，从知识库条目中提取问答内容。
func extractWikiQA(item *wikiMng.CollectionSearchResponseItem) (question, answer, content string) {
	if item == nil {
		return "", "", ""
	}
	content = strings.TrimSpace(item.Content)
	question = strings.TrimSpace(item.OriginalQuestion)
	answer = ""

	if len(item.TableChunkFields) > 0 {
		for _, field := range item.TableChunkFields {
			name := strings.ToUpper(strings.TrimSpace(field.FieldName))
			value := strings.TrimSpace(interfaceToString(field.FieldValue))
			if value == "" {
				continue
			}
			switch name {
			case "Q", "QUESTION", "问题":
				question = value
			case "A", "ANSWER", "答案":
				answer = value
			default:
				if content == "" {
					content = value
				}
			}
		}
	}

	if question == "" || answer == "" {
		q2, a2 := splitDocStr(item.Content)
		if question == "" {
			question = strings.TrimSpace(q2)
		}
		if answer == "" {
			answer = strings.TrimSpace(a2)
		}
	}
	if question == "" {
		question = strings.TrimSpace(item.ChunkTitle)
	}
	if question == "" {
		question = content
	}
	if answer == "" {
		answer = content
	}
	return
}

// interfaceToString 是辅助函数，将通用类型转换为字符串。
func interfaceToString(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return ""
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 64)
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	default:
		return fmt.Sprintf("%v", val)
	}
}

const questionPrefix = "Q:"
const answerPrefix = "A:"

// splitDocStr 是文本拆分函数，从带前缀的内容中拆解问题与答案。
func splitDocStr(inputStr string) (question, answer string) {
	qIdx := strings.Index(inputStr, questionPrefix)
	if qIdx == -1 {
		fmt.Println("未找到Question前缀")
		return
	}
	qStart := qIdx + len(questionPrefix)
	qEnd := strings.Index(inputStr[qStart:], "\n")
	if qEnd == -1 {
		question = strings.TrimSpace(inputStr[qStart:])
	} else {
		question = strings.TrimSpace(inputStr[qStart : qStart+qEnd])
	}

	aIdx := strings.Index(inputStr, answerPrefix)
	if aIdx == -1 {
		fmt.Println("未找到Answer前缀")
		return
	}
	aStart := aIdx + len(answerPrefix)
	answer = strings.TrimSpace(inputStr[aStart:])

	return
}
