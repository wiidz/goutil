package configStruct

import (
	"time"

	"gorm.io/gorm/logger"
)

// AppProfile App资料
type AppProfile struct {
	No      string `mapstructure:"no"`
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`

	// Domain string `mapstructure:"domain"` // 外部域名地址
	// Host   string `mapstructure:"host"`   // 监听的地址（如0.0.0.0，127.0.0.1等）
	// Port   string `mapstructure:"port"`   // 监听的端口号

	Debug bool `mapstructure:"debug"` // 是否调试模式
}

// 一个AppMng中可能包含多个Server，例如(Client和Console)，每个Server可能对应一个不同的域名和端口
type HttpServerConfig struct {
	Label  string `mapstructure:"label"`  // 标签（例如client，console等）
	Domain string `mapstructure:"domain"` // 外部域名地址
	Host   string `mapstructure:"host"`   // 监听的地址（如0.0.0.0，127.0.0.1等）
	Port   string `mapstructure:"port"`   // 监听的端口号
}

type ViperConfig struct {
	DirPath  string `mapstructure:"dir_path"`  // 例如 ./configs
	FileName string `mapstructure:"file_name"` // 例如 config
	FileType string `mapstructure:"file_type"` // 通常我们用yaml，不需要带点
}
type HttpConfig struct {
	IP   string `mapstructure:"ip"`
	Port string `mapstructure:"port"`
}
type RepoConfig struct {
	DSN         string `mapstructure:"dsn"`
	AutoMigrate string `mapstructure:"auto_migrate"`
}

// BaseConfig 参数
type BaseConfig struct {
	Profile  *AppProfile    `mapstructure:"profile"`
	Location *time.Location `gorm:"-" json:"-" mapstructure:"location"` // 时区

	HttpServerConfig []*HttpServerConfig `mapstructure:"http_server_config"` // http服务器设定

	MysqlConfig    *MysqlConfig    `mapstructure:"mysql_config"`    // 数据库设定
	PostgresConfig *PostgresConfig `mapstructure:"postgres_config"` // Postgres 设定
	RedisConfig    *RedisConfig    `mapstructure:"redis_config"`    // redis设定

	EsConfig       *EsConfig       `mapstructure:"es_config"`       // es设定
	RabbitMQConfig *RabbitMQConfig `mapstructure:"rabbitmq_config"` // es设定

	WechatMiniConfig  *WechatMiniConfig  `mapstructure:"wechat_mini_config"`   // 小程序设定
	WechatOaConfig    *WechatOaConfig    `mapstructure:"wechat_oa_config"`     // 公众号设定
	WechatOpenConfig  *WechatOpenConfig  `mapstructure:"wechat_open_config"`   // 开放平台设定
	WechatPayConfigV3 *WechatPayConfigV3 `mapstructure:"wechat_pay_config_v3"` // V3微信支付设定
	WechatPayConfigV2 *WechatPayConfigV2 `mapstructure:"wechat_pay_config_v2"` // V2微信支付设定
	AliPayConfig      *AliPayConfig      `mapstructure:"ali_pay_config"`       // 支付宝设定

	AliApiConfig *AliApiConfig `mapstructure:"ali_api_config"` // 阿里云APi市场设定
	AliSmsConfig *AliSmsConfig `mapstructure:"ali_sms_config"` // 阿里云短信服务设定
	AliIotConfig *AliIotConfig `mapstructure:"ali_iot_config"` // 阿里云物联网市场设定
	AliOssConfig *AliOssConfig `mapstructure:"ali_oss_config"` // 阿里云oss对象存储设定
	AmapConfig   *AmapConfig   `mapstructure:"amap_config"`    // 高德地图设定

	YunxinConfig *YunxinConfig `mapstructure:"yunxin_config"` // 网易云信设定
	XFYunConfig  *XFYunConfig  `mapstructure:"xfyun_config"`  // 科大讯飞配置
}

// MysqlConfig mysql数据库参数
type MysqlConfig struct {
	Host             string `gorm:"column:db_host" json:"host" mapstructure:"host"`
	Port             string `gorm:"column:db_host" json:"port" mapstructure:"port"`
	Username         string `gorm:"column:db_account" json:"username" mapstructure:"username"`
	Password         string `gorm:"column:db_password" json:"password" mapstructure:"password"`
	DbName           string `gorm:"column:db_name" json:"db_name" mapstructure:"db_name"`
	Charset          string `mapstructure:"charset"`
	Collation        string `mapstructure:"collation"`
	MaxOpenConns     int    `json:"max_open_conns" mapstructure:"max_open_conns"`                                          // 默认10
	MaxIdle          int    `json:"max_idle" mapstructure:"max_idle"`                                                      // 默认5
	MaxLifeTime      int    `json:"max_life_time" mapstructure:"max_life_time"`                                            // 最长生命周期（秒） 默认60
	SettingTableName string `gorm:"column:setting_table_name" json:"setting_table_name" mapstructure:"setting_table_name"` // 设置表的表名
	TimeZone         string `mapstructure:"time_zone"`                                                                     // 时区
	ParseTime        bool   `mapstructure:"parse_time"`

	Logger logger.Interface `mapstructure:"logger"`
}

// EsConfig elastic search 设置
type EsConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     string `json:"port" mapstructure:"port"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

// RedisConfig redis服务器设置
type RedisConfig struct {
	Host        string `json:"host" mapstructure:"host"`
	Port        string `json:"port" mapstructure:"port"`
	Username    string `json:"username" mapstructure:"username"`
	Password    string `json:"password" mapstructure:"password"`
	IdleTimeout int    `json:"idle_timeout" mapstructure:"idle_timeout" default:"60"` // 默认60
	Database    int    `json:"datebase" mapstructure:"datebase" default:"0"`          // 默认0
	MaxActive   int    `json:"max_active" mapstructure:"max_active" default:"10"`     // 默认10
	MaxIdle     int    `json:"max_idle" mapstructure:"max_idle" default:"10"`         // 默认10
}

// DBType 数据库类型
type DBType string

const (
	DBTypePostgres DBType = "postgres" // PostgreSQL
	DBTypeMysql    DBType = "mysql"    // MySQL
)

// DBConfig 统一的数据库配置（支持 PostgreSQL 和 MySQL）
type DBConfig struct {
	// 数据库类型（postgres 或 mysql）
	Type DBType `json:"type" mapstructure:"type"`

	// DSN 连接字符串（优先级最高，如果设置了 DSN，将忽略下面的 Host/Port 等字段）
	// PostgreSQL: "postgres://user:password@host:port/dbname?sslmode=disable"
	// MySQL: "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	DSN string `json:"dsn" mapstructure:"dsn"`

	// 连接参数（当 DSN 为空时使用，主要用于 MySQL）
	Host     string `json:"host" mapstructure:"host"`
	Port     string `json:"port" mapstructure:"port"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	DbName   string `json:"db_name" mapstructure:"db_name"`

	// MySQL 特有参数（当 Type 为 mysql 且 DSN 为空时使用）
	Charset   string `json:"charset" mapstructure:"charset"`       // 默认 utf8mb4
	Collation string `json:"collation" mapstructure:"collation"`   // 默认 utf8mb4_unicode_ci
	TimeZone  string `json:"time_zone" mapstructure:"time_zone"`   // 默认 Asia/Shanghai
	ParseTime bool   `json:"parse_time" mapstructure:"parse_time"` // 默认 true

	// 连接池参数（两种数据库通用）
	ConnMaxIdle      int              `json:"conn_max_idle" mapstructure:"conn_max_idle"`           // 最大空闲连接数
	ConnMaxOpen      int              `json:"conn_max_open" mapstructure:"conn_max_open"`           // 最大打开连接数
	ConnMaxLifetime  time.Duration    `json:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`   // 连接最大生命周期
	SettingTableName string           `json:"setting_table_name" mapstructure:"setting_table_name"` // 设置表的表名
	Logger           logger.Interface `json:"logger" mapstructure:"logger"`
}

// PostgresConfig PostgreSQL 数据库配置（保留用于向后兼容）
// Deprecated: 建议使用 DBConfig 替代
type PostgresConfig struct {
	DSN              string           `json:"dsn" mapstructure:"dsn"`
	ConnMaxIdle      int              `json:"conn_max_idle" mapstructure:"conn_max_idle"`
	ConnMaxOpen      int              `json:"conn_max_open" mapstructure:"conn_max_open"`
	ConnMaxLifetime  time.Duration    `json:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	SettingTableName string           `json:"setting_table_name" mapstructure:"setting_table_name"` // 设置表的表名
	Logger           logger.Interface `json:"logger" mapstructure:"logger"`
}

// RabbitMQ Exchange 类型
type ExchangeType string

const (
	Fanout          ExchangeType = "fanout"
	Direct          ExchangeType = "direct"
	Topic           ExchangeType = "topic"
	XDelayedMessage ExchangeType = "x-delayed-message"
	DeadLetterDelay ExchangeType = "dead_letter_delay"
)

// RabbitMQConfig rabbit mq配置
type RabbitMQConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

// WechatMiniConfig 微信小程序参数
type WechatMiniConfig struct {
	AppID     string `gorm:"column:wechat_mini_app_id" json:"wechat_mini_app_id" mapstructure:"wechat_mini_app_id"`
	AppSecret string `gorm:"column:wechat_mini_app_secret" json:"wechat_mini_app_secret" mapstructure:"wechat_mini_app_secret"`
}

// WechatOaConfig 微信公众号参数
type WechatOaConfig struct {
	AppID          string `gorm:"column:wechat_oa_app_id" json:"wechat_oa_app_id" mapstructure:"wechat_oa_app_id"`
	AppSecret      string `gorm:"column:wechat_oa_app_secret" json:"wechat_oa_app_secret" mapstructure:"wechat_oa_app_secret"`
	Token          string `gorm:"column:token" json:"token" mapstructure:"token"`
	EncodingAESKey string `gorm:"column:encoding_aes_key" json:"encoding_aes_key" mapstructure:"encoding_aes_key"`
}

// WechatOpenConfig 微信开放平台参数
type WechatOpenConfig struct {
	AppID     string `gorm:"column:wechat_oa_app_id" json:"wechat_oa_app_id" mapstructure:"wechat_oa_app_id"`
	AppSecret string `gorm:"column:wechat_oa_app_secret" json:"wechat_oa_app_secret" mapstructure:"wechat_oa_app_secret"`
}

// WechatPayConfigV3 V3微信支付参数
type WechatPayConfigV3 struct {
	AppID                     string `gorm:"column:wechat_pay_app_id" json:"app_id" mapstructure:"app_id"`                                                        //【微信支付】appID
	ApiKeyV3                  string `gorm:"column:wechat_api_key_v3" json:"api_key_v3" mapstructure:"api_key_v3"`                                                //【微信支付】apiKey,apiV3Key（v3）
	MchID                     string `gorm:"column:wechat_pay_mch_id" json:"mch_id" mapstructure:"mch_id"`                                                        //【微信支付】商户ID 或者服务商模式的 sp_mchid
	CertURI                   string `gorm:"column:wechat_pay_cert_uri" json:"cert_uri" mapstructure:"cert_uri"`                                                  //【微信支付】公钥文件
	KeyURI                    string `gorm:"column:wechat_pay_key_uri" json:"key_uri" mapstructure:"key_uri"`                                                     //【微信支付】私钥文件
	CertSerialNo              string `gorm:"column:cert_serial_mo" json:"cert_serial_mo" mapstructure:"cert_serial_mo"`                                           //【微信支付】证书序列号（V3使用）
	NotifyURL                 string `gorm:"column:notify_url" json:"notify_url" mapstructure:"notify_url"`                                                       // 【微信支付】支付回调地址
	RefundNotifyURL           string `gorm:"column:refund_notify_url" json:"refund_notify_url" mapstructure:"refund_notify_url"`                                  // 【微信支付】退款回调地址
	MerchantTransferNotifyURL string `gorm:"column:merchant_transfer_notify_url" json:"merchant_transfer_notify_url" mapstructure:"merchant_transfer_notify_url"` // 【微信支付】商家转账回调地址
	Debug                     bool   `gorm:"column:debug" json:"debug" mapstructure:"debug"`                                                                      // 【微信支付】是否是调试模式
	PEMCertContent            string `gorm:"column:pem_cert_content" json:"pem_cert_content" mapstructure:"pem_cert_content"`                                     //【微信支付】证书pem格式（apiclient_cert.pem） 从apiclient_cert.p12中导出证书部分的文件，为pem格式，请妥善保管不要泄漏和被他人复制 部分开发语言和环境，不能直接使用p12文件，而需要使用pem，所以为了方便您使用，已为您直接提供
	PEMPrivateKeyContent      string `gorm:"column:pem_private_key_content" json:"pem_private_key_content" mapstructure:"pem_private_key_content"`                //【微信支付】证书密钥pem格式（apiclient_key.pem） 从apiclient_cert.p12中导出密钥部分的文件，为pem格式 部分开发语言和环境，不能直接使用p12文件，而需要使用pem，所以为了方便您使用，已为您直接提供
	//PEMPublicKeyContent  string `gorm:"column:pem_public_key_content" json:"pem_public_key_content"`   //【微信支付】证书公钥pem格式(我们手动生成的)；；新：：：：：不用我们去维护公钥！！！
}

// WechatPayConfigV2 V2微信支付参数
type WechatPayConfigV2 struct {
	AppID           string `gorm:"column:wechat_pay_app_id" json:"app_id" mapstructure:"app_id"`                          //【微信支付】appID
	ApiKey          string `gorm:"column:wechat_api_key" json:"api_key" mapstructure:"api_key"`                           //【微信支付】apiKey（v2）
	MchID           string `gorm:"column:wechat_pay_mch_id" json:"mch_id" mapstructure:"mch_id"`                          //【微信支付】商户ID 或者服务商模式的 sp_mchid
	CertURI         string `gorm:"column:wechat_pay_cert_uri" json:"cert_uri" mapstructure:"cert_uri"`                    //【微信支付】公钥文件
	KeyURI          string `gorm:"column:wechat_pay_key_uri" json:"key_uri" mapstructure:"key_uri"`                       //【微信支付】私钥文件
	CertSerialNo    string `gorm:"column:cert_serial_mo" json:"cert_serial_mo" mapstructure:"cert_serial_mo"`             //【微信支付】证书序列号（V3使用）
	NotifyURL       string `gorm:"column:notify_url" json:"notify_url" mapstructure:"notify_url"`                         // 【微信支付】支付回调地址
	RefundNotifyURL string `gorm:"column:refund_notify_url" json:"refund_notify_url" mapstructure:"refund_notify_url"`    // 【微信支付】退款回调地址
	Debug           bool   `gorm:"column:debug" json:"debug" mapstructure:"debug"`                                        // 【微信支付】是否是调试模式
	P12CertFilePath string `gorm:"column:p12_cert_file_path" json:"p12_cert_file_path" mapstructure:"p12_cert_file_path"` // apiclient_cert.p12的路径
	//PEMCertContent       string `gorm:"column:pem_cert_content" json:"pem_cert_content"`               //【微信支付】证书pem格式（apiclient_cert.pem） 从apiclient_cert.p12中导出证书部分的文件，为pem格式，请妥善保管不要泄漏和被他人复制 部分开发语言和环境，不能直接使用p12文件，而需要使用pem，所以为了方便您使用，已为您直接提供
	//PEMPrivateKeyContent string `gorm:"column:pem_private_key_content" json:"pem_private_key_content"` //【微信支付】证书密钥pem格式（apiclient_key.pem） 从apiclient_cert.p12中导出密钥部分的文件，为pem格式 部分开发语言和环境，不能直接使用p12文件，而需要使用pem，所以为了方便您使用，已为您直接提供
	//PEMPublicKeyContent  string `gorm:"column:pem_public_key_content" json:"pem_public_key_content"`   //【微信支付】证书公钥pem格式(我们手动生成的)；；新：：：：：不用我们去维护公钥！！！
}

// AliPayConfig 支付宝参数
type AliPayConfig struct {
	AppID      string `gorm:"column:alipay_app_id" json:"alipay_app_id" mapstructure:"alipay_app_id"`                //【支付宝】appID
	PrivateKey string `gorm:"column:alipay_private_key" json:"alipay_private_key" mapstructure:"alipay_private_key"` //【支付宝】密钥（PKCS1）

	AppCertPublicKey string `gorm:"column:app_cert_public_key" json:"app_cert_public_key" mapstructure:"app_cert_public_key"`
	RootCert         string `gorm:"column:root_cert" json:"root_cert" mapstructure:"root_cert"`
	CertPublicKey    string `gorm:"column:cert_public_key" json:"cert_public_key" mapstructure:"cert_public_key"`

	NotifyURL string `mapstructure:"notify_url"`
	Debug     bool   `mapstructure:"debug"`
}

type ProjectConfig interface {
	Build(baseConfig *BaseConfig) error // 构建参数
}

// AliApiConfig 阿里云市场提供的服务的基本配置
type AliApiConfig struct {
	AppCode   string `mapstructure:"app_code"`
	AppKey    string `mapstructure:"app_key"`
	AppSecret string `mapstructure:"app_secret"`
}

// AliOssConfig oss参数
type AliOssConfig struct {
	AccessKeyID     string `gorm:"column:oss_access_key_id;type:varchar(128)" json:"oss_access_key_id" mapstructure:"oss_access_key_id"`
	AccessKeySecret string `gorm:"column:oss_access_key_secret;type:varchar(128)" json:"oss_access_key_secret" mapstructure:"oss_access_key_secret"`
	Host            string `gorm:"column:oss_host;type:varchar(128)" json:"oss_host" mapstructure:"oss_host"`
	EndPoint        string `gorm:"column:oss_end_point;type:varchar(128)" json:"oss_end_point" mapstructure:"oss_end_point"`
	BucketName      string `gorm:"column:oss_bucket_name;type:varchar(128)" json:"oss_bucket_name" mapstructure:"oss_bucket_name"`
	ARN             string `gorm:"arn" json:"arn" mapstructure:"arn"`
	ExpireTime      int64  `mapstructure:"expire_time"`
}

// AliSmsConfig 阿里云短信服务的配置
type AliSmsConfig struct {
	AccessKeyID     string `gorm:"column:oss_access_key_id;type:varchar(128)" json:"oss_access_key_id" mapstructure:"oss_access_key_id"`
	AccessKeySecret string `gorm:"column:oss_access_key_secret;type:varchar(128)" json:"oss_access_key_secret" mapstructure:"oss_access_key_secret"`
}

// AliRamConfig 阿里云RAM访问控制的账号和密码
type AliRamConfig struct {
	AccessKeyID     string `gorm:"column:oss_access_key_id;type:varchar(128)" json:"oss_access_key_id" mapstructure:"oss_access_key_id"`
	AccessKeySecret string `gorm:"column:oss_access_key_secret;type:varchar(128)" json:"oss_access_key_secret" mapstructure:"oss_access_key_secret"`
}

// AliIotConfig 阿里云物联网的基本配置（每个实例单独放）
// 因为一个项目用的服务器基本上是一个区域，一个账户，所以以下属性是公用的
type AliIotConfig struct {
	AccessKeyID     string `gorm:"column:oss_access_key_id;type:varchar(128)" json:"oss_access_key_id" mapstructure:"oss_access_key_id"`
	AccessKeySecret string `gorm:"column:oss_access_key_secret;type:varchar(128)" json:"oss_access_key_secret" mapstructure:"oss_access_key_secret"`
	EndPoint        string `gorm:"end_point;type:varchar(128)" json:"end_point" mapstructure:"end_point"`
	RegionID        string `gorm:"region_id;type:varchar(128)" json:"region_id" mapstructure:"region_id"`
}

// YunxinConfig 网易云信
type YunxinConfig struct {
	AppKey    string `gorm:"column:app_key;type:varchar(128)" json:"app_key" mapstructure:"app_key"`          // 【云信】密钥
	AppSecret string `gorm:"column:app_secret;type:varchar(128)" json:"app_secret" mapstructure:"app_secret"` // 【云信】密钥
	CCURL     string `gorm:"column:cc_url;type:varchar(128)" json:"cc_url" mapstructure:"cc_url"`             // 信息抄送地址
}

// AmapConfig 高德地图配置
type AmapConfig struct {
	Key string `mapstructure:"key"`
}

type TcpConfig struct {
	IP           string        `json:"ip" mapstructure:"ip"`
	Port         int           `json:"port" mapstructure:"port"`
	ReadTimeOut  time.Duration `mapstructure:"read_timeout"`  // 读取超时时间
	WriteTimeOut time.Duration `mapstructure:"write_timeout"` // 写入超时时间
}

type KookConfig struct {
	Debug bool `mapstructure:"debug"`

	GuildID     string `mapstructure:"guild_id"`
	Token       string `mapstructure:"token"`
	EncryptKey  string `mapstructure:"encrypt_key"`
	VerifyToken string `mapstructure:"verify_token"`
	CallbackURL string `mapstructure:"callback_url"`
	RobotID     string `mapstructure:"robot_id"`

	VerifiedPlayerRoleID uint64 `mapstructure:"verified_player_role_id"` // 认证角色ID
	BanRoleID            uint64 `mapstructure:"ban_role_id"`             // 禁言角色ID
}

// VolcengineConfig 火山引擎
type VolcengineConfig struct {
	Debug bool `mapstructure:"debug"`

	ApiKey string `mapstructure:"api_key"`
}

// XFYunConfig 科大讯飞大模型配置
type XFYunConfig struct {
	Debug bool `mapstructure:"debug"`

	AppID     string `mapstructure:"app_id"`
	ApiKey    string `mapstructure:"api_key"`
	ApiSecret string `mapstructure:"api_secret"`

	Scheme string `mapstructure:"scheme"` // 默认: wss
	Host   string `mapstructure:"host"`   // 默认: iat.xf-yun.com
	Path   string `mapstructure:"path"`   // 默认: /v1
}

// Config 配置信息
type WikiConfig struct {
	AccessKeyID string `mapstructure:"access_key_id"` // Access Key
	SecretKey   string `mapstructure:"secret_key"`    // Secret Key
	ApiKey      string `mapstructure:"api_key"`       // API Key

	StreamTimeout time.Duration `mapstructure:"stream_timeout"` // 流式超时
	SimpleTimeout time.Duration `mapstructure:"simple_timeout"` // 单此请求
}
