package aiMng

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/wiidz/goutil/helpers/typeHelper"
)

const productContextPrefix = "[产品上下文]"

// sessionRedisKey 是一个构造函数，用于生成会话缓存对应的 Redis 键。
func (m *Manager) sessionRedisKey(sessionID string) string {
	return fmt.Sprintf("chat:session:%s", sessionID)
}

// loadSessionConversation 是会话读取函数，用于从 Redis 中取回历史会话消息。
func (m *Manager) loadSessionConversation(ctx context.Context, sessionID string) ([]*model.ChatCompletionMessage, error) {
	if m == nil || m.redis == nil || sessionID == "" {
		return []*model.ChatCompletionMessage{}, nil
	}

	raw, err := m.redis.GetString(ctx, m.sessionRedisKey(sessionID))
	if err != nil {
		return nil, err
	}
	if raw == "" {
		return []*model.ChatCompletionMessage{}, nil
	}

	var payload sessionPayload
	if err = json.Unmarshal([]byte(raw), &payload); err == nil && len(payload.Entries) > 0 {
		messages := make([]*model.ChatCompletionMessage, 0, len(payload.Entries))
		timestamps := make([]messageTimestamp, 0, len(payload.Entries))
		for _, entry := range payload.Entries {
			messages = append(messages, entry.Message)
			timestamps = append(timestamps, m.normalizeTimestamp(entry.SavedAt, entry.SavedUnix))
		}
		m.recordSessionMessageTimes(sessionID, timestamps)
		if len(timestamps) > 0 && timestamps[len(timestamps)-1].Value != "" {
			m.recordSessionSavedAt(sessionID, timestamps[len(timestamps)-1].Value)
		}
		return messages, nil
	}

	var legacyPayload legacySessionPayload
	if err = json.Unmarshal([]byte(raw), &legacyPayload); err == nil && len(legacyPayload.Messages) > 0 {
		messages := legacyPayload.Messages
		ts := m.normalizeTimestamp(legacyPayload.SavedAt, legacyPayload.SavedUnix)
		if ts.Value != "" {
			m.recordSessionSavedAt(sessionID, ts.Value)
		}
		if len(messages) > 0 && ts.Value != "" {
			timestamps := make([]messageTimestamp, len(messages))
			for i := range timestamps {
				timestamps[i] = ts
			}
			m.recordSessionMessageTimes(sessionID, timestamps)
		}
		return messages, nil
	}

	var legacy []*model.ChatCompletionMessage
	if err = json.Unmarshal([]byte(raw), &legacy); err != nil {
		return nil, err
	}
	// legacy format without timestamps
	m.recordSessionMessageTimes(sessionID, nil)
	return legacy, nil
}

// saveSessionConversation 是将会话数据写入 Redis 的持久化函数。
func (m *Manager) saveSessionConversation(ctx context.Context, sessionID string, messages []*model.ChatCompletionMessage) error {
	if m == nil || m.redis == nil || sessionID == "" {
		return nil
	}

	existing := m.sessionMessageTimesFor(sessionID)
	if len(existing) > len(messages) {
		existing = existing[len(existing)-len(messages):]
	}
	newTimes := make([]messageTimestamp, len(messages))
	copyCount := copy(newTimes, existing)
	for i := 0; i < len(messages); i++ {
		if i < copyCount && newTimes[i].Value != "" {
			continue
		}
		now := time.Now()
		newTimes[i] = messageTimestamp{
			Value: now.Format(time.DateTime),
			Unix:  now.Unix(),
		}
	}

	now := time.Now()
	entries := make([]sessionEntry, len(messages))
	for i, msg := range messages {
		ts := newTimes[i]
		if ts.Value == "" {
			ts = messageTimestamp{Value: now.Format(time.DateTime), Unix: now.Unix()}
			newTimes[i] = ts
		}
		entries[i] = sessionEntry{
			Message:   msg,
			SavedAt:   ts.Value,
			SavedUnix: ts.Unix,
		}
	}
	payload := sessionPayload{Entries: entries}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := m.redis.Set(ctx, m.sessionRedisKey(sessionID), data, m.config.SessionTTL); err != nil {
		return err
	}
	m.recordSessionMessageTimes(sessionID, newTimes)
	if len(newTimes) > 0 && newTimes[len(newTimes)-1].Value != "" {
		m.recordSessionSavedAt(sessionID, newTimes[len(newTimes)-1].Value)
	} else {
		m.recordSessionSavedAt(sessionID, now.Format(time.DateTime))
	}
	return nil
}

// trimConversationHistory 是会话整理函数，确保消息数量在设定上限内。
func (m *Manager) trimConversationHistory(messages []*model.ChatCompletionMessage) []*model.ChatCompletionMessage {
	if len(messages) <= m.config.MaxSessionMessages {
		return messages
	}

	if len(messages) == 0 {
		return messages
	}

	first := messages[0]
	if first != nil && first.Role == model.ChatMessageRoleSystem {
		capacity := m.config.MaxSessionMessages
		trimmed := make([]*model.ChatCompletionMessage, 0, capacity)
		trimmed = append(trimmed, first)

		remaining := messages[1:]
		if len(remaining) > capacity-1 {
			remaining = remaining[len(remaining)-(capacity-1):]
		}
		trimmed = append(trimmed, remaining...)
		return trimmed
	}

	start := len(messages) - m.config.MaxSessionMessages
	if start < 0 {
		start = 0
	}
	return append([]*model.ChatCompletionMessage(nil), messages[start:]...)
}

// debugPrintSession 是调试辅助函数，用于打印会话状态信息。
func (m *Manager) debugPrintSession(sessionID, label string, messages []*model.ChatCompletionMessage) {
	if m == nil {
		return
	}
	displaySessionID := sessionID
	if displaySessionID == "" {
		displaySessionID = "unknown"
	}
	lastSaved := m.sessionSavedAt(sessionID)
	if lastSaved == "" {
		lastSaved = "unknown"
	}
	product := sessionProductSnapshotFor(sessionID)
	infoParts := make([]string, 0, 5)
	if product.Primary != "" {
		infoParts = append(infoParts, product.Primary)
	}
	if product.SpuID != 0 {
		infoParts = append(infoParts, fmt.Sprintf("spu:%d", product.SpuID))
	}
	if product.SkuCount > 0 {
		infoParts = append(infoParts, fmt.Sprintf("sku_count:%d", product.SkuCount))
	}
	if len(product.AttributeNames) > 0 {
		infoParts = append(infoParts, fmt.Sprintf("attrs:%s", strings.Join(product.AttributeNames, ", ")))
	}
	if len(product.Candidates) > 0 {
		infoParts = append(infoParts, fmt.Sprintf("candidates:%s", strings.Join(product.Candidates, ", ")))
	}
	if len(infoParts) == 0 {
		infoParts = append(infoParts, "unknown")
	}
	productInfo := strings.Join(infoParts, " | ")
	nowStr := time.Now().Format(time.DateTime)
	fmt.Printf("\n---------- [%s] %s [session:%s] (%d messages, last saved at %s, product %s)\n", nowStr, label, displaySessionID, len(messages), lastSaved, productInfo)

	const maxShown = 6
	start := 0
	if len(messages) > maxShown {
		start = len(messages) - maxShown
	}
	messageTimes := m.sessionMessageTimesFor(sessionID)
	for i := start; i < len(messages); i++ {
		msg := messages[i]
		if msg == nil {
			fmt.Printf("  #%d <nil>\n", i)
			continue
		}
		content := m.messageSummary(msg)
		tsStr := "unknown"
		if i < len(messageTimes) {
			ts := messageTimes[i]
			if ts.Value != "" {
				tsStr = ts.Value
			} else if ts.Unix != 0 {
				tsStr = time.Unix(ts.Unix, 0).Format(time.DateTime)
			}
		}
		fmt.Printf("  #%s %s [%s]: %s\n", fmt.Sprintf("%-*s", 3, typeHelper.Int2Str(i)), tsStr, fmt.Sprintf("%-*s", 9, msg.Role), content)
	}
}

type sessionPayload struct {
	Entries []sessionEntry `json:"entries"`
}

type sessionEntry struct {
	Message   *model.ChatCompletionMessage `json:"message"`
	SavedAt   string                       `json:"saved_at,omitempty"`
	SavedUnix int64                        `json:"saved_unix,omitempty"`
}

type legacySessionPayload struct {
	Messages  []*model.ChatCompletionMessage `json:"messages"`
	SavedAt   string                         `json:"saved_at"`
	SavedUnix int64                          `json:"saved_unix"`
}

var (
	sessionSavedAtMu      sync.RWMutex
	sessionSavedMap       = make(map[string]string)
	sessionMessageTimesMu sync.RWMutex
	sessionMessageTimes   = make(map[string][]messageTimestamp)
	sessionProductMu      sync.RWMutex
	sessionProductMap     = make(map[string]sessionProductSnapshot)
)

type messageTimestamp struct {
	Value string
	Unix  int64
}

type sessionProductSnapshot struct {
	Primary        string
	Candidates     []string
	SpuID          uint64
	SkuCount       int
	AttributeNames []string
}

// recordSessionSavedAt 是记录函数，保存最新的会话持久化时间。
func (m *Manager) recordSessionSavedAt(sessionID, savedAt string) {
	if sessionID == "" || savedAt == "" {
		return
	}
	sessionSavedAtMu.Lock()
	sessionSavedMap[sessionID] = savedAt
	sessionSavedAtMu.Unlock()
}

// sessionSavedAt 是查询函数，返回会话最近一次持久化的时间。
func (m *Manager) sessionSavedAt(sessionID string) string {
	if sessionID == "" {
		return ""
	}
	if m == nil {
		return ""
	}
	sessionSavedAtMu.RLock()
	defer sessionSavedAtMu.RUnlock()
	return sessionSavedMap[sessionID]
}

// recordSessionMessageTimes 是记录函数，用于保存会话中各消息的时间戳。
func (m *Manager) recordSessionMessageTimes(sessionID string, times []messageTimestamp) {
	if sessionID == "" {
		return
	}
	sessionMessageTimesMu.Lock()
	if len(times) == 0 {
		delete(sessionMessageTimes, sessionID)
	} else {
		copied := make([]messageTimestamp, len(times))
		copy(copied, times)
		sessionMessageTimes[sessionID] = copied
	}
	sessionMessageTimesMu.Unlock()
}

// sessionMessageTimesFor 是查询函数，获取特定会话的消息时间戳列表。
func (m *Manager) sessionMessageTimesFor(sessionID string) []messageTimestamp {
	if sessionID == "" {
		return nil
	}
	sessionMessageTimesMu.RLock()
	defer sessionMessageTimesMu.RUnlock()
	src := sessionMessageTimes[sessionID]
	if len(src) == 0 {
		return nil
	}
	dst := make([]messageTimestamp, len(src))
	copy(dst, src)
	return dst
}

// recordSessionProductSnapshot 是记录函数，用于保存会话中的产品快照信息。
func recordSessionProductSnapshot(sessionID string, snapshot sessionProductSnapshot) {
	if sessionID == "" {
		return
	}
	copySnap := sessionProductSnapshot{
		Primary:  snapshot.Primary,
		SpuID:    snapshot.SpuID,
		SkuCount: snapshot.SkuCount,
	}
	if len(snapshot.Candidates) > 0 {
		copySnap.Candidates = cloneStringSlice(snapshot.Candidates)
	}
	if len(snapshot.AttributeNames) > 0 {
		copySnap.AttributeNames = cloneStringSlice(snapshot.AttributeNames)
	}
	sessionProductMu.Lock()
	sessionProductMap[sessionID] = copySnap
	sessionProductMu.Unlock()
}

// recordSessionProductSnapshotFromContext 是根据上下文生成产品快照的辅助函数。
func recordSessionProductSnapshotFromContext(sessionID string, ctx *sessionProductContext) {
	if sessionID == "" {
		return
	}
	if ctx == nil {
		clearSessionProductSnapshot(sessionID)
		return
	}
	snapshot := sessionProductSnapshot{}
	if ctx.Primary != nil {
		snapshot.Primary = describeProductCandidate(*ctx.Primary)
	}
	if ctx.SpuID != 0 {
		snapshot.SpuID = ctx.SpuID
	}
	if ctx.SkuCount != 0 {
		snapshot.SkuCount = ctx.SkuCount
	}
	if len(ctx.Candidates) > 0 {
		snapshot.Candidates = make([]string, 0, len(ctx.Candidates))
		for _, c := range ctx.Candidates {
			snapshot.Candidates = append(snapshot.Candidates, describeProductCandidate(c))
		}
	}
	if len(ctx.AttributeNames) > 0 {
		snapshot.AttributeNames = cloneStringSlice(ctx.AttributeNames)
	}
	recordSessionProductSnapshot(sessionID, snapshot)
}

// clearSessionProductSnapshot 是清理函数，用于删除会话对应的产品快照。
func clearSessionProductSnapshot(sessionID string) {
	if sessionID == "" {
		return
	}
	sessionProductMu.Lock()
	delete(sessionProductMap, sessionID)
	sessionProductMu.Unlock()
}

// sessionProductSnapshotFor 是查询函数，返回指定会话的产品快照数据。
func sessionProductSnapshotFor(sessionID string) sessionProductSnapshot {
	if sessionID == "" {
		return sessionProductSnapshot{}
	}
	sessionProductMu.RLock()
	defer sessionProductMu.RUnlock()
	if snap, ok := sessionProductMap[sessionID]; ok {
		result := sessionProductSnapshot{
			Primary:  snap.Primary,
			SpuID:    snap.SpuID,
			SkuCount: snap.SkuCount,
		}
		if len(snap.Candidates) > 0 {
			result.Candidates = cloneStringSlice(snap.Candidates)
		}
		if len(snap.AttributeNames) > 0 {
			result.AttributeNames = cloneStringSlice(snap.AttributeNames)
		}
		return result
	}
	return sessionProductSnapshot{}
}

// cloneStringSlice 是工具函数，用于复制字符串切片。
func cloneStringSlice(src []string) []string {
	if len(src) == 0 {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

// formatProductSnapshotMessage 是格式化函数，将产品快照转换为系统消息文本。
func formatProductSnapshotMessage(snap sessionProductSnapshot) string {
	if snap.Primary == "" && len(snap.Candidates) == 0 && snap.SpuID == 0 && snap.SkuCount == 0 && len(snap.AttributeNames) == 0 {
		return ""
	}
	builder := strings.Builder{}
	builder.WriteString(productContextPrefix)
	builder.WriteString("\n- 主产品: ")
	if snap.Primary != "" {
		builder.WriteString(snap.Primary)
	} else {
		builder.WriteString("unknown")
	}
	if snap.SpuID != 0 {
		builder.WriteString(fmt.Sprintf("\n- SPU ID: %d", snap.SpuID))
	}
	if snap.SkuCount > 0 {
		builder.WriteString(fmt.Sprintf("\n- SKU 数量: %d", snap.SkuCount))
	}
	if len(snap.AttributeNames) > 0 {
		builder.WriteString("\n- 关键属性: ")
		builder.WriteString(strings.Join(snap.AttributeNames, ", "))
	}
	if len(snap.Candidates) > 0 {
		builder.WriteString("\n- 其它候选: ")
		builder.WriteString(strings.Join(snap.Candidates, "; "))
	}
	return builder.String()
}

// appendProductContextMessage 是上下文补充函数，为会话追加产品上下文信息。
func (m *Manager) appendProductContextMessage(ctx context.Context, sessionID string, conversation []*model.ChatCompletionMessage) []*model.ChatCompletionMessage {
	if m == nil || sessionID == "" {
		return conversation
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if _, err := m.loadSessionProductContext(ctx, sessionID); err != nil {
		fmt.Printf("load session product context error: %v\n", err)
	}
	snap := sessionProductSnapshotFor(sessionID)
	msg := formatProductSnapshotMessage(snap)
	if msg == "" {
		return conversation
	}
	if len(conversation) > 0 {
		last := conversation[len(conversation)-1]
		if last != nil && last.Role == model.ChatMessageRoleSystem && last.Content != nil && last.Content.StringValue != nil {
			if strings.HasPrefix(*last.Content.StringValue, productContextPrefix) && *last.Content.StringValue == msg {
				return conversation
			}
		}
	}
	conversation = append(conversation, m.newChatMessage(sessionID, model.ChatMessageRoleSystem, msg, nil))
	return conversation
}

// describeProductCandidate 是描述函数，生成产品候选的文本摘要。
func describeProductCandidate(c productCandidate) string {
	parts := make([]string, 0, 4)
	if main := strings.TrimSpace(c.StandardName); main != "" {
		parts = append(parts, main)
	} else if alt := strings.TrimSpace(c.Candidate); alt != "" {
		parts = append(parts, alt)
	}
	if alias := strings.TrimSpace(c.Candidate); alias != "" && !strings.EqualFold(alias, c.StandardName) {
		parts = append(parts, fmt.Sprintf("alias:%s", alias))
	}
	if search := strings.TrimSpace(c.SearchUsed); search != "" && !strings.EqualFold(search, c.StandardName) && !strings.EqualFold(search, c.Candidate) {
		parts = append(parts, fmt.Sprintf("search:%s", search))
	}
	if src := strings.TrimSpace(c.Source); src != "" {
		parts = append(parts, fmt.Sprintf("source:%s", src))
	}
	if len(parts) == 0 {
		return "unknown"
	}
	return strings.Join(parts, " | ")
}

// candidateKey 是键生成函数，为产品候选生成去重标识。
func candidateKey(c productCandidate) string {
	if key := strings.TrimSpace(c.StandardName); key != "" {
		return strings.ToLower(key)
	}
	if key := strings.TrimSpace(c.Candidate); key != "" {
		return strings.ToLower(key)
	}
	if key := strings.TrimSpace(c.SearchUsed); key != "" {
		return strings.ToLower(key)
	}
	return ""
}

// mergeCandidate 是合并函数，用于合并两个候选产品的信息。
func mergeCandidate(dst *productCandidate, src productCandidate) bool {
	if dst == nil {
		return false
	}
	changed := false
	if strings.TrimSpace(dst.StandardName) == "" && strings.TrimSpace(src.StandardName) != "" {
		dst.StandardName = strings.TrimSpace(src.StandardName)
		changed = true
	}
	if strings.TrimSpace(dst.Candidate) == "" && strings.TrimSpace(src.Candidate) != "" {
		dst.Candidate = strings.TrimSpace(src.Candidate)
		changed = true
	}
	if strings.TrimSpace(dst.SearchUsed) == "" && strings.TrimSpace(src.SearchUsed) != "" {
		dst.SearchUsed = strings.TrimSpace(src.SearchUsed)
		changed = true
	}
	if src.Source != "" {
		if dst.Source == "" {
			dst.Source = src.Source
			changed = true
		} else if !strings.Contains(dst.Source, src.Source) {
			dst.Source = dst.Source + "," + src.Source
			changed = true
		}
	}
	return changed
}

// normalizeProductCandidates 是去重规范化函数，整理候选产品列表。
func normalizeProductCandidates(candidates []productCandidate) []productCandidate {
	if len(candidates) == 0 {
		return nil
	}
	order := make([]string, 0, len(candidates))
	result := make(map[string]*productCandidate, len(candidates))
	fallback := 0
	for _, cand := range candidates {
		key := candidateKey(cand)
		if key == "" {
			key = fmt.Sprintf("_%d", fallback)
			fallback++
		}
		if existing, ok := result[key]; ok {
			if mergeCandidate(existing, cand) {
				// merged
			}
			continue
		}
		cp := cand
		result[key] = &cp
		order = append(order, key)
	}
	normalized := make([]productCandidate, 0, len(order))
	for _, key := range order {
		normalized = append(normalized, *result[key])
	}
	return normalized
}

// addCandidate 是追加函数，将候选产品加入列表并处理去重。
func addCandidate(list *[]productCandidate, cand productCandidate) bool {
	if list == nil {
		return false
	}
	key := candidateKey(cand)
	if key != "" {
		for i := range *list {
			if candidateKey((*list)[i]) == key {
				return mergeCandidate(&(*list)[i], cand)
			}
		}
	}
	*list = append(*list, cand)
	return true
}

// removeCandidateByKey 是删除函数，按键移除候选产品。
func removeCandidateByKey(list []productCandidate, key string) []productCandidate {
	if key == "" || len(list) == 0 {
		return list
	}
	result := make([]productCandidate, 0, len(list))
	for _, cand := range list {
		if candidateKey(cand) == key {
			continue
		}
		result = append(result, cand)
	}
	return result
}

// getCandidateByKey 是查询函数，从上下文中按键获取候选产品。
func getCandidateByKey(ctx *sessionProductContext, key string) (productCandidate, bool) {
	if ctx == nil || key == "" {
		return productCandidate{}, false
	}
	if ctx.Primary != nil && candidateKey(*ctx.Primary) == key {
		return *ctx.Primary, true
	}
	for _, cand := range ctx.Candidates {
		if candidateKey(cand) == key {
			return cand, true
		}
	}
	return productCandidate{}, false
}

// applyAssistantCandidates 是上下文同步函数，将助手返回的候选并入上下文。
func applyAssistantCandidates(ctx *sessionProductContext, newCandidates []productCandidate) bool {
	if ctx == nil || len(newCandidates) == 0 {
		return false
	}
	normalized := normalizeProductCandidates(newCandidates)
	if len(normalized) == 0 {
		return false
	}
	changed := false
	if ctx.Primary == nil {
		if len(normalized) == 1 {
			candidate := normalized[0]
			ctx.Primary = &candidate
			ctx.Candidates = removeCandidateByKey(ctx.Candidates, candidateKey(candidate))
			changed = true
		} else {
			for _, cand := range normalized {
				if addCandidate(&ctx.Candidates, cand) {
					changed = true
				}
			}
			if changed {
				ctx.Candidates = normalizeProductCandidates(ctx.Candidates)
			}
		}
		return changed
	}

	primaryKey := candidateKey(*ctx.Primary)
	for _, cand := range normalized {
		key := candidateKey(cand)
		if key != "" && key == primaryKey {
			if mergeCandidate(ctx.Primary, cand) {
				changed = true
			}
			continue
		}
		if addCandidate(&ctx.Candidates, cand) {
			changed = true
		}
	}
	if len(ctx.Candidates) > 0 {
		ctx.Candidates = normalizeProductCandidates(ctx.Candidates)
		if primaryKey != "" {
			ctx.Candidates = removeCandidateByKey(ctx.Candidates, primaryKey)
		}
	}
	return changed
}

// updateProductContextFromAssistant 是上下文更新函数，用于解析助手回复并刷新产品信息。
func (m *Manager) updateProductContextFromAssistant(ctx context.Context, sessionID string, reply string) {
	if m == nil || sessionID == "" {
		return
	}
	reply = strings.TrimSpace(reply)
	if reply == "" {
		return
	}
	tokens := extractSkuCandidates(reply)
	if len(tokens) == 0 {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var productCtx *sessionProductContext
	var err error
	if m.redis != nil {
		productCtx, err = m.loadSessionProductContext(ctx, sessionID)
		if err != nil {
			fmt.Printf("load session product context error: %v\n", err)
		}
	}
	if productCtx == nil {
		productCtx = &sessionProductContext{}
	}

	seen := make(map[string]struct{}, len(tokens))
	candidates := make([]productCandidate, 0, len(tokens))
	for _, token := range tokens {
		normalized := normalizeSkuCandidate(token)
		if normalized == "" {
			normalized = strings.ToUpper(strings.TrimSpace(token))
		}
		normalized = strings.TrimSpace(normalized)
		if normalized == "" {
			continue
		}
		key := strings.ToLower(normalized)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		cand := productCandidate{
			Candidate: normalized,
			Source:    "assistant_reply",
		}
		if existing, ok := getCandidateByKey(productCtx, key); ok {
			cand = existing
			mergeCandidate(&cand, productCandidate{Source: "assistant_reply"})
		} else {
			if result, err := m.lookupSkuRecords(ctx, normalized, 0); err == nil && result != nil && len(result.Records) > 0 {
				cand.StandardName = strings.TrimSpace(result.StandardName)
				cand.SearchUsed = result.SearchUsed
			} else if err != nil {
				fmt.Printf("lookup sku records via assistant reply error (%s): %v\n", normalized, err)
			}
		}
		candidates = append(candidates, cand)
	}

	if len(candidates) == 0 {
		return
	}

	if !applyAssistantCandidates(productCtx, candidates) {
		return
	}

	productCtx.LastQuery = ""
	if err := m.saveSessionProductContext(ctx, sessionID, productCtx); err != nil {
		fmt.Printf("save session product context error: %v\n", err)
	}
}

// normalizeTimestamp 是时间标准化函数，将字符串和时间戳对齐。
func (m *Manager) normalizeTimestamp(value string, unix int64) messageTimestamp {
	ts := messageTimestamp{
		Value: value,
		Unix:  unix,
	}
	if ts.Value == "" && ts.Unix != 0 {
		ts.Value = time.Unix(ts.Unix, 0).Format(time.DateTime)
	}
	if ts.Unix == 0 && ts.Value != "" {
		if parsed, err := time.Parse(time.DateTime, ts.Value); err == nil {
			ts.Unix = parsed.Unix()
		}
	}
	return ts
}

// messageSummary 是调试函数，用于为消息生成日志摘要。
func (m *Manager) messageSummary(msg *model.ChatCompletionMessage) string {
	if msg == nil || msg.Content == nil {
		return "<empty>"
	}
	if msg.Content.StringValue != nil {
		return m.truncateForLog(*msg.Content.StringValue)
	}
	if len(msg.Content.ListValue) > 0 {
		return fmt.Sprintf("[content parts: %d]", len(msg.Content.ListValue))
	}
	return "<empty>"
}

// truncateForLog 是日志辅助函数，对长文本进行截断。
func (m *Manager) truncateForLog(text string) string {
	const limit = 120
	if len(text) <= limit {
		return text
	}
	return text[:limit] + "..."
}
