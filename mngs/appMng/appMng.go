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

const (
	defaultCacheTTL   = 30 * time.Minute
	cacheCleanupCycle = 5 * time.Minute
)

// Loader 定义了如何装载应用配置的接口。
// 返回 Result，其中至少要包含 BaseConfig。
type Loader interface {
	Load(ctx context.Context) (*Result, error)
}

// LoaderFunc 便于使用函数式实现 Loader。
type LoaderFunc func(ctx context.Context) (*Result, error)

// Load 实现 Loader 接口。
func (f LoaderFunc) Load(ctx context.Context) (*Result, error) { return f(ctx) }

// Result 是 Loader 返回的数据结构。
type Result struct {
	BaseConfig *configStruct.BaseConfig
	Mysql      *mysqlMng.MysqlMng
	Postgres   *psqlMng.Manager
}

// Options 描述了创建/获取 AppMng 时的参数。
type Options struct {
	ID            string
	Loader        Loader
	ProjectConfig configStruct.ProjectConfig
	CheckStart    *configStruct.CheckStart
	CacheTTL      time.Duration
}

// AppMng 表示一个应用实例，封装了基础配置以及资源句柄。
type AppMng struct {
	ID            string
	BaseConfig    *configStruct.BaseConfig
	ProjectConfig configStruct.ProjectConfig

	Mysql    *mysqlMng.MysqlMng
	Postgres *psqlMng.Manager
}

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

// Get 根据 Options 获取（或创建）AppMng。
func (m *Manager) Get(ctx context.Context, opts Options) (*AppMng, error) {
	if opts.Loader == nil {
		return nil, errors.New("appMng: loader is required")
	}
	if opts.ProjectConfig == nil {
		return nil, errors.New("appMng: project config is required")
	}
	if opts.CheckStart == nil {
		opts.CheckStart = &configStruct.CheckStart{}
	}

	key := opts.ID
	if key == "" {
		key = "default"
	}

	m.mu.RLock()
	if cached, ok := m.cache.Get(key); ok {
		m.mu.RUnlock()
		return cached.(*AppMng), nil
	}
	m.mu.RUnlock()

	res, err := opts.Loader.Load(ctx)
	if err != nil {
		return nil, err
	}
	if res == nil || res.BaseConfig == nil {
		return nil, errors.New("appMng: loader returned empty base config")
	}

	app := &AppMng{
		ID:            key,
		BaseConfig:    res.BaseConfig,
		ProjectConfig: opts.ProjectConfig,
		Mysql:         res.Mysql,
		Postgres:      res.Postgres,
	}

	if opts.CheckStart.Mysql && app.Mysql == nil && app.BaseConfig.MysqlConfig != nil {
		app.Mysql, err = mysqlMng.NewMysqlMng(app.BaseConfig.MysqlConfig, nil)
		if err != nil {
			return nil, fmt.Errorf("appMng: init mysql failed: %w", err)
		}
	}

	if opts.CheckStart.Postgres && app.Postgres == nil && app.BaseConfig.PostgresConfig != nil {
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

	if opts.CheckStart.Redis && app.BaseConfig.RedisConfig != nil {
		if err := redisMng.Init(ctx, app.BaseConfig.RedisConfig); err != nil {
			return nil, fmt.Errorf("appMng: init redis failed: %w", err)
		}
	}

	if opts.CheckStart.Es && app.BaseConfig.EsConfig != nil {
		if err := esMng.Init(app.BaseConfig.EsConfig); err != nil {
			return nil, fmt.Errorf("appMng: init es failed: %w", err)
		}
	}

	if opts.CheckStart.RabbitMQ && app.BaseConfig.RabbitMQConfig != nil {
		if err := amqpMng.Init(app.BaseConfig.RabbitMQConfig); err != nil {
			return nil, fmt.Errorf("appMng: init rabbitmq failed: %w", err)
		}
	}

	if err := app.ProjectConfig.Build(app.BaseConfig); err != nil {
		return nil, fmt.Errorf("appMng: project build failed: %w", err)
	}

	ttl := opts.CacheTTL
	if ttl <= 0 {
		ttl = defaultCacheTTL
	}

	m.mu.Lock()
	m.cache.Set(key, app, ttl)
	m.mu.Unlock()

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

// GetSingletonAppMng 兼容旧接口，通过 MySQL 表装载配置。
func GetSingletonAppMng(appID uint64, mysqlConfig *configStruct.MysqlConfig, projectConfig configStruct.ProjectConfig, checkStart *configStruct.CheckStart) (*AppMng, error) {
	if mysqlConfig == nil {
		return nil, errors.New("appMng: mysql config is required")
	}

	loader := LoaderFunc(func(ctx context.Context) (*Result, error) {
		mysqlMngInst, err := mysqlMng.NewMysqlMng(mysqlConfig, nil)
		if err != nil {
			return nil, err
		}

		conn := mysqlMngInst.GetConn()
		var rows []*DbSettingRow
		if err := conn.WithContext(ctx).Table(mysqlConfig.SettingTableName).Where("belonging = ?", "system").Find(&rows).Error; err != nil {
			return nil, err
		}

		baseConfig := buildBaseConfig(rows)
		baseConfig.MysqlConfig = mysqlConfig

		return &Result{
			BaseConfig: baseConfig,
			Mysql:      mysqlMngInst,
		}, nil
	})

	opts := Options{
		ID:            fmt.Sprintf("%d", appID),
		Loader:        loader,
		ProjectConfig: projectConfig,
		CheckStart:    checkStart,
	}

	return defaultManager.Get(context.Background(), opts)
}
