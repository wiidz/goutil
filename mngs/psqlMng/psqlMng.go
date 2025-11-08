package psqlMng

import (
	"context"
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewMng creates a PostgreSQL manager
func NewMng(cfg *Config) (*Manager, error) {
	if cfg == nil || cfg.DSN == "" {
		return nil, errors.New("psqlMng: DSN is required")
	}
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: cfg.Logger,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err == nil {
		if cfg.ConnMaxIdle > 0 {
			sqlDB.SetMaxIdleConns(cfg.ConnMaxIdle)
		}
		if cfg.ConnMaxOpen > 0 {
			sqlDB.SetMaxOpenConns(cfg.ConnMaxOpen)
		}
		if cfg.ConnMaxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
		}
	}
	return &Manager{db: db}, nil
}

// AutoMigrate runs gorm automigrate on provided models
func (m *Manager) AutoMigrate(models ...interface{}) error {
	if m.db == nil {
		return nil
	}
	return m.db.AutoMigrate(models...)
}

// Ping checks db connectivity
func (m *Manager) Ping() error {
	if m.db == nil {
		return errors.New("psqlMng: db is nil")
	}
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return sqlDB.PingContext(ctx)
}

// Close closes the underlying connection pool
func (m *Manager) Close() error {
	if m.db == nil {
		return nil
	}
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
