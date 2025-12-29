package aiMng

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/wiidz/goutil/mngs/eyunMng/msgSend"
)

type directSkuContext struct {
	Candidate    string                         `json:"candidate"`
	StandardName string                         `json:"standard_name"`
	SearchUsed   string                         `json:"search_used"`
	SpuID        uint64                         `json:"spu_id,omitempty"`
	Total        int64                          `json:"total"`
	AllRecords   []ProductPriceRecord `json:"all_records,omitempty"`
	Records      []ProductPriceRecord `json:"records"`
	LastQuery    string                         `json:"last_query,omitempty"`
	UpdatedAt    int64                          `json:"updated_at"`
}

type productCandidate struct {
	StandardName string `json:"standard_name,omitempty"`
	Candidate    string `json:"candidate,omitempty"`
	SearchUsed   string `json:"search_used,omitempty"`
	Source       string `json:"source,omitempty"`
}

// directSkuContextRedisKey 构造 direct SKU 上下文对应的 Redis 键。
func directSkuContextRedisKey(sessionID string) string {
	return fmt.Sprintf("chat:sku:%s", sessionID)
}

// loadDirectSkuContext 从 Redis 读取 direct SKU 会话上下文。
func (m *Manager) loadDirectSkuContext(ctx context.Context, sessionID string) (*directSkuContext, error) {
	if m == nil || m.redis == nil || sessionID == "" {
		return nil, nil
	}
	raw, err := m.redis.GetString(ctx, directSkuContextRedisKey(sessionID))
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var payload directSkuContext
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

// saveDirectSkuContext 将 direct SKU 会话上下文写入 Redis。
func (m *Manager) saveDirectSkuContext(ctx context.Context, sessionID string, payload *directSkuContext) error {
	if m == nil || m.redis == nil || sessionID == "" || payload == nil {
		return nil
	}
	payload.UpdatedAt = time.Now().Unix()
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return m.redis.Set(ctx, directSkuContextRedisKey(sessionID), data, m.config.DirectSkuContextTTL)
}

// clearDirectSkuContext 清空 direct SKU 会话上下文缓存。
func (m *Manager) clearDirectSkuContext(ctx context.Context, sessionID string) error {
	if m == nil || m.redis == nil || sessionID == "" {
		return nil
	}
	return m.redis.Set(ctx, directSkuContextRedisKey(sessionID), "", 0)
}

// sessionProductContextRedisKey 构造产品上下文对应的 Redis 键。
func sessionProductContextRedisKey(sessionID string) string {
	return fmt.Sprintf("chat:product:%s", sessionID)
}

type sessionProductContext struct {
	Primary          *productCandidate   `json:"primary,omitempty"`
	Candidates       []productCandidate  `json:"candidates,omitempty"`
	SpuID            uint64              `json:"spu_id,omitempty"`
	SkuCount         int                 `json:"sku_count,omitempty"`
	AttributeNames   []string            `json:"attribute_names,omitempty"`
	AttributeOptions map[string][]string `json:"attribute_options,omitempty"`
	LastQuery        string              `json:"last_query,omitempty"`
	UpdatedAt        int64               `json:"updated_at"`
}

var productAttributeCache sync.Map

type attrNameRow struct {
	AttrName string `gorm:"column:attr_name"`
}

// cloneStrings 创建字符串切片的副本，避免共享底层数组。
func cloneStrings(src []string) []string {
	if len(src) == 0 {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

func stringSliceContains(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

// stringSliceEqual 判断两个字符串切片内容是否一致。
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// resolveSPUID 从产品记录中推导出 SPU ID。
func resolveSPUID(records []ProductPriceRecord) uint64 {
	for _, record := range records {
		if record.SpuID != 0 {
			return record.SpuID
		}
	}
	return 0
}

// productAttributeNames 查询指定 SPU 对应的属性名称集合。
func (m *Manager) productAttributeNames(ctx context.Context, spuID uint64) ([]string, error) {
	if spuID == 0 || m == nil || m.db == nil {
		return nil, nil
	}
	if cached, ok := productAttributeCache.Load(spuID); ok {
		if names, ok := cached.([]string); ok {
			return cloneStrings(names), nil
		}
	}

	if ctx == nil {
		ctx = context.Background()
	}

	rows := make([]attrNameRow, 0)
	query := m.db.WithContext(ctx).
		Table("a_spu_attr AS spa").
		Select("k.attr_name").
		Joins("LEFT JOIN a_attr_key AS k ON spa.attribute_key_id = k.id").
		Where("spa.spu_id = ?", spuID).
		Order("spa.sequence ASC")

	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	names := make([]string, 0, len(rows))
	seen := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(row.AttrName)
		if name == "" {
			continue
		}
		lower := strings.ToLower(name)
		if _, ok := seen[lower]; ok {
			continue
		}
		seen[lower] = struct{}{}
		names = append(names, name)
	}

	productAttributeCache.Store(spuID, cloneStrings(names))
	return names, nil
}

func (m *Manager) buildAttributeOptions(ctx context.Context, spuID uint64) (map[string][]string, error) {
	if spuID == 0 || m == nil || m.db == nil {
		return nil, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	type cacheRow struct {
		SKUID     int64  `gorm:"column:id"`
		AttrNames []byte `gorm:"column:attr_value_cache_names"`
	}
	rows := make([]cacheRow, 0)
	if err := m.db.WithContext(ctx).
		Table("a_sku").
		Select("id, attr_value_cache_names").
		Where("spu_id = ?", spuID).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	missing := make([]int64, 0)

	for _, row := range rows {
		if len(row.AttrNames) == 0 {
			missing = append(missing, row.SKUID)
			continue
		}
		var cache map[string]string
		if err := json.Unmarshal(row.AttrNames, &cache); err != nil || len(cache) == 0 {
			missing = append(missing, row.SKUID)
			continue
		}
		for name, value := range cache {
			name = strings.TrimSpace(name)
			value = strings.TrimSpace(value)
			if name == "" || value == "" {
				continue
			}
			list := result[name]
			if !stringSliceContains(list, value) {
				result[name] = append(list, value)
			}
		}
	}

	if len(missing) > 0 {
		type attrRow struct {
			AttrName  string `gorm:"column:attr_name"`
			AttrValue string `gorm:"column:attr_value"`
		}
		attrRows := make([]attrRow, 0)
		if err := m.db.WithContext(ctx).
			Table("a_sku_attr_map AS m").
			Select("k.attr_name, m.attr_value").
			Joins("LEFT JOIN a_attr_key AS k ON m.attr_key_id = k.id").
			Where("m.sku_id IN ?", missing).
			Find(&attrRows).Error; err != nil {
			return nil, err
		}
		for _, row := range attrRows {
			name := strings.TrimSpace(row.AttrName)
			value := strings.TrimSpace(row.AttrValue)
			if name == "" || value == "" {
				continue
			}
			list := result[name]
			if !stringSliceContains(list, value) {
				result[name] = append(list, value)
			}
		}
	}

	for key, list := range result {
		if len(list) > 1 {
			sort.Strings(list)
			result[key] = list
		}
	}
	return result, nil
}

// ensureProductAttributes 确保会话产品上下文同步最新的属性名称。
func (m *Manager) ensureProductAttributes(ctx context.Context, productCtx *sessionProductContext, stored *directSkuContext) []string {
	var names []string
	if productCtx != nil && len(productCtx.AttributeNames) > 0 {
		names = cloneStrings(productCtx.AttributeNames)
	}

	if len(names) == 0 && stored != nil {
		spuID := stored.SpuID
		if spuID == 0 {
			spuID = resolveSPUID(stored.AllRecords)
			if spuID == 0 {
				spuID = resolveSPUID(stored.Records)
			}
		}
		if spuID != 0 {
			if fetched, err := m.productAttributeNames(ctx, spuID); err == nil && len(fetched) > 0 {
				names = fetched
			} else if err != nil {
				fmt.Printf("load product attribute names error: %v\n", err)
			}
		}
	}

	if len(names) == 0 && stored != nil {
		nameSet := make(map[string]struct{})
		for _, record := range stored.AllRecords {
			for key := range record.Attributes {
				if strings.TrimSpace(key) == "" {
					continue
				}
				nameSet[key] = struct{}{}
			}
		}
		if len(nameSet) == 0 {
			for _, record := range stored.Records {
				for key := range record.Attributes {
					if strings.TrimSpace(key) == "" {
						continue
					}
					nameSet[key] = struct{}{}
				}
			}
		}
		if len(nameSet) > 0 {
			names = make([]string, 0, len(nameSet))
			for key := range nameSet {
				names = append(names, key)
			}
			sort.Strings(names)
		}
	}

	if productCtx != nil && len(names) > 0 {
		if !stringSliceEqual(productCtx.AttributeNames, names) {
			productCtx.AttributeNames = cloneStrings(names)
		}
	}

	if productCtx != nil && productCtx.SpuID != 0 {
		if len(productCtx.AttributeOptions) == 0 {
			if options, err := m.buildAttributeOptions(ctx, productCtx.SpuID); err == nil && len(options) > 0 {
				productCtx.AttributeOptions = options
			} else if err != nil {
				fmt.Printf("load product attribute options error: %v\n", err)
			}
		}
	}
	return names
}

// generateAttributeSynonyms 为属性名称生成匹配时可使用的同义词。
func generateAttributeSynonyms(attr string) []string {
	attr = strings.TrimSpace(attr)
	if attr == "" {
		return nil
	}
	targets := make([]string, 0, 8)
	add := func(val string) {
		val = strings.TrimSpace(val)
		if val == "" {
			return
		}
		targets = append(targets, val)
	}
	add(attr)

	splitters := func(r rune) bool {
		switch r {
		case ' ', '\t', '\n', '\r', '/', '\\', '-', '_', ',', '，', '、', '|', '（', '）', '(', ')', '【', '】', '[', ']', '：', ':':
			return true
		default:
			return false
		}
	}
	parts := strings.FieldsFunc(attr, splitters)
	for _, part := range parts {
		add(part)
	}

	unique := make([]string, 0, len(targets))
	seen := make(map[string]struct{}, len(targets))
	for _, val := range targets {
		key := strings.ToLower(val)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		unique = append(unique, val)
	}
	return unique
}

// detectAttributeKeyword 在用户输入中识别关注的属性名称。
func detectAttributeKeyword(userText string, attrNames []string, records []ProductPriceRecord) string {
	if strings.TrimSpace(userText) == "" {
		return ""
	}
	lowerText := strings.ToLower(userText)

	attrSet := make(map[string]struct{})
	for _, name := range attrNames {
		if trimmed := strings.TrimSpace(name); trimmed != "" {
			attrSet[trimmed] = struct{}{}
		}
	}
	for _, record := range records {
		for key := range record.Attributes {
			if trimmed := strings.TrimSpace(key); trimmed != "" {
				attrSet[trimmed] = struct{}{}
			}
		}
	}
	if len(attrSet) == 0 {
		return ""
	}
	keys := make([]string, 0, len(attrSet))
	for key := range attrSet {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, attr := range keys {
		synonyms := generateAttributeSynonyms(attr)
		for _, synonym := range synonyms {
			if synonym == "" {
				continue
			}
			synLower := strings.ToLower(synonym)
			if strings.Contains(lowerText, synLower) || strings.Contains(userText, synonym) {
				return attr
			}
		}
	}
	return ""
}

// collectAttributeValues 汇总产品记录中某属性的全部可能取值。
func collectAttributeValues(records []ProductPriceRecord, key string) []string {
	if key == "" {
		return nil
	}
	seen := make(map[string]struct{})
	values := make([]string, 0, len(records))
	for _, record := range records {
		if record.Attributes == nil {
			continue
		}
		val := strings.TrimSpace(record.Attributes[key])
		if val == "" {
			continue
		}
		if _, ok := seen[val]; ok {
			continue
		}
		seen[val] = struct{}{}
		values = append(values, val)
	}
	if len(values) == 0 {
		return nil
	}
	sort.Strings(values)
	return values
}

// loadSessionProductContext 从 Redis 加载会话级的产品上下文。
func (m *Manager) loadSessionProductContext(ctx context.Context, sessionID string) (*sessionProductContext, error) {
	if m == nil || m.redis == nil || sessionID == "" {
		return nil, nil
	}
	raw, err := m.redis.GetString(ctx, sessionProductContextRedisKey(sessionID))
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var payload sessionProductContext
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, err
	}
	recordSessionProductSnapshotFromContext(sessionID, &payload)
	return &payload, nil
}

// saveSessionProductContext 将会话级产品上下文写入 Redis。
func (m *Manager) saveSessionProductContext(ctx context.Context, sessionID string, payload *sessionProductContext) error {
	if m == nil || m.redis == nil || sessionID == "" || payload == nil {
		return nil
	}
	payload.UpdatedAt = time.Now().Unix()
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if err := m.redis.Set(ctx, sessionProductContextRedisKey(sessionID), data, m.config.SessionTTL); err != nil {
		return err
	}
	recordSessionProductSnapshotFromContext(sessionID, payload)
	return nil
}

// clearSessionProductContext 删除会话级产品上下文缓存。
func (m *Manager) clearSessionProductContext(ctx context.Context, sessionID string) error {
	if m == nil || m.redis == nil || sessionID == "" {
		return nil
	}
	clearSessionProductSnapshot(sessionID)
	return m.redis.Set(ctx, sessionProductContextRedisKey(sessionID), "", 0)
}

// productFocusLock 表示待锁定的单个产品候选及其检索结果。
type productFocusLock struct {
	Query  string
	Result *skuLookupResult
}

// applyProductFocusFromLookup 根据工具识别结果写入产品上下文与直接检索缓存。
func (m *Manager) applyProductFocusFromLookup(ctx context.Context, sessionID string, locks []productFocusLock) error {
	if m == nil || sessionID == "" || len(locks) == 0 {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	productCtx, err := m.loadSessionProductContext(ctx, sessionID)
	if err != nil {
		return err
	}
	if productCtx == nil {
		productCtx = &sessionProductContext{}
	}

	type focusEntry struct {
		lock      productFocusLock
		candidate productCandidate
	}
	entries := make([]focusEntry, 0, len(locks))
	for _, lock := range locks {
		if lock.Result == nil || len(lock.Result.Records) == 0 {
			continue
		}
		candidate := productCandidate{
			StandardName: strings.TrimSpace(lock.Result.StandardName),
			Candidate:    strings.TrimSpace(lock.Query),
			SearchUsed:   lock.Result.SearchUsed,
			Source:       "identify_product_focus",
		}
		entries = append(entries, focusEntry{
			lock:      lock,
			candidate: candidate,
		})
	}
	if len(entries) == 0 {
		return nil
	}

	newCandidates := make([]productCandidate, 0, len(entries))
	for _, entry := range entries {
		newCandidates = append(newCandidates, entry.candidate)
	}
	applyAssistantCandidates(productCtx, newCandidates)

	var primaryEntry focusEntry
	if productCtx.Primary != nil {
		targetKey := candidateKey(*productCtx.Primary)
		for _, entry := range entries {
			if targetKey != "" && candidateKey(entry.candidate) == targetKey {
				primaryEntry = entry
				break
			}
			if strings.TrimSpace(entry.candidate.StandardName) != "" &&
				strings.EqualFold(strings.TrimSpace(entry.candidate.StandardName), strings.TrimSpace(productCtx.Primary.StandardName)) {
				primaryEntry = entry
				break
			}
		}
	}
	if primaryEntry.lock.Result == nil {
		primaryEntry = entries[0]
		if productCtx.Primary == nil {
			c := primaryEntry.candidate
			productCtx.Primary = &c
		} else {
			mergeCandidate(productCtx.Primary, primaryEntry.candidate)
		}
	}

	productCtx.LastQuery = strings.TrimSpace(primaryEntry.lock.Query)

	if primaryEntry.lock.Result != nil {
		if productCtx.SpuID == 0 {
			productCtx.SpuID = resolveSPUID(primaryEntry.lock.Result.Records)
		}
		skuCount := int(primaryEntry.lock.Result.Total)
		if skuCount == 0 {
			skuCount = len(primaryEntry.lock.Result.Records)
		}
		productCtx.SkuCount = skuCount

		direct := &directSkuContext{
			Candidate:    strings.TrimSpace(primaryEntry.candidate.Candidate),
			StandardName: strings.TrimSpace(primaryEntry.candidate.StandardName),
			SearchUsed:   primaryEntry.lock.Result.SearchUsed,
			SpuID:        resolveSPUID(primaryEntry.lock.Result.Records),
			Total:        primaryEntry.lock.Result.Total,
			AllRecords:   primaryEntry.lock.Result.Records,
			Records:      primaryEntry.lock.Result.Records,
			LastQuery:    strings.TrimSpace(primaryEntry.lock.Query),
		}
		productCtx.AttributeNames = m.ensureProductAttributes(ctx, productCtx, direct)
		if err := m.saveDirectSkuContext(ctx, sessionID, direct); err != nil {
			return err
		}
		if m.redis == nil {
			recordSessionProductSnapshotFromContext(sessionID, productCtx)
		}
	}

	if err := m.saveSessionProductContext(ctx, sessionID, productCtx); err != nil {
		return err
	}
	if m.redis == nil {
		recordSessionProductSnapshotFromContext(sessionID, productCtx)
	}

	return nil
}

func (m *Manager) maybeHandleDirectSkuQuery(
	ctx context.Context,
	sessionID string,
	conversation []*model.ChatCompletionMessage,
	userText string,
) ([]msgSend.MessagePayload, bool, error) {
	trimmedText := strings.TrimSpace(userText)
	if trimmedText == "" {
		return nil, false, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var storedContext *directSkuContext
	if cached, err := m.loadDirectSkuContext(ctx, sessionID); err != nil {
		fmt.Printf("load direct sku context error: %v\n", err)
	} else {
		storedContext = cached
		if cached != nil {
			if storedContext.SpuID == 0 {
				if id := resolveSPUID(storedContext.AllRecords); id != 0 {
					storedContext.SpuID = id
				} else if id := resolveSPUID(storedContext.Records); id != 0 {
					storedContext.SpuID = id
				}
			}
			recordSessionProductSnapshot(sessionID, sessionProductSnapshot{
				Primary: describeProductCandidate(productCandidate{
					StandardName: cached.StandardName,
					Candidate:    cached.Candidate,
					SearchUsed:   cached.SearchUsed,
					Source:       "direct_sku_cache",
				}),
			})
		}
	}

	productContext, err := m.loadSessionProductContext(ctx, sessionID)
	if err != nil {
		fmt.Printf("load session product context error: %v\n", err)
	}

	if storedContext != nil {
		if messages, handled, err := m.handleDirectSkuFollowUp(ctx, sessionID, conversation, storedContext, trimmedText); handled {
			return messages, true, err
		}
	}

	if storedContext == nil && productContext != nil && productContext.Primary != nil {
		name := strings.TrimSpace(productContext.Primary.StandardName)
		if name == "" {
			name = strings.TrimSpace(productContext.Primary.Candidate)
		}
		if name != "" {
			result, err := m.lookupSkuRecords(ctx, name, 0)
			if err != nil {
				fmt.Printf("lookup sku records via product context error: %v\n", err)
			} else if result != nil && len(result.Records) > 0 {
				candidate := strings.TrimSpace(productContext.Primary.Candidate)
				if candidate == "" {
					candidate = name
				}
				storedContext = &directSkuContext{
					Candidate:    candidate,
					StandardName: result.StandardName,
					SearchUsed:   result.SearchUsed,
					Total:        result.Total,
					AllRecords:   result.Records,
					Records:      result.Records,
					LastQuery:    productContext.LastQuery,
				}
				if messages, handled, err := m.handleDirectSkuFollowUp(ctx, sessionID, conversation, storedContext, trimmedText); handled {
					return messages, true, err
				}
			}
		}
	}

	candidates := extractSkuCandidates(trimmedText)
	if len(candidates) == 0 {
		return m.handlePrimaryFollowUp(ctx, sessionID, conversation, productContext, trimmedText)
	}
	if !hasSkuIntentKeyword(trimmedText) && !(len(candidates) == 1 && len([]rune(trimmedText)) <= 16) {
		if messages, handled, err := m.handlePrimaryFollowUp(ctx, sessionID, conversation, productContext, trimmedText); handled || err != nil {
			return messages, handled, err
		}
		return nil, false, nil
	}

	for _, candidate := range candidates {
		result, err := m.lookupSkuRecords(ctx, candidate, 0)
		if err != nil {
			return nil, true, err
		}
		if result == nil || len(result.Records) == 0 {
			continue
		}

		reply := composeSkuReply(result, candidate)
		if reply == "" {
			continue
		}

		conversation = append(conversation, m.newChatMessage(sessionID, model.ChatMessageRoleAssistant, reply, nil))

		contextPayload := &directSkuContext{
			Candidate:    candidate,
			StandardName: result.StandardName,
			SearchUsed:   result.SearchUsed,
			SpuID:        resolveSPUID(result.Records),
			Total:        result.Total,
			AllRecords:   result.Records,
			Records:      result.Records,
			LastQuery:    trimmedText,
		}
		if err := m.saveDirectSkuContext(ctx, sessionID, contextPayload); err != nil {
			fmt.Printf("save direct sku context error: %v\n", err)
		}
		productPayload := &sessionProductContext{
			Primary: &productCandidate{
				StandardName: result.StandardName,
				Candidate:    candidate,
				SearchUsed:   result.SearchUsed,
				Source:       "direct_sku",
			},
			SpuID:     contextPayload.SpuID,
			SkuCount:  int(result.Total),
			LastQuery: trimmedText,
		}
		productPayload.AttributeNames = m.ensureProductAttributes(ctx, productPayload, contextPayload)
		if err := m.saveSessionProductContext(ctx, sessionID, productPayload); err != nil {
			fmt.Printf("save session product context error: %v\n", err)
		}
		productContext = productPayload
		recordSessionProductSnapshotFromContext(sessionID, productPayload)
		conversation = m.appendProductContextMessage(ctx, sessionID, conversation)
		m.debugPrintSession(sessionID, fmt.Sprintf("direct sku reply (%s)", candidate), conversation)
		if saveErr := m.saveSessionConversation(ctx, sessionID, m.trimConversationHistory(conversation)); saveErr != nil {
			fmt.Printf("save session conversation error: %v\n", saveErr)
		}

		messages := make([]msgSend.MessagePayload, 0, 1+len(result.Images))
		messages = append(messages, &msgSend.TextMessage{
			WcID:    sessionID,
			Content: reply,
		})
		for _, url := range result.Images {
			messages = append(messages, &msgSend.ImageMessage{
				WcID: sessionID,
				URL:  url,
			})
		}
		return messages, true, nil
	}

	return nil, false, nil
}

// handlePrimaryFollowUp 在已有主产品时补充检索并生成回答。
func (m *Manager) handlePrimaryFollowUp(
	ctx context.Context,
	sessionID string,
	conversation []*model.ChatCompletionMessage,
	productContext *sessionProductContext,
	userText string,
) ([]msgSend.MessagePayload, bool, error) {
	if m == nil || sessionID == "" || productContext == nil || productContext.Primary == nil {
		return nil, false, nil
	}

	name := strings.TrimSpace(productContext.Primary.StandardName)
	if name == "" {
		name = strings.TrimSpace(productContext.Primary.Candidate)
	}
	if name == "" {
		return nil, false, nil
	}

	result, err := m.lookupSkuRecords(ctx, name, 0)
	if err != nil {
		fmt.Printf("lookup sku records via primary product error: %v\n", err)
		return nil, false, nil
	}
	if result == nil || len(result.Records) == 0 {
		return nil, false, nil
	}

	candidate := strings.TrimSpace(productContext.Primary.Candidate)
	if candidate == "" {
		candidate = name
	}

	stored := &directSkuContext{
		Candidate:    candidate,
		StandardName: name,
		SearchUsed:   result.SearchUsed,
		SpuID:        resolveSPUID(result.Records),
		Total:        result.Total,
		AllRecords:   result.Records,
		Records:      result.Records,
		LastQuery:    productContext.LastQuery,
	}

	skuCount := int(result.Total)
	if skuCount == 0 {
		skuCount = len(result.Records)
	}
	if productContext.SkuCount != skuCount {
		productContext.SkuCount = skuCount
	}
	if stored.SpuID != 0 && productContext.SpuID != stored.SpuID {
		productContext.SpuID = stored.SpuID
	}

	attrNames := m.ensureProductAttributes(ctx, productContext, stored)

	if messages, handled, err := m.handleDirectSkuFollowUp(ctx, sessionID, conversation, stored, userText); handled || err != nil {
		return messages, handled, err
	}

	return m.handleAttributeInquiry(ctx, sessionID, conversation, stored, productContext, attrNames, userText)
}

// handleDirectSkuFollowUp 处理针对 direct SKU 列表的深入查询或筛选。
func (m *Manager) handleDirectSkuFollowUp(
	ctx context.Context,
	sessionID string,
	conversation []*model.ChatCompletionMessage,
	stored *directSkuContext,
	userText string,
) ([]msgSend.MessagePayload, bool, error) {
	if stored == nil {
		return nil, false, nil
	}
	sourceRecords := stored.Records
	if len(stored.AllRecords) > 0 {
		sourceRecords = stored.AllRecords
	}
	if len(sourceRecords) == 0 {
		return nil, false, nil
	}

	indices := extractRecordSelectionIndexes(userText, len(sourceRecords))
	var matches []ProductPriceRecord
	if len(indices) > 0 {
		matches = make([]ProductPriceRecord, 0, len(indices))
		seen := make(map[int]struct{}, len(indices))
		for _, idx := range indices {
			if idx < 0 || idx >= len(sourceRecords) {
				continue
			}
			if _, ok := seen[idx]; ok {
				continue
			}
			seen[idx] = struct{}{}
			matches = append(matches, sourceRecords[idx])
		}
		if len(matches) == 0 {
			return nil, false, nil
		}
	} else {
		tokens := extractMeaningfulTokens(userText)
		if len(tokens) == 0 {
			return nil, false, nil
		}
		matches = filterRecordsByTokens(sourceRecords, tokens)
		if len(matches) == 0 {
			return nil, false, nil
		}
		if len(matches) == len(stored.Records) && strings.EqualFold(strings.TrimSpace(stored.LastQuery), strings.TrimSpace(userText)) {
			return nil, false, nil
		}
	}

	reply := composeSkuFollowUpReply(stored, matches, userText)
	if reply == "" {
		return nil, false, nil
	}

	conversation = append(conversation, m.newChatMessage(sessionID, model.ChatMessageRoleAssistant, reply, nil))

	if m != nil {
		updated := &directSkuContext{
			Candidate:    stored.Candidate,
			StandardName: stored.StandardName,
			SearchUsed:   stored.SearchUsed,
			SpuID:        stored.SpuID,
			Total:        stored.Total,
			AllRecords:   stored.AllRecords,
			Records:      matches,
			LastQuery:    userText,
		}
		if len(updated.AllRecords) == 0 {
			updated.AllRecords = sourceRecords
		}
		if err := m.saveDirectSkuContext(ctx, sessionID, updated); err != nil {
			fmt.Printf("save direct sku context error: %v\n", err)
		}
		productPayload := &sessionProductContext{
			Primary: &productCandidate{
				StandardName: updated.StandardName,
				Candidate:    updated.Candidate,
				SearchUsed:   updated.SearchUsed,
				Source:       "direct_sku_follow_up",
			},
			SpuID:     updated.SpuID,
			SkuCount:  int(updated.Total),
			LastQuery: userText,
		}
		productPayload.AttributeNames = m.ensureProductAttributes(ctx, productPayload, updated)
		if err := m.saveSessionProductContext(ctx, sessionID, productPayload); err != nil {
			fmt.Printf("save session product context error: %v\n", err)
		}
		recordSessionProductSnapshotFromContext(sessionID, productPayload)
	}
	conversation = m.appendProductContextMessage(ctx, sessionID, conversation)
	m.debugPrintSession(sessionID, "direct sku follow-up reply", conversation)
	if saveErr := m.saveSessionConversation(ctx, sessionID, m.trimConversationHistory(conversation)); saveErr != nil {
		fmt.Printf("save session conversation error: %v\n", saveErr)
	}

	messages := make([]msgSend.MessagePayload, 0, 1+len(matches))
	messages = append(messages, &msgSend.TextMessage{
		WcID:    sessionID,
		Content: reply,
	})
	if thumbs := collectThumbnailsFromRecords(matches); len(thumbs) > 0 {
		for _, url := range thumbs {
			messages = append(messages, &msgSend.ImageMessage{
				WcID: sessionID,
				URL:  url,
			})
		}
	}

	return messages, true, nil
}

// handleAttributeInquiry 汇总用户关注属性的可选值并生成回复。
func (m *Manager) handleAttributeInquiry(
	ctx context.Context,
	sessionID string,
	conversation []*model.ChatCompletionMessage,
	stored *directSkuContext,
	productCtx *sessionProductContext,
	attrNames []string,
	userText string,
) ([]msgSend.MessagePayload, bool, error) {
	if stored == nil || len(stored.AllRecords) == 0 {
		return nil, false, nil
	}
	if len(attrNames) == 0 {
		attrNames = m.ensureProductAttributes(ctx, productCtx, stored)
	}
	attrKey := detectAttributeKeyword(userText, attrNames, stored.AllRecords)
	if attrKey == "" {
		return nil, false, nil
	}

	var values []string
	if productCtx != nil {
		if opts, ok := productCtx.AttributeOptions[attrKey]; ok && len(opts) > 0 {
			values = cloneStrings(opts)
		}
	}
	if len(values) == 0 {
		values = collectAttributeValues(stored.AllRecords, attrKey)
	}
	if len(values) == 0 {
		return nil, false, nil
	}

	productName := strings.TrimSpace(stored.StandardName)
	if productName == "" {
		productName = strings.TrimSpace(stored.SearchUsed)
	}
	if productName == "" {
		productName = strings.TrimSpace(stored.Candidate)
	}
	if productName == "" {
		productName = "该产品"
	}

	lines := make([]string, 0, len(values)+3)
	lines = append(lines, fmt.Sprintf("%s 当前可选的%s规格如下：", productName, attrKey))
	for _, val := range values {
		lines = append(lines, fmt.Sprintf("- %s", val))
	}
	lines = append(lines, "如需其中任意孔径对应的额定电流、精度或其它参数，请告诉我，我再帮您确认。")

	reply := strings.Join(lines, "\n")
	conversation = append(conversation, m.newChatMessage(sessionID, model.ChatMessageRoleAssistant, reply, nil))

	stored.LastQuery = userText
	if stored.SpuID == 0 {
		stored.SpuID = resolveSPUID(stored.AllRecords)
	}
	stored.Records = stored.AllRecords
	if err := m.saveDirectSkuContext(ctx, sessionID, stored); err != nil {
		fmt.Printf("save direct sku context error: %v\n", err)
	}

	if productCtx == nil {
		productCtx = &sessionProductContext{}
	}
	primary := productCandidate{
		StandardName: stored.StandardName,
		Candidate:    stored.Candidate,
		SearchUsed:   stored.SearchUsed,
		Source:       "attribute_summary",
	}
	if productCtx.Primary == nil {
		c := primary
		productCtx.Primary = &c
	} else {
		mergeCandidate(productCtx.Primary, primary)
	}
	if stored.SpuID != 0 && productCtx.SpuID != stored.SpuID {
		productCtx.SpuID = stored.SpuID
	}
	skuCount := len(stored.AllRecords)
	if skuCount == 0 {
		skuCount = len(stored.Records)
	}
	if skuCount == 0 && stored.Total > 0 {
		skuCount = int(stored.Total)
	}
	if skuCount > 0 {
		productCtx.SkuCount = skuCount
	}
	productCtx.AttributeNames = m.ensureProductAttributes(ctx, productCtx, stored)
	productCtx.LastQuery = userText
	if err := m.saveSessionProductContext(ctx, sessionID, productCtx); err != nil {
		fmt.Printf("save session product context error: %v\n", err)
	}
	conversation = m.appendProductContextMessage(ctx, sessionID, conversation)
	m.debugPrintSession(sessionID, fmt.Sprintf("attribute summary (%s)", attrKey), conversation)
	if err := m.saveSessionConversation(ctx, sessionID, m.trimConversationHistory(conversation)); err != nil {
		fmt.Printf("save session conversation error: %v\n", err)
	}

	messages := []msgSend.MessagePayload{
		&msgSend.TextMessage{
			WcID:    sessionID,
			Content: reply,
		},
	}
	return messages, true, nil
}

// extractRecordSelectionIndexes 从文本解析用户提到的序号列表。
func extractRecordSelectionIndexes(text string, max int) []int {
	if max <= 0 {
		return nil
	}
	matches := selectionIndexPattern.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return nil
	}
	result := make([]int, 0, len(matches))
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		num, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		idx := num - 1
		if idx >= 0 && idx < max {
			result = append(result, idx)
		}
	}
	return result
}

// extractMeaningfulTokens 拆分并收集可用于匹配的关键 token。
func extractMeaningfulTokens(text string) []string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return nil
	}
	tokens := make([]string, 0, 4)
	seen := make(map[string]struct{})

	addToken := func(token string) {
		t := strings.TrimSpace(token)
		if t == "" {
			return
		}
		lower := strings.ToLower(t)
		if _, skip := skuFollowUpStopWords[lower]; skip {
			return
		}
		if !shouldKeepToken(t) {
			return
		}
		if _, ok := seen[lower]; ok {
			return
		}
		seen[lower] = struct{}{}
		tokens = append(tokens, t)
	}

	addToken(trimmed)
	splitTokens := strings.FieldsFunc(trimmed, func(r rune) bool {
		switch r {
		case ' ', '\t', '\n', '\r', ',', '，', '。', ';', '；', '、', '|':
			return true
		default:
			return false
		}
	})
	for _, tok := range splitTokens {
		addToken(tok)
		for _, sub := range extractAlphaNumSequences(tok) {
			addToken(sub)
		}
	}

	return tokens
}

// shouldKeepToken 判断 token 是否值得保留参与匹配。
func shouldKeepToken(token string) bool {
	if token == "" {
		return false
	}
	runes := []rune(token)
	if len(runes) > 1 {
		return true
	}
	r := runes[0]
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return true
	}
	if unicode.Is(unicode.Han, r) {
		return true
	}
	return false
}

// extractAlphaNumSequences 提取 token 内按顺序出现的字母数字片段。
func extractAlphaNumSequences(token string) []string {
	if token == "" {
		return nil
	}
	result := make([]string, 0, 2)
	var builder strings.Builder
	flushBuilder := func() {
		if builder.Len() == 0 {
			return
		}
		result = append(result, builder.String())
		builder.Reset()
	}
	for _, r := range token {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
			continue
		}
		flushBuilder()
	}
	flushBuilder()
	return result
}

// filterRecordsByTokens 根据 token 匹配得分筛选最相关的产品记录。
func filterRecordsByTokens(records []ProductPriceRecord, tokens []string) []ProductPriceRecord {
	if len(records) == 0 || len(tokens) == 0 {
		return nil
	}
	bestScore := 0
	scores := make([]int, len(records))
	for i, record := range records {
		score := computeRecordTokenScore(record, tokens)
		if score == 0 {
			continue
		}
		scores[i] = score
		if score > bestScore {
			bestScore = score
		}
	}
	if bestScore == 0 {
		return nil
	}
	result := make([]ProductPriceRecord, 0, len(records))
	for i, score := range scores {
		if score == bestScore && score > 0 {
			result = append(result, records[i])
		}
	}
	return result
}

// buildRecordSearchCorpus 构造产品记录的检索语料信息。
func buildRecordSearchCorpus(record ProductPriceRecord) recordSearchCorpus {
	parts := make([]string, 0, len(record.Attributes)*2+6)
	if display := strings.TrimSpace(record.DisplayName); display != "" {
		parts = append(parts, display)
	}
	if code := strings.TrimSpace(record.SKUCode); code != "" {
		parts = append(parts, code)
	}
	if len(record.Attributes) > 0 {
		keys := make([]string, 0, len(record.Attributes))
		for key := range record.Attributes {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			val := strings.TrimSpace(record.Attributes[key])
			if val == "" {
				continue
			}
			parts = append(parts, key, val)
		}
	}
	if record.PricePromot != nil {
		parts = append(parts, fmt.Sprintf("%.2f", *record.PricePromot))
	}
	if record.PriceIn != nil {
		parts = append(parts, fmt.Sprintf("%.2f", *record.PriceIn))
	}
	raw := strings.TrimSpace(strings.Join(parts, " "))
	lower := strings.ToLower(raw)
	upper := strings.ToUpper(raw)
	normalized := ""
	if lower != "" {
		normalized = normalizeNameForMatch(lower)
	}
	return recordSearchCorpus{
		raw:        raw,
		lower:      lower,
		upper:      upper,
		normalized: normalized,
	}
}

// computeRecordTokenScore 计算产品记录命中 token 的总得分。
func computeRecordTokenScore(record ProductPriceRecord, tokens []string) int {
	if len(tokens) == 0 {
		return 0
	}
	corpus := buildRecordSearchCorpus(record)
	if corpus.raw == "" {
		return 0
	}
	score := 0
	for _, token := range tokens {
		if token == "" {
			continue
		}
		if matchTokenInCorpus(token, corpus) {
			score++
		}
	}
	return score
}

// matchTokenInCorpus 判断 token 是否在语料中出现。
func matchTokenInCorpus(token string, corpus recordSearchCorpus) bool {
	t := strings.TrimSpace(token)
	if t == "" {
		return false
	}
	upper := strings.ToUpper(t)
	if upper != "" && strings.Contains(corpus.upper, upper) {
		return true
	}
	lower := strings.ToLower(t)
	if lower != "" && strings.Contains(corpus.lower, lower) {
		return true
	}
	if strings.Contains(corpus.raw, t) {
		return true
	}
	if corpus.normalized != "" {
		if normalized := normalizeNameForMatch(t); normalized != "" && strings.Contains(corpus.normalized, normalized) {
			return true
		}
	}
	return false
}

// composeSkuFollowUpReply 按照匹配结果生成针对 SKU 的详细回复。
func composeSkuFollowUpReply(ctxData *directSkuContext, records []ProductPriceRecord, query string) string {
	if len(records) == 0 {
		return ""
	}
	name := strings.TrimSpace(ctxData.StandardName)
	if name == "" {
		name = strings.TrimSpace(ctxData.Candidate)
	}
	queryTrimmed := strings.TrimSpace(query)

	if len(records) == 1 {
		record := records[0]
		lines := make([]string, 0, 5)
		if name != "" {
			if queryTrimmed != "" {
				lines = append(lines, fmt.Sprintf("%s 中符合“%s”的规格：", name, queryTrimmed))
			} else {
				lines = append(lines, fmt.Sprintf("%s 的详细规格：", name))
			}
		} else if queryTrimmed != "" {
			lines = append(lines, fmt.Sprintf("符合“%s”的规格：", queryTrimmed))
		}
		lines = append(lines, summarizeProductRecordDetailed(record))
		lines = append(lines, "如需报价或其它参数，请继续告诉我。")
		return strings.Join(lines, "\n")
	}

	lines := make([]string, 0, len(records)+4)
	if name != "" {
		if queryTrimmed != "" {
			lines = append(lines, fmt.Sprintf("%s 中符合“%s”的共有 %d 个规格：", name, queryTrimmed, len(records)))
		} else {
			lines = append(lines, fmt.Sprintf("%s 当前可选的 %d 个规格：", name, len(records)))
		}
	} else if queryTrimmed != "" {
		lines = append(lines, fmt.Sprintf("符合“%s”的共有 %d 个规格：", queryTrimmed, len(records)))
	} else {
		lines = append(lines, fmt.Sprintf("共有 %d 个候选规格：", len(records)))
	}
	for idx, record := range records {
		lines = append(lines, fmt.Sprintf("%d. %s", idx+1, summarizeProductRecord(record)))
	}
	lines = append(lines, "如需其中某个规格的更详细信息，请告诉我序号或继续补充需求。")
	return strings.Join(lines, "\n")
}

// extractSkuCandidates 从文本提取疑似产品型号的 token。
func extractSkuCandidates(text string) []string {
	raw := strings.TrimSpace(text)
	if raw == "" {
		return nil
	}
	matches := skuCodePattern.FindAllString(raw, -1)
	if len(matches) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(matches))
	result := make([]string, 0, len(matches))
	for _, match := range matches {
		token := normalizeSkuCandidate(match)
		if token == "" {
			continue
		}
		if _, ok := seen[token]; ok {
			continue
		}
		seen[token] = struct{}{}
		result = append(result, token)
	}
	return result
}

// normalizeSkuCandidate 对产品型号 token 做标准化处理。
func normalizeSkuCandidate(token string) string {
	t := strings.ToUpper(strings.TrimSpace(token))
	if t == "" {
		return ""
	}
	replacer := strings.NewReplacer(
		"（", "",
		"）", "",
		"(", "",
		")", "",
		"，", "",
		"。", "",
		",", "",
		"！", "",
		"!", "",
		"？", "",
		"?", "",
		"；", "",
		";", "",
		"：", "",
		":", "",
		"×", "X",
		"–", "-",
		"—", "-",
		"－", "-",
		" ", "",
	)
	t = replacer.Replace(t)
	t = strings.Trim(t, ".-")
	if t == "" {
		return ""
	}
	return t
}

// hasSkuIntentKeyword 判断文本是否包含明显的型号/规格意图。
func hasSkuIntentKeyword(text string) bool {
	for _, kw := range skuIntentKeywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	lower := strings.ToLower(text)
	for _, kw := range skuIntentKeywordsLower {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

type skuLookupResult struct {
	StandardName string
	Records      []ProductPriceRecord
	Total        int64
	SearchUsed   string
	Images       []string
}

type recordSearchCorpus struct {
	raw        string
	lower      string
	upper      string
	normalized string
}

// lookupSkuRecords 调用数据库查询匹配的产品记录。
func (m *Manager) lookupSkuRecords(ctx context.Context, candidate string, limit int) (*skuLookupResult, error) {
	searchName := candidate
	if m.wikiMng != nil {
		if normalizations, _, _, err := m.queryWikiNameNormalizations(ctx, m.wikiMng, candidate, m.config.DefaultNameNormalizationLimit); err == nil {
			if selected, ok := selectNormalizationCandidate(normalizations, candidate); ok && strings.TrimSpace(selected.StandardName) != "" {
				searchName = strings.TrimSpace(selected.StandardName)
			}
		} else {
			fmt.Printf("normalize product name error (%s): %v\n", candidate, err)
		}
	}

	attempts := make([]string, 0, 2)
	seen := make(map[string]struct{}, 2)
	addAttempt := func(value string) {
		val := strings.TrimSpace(value)
		if val == "" {
			return
		}
		key := strings.ToLower(val)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		attempts = append(attempts, val)
	}
	addAttempt(searchName)
	addAttempt(candidate)

	var lastErr error
	for _, attempt := range attempts {
		if m.productPriceQuery == nil {
			return nil, errors.New("产品价格查询接口未初始化")
		}
		records, total, _, err := m.productPriceQuery.QueryProductPrices(ctx, ProductPriceQueryParams{
			SearchText: attempt,
			Limit:      limit,
		})
		if err != nil {
			lastErr = err
			continue
		}
		if total == 0 || len(records) == 0 {
			continue
		}

		result := &skuLookupResult{
			StandardName: strings.TrimSpace(searchName),
			Records:      records,
			Total:        total,
			SearchUsed:   attempt,
			Images:       collectThumbnailsFromRecords(records),
		}
		return result, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, nil
}

// selectNormalizationCandidate 从知识库标准化结果中挑选最合适的条目。
func selectNormalizationCandidate(results []wikiNameNormalization, candidate string) (wikiNameNormalization, bool) {
	if len(results) == 0 {
		return wikiNameNormalization{}, false
	}
	normalizedCandidate := normalizeNameForMatch(candidate)
	for _, item := range results {
		if len(item.MatchedValues) > 0 {
			return item, true
		}
		if normalizeNameForMatch(item.StandardName) == normalizedCandidate && strings.TrimSpace(item.StandardName) != "" {
			return item, true
		}
	}
	return results[0], true
}

// composeSkuReply 根据查询结果拼装首轮回复。
func composeSkuReply(result *skuLookupResult, original string) string {
	if result == nil || len(result.Records) == 0 {
		return ""
	}
	nameForDisplay := strings.TrimSpace(result.StandardName)
	if nameForDisplay == "" {
		nameForDisplay = strings.TrimSpace(result.SearchUsed)
	}
	if nameForDisplay == "" {
		nameForDisplay = strings.TrimSpace(original)
	}

	lines := make([]string, 0, len(result.Records)+4)
	if !strings.EqualFold(strings.TrimSpace(original), nameForDisplay) {
		lines = append(lines, fmt.Sprintf("“%s”对应的标准型号是 %s。", strings.TrimSpace(original), nameForDisplay))
	}
	lines = append(lines, fmt.Sprintf("%s 系列目前找到 %d 个规格，我先列出前 %d 条：", nameForDisplay, result.Total, len(result.Records)))

	for idx, record := range result.Records {
		lines = append(lines, fmt.Sprintf("%d. %s", idx+1, summarizeProductRecord(record)))
	}

	if result.Total > int64(len(result.Records)) {
		lines = append(lines, fmt.Sprintf("仅展示前 %d 条，如需其它规格请告诉我。", len(result.Records)))
	}
	lines = append(lines, "可以告诉我期望的孔径、电流等参数，我再帮您确认最合适的型号。")

	return strings.Join(lines, "\n")
}

// summarizeProductRecord 输出产品记录的简短摘要。
func summarizeProductRecord(record ProductPriceRecord) string {
	parts := make([]string, 0, 4)
	if display := strings.TrimSpace(record.DisplayName); display != "" {
		parts = append(parts, display)
	} else if attrSummary := summarizeProductAttributes(record.Attributes); attrSummary != "" {
		parts = append(parts, attrSummary)
	}
	if code := strings.TrimSpace(record.SKUCode); code != "" {
		parts = append(parts, fmt.Sprintf("SKU: %s", code))
	}
	if record.PricePromot != nil {
		price := *record.PricePromot
		parts = append(parts, formatPrice("参考价", price, record.Currency))
	} else if record.PriceIn != nil {
		price := *record.PriceIn
		parts = append(parts, formatPrice("含税价", price, record.Currency))
	}
	return strings.Join(parts, "，")
}

// summarizeProductAttributes 汇总产品属性为短语形式。
func summarizeProductAttributes(attrs map[string]string) string {
	if len(attrs) == 0 {
		return ""
	}
	keys := make([]string, 0, len(attrs))
	for key := range attrs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		val := strings.TrimSpace(attrs[key])
		if val == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s:%s", key, val))
	}
	return strings.Join(parts, " / ")
}

// summarizeProductRecordDetailed 输出产品记录的详细描述。
func summarizeProductRecordDetailed(record ProductPriceRecord) string {
	lines := make([]string, 0, 6)
	if display := strings.TrimSpace(record.DisplayName); display != "" {
		lines = append(lines, display)
	}
	if code := strings.TrimSpace(record.SKUCode); code != "" {
		lines = append(lines, fmt.Sprintf("SKU: %s", code))
	}
	if attrDetail := formatAttributeDetails(record.Attributes); attrDetail != "" {
		lines = append(lines, attrDetail)
	}
	if record.PricePromot != nil {
		lines = append(lines, formatPrice("优惠价", *record.PricePromot, record.Currency))
	} else if record.PriceIn != nil {
		lines = append(lines, formatPrice("参考价", *record.PriceIn, record.Currency))
	}
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n")
}

// formatAttributeDetails 以多行格式展示产品属性详情。
func formatAttributeDetails(attrs map[string]string) string {
	if len(attrs) == 0 {
		return ""
	}
	keys := make([]string, 0, len(attrs))
	for key := range attrs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	lines := make([]string, 0, len(keys)+1)
	lines = append(lines, "关键参数：")
	for _, key := range keys {
		val := strings.TrimSpace(attrs[key])
		if val == "" {
			continue
		}
		lines = append(lines, fmt.Sprintf("  - %s：%s", key, val))
	}
	if len(lines) == 1 {
		return ""
	}
	return strings.Join(lines, "\n")
}

// formatPrice 统一格式化价格信息。
func formatPrice(label string, price float64, currency string) string {
	value := fmt.Sprintf("%.2f", price)
	ccy := strings.TrimSpace(currency)
	if ccy != "" {
		return fmt.Sprintf("%s %s %s", label, ccy, value)
	}
	return fmt.Sprintf("%s %s元", label, value)
}

// collectThumbnailsFromRecords 提取产品记录里的缩略图链接去重集合。
func collectThumbnailsFromRecords(records []ProductPriceRecord) []string {
	if len(records) == 0 {
		return nil
	}
	seen := make(map[string]struct{})
	result := make([]string, 0, len(records))
	for _, record := range records {
		url := strings.TrimSpace(record.Thumbnail)
		if url == "" {
			continue
		}
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}
		result = append(result, url)
	}
	return result
}
