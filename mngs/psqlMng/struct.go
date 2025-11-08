package psqlMng

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config for PostgreSQL manager
type Config struct {
	DSN             string           // postgres DSN
	ConnMaxIdle     int              // idle connections
	ConnMaxOpen     int              // open connections
	ConnMaxLifetime time.Duration    // connection max lifetime
	Logger          logger.Interface `mapstructure:"logger"`
}

// Manager wraps a gorm DB for PostgreSQL
type Manager struct {
	db *gorm.DB
}

// DB returns underlying *gorm.DB
func (m *Manager) DB() *gorm.DB { return m.db }

// WithTx executes fn in a transaction
func (m *Manager) WithTx(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error {
	if m.db == nil {
		return fn(ctx, nil)
	}
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error { return fn(ctx, tx) })
}
