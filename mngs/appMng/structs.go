package appMng

import (
	"fmt"
	"time"

	"github.com/wiidz/goutil/mngs/amqpMng"
	"github.com/wiidz/goutil/mngs/esMng"
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
}

// ConfigSource 配置来源类型
type ConfigSource string

const (
	SourceYAML     ConfigSource = "yaml"     // 从 YAML 文件加载
	SourceDatabase ConfigSource = "database" // 从数据库加载
)

// ConfigKey 统一描述配置键与中文名
type ConfigKey struct {
	Key     string
	CnLabel string
}

// ConfigKeys 将键名与中文名集中定义，主键与子键统一入口
var ConfigKeys = struct {
	// YAML/主配置键
	Profile     ConfigKey
	Location    ConfigKey
	Redis       ConfigKey
	Es          ConfigKey
	RabbitMQ    ConfigKey
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
	Profile:     ConfigKey{Key: "profile", CnLabel: "基础配置"},
	Location:    ConfigKey{Key: "location", CnLabel: "时区/位置"},
	Redis:       ConfigKey{Key: "redis", CnLabel: "Redis"},
	Es:          ConfigKey{Key: "es", CnLabel: "Elasticsearch"},
	RabbitMQ:    ConfigKey{Key: "rabbit_mq", CnLabel: "RabbitMQ"},
	Postgres:    ConfigKey{Key: "postgres", CnLabel: "PostgreSQL"},
	Mysql:       ConfigKey{Key: "mysql", CnLabel: "MySQL"},
	WechatMini:  ConfigKey{Key: "wechat_mini", CnLabel: "微信小程序"},
	WechatOa:    ConfigKey{Key: "wechat_oa", CnLabel: "微信公众号"},
	WechatOpen:  ConfigKey{Key: "wechat_open", CnLabel: "微信开放平台"},
	WechatPayV3: ConfigKey{Key: "wechat_pay_v3", CnLabel: "微信支付 V3"},
	WechatPayV2: ConfigKey{Key: "wechat_pay_v2", CnLabel: "微信支付 V2"},
	AliOss:      ConfigKey{Key: "ali_oss", CnLabel: "阿里云 OSS"},
	AliPay:      ConfigKey{Key: "ali_pay", CnLabel: "支付宝"},
	AliApi:      ConfigKey{Key: "ali_api", CnLabel: "阿里云 API"},
	AliSms:      ConfigKey{Key: "ali_sms", CnLabel: "阿里云短信"},
	AliIot:      ConfigKey{Key: "ali_iot", CnLabel: "阿里云 IoT"},
	Amap:        ConfigKey{Key: "amap", CnLabel: "高德地图"},
	Yunxin:      ConfigKey{Key: "yunxin", CnLabel: "网易云信"},
	Volcengine:  ConfigKey{Key: "volcengine", CnLabel: "火山引擎"},
	LocationTZ:  ConfigKey{Key: "location.timezone", CnLabel: "时区"},
	App:         ConfigKey{Key: "app", CnLabel: "应用配置"},
	HttpServer:  ConfigKey{Key: "http_server", CnLabel: "服务配置"},
	TimeZone:    ConfigKey{Key: "time_zone", CnLabel: "时区"},
	Wechat:      ConfigKey{Key: "wechat", CnLabel: "微信配置"},
	Ali:         ConfigKey{Key: "ali", CnLabel: "阿里云配置"},
	Netease:     ConfigKey{Key: "netease", CnLabel: "网易配置"},

	WechatMiniFlag:  ConfigKey{Key: "mini", CnLabel: "微信小程序"},
	WechatOaFlag:    ConfigKey{Key: "oa", CnLabel: "微信公众号"},
	WechatOpenFlag:  ConfigKey{Key: "open", CnLabel: "微信开放平台"},
	WechatPayV3Flag: ConfigKey{Key: "pay_v3", CnLabel: "微信支付 V3"},
	WechatPayV2Flag: ConfigKey{Key: "pay_v2", CnLabel: "微信支付 V2"},
	AliPayFlag:      ConfigKey{Key: "pay", CnLabel: "支付宝"},
	AliOssFlag:      ConfigKey{Key: "oss", CnLabel: "阿里云 OSS"},
	AliApiFlag:      ConfigKey{Key: "api", CnLabel: "阿里云 API"},
	AliSmsFlag:      ConfigKey{Key: "sms", CnLabel: "阿里云短信"},
	AliIotFlag:      ConfigKey{Key: "iot", CnLabel: "阿里云 IoT"},
	AliAmapFlag:     ConfigKey{Key: "amap", CnLabel: "高德地图"},
	YunxinFlag:      ConfigKey{Key: "yunxin", CnLabel: "网易云信"},
}

// ConfigKeys 配置键与中文名的组合，便于直接引用
var ConfigKeyList = []ConfigKey{
	ConfigKeys.Profile,
	ConfigKeys.Location,
	ConfigKeys.LocationTZ,
	ConfigKeys.Redis,
	ConfigKeys.Es,
	ConfigKeys.RabbitMQ,
	ConfigKeys.Postgres,
	ConfigKeys.Mysql,
	ConfigKeys.WechatMini,
	ConfigKeys.WechatOa,
	ConfigKeys.WechatOpen,
	ConfigKeys.WechatPayV3,
	ConfigKeys.WechatPayV2,
	ConfigKeys.AliOss,
	ConfigKeys.AliPay,
	ConfigKeys.AliApi,
	ConfigKeys.AliSms,
	ConfigKeys.AliIot,
	ConfigKeys.Amap,
	ConfigKeys.Yunxin,
	ConfigKeys.Volcengine,
	ConfigKeys.App,
	ConfigKeys.HttpServer,
	ConfigKeys.TimeZone,
	ConfigKeys.Wechat,
	ConfigKeys.Ali,
	ConfigKeys.Netease,
}

// ConfigFlagKeys 子键与中文名的组合
var ConfigFlagKeyList = []ConfigKey{
	ConfigKeys.WechatMiniFlag,
	ConfigKeys.WechatOaFlag,
	ConfigKeys.WechatOpenFlag,
	ConfigKeys.WechatPayV3Flag,
	ConfigKeys.WechatPayV2Flag,
	ConfigKeys.AliPayFlag,
	ConfigKeys.AliOssFlag,
	ConfigKeys.AliApiFlag,
	ConfigKeys.AliSmsFlag,
	ConfigKeys.AliIotFlag,
	ConfigKeys.AliAmapFlag,
	ConfigKeys.YunxinFlag,
}

// ConfigKeyDisplayNames 用于将配置键映射到可读的中文名，便于提示
var ConfigKeyDisplayNames = func() map[string]string {
	m := make(map[string]string, len(ConfigKeyList))
	for _, item := range ConfigKeyList {
		m[item.Key] = item.CnLabel
	}
	return m
}()

// ConfigFlagDisplayNames 用于将子键映射到可读的中文名
var ConfigFlagDisplayNames = func() map[string]string {
	m := make(map[string]string, len(ConfigFlagKeyList))
	for _, item := range ConfigFlagKeyList {
		m[item.Key] = item.CnLabel
	}
	return m
}()

// GetKeyDisplayName 返回配置键的中文名，若不存在则返回原键
func GetKeyDisplayName(key string) string {
	if name, ok := ConfigKeyDisplayNames[key]; ok {
		return name
	}
	return key
}

// GetFlagDisplayName 返回子键的中文名，若不存在则返回原键
func GetFlagDisplayName(flag string) string {
	if name, ok := ConfigFlagDisplayNames[flag]; ok {
		return name
	}
	return flag
}

// InitialConfig 初始配置，在应用构建之初传入
type InitialConfig struct {

	// HttpServer标签列表，用于区分不同的HttpServer，例如(client和console)
	HttpServerLabels []string `mapstructure:"http_server_labels"`
}

// ConfigSourceStrategy 配置来源策略，定义每个配置项应该从哪个来源加载
type ConfigSourceStrategy struct {

	// Profile 和 Location 配置来源（通常从数据库或第一个 YAML 文件）
	Profile  ConfigSource `mapstructure:"profile"`  // Profile 配置来源
	Location ConfigSource `mapstructure:"location"` // Location 配置来源

	HttpServer ConfigSource `mapstructure:"http_server"` // HttpServer 配置来源

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
	Yunxin     ConfigSource `mapstructure:"yunxin"`     // 网易云信配置来源
	Volcengine ConfigSource `mapstructure:"volcengine"` // 火山引擎配置来源

	// 自定义配置项策略（用于 ProjectConfig 扩展）
	// key: 配置项名称（如 "my_custom_config"）
	// value: 配置来源（SourceDatabase 或 SourceYAML）
	Custom map[string]ConfigSource `mapstructure:"custom"`
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
	return fmt.Sprintf("加载 %s 配置失败: %v", configName, err)
}

// configLoadSuccess 加载成功
func configLoadSuccess(configName string) string {
	return fmt.Sprintf("加载 %s 配置成功", configName)
}
