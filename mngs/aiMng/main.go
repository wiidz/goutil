package aiMng

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wiidz/goutil/mngs/eyunMng/msgSend"
)

// SearchWithProductPrice 是 Manager 封装的对话检索函数，结合数据库与知识库为会话生成回复消息。
func (m *Manager) SearchWithProductPrice(
	ctx context.Context,
	sessionID string,
	prompt string,
	content string,
	messageTime *time.Time,
	returnTokenUsage bool,
) ([]msgSend.MessagePayload, error) {
	return m.searchWithProductPrice(ctx, sessionID, prompt, content, messageTime, returnTokenUsage)
}

// Search 是 Manager 暴露的基础对话函数，基于提示词执行一次简单的模型交互。
func (m *Manager) Search(
	ctx *gin.Context,
	prompt string,
	content string,
	returnTokenUsage bool,
) (string, error) {
	return m.aiSearch(ctx, prompt, content, returnTokenUsage)
}

// ClearSessionCache 是 Manager 的缓存管理函数，用于主动清理指定会话的上下文缓存。
func (m *Manager) ClearSessionCache(ctx context.Context, sessionID string) error {
	err1 := m.clearDirectSkuContext(ctx, sessionID)
	err2 := m.clearSessionProductContext(ctx, sessionID)
	clearSessionProductSnapshot(sessionID)
	if err1 != nil {
		return err1
	}
	return err2
}

