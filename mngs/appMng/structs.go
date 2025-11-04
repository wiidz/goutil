package appMng

import (
	"sync"
	"time"

	"github.com/wiidz/goutil/mngs/esMng"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"github.com/wiidz/goutil/mngs/psqlMng"
	"github.com/wiidz/goutil/mngs/redisMng"
	"github.com/wiidz/goutil/structs/configStruct"

	"context"
	"github.com/patrickmn/go-cache"
)

/******sql******
CREATE TABLE `u_setting` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `kind` tinyint(4) DEFAULT NULL COMMENT '类别，1=一般设定，2=页面配置',
  `belonging` varchar(128) DEFAULT NULL COMMENT '类别',
  `name` varchar(24) DEFAULT NULL COMMENT '名称',
  `flag_1` varchar(128) DEFAULT NULL COMMENT '【属性】补充的一个标识符1',
  `flag_2` varchar(128) DEFAULT NULL COMMENT '【属性】补充的一个标识符2',
  `value` text COMMENT '值',
  `value_1` text COMMENT '值-2',
  `value_2` text COMMENT '值-1',
  `tips` varchar(255) DEFAULT NULL COMMENT '说明',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '【时间】创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '【时间】最后修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '【时间】删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8
******sql******/
// SettingDbRow [...]
type DbSettingRow struct {
	ID        uint64    `gorm:"primary_key;column:id;type:int(11);not null" json:"id"` // sku编号
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"`    // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updated_at"`    // 修改时间
	Kind      int8      `gorm:"column:kind;type:tinyint(4)" json:"kind"`               // 类别，1=一般设定，2=页面配置
	Belonging string    `gorm:"column:belonging;type:varchar(128)" json:"belonging"`   // 类别
	Name      string    `gorm:"column:name;type:varchar(24)" json:"name"`              // 名称
	Flag1     string    `gorm:"column:flag_1;type:varchar(128)" json:"flag_1"`         // 【属性】补充的一个标识符1
	Flag2     string    `gorm:"column:flag_2;type:varchar(128)" json:"flag_2"`         // 【属性】补充的一个标识符2
	Value1    string    `gorm:"column:value_1;type:text" json:"value_1"`               // 值-1
	Value2    string    `gorm:"column:value_2;type:text" json:"value_2"`               // 值-2
	Value3    string    `gorm:"column:value_3;type:text" json:"value_3"`               // 值-2
	Tips      string    `gorm:"column:tips;type:varchar(255)" json:"tips"`             // 说明
}

// SettingPage 页面设置（带json decode）
type SettingPage struct {
	Kind        int8        `gorm:"column:kind;type:tinyint(4)" json:"-"`        // 类别，1=一般设定，2=页面配置
	Belonging   string      `gorm:"column:belonging;type:varchar(128)" json:"-"` // 类别
	Name        string      `gorm:"column:name;type:varchar(24)" json:"name"`    // 名称
	Value       string      `gorm:"column:value;type:text" json:"-"`             // 值
	ValueParsed interface{} `gorm:"-" json:"value"`                              // 值
}

// Loader 定义了如何装载应用配置的接口。
// 这个是比较复杂逻辑时使用的，为了拓展功能使用
// 返回 Result，其中至少要包含 BaseConfig。
type Loader interface {
	Load(ctx context.Context) (*LoaderResult, error)
	// ... 具体结构体中增加其他的验证、远程拉取等方法
}

// LoaderFunc 便于使用函数式实现 Loader。
type LoaderFunc func(ctx context.Context) (*LoaderResult, error)

// Load 实现 Loader 接口。
func (f LoaderFunc) Load(ctx context.Context) (*LoaderResult, error) { return f(ctx) }

// LoaderResult 是 Loader 返回的数据结构。
type LoaderResult struct {
	BaseConfig *configStruct.BaseConfig
	Mysql      *mysqlMng.MysqlMng
	Postgres   *psqlMng.Manager
	Redis      *redisMng.RedisMng
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
	ID string

	HttpConfig    *configStruct.HttpConfig
	BaseConfig    *configStruct.BaseConfig
	ProjectConfig configStruct.ProjectConfig

	Mysql    *mysqlMng.MysqlMng
	Postgres *psqlMng.Manager
	Redis    *redisMng.RedisMng
	Es       *esMng.EsMng
}

// Manager 负责缓存和复用 AppMng 实例。
type Manager struct {
	mu    sync.RWMutex
	cache *cache.Cache
}
