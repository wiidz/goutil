package appMng

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/wiidz/goutil/mngs/amqpMng"
	"github.com/wiidz/goutil/mngs/esMng"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"github.com/wiidz/goutil/mngs/psqlMng"
	"github.com/wiidz/goutil/mngs/redisMng"
	"github.com/wiidz/goutil/structs/configStruct"
	"gorm.io/gorm"

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
	ID        uint64 `gorm:"primary_key;column:id;type:int(11);not null" json:"id"` // 主键
	Kind      int8   `gorm:"column:kind;type:tinyint(4)" json:"kind"`               // 类别，1=一般设定，2=页面配置
	Belonging string `gorm:"column:belonging;type:varchar(128)" json:"belonging"`   // 类别
	Name      string `gorm:"column:name;type:varchar(24)" json:"name"`              // 名称

	Flag1  string `gorm:"column:flag_1;type:varchar(128)" json:"flag_1"` // 【属性】补充的一个标识符1
	Flag2  string `gorm:"column:flag_2;type:varchar(128)" json:"flag_2"` // 【属性】补充的一个标识符2
	Value1 string `gorm:"column:value_1;type:text" json:"value_1"`       // 值-1
	Value2 string `gorm:"column:value_2;type:text" json:"value_2"`       // 值-2
	Value3 string `gorm:"column:value_3;type:text" json:"value_3"`       // 值-2

	Tips string `gorm:"column:tips;type:varchar(255)" json:"tips"` // 说明

	CreatedAt time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updated_at"` // 修改时间
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
	Es         *esMng.EsMng
	RabbitMQ   *amqpMng.RabbitMQ
}

// Options 描述了创建/获取 AppMng 时的参数。
type Options struct {
	ID            string
	Loader        Loader
	ProjectConfig configStruct.ProjectConfig
	CacheTTL      time.Duration
}

// AppMng 表示一个应用实例，封装了基础配置以及资源句柄。
type AppMng struct {
	ID string

	BaseConfig    *configStruct.BaseConfig
	ProjectConfig configStruct.ProjectConfig

	Mysql    *mysqlMng.MysqlMng
	Postgres *psqlMng.Manager
	Redis    *redisMng.RedisMng
	Es       *esMng.EsMng
	RabbitMQ *amqpMng.RabbitMQ
}

// Manager 负责缓存和复用 AppMng 实例。
type Manager struct {
	mu    sync.RWMutex
	cache *cache.Cache
}

// ConfigSource 配置来源类型
type ConfigSource string

const (
	SourceYAML     ConfigSource = "yaml"     // 从 YAML 文件加载
	SourceDatabase ConfigSource = "database" // 从数据库加载
)

// 配置键名常量（用于 YAML 和数据库配置）
const (
	// YAML 配置键（用于 UnmarshalKey）
	ConfigKeyProfile     = "profile"           // Profile 配置键
	ConfigKeyLocation    = "location"          // Location 配置键
	ConfigKeyRedis       = "redis"             // Redis 配置键
	ConfigKeyEs          = "es"                // Elasticsearch 配置键
	ConfigKeyRabbitMQ    = "rabbit_mq"         // RabbitMQ 配置键
	ConfigKeyPostgres    = "postgres"          // PostgreSQL 配置键
	ConfigKeyMysql       = "mysql"             // MySQL 配置键
	ConfigKeyWechatMini  = "wechat_mini"       // 微信小程序配置键
	ConfigKeyWechatOa    = "wechat_oa"         // 微信公众号配置键
	ConfigKeyWechatOpen  = "wechat_open"       // 微信开放平台配置键
	ConfigKeyWechatPayV3 = "wechat_pay_v3"     // 微信支付 V3 配置键
	ConfigKeyWechatPayV2 = "wechat_pay_v2"     // 微信支付 V2 配置键
	ConfigKeyAliOss      = "ali_oss"           // 阿里云 OSS 配置键
	ConfigKeyAliPay      = "ali_pay"           // 支付宝配置键
	ConfigKeyAliApi      = "ali_api"           // 阿里云 API 配置键
	ConfigKeyAliSms      = "ali_sms"           // 阿里云短信配置键
	ConfigKeyAliIot      = "ali_iot"           // 阿里云 IoT 配置键
	ConfigKeyAmap        = "amap"              // 高德地图配置键
	ConfigKeyYunxin      = "yunxin"            // 网易云信配置键
	ConfigKeyLocationTZ  = "location.timezone" // Location 时区配置键

	// 数据库配置键（用于 GetValueFromRow 的 name 参数）
	ConfigKeyApp      = "app"       // 应用配置键
	ConfigKeyTimeZone = "time_zone" // 时区配置键
	ConfigKeyWechat   = "wechat"    // 微信配置键
	ConfigKeyAli      = "ali"       // 阿里配置键
	ConfigKeyNetease  = "netease"   // 网易配置键

	// 数据库配置子键（用于 GetValueFromRow 的 flag1 参数）
	ConfigKeyWechatMiniFlag  = "mini"   // 微信小程序子键
	ConfigKeyWechatOaFlag    = "oa"     // 微信公众号子键
	ConfigKeyWechatOpenFlag  = "open"   // 微信开放平台子键
	ConfigKeyWechatPayV3Flag = "pay_v3" // 微信支付 V3 子键
	ConfigKeyWechatPayV2Flag = "pay_v2" // 微信支付 V2 子键
	ConfigKeyAliPayFlag      = "pay"    // 支付宝子键
	ConfigKeyAliOssFlag      = "oss"    // 阿里云 OSS 子键
	ConfigKeyAliApiFlag      = "api"    // 阿里云 API 子键
	ConfigKeyAliSmsFlag      = "sms"    // 阿里云短信子键
	ConfigKeyAliIotFlag      = "iot"    // 阿里云 IoT 子键
	ConfigKeyAliAmapFlag     = "amap"   // 高德地图子键
	ConfigKeyYunxinFlag      = "yunxin" // 网易云信子键
)

const (
	defaultCacheTTL   = 30 * time.Minute
	cacheCleanupCycle = 5 * time.Minute
)

// InitialConfig 初始配置，在应用构建之初传入
// 如果项目涉及数据库连接，必须包含数据库连接信息
type InitialConfig struct {
	// 数据库连接配置（统一的数据库配置，支持 PostgreSQL 和 MySQL）
	DB *configStruct.DBConfig `mapstructure:"db"`

	// 配置表名（从数据库加载配置时使用的表名）
	// 优先级：InitialConfig.SettingTableName > DB.SettingTableName > 默认值 "a_setting"
	SettingTableName string `mapstructure:"setting_table_name"`

	// YAML 配置文件列表（支持多个 YAML 文件）
	YAMLFiles []*configStruct.ViperConfig `mapstructure:"yaml_files"`
}

// ConfigSourceStrategy 配置来源策略，定义每个配置项应该从哪个来源加载
type ConfigSourceStrategy struct {

	// Profile 和 Location 配置来源（通常从数据库或第一个 YAML 文件）
	Profile  ConfigSource `mapstructure:"profile"`  // Profile 配置来源
	Location ConfigSource `mapstructure:"location"` // Location 配置来源

	// 存储相关配置
	Redis    ConfigSource `mapstructure:"redis"`    // Redis 配置来源
	Es       ConfigSource `mapstructure:"es"`       // Elasticsearch 配置来源
	RabbitMQ ConfigSource `mapstructure:"rabbitmq"` // RabbitMQ 配置来源
	Postgres ConfigSource `mapstructure:"postgres"` // PostgreSQL 配置来源
	Mysql    ConfigSource `mapstructure:"mysql"`    // MySQL 配置来源

	// 微信相关配置
	WechatMini  ConfigSource `mapstructure:"wechat_mini"`   // 微信小程序配置来源
	WechatOa    ConfigSource `mapstructure:"wechat_oa"`     // 微信公众号配置来源
	WechatOpen  ConfigSource `mapstructure:"wechat_open"`   // 微信开放平台配置来源
	WechatPayV3 ConfigSource `mapstructure:"wechat_pay_v3"` // 微信支付 V3 配置来源
	WechatPayV2 ConfigSource `mapstructure:"wechat_pay_v2"` // 微信支付 V2 配置来源

	// 阿里相关配置
	AliOss ConfigSource `mapstructure:"ali_oss"` // 阿里云 OSS 配置来源
	AliPay ConfigSource `mapstructure:"ali_pay"` // 支付宝配置来源
	AliApi ConfigSource `mapstructure:"ali_api"` // 阿里云 API 配置来源
	AliSms ConfigSource `mapstructure:"ali_sms"` // 阿里云短信配置来源
	AliIot ConfigSource `mapstructure:"ali_iot"` // 阿里云 IoT 配置来源
	Amap   ConfigSource `mapstructure:"amap"`    // 高德地图配置来源

	// 其他配置
	Yunxin ConfigSource `mapstructure:"yunxin"` // 网易云信配置来源

}

// ConfigBuilder 配置构建器，支持从多个来源加载配置
type ConfigBuilder struct {
	initialConfig *InitialConfig
	strategy      *ConfigSourceStrategy
	db            *gorm.DB       // 数据库连接（如果使用数据库）
	yamlVipers    []*viper.Viper // YAML 配置实例列表
}

// errConfigFromDatabaseEmpty 生成从数据库加载配置失败的错误信息（数据库配置行为空）
func errConfigFromDatabaseEmpty(configName string) string {
	return fmt.Sprintf("策略要求从数据库加载%s配置，但数据库配置行为空", configName)
}

// errConfigFromYAMLNotInit 生成从 YAML 加载配置失败的错误信息（YAML 配置未初始化）
func errConfigFromYAMLNotInit(configName string) string {
	return fmt.Sprintf("策略要求从 YAML 加载%s配置，但 YAML 配置未初始化", configName)
}

// errConfigLoadFailed 生成从 YAML 加载配置失败的错误信息（YAML 配置未初始化）
func errConfigLoadFailed(configName string, err error) string {
	return fmt.Sprintf("加载 %s 配置失败: %w", configName, err)
}

// configLoadSuccess 加载成功
func configLoadSuccess(configName string) string {
	return fmt.Sprintf("加载 %s 配置成功", configName)
}
