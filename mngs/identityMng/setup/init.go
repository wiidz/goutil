package setup

import (
	sagin "github.com/click33/sa-token-go/integrations/gin"
	"github.com/click33/sa-token-go/storage/memory"
)

// Init initializes Sa-Token manager (required by stputil) using integrations/gin + memory storage.
func Init() {
	storage := memory.NewStorage()
	cfg := sagin.DefaultConfig()
	manager := sagin.NewManager(storage, cfg)
	sagin.SetManager(manager)
}
