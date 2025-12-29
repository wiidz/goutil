package appMng

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/wiidz/goutil/mngs/amqpMng"
	"github.com/wiidz/goutil/mngs/esMng"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"github.com/wiidz/goutil/mngs/psqlMng"
	"github.com/wiidz/goutil/mngs/redisMng"
)

const (
	defaultCacheTTL   = 30 * time.Minute
	cacheCleanupCycle = 5 * time.Minute
)

var defaultManager = NewManager()

// NewManager 创建一个新的 Manager。
func NewManager() *Manager {
	return &Manager{cache: cache.New(defaultCacheTTL, cacheCleanupCycle)}
}

// Get 根据 Options 获取（或创建）AppMng。
func (m *Manager) Get(ctx context.Context, opts *Options) (*AppMng, error) {

	//【1】参数验证
	if opts.Loader == nil {
		opts.Loader = DefaultLoader()
	}

	key := opts.ID
	if key == "" {
		key = "default"
	}

	//【2】从缓存中读取
	m.mu.RLock()
	if cached, ok := m.cache.Get(key); ok {
		m.mu.RUnlock()
		return cached.(*AppMng), nil
	}
	m.mu.RUnlock()

	//【3】构建参数
	res, err := opts.Loader.Load(ctx)
	if err != nil {
		return nil, err
	}
	if res == nil || res.BaseConfig == nil {
		return nil, errors.New("appMng: loader returned empty base config")
	}

	//【4】初始化
	app := &AppMng{
		ID:            key,
		BaseConfig:    res.BaseConfig,
		ProjectConfig: opts.ProjectConfig,
		Mysql:         res.Mysql,
		Postgres:      res.Postgres,
		Redis:         res.Redis,
	}

	//【5】启动检查
	// 根据 BaseConfig 中的配置自动初始化服务
	// 注意：配置的有效性验证应该在 Loader 中完成（例如 ConfigBuilder.Build）
	// 如果 BaseConfig 中有配置，则自动初始化对应的服务
	// 这样避免了重复配置：如果 ConfigSourceStrategy 中约定了配置必须加载，则自动初始化服务

	// MySQL
	if app.BaseConfig.MysqlConfig != nil {
		if res.Mysql != nil {
			app.Mysql = res.Mysql
		} else {
			app.Mysql, err = mysqlMng.NewMysqlMng(app.BaseConfig.MysqlConfig, nil)
			if err != nil {
				return nil, fmt.Errorf("appMng: init mysql failed: %w", err)
			}
		}
	}

	// PostgreSQL
	if app.BaseConfig.PostgresConfig != nil {
		if res.Postgres != nil {
			app.Postgres = res.Postgres
		} else {
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
	}

	// Redis
	if app.BaseConfig.RedisConfig != nil {
		if res.Redis != nil {
			app.Redis = res.Redis
		} else {
			if app.Redis, err = redisMng.NewRedisMng(ctx, app.BaseConfig.RedisConfig); err != nil {
				return nil, fmt.Errorf("appMng: init redis failed: %w", err)
			}
		}
	}

	// Elasticsearch
	if app.BaseConfig.EsConfig != nil {
		if res.Es != nil {
			app.Es = res.Es
		} else {
			if err = esMng.Init(app.BaseConfig.EsConfig); err != nil {
				return nil, fmt.Errorf("appMng: init es failed: %w", err)
			}
		}
	}

	// RabbitMQ
	if app.BaseConfig.RabbitMQConfig != nil {
		if res.RabbitMQ != nil {
			app.RabbitMQ = res.RabbitMQ
		} else {
			if err = amqpMng.Init(app.BaseConfig.RabbitMQConfig); err != nil {
				return nil, fmt.Errorf("appMng: init rabbitmq failed: %w", err)
			}
		}
	}

	//【6】项目独特配置
	if app.ProjectConfig != nil {
		if err = app.ProjectConfig.Build(app.BaseConfig); err != nil {
			return nil, fmt.Errorf("appMng: project build failed: %w", err)
		}
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
