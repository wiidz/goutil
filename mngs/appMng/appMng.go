package appMng

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/wiidz/goutil/mngs/amqpMng"
	"github.com/wiidz/goutil/mngs/esMng"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"github.com/wiidz/goutil/mngs/psqlMng"
	"github.com/wiidz/goutil/mngs/redisMng"
	"github.com/wiidz/goutil/structs/configStruct"
)

// Manager 负责缓存和复用 AppMng 实例。
type Manager struct {
	mu    sync.RWMutex
	cache *cache.Cache
}

var defaultManager = NewManager()

// NewManager 创建一个新的 Manager。
func NewManager() *Manager {
	return &Manager{cache: cache.New(defaultCacheTTL, cacheCleanupCycle)}
}

// NewApp 直接构建一个 AppMng，不使用缓存。
func NewApp(ctx context.Context, builder *ConfigBuilder, projectCfg configStruct.ProjectConfig) (*AppMng, error) {
	m := NewManager()
	return m.build(ctx, "default", builder, projectCfg, 0, false)
}

// Get 从缓存获取或构建新的 AppMng。
func (m *Manager) Get(ctx context.Context, opts *Options) (*AppMng, error) {
	if opts == nil || opts.Builder == nil {
		return nil, errors.New("appMng: builder is nil")
	}

	key := opts.ID
	if key == "" {
		key = "default"
	}

	// 读取缓存
	m.mu.RLock()
	if cached, ok := m.cache.Get(key); ok {
		m.mu.RUnlock()
		return cached.(*AppMng), nil
	}
	m.mu.RUnlock()

	app, err := m.build(ctx, key, opts.Builder, opts.ProjectConfig, opts.CacheTTL, true)
	if err != nil {
		return nil, err
	}
	return app, nil
}

// 构建实例，可选择写入缓存
func (m *Manager) build(ctx context.Context, key string, builder *ConfigBuilder, projectCfg configStruct.ProjectConfig, ttl time.Duration, cacheResult bool) (*AppMng, error) {
	baseCfg, err := builder.Build(ctx)
	if err != nil {
		return nil, err
	}
	if baseCfg == nil {
		return nil, errors.New("appMng: builder returned empty base config")
	}

	app := &AppMng{
		ID:            key,
		BaseConfig:    baseCfg,
		ProjectConfig: projectCfg,
	}

	// 初始化依赖（若 loader 未提供实例则自行初始化）
	if app.BaseConfig.MysqlConfig != nil && app.Mysql == nil {
		app.Mysql, err = mysqlMng.NewMysqlMng(app.BaseConfig.MysqlConfig, nil)
		if err != nil {
			return nil, fmt.Errorf("appMng: init mysql failed: %w", err)
		}
	}
	if app.BaseConfig.PostgresConfig != nil && app.Postgres == nil {
		app.Postgres, err = psqlMng.NewMng(&psqlMng.Config{
			DSN:             app.BaseConfig.PostgresConfig.DSN,
			ConnMaxIdle:     app.BaseConfig.PostgresConfig.ConnMaxIdle,
			ConnMaxOpen:     app.BaseConfig.PostgresConfig.ConnMaxOpen,
			ConnMaxLifetime: app.BaseConfig.PostgresConfig.ConnMaxLifetime,
		})
		if err != nil {
			return nil, fmt.Errorf("appMng: init postgres failed: %w", err)
		}
	}
	if app.BaseConfig.RedisConfig != nil && app.Redis == nil {
		if app.Redis, err = redisMng.NewRedisMng(ctx, app.BaseConfig.RedisConfig); err != nil {
			return nil, fmt.Errorf("appMng: init redis failed: %w", err)
		}
	}
	if app.BaseConfig.EsConfig != nil && app.Es == nil {
		if err = esMng.Init(app.BaseConfig.EsConfig); err != nil {
			return nil, fmt.Errorf("appMng: init es failed: %w", err)
		}
	}
	if app.BaseConfig.RabbitMQConfig != nil && app.RabbitMQ == nil {
		if err = amqpMng.Init(app.BaseConfig.RabbitMQConfig); err != nil {
			return nil, fmt.Errorf("appMng: init rabbitmq failed: %w", err)
		}
	}

	// 项目级构建
	if app.ProjectConfig != nil {
		if err = app.ProjectConfig.Build(app.BaseConfig); err != nil {
			return nil, fmt.Errorf("appMng: project build failed: %w", err)
		}
	}

	if cacheResult {
		if ttl <= 0 {
			ttl = defaultCacheTTL
		}
		m.mu.Lock()
		m.cache.Set(key, app, ttl)
		m.mu.Unlock()
	}
	return app, nil
}

// Invalidate 清除缓存，迫使下次重新加载。
func (m *Manager) Invalidate(id string) {
	if id == "" {
		id = "default"
	}
	m.mu.Lock()
	m.cache.Delete(id)
	m.mu.Unlock()
}

// DefaultManager 返回全局的默认 Manager。
func DefaultManager() *Manager { return defaultManager }
