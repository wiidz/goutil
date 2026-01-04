package appMng

import (
	"time"

	"github.com/wiidz/goutil/mngs/amqpMng"
	"github.com/wiidz/goutil/mngs/esMng"
	"github.com/wiidz/goutil/mngs/identityMng"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"github.com/wiidz/goutil/mngs/psqlMng"
	"github.com/wiidz/goutil/mngs/redisMng"
	"github.com/wiidz/goutil/structs/configStruct"
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

// AppMng 表示一个应用实例，封装了基础配置以及资源句柄。
type AppMng struct {
	ID string

	BaseConfig    *configStruct.BaseConfig
	ProjectConfig ProjectConfig

	Repos struct {
		Mysql    *mysqlMng.MysqlMng
		Postgres *psqlMng.Manager
		Redis    *redisMng.RedisMng
		Es       *esMng.EsMng
		RabbitMQ *amqpMng.RabbitMQ
	}

	IdMng *identityMng.IdentityMng
}

// ConfigSource 配置来源类型
type ConfigSource string

const (
	SourceYAML     ConfigSource = "yaml"     // 从 YAML 文件加载
	SourceDatabase ConfigSource = "database" // 从数据库加载
)

// ConfigKey 统一描述配置键与中文名
type ConfigKey struct {
	Key     string // 配置键（用于数据库/YAML，首字母大写，如 "Redis"）
	CnLabel string // 中文标签（用于显示）
}

// NewConfigKey 创建配置键
// key: 配置键（首字母大写，如 "Redis"），直接作为字段名使用
// cnLabel: 中文标签（如 "Redis"）
func NewConfigKey(key, cnLabel string) ConfigKey {
	return ConfigKey{
		Key:     key,
		CnLabel: cnLabel,
	}
}

// ConfigKeyRegistry 配置键注册表（使用泛型管理，便于扩展）
type ConfigKeyRegistry struct {
	keys map[string]ConfigKey
}

// NewConfigKeyRegistry 创建新的配置键注册表
func NewConfigKeyRegistry() *ConfigKeyRegistry {
	return &ConfigKeyRegistry{
		keys: make(map[string]ConfigKey),
	}
}

// Register 注册配置键
func (r *ConfigKeyRegistry) Register(name string, key ConfigKey) {
	r.keys[name] = key
}

// Get 获取配置键
func (r *ConfigKeyRegistry) Get(name string) (ConfigKey, bool) {
	key, ok := r.keys[name]
	return key, ok
}

// All 获取所有配置键
func (r *ConfigKeyRegistry) All() map[string]ConfigKey {
	result := make(map[string]ConfigKey, len(r.keys))
	for k, v := range r.keys {
		result[k] = v
	}
	return result
}

// ConfigKeys 将键名与中文名集中定义，主键与子键统一入口
// 使用泛型辅助函数简化创建，保持向后兼容
var ConfigKeys = struct {
	// YAML/主配置键
	Profile     ConfigKey
	Location    ConfigKey
	Redis       ConfigKey
	Es          ConfigKey
	Rabbitmq    ConfigKey
	Postgres    ConfigKey
	Mysql       ConfigKey
	WechatMini  ConfigKey
	WechatOa    ConfigKey
	WechatOpen  ConfigKey
	WechatPayV3 ConfigKey
	WechatPayV2 ConfigKey
	AliOss      ConfigKey
	AliPay      ConfigKey
	AliApi      ConfigKey
	AliSms      ConfigKey
	AliIot      ConfigKey
	Amap        ConfigKey
	Yunxin      ConfigKey
	Volcengine  ConfigKey
	LocationTZ  ConfigKey
	App         ConfigKey
	HttpServer  ConfigKey
	TimeZone    ConfigKey
	Wechat      ConfigKey
	Ali         ConfigKey
	Netease     ConfigKey

	// 子键
	WechatMiniFlag  ConfigKey
	WechatOaFlag    ConfigKey
	WechatOpenFlag  ConfigKey
	WechatPayV3Flag ConfigKey
	WechatPayV2Flag ConfigKey
	AliPayFlag      ConfigKey
	AliOssFlag      ConfigKey
	AliApiFlag      ConfigKey
	AliSmsFlag      ConfigKey
	AliIotFlag      ConfigKey
	AliAmapFlag     ConfigKey
	YunxinFlag      ConfigKey
}{
	// 使用泛型辅助函数简化创建
	Profile:     NewConfigKey("Profile", "基础配置"),
	Location:    NewConfigKey("Location", "时区/位置"),
	Redis:       NewConfigKey("Redis", "Redis"),
	Es:          NewConfigKey("Es", "Elasticsearch"),
	Rabbitmq:    NewConfigKey("Rabbitmq", "RabbitMQ"),
	Postgres:    NewConfigKey("Postgres", "PostgreSQL"),
	Mysql:       NewConfigKey("Mysql", "MySQL"),
	WechatMini:  NewConfigKey("WechatMini", "微信小程序"),
	WechatOa:    NewConfigKey("WechatOa", "微信公众号"),
	WechatOpen:  NewConfigKey("WechatOpen", "微信开放平台"),
	WechatPayV3: NewConfigKey("WechatPayV3", "微信支付 V3"),
	WechatPayV2: NewConfigKey("WechatPayV2", "微信支付 V2"),
	AliOss:      NewConfigKey("AliOss", "阿里云 OSS"),
	AliPay:      NewConfigKey("AliPay", "支付宝"),
	AliApi:      NewConfigKey("AliApi", "阿里云 API"),
	AliSms:      NewConfigKey("AliSms", "阿里云短信"),
	AliIot:      NewConfigKey("AliIot", "阿里云 IoT"),
	Amap:        NewConfigKey("Amap", "高德地图"),
	Yunxin:      NewConfigKey("Yunxin", "网易云信"),
	Volcengine:  NewConfigKey("Volcengine", "火山引擎"),
	LocationTZ:  NewConfigKey("LocationTZ", "时区"),
	App:         NewConfigKey("App", "应用配置"),
	HttpServer:  NewConfigKey("HttpServer", "服务配置"),
	TimeZone:    NewConfigKey("TimeZone", "时区"),
	Wechat:      NewConfigKey("Wechat", "微信配置"),
	Ali:         NewConfigKey("Ali", "阿里云配置"),
	Netease:     NewConfigKey("Netease", "网易配置"),

	WechatMiniFlag:  NewConfigKey("mini", "微信小程序"),
	WechatOaFlag:    NewConfigKey("oa", "微信公众号"),
	WechatOpenFlag:  NewConfigKey("open", "微信开放平台"),
	WechatPayV3Flag: NewConfigKey("pay_v3", "微信支付 V3"),
	WechatPayV2Flag: NewConfigKey("pay_v2", "微信支付 V2"),
	AliPayFlag:      NewConfigKey("pay", "支付宝"),
	AliOssFlag:      NewConfigKey("oss", "阿里云 OSS"),
	AliApiFlag:      NewConfigKey("api", "阿里云 API"),
	AliSmsFlag:      NewConfigKey("sms", "阿里云短信"),
	AliIotFlag:      NewConfigKey("iot", "阿里云 IoT"),
	AliAmapFlag:     NewConfigKey("amap", "高德地图"),
	YunxinFlag:      NewConfigKey("yunxin", "网易云信"),
}

// ConfigSourceStrategy 配置来源策略，定义每个配置项应该从哪个来源加载
type ConfigSourceStrategy struct {

	// Profile 和 Location 配置来源（通常从数据库或第一个 YAML 文件）
	Profile  ConfigSource `mapstructure:"Profile"`  // Profile 配置来源
	Location ConfigSource `mapstructure:"Location"` // Location 配置来源

	HttpServer ConfigSource `mapstructure:"HttpServer"` // HttpServer 配置来源

	// 存储相关配置
	Redis    ConfigSource `mapstructure:"Redis"`    // Redis 配置来源
	Es       ConfigSource `mapstructure:"Es"`       // Elasticsearch 配置来源
	Rabbitmq ConfigSource `mapstructure:"Rabbitmq"` // RabbitMQ 配置来源
	Postgres ConfigSource `mapstructure:"Postgres"` // PostgreSQL 配置来源
	Mysql    ConfigSource `mapstructure:"Mysql"`    // MySQL 配置来源

	// 微信相关配置
	// WechatMini  ConfigSource `mapstructure:"wechat_mini"`   // 微信小程序配置来源
	// WechatOa    ConfigSource `mapstructure:"wechat_oa"`     // 微信公众号配置来源
	// WechatOpen  ConfigSource `mapstructure:"wechat_open"`   // 微信开放平台配置来源
	// WechatPayV3 ConfigSource `mapstructure:"wechat_pay_v3"` // 微信支付 V3 配置来源
	// WechatPayV2 ConfigSource `mapstructure:"wechat_pay_v2"` // 微信支付 V2 配置来源

	// // 阿里相关配置
	// AliOss ConfigSource `mapstructure:"ali_oss"` // 阿里云 OSS 配置来源
	// AliPay ConfigSource `mapstructure:"ali_pay"` // 支付宝配置来源
	// AliApi ConfigSource `mapstructure:"ali_api"` // 阿里云 API 配置来源
	// AliSms ConfigSource `mapstructure:"ali_sms"` // 阿里云短信配置来源
	// AliIot ConfigSource `mapstructure:"ali_iot"` // 阿里云 IoT 配置来源
	// Amap   ConfigSource `mapstructure:"amap"`    // 高德地图配置来源

	// // 其他配置
	// Yunxin     ConfigSource `mapstructure:"yunxin"`     // 网易云信配置来源
	// Volcengine ConfigSource `mapstructure:"volcengine"` // 火山引擎配置来源

	// 自定义配置项策略（用于 ProjectConfig 扩展）
	// key: 配置项名称（如 "my_custom_config"）
	// value: 配置来源（SourceDatabase 或 SourceYAML）
	Custom map[string]ConfigSource `mapstructure:"custom"`
}
