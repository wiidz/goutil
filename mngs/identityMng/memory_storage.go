package identityMng

import (
	"github.com/click33/sa-token-go/core/adapter"
	"github.com/click33/sa-token-go/storage/memory"
)

// NewMemoryStorage 返回内存存储实现（直接使用官方 memory storage）
func NewMemoryStorage() adapter.Storage {
	return memory.NewStorage()
}
