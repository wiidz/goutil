package repoMng

import (
	"context"
	"errors"
	"sync"

	"gorm.io/gorm"
)

// Manager holds named *gorm.DB instances and provides typed repo access.
type Manager struct {
	mu  sync.RWMutex
	dbs map[string]*gorm.DB
	def string
}

func NewManager() *Manager { return &Manager{dbs: make(map[string]*gorm.DB)} }

// SetupDefault sets the default database.
func (m *Manager) SetupDefault(db *gorm.DB) {
	m.mu.Lock()
	m.dbs["default"] = db
	m.def = "default"
	m.mu.Unlock()
}

// Register registers a named database.
func (m *Manager) Register(name string, db *gorm.DB) {
	m.mu.Lock()
	m.dbs[name] = db
	if m.def == "" {
		m.def = name
	}
	m.mu.Unlock()
}

// Default returns a Set bound to the default DB.
func (m *Manager) Default() *Set { return m.For(m.def) }

// For returns a Set bound to the named DB.
func (m *Manager) For(name string) *Set {
	m.mu.RLock()
	db := m.dbs[name]
	m.mu.RUnlock()
	if db == nil {
		return &Set{db: nil}
	}
	return &Set{db: db}
}

// InTx runs fn in a transaction for the named DB.
func (m *Manager) InTx(name string, ctx context.Context, fn func(ctx context.Context, s *Set) error) error {
	s := m.For(name)
	if s.db == nil {
		return errors.New("repoMng: db not found: " + name)
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error { return fn(ctx, s.BindTx(tx)) })
}

// Set is a typed-access view on top of a *gorm.DB.
type Set struct{ db *gorm.DB }

func (s *Set) DB() *gorm.DB { return s.db }

func (s *Set) BindTx(tx *gorm.DB) *Set { return &Set{db: tx} }
