package aiMng

import (
	"errors"

	"github.com/wiidz/goutil/mngs/psqlMng"
	"github.com/wiidz/goutil/mngs/redisMng"
	"github.com/wiidz/goutil/mngs/volcengineMng"
	"github.com/wiidz/goutil/mngs/volcengineMng/wikiMng"
	"github.com/wiidz/goutil/structs/configStruct"
	"gorm.io/gorm"
)

// Manager 负责承载 AI 对话及产品检索的核心流程，避免聊天服务过于臃肿。
type Manager struct {
	config           *configStruct.AiMngConfig
	aiMng            *volcengineMng.VolcengineMng
	dbMgr            *psqlMng.Manager
	db               *gorm.DB
	wikiMng          *wikiMng.WikiMng
	redis            *redisMng.RedisMng
	productPriceQuery ProductPriceQuery // 产品价格查询接口
}

// ManagerOption 用于在创建 Manager 时注入依赖。
type ManagerOption func(*managerOptions)

type managerOptions struct {
	config           *configStruct.AiMngConfig
	aiMng            *volcengineMng.VolcengineMng
	dbMgr            *psqlMng.Manager
	db               *gorm.DB
	wikiMng          *wikiMng.WikiMng
	redis            *redisMng.RedisMng
	productPriceQuery ProductPriceQuery
}

var (
	errNilConfig  = errors.New("AI 管理器配置未提供")
	errNilAIMng   = errors.New("AI 管理器未提供")
	errNilDB      = errors.New("数据库连接未就绪")
	errNilRedis   = errors.New("Redis 管理器未初始化")
	errNilWikiMng = errors.New("知识库管理器未初始化")
)

// NewManager 创建并返回一个新的 AI 搜索管理器实例。
func NewManager(config *configStruct.AiMngConfig, aiMng *volcengineMng.VolcengineMng, opts ...ManagerOption) (*Manager, error) {
	if config == nil {
		return nil, errNilConfig
	}

	cfg := managerOptions{
		config: config,
		aiMng:  aiMng,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	if cfg.db == nil {
		if cfg.dbMgr != nil {
			cfg.db = cfg.dbMgr.DB()
		}
	}

	if cfg.aiMng == nil {
		return nil, errNilAIMng
	}
	if cfg.db == nil {
		return nil, errNilDB
	}
	if cfg.redis == nil {
		return nil, errNilRedis
	}
	if cfg.wikiMng == nil {
		return nil, errNilWikiMng
	}

	return &Manager{
		config:           cfg.config,
		aiMng:            cfg.aiMng,
		dbMgr:            cfg.dbMgr,
		db:               cfg.db,
		wikiMng:          cfg.wikiMng,
		redis:            cfg.redis,
		productPriceQuery: cfg.productPriceQuery,
	}, nil
}

// WithDB 指定自定义的数据库连接。
func WithDB(db *gorm.DB) ManagerOption {
	return func(opts *managerOptions) {
		if db != nil {
			opts.db = db
		}
	}
}

// WithDBManager 指定数据库管理器，并自动提取 *gorm.DB。
func WithDBManager(dbMgr *psqlMng.Manager) ManagerOption {
	return func(opts *managerOptions) {
		if dbMgr != nil {
			opts.dbMgr = dbMgr
			opts.db = dbMgr.DB()
		}
	}
}

// WithRedis 指定 Redis 管理器。
func WithRedis(redis *redisMng.RedisMng) ManagerOption {
	return func(opts *managerOptions) {
		if redis != nil {
			opts.redis = redis
		}
	}
}

// WithWiki 指定知识库管理器。
func WithWiki(wikiMgr *wikiMng.WikiMng) ManagerOption {
	return func(opts *managerOptions) {
		if wikiMgr != nil {
			opts.wikiMng = wikiMgr
		}
	}
}

// WithProductPriceQuery 指定产品价格查询接口。
func WithProductPriceQuery(query ProductPriceQuery) ManagerOption {
	return func(opts *managerOptions) {
		if query != nil {
			opts.productPriceQuery = query
		}
	}
}

