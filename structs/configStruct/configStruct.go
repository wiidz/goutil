package configStruct

import "time"

// AppProfile App资料
type AppProfile struct {
	Host    string // host地址
	Debug   bool   // 是否调试模式
	No      string
	Name    string
	Version string
}

// CheckStart 约定启动时要验证的项目
type CheckStart struct {
	Mysql    bool //
	Redis    bool
	Es       bool
	RabbitMQ bool
}

// BaseConfig 参数
type BaseConfig struct {
	Profile *AppProfile

	Location *time.Location `gorm:"-" json:"-"` // 时区

	MysqlConfig    *MysqlConfig    // 数据库设定
	RedisConfig    *RedisConfig    // redis设定
	OssConfig      *OssConfig      // oss对象存储设定
	EsConfig       *EsConfig       // es设定
	RabbitMQConfig *RabbitMQConfig // es设定

	WechatMiniConfig *WechatMiniConfig // 小程序设定
	WechatOaConfig   *WechatOaConfig   // 公众号设定
	WechatOpenConfig *WechatOpenConfig // 开放平台设定
	WechatPayConfig  *WechatPayConfig  // 微信支付设定
	AliPayConfig     *AliPayConfig     // 支付宝设定

	AliApiConfig *AliApiConfig // 阿里云APi市场设定
	AliSmsConfig *AliSmsConfig // 阿里云APi市场设定
}

// MysqlConfig mysql数据库参数
type MysqlConfig struct {
	Host             string `gorm:"column:db_host" json:"host"`
	Port             string `gorm:"column:db_host" json:"port"`
	Username         string `gorm:"column:db_account" json:"username"`
	Password         string `gorm:"column:db_password" json:"password"`
	DbName           string `gorm:"column:db_name" json:"db_name"`
	Charset          string
	Collation        string
	MaxOpenConns     int    `json:"max_open_conns"`                                      // 默认10
	MaxIdle          int    `json:"max_idle"`                                            // 默认5
	MaxLifeTime      int    `json:"max_life_time"`                                       // 最长生命周期（秒） 默认60
	SettingTableName string `gorm:"column:setting_table_name" json:"setting_table_name"` // 设置表的表名
	TimeZone         string // 时区
	ParseTime        bool
}

// EsConfig elastic search 设置
type EsConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// RedisConfig redis服务器设置
type RedisConfig struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	IdleTimeout int    `json:"idle_timeout"` // 默认60
	Database    int    `json:"datebase"`     // 默认0
	MaxActive   int    `json:"max_active"`   // 默认10
	MaxIdle     int    `json:"max_idle"`     // 默认10
}

// RabbitMQConfig rabbit mq配置
type RabbitMQConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// OssConfig oss参数
type OssConfig struct {
	AccessKeyID     string `gorm:"column:oss_access_key_id;type:varchar(128)" json:"oss_access_key_id"`         // 【OSS】密钥ID
	AccessKeySecret string `gorm:"column:oss_access_key_secret;type:varchar(128)" json:"oss_access_key_secret"` // 【OSS】密钥
	Host            string `gorm:"column:oss_host;type:varchar(128)" json:"oss_host"`                           // 【OSS】域名
	EndPoint        string `gorm:"column:oss_end_point;type:varchar(128)" json:"oss_end_point"`                 // 【OSS】端
	BucketName      string `gorm:"column:oss_bucket_name;type:varchar(128)" json:"oss_bucket_name"`             // 【OSS】bucket名称
	ExpireTime      int64  // 上传策略Policy的失效时间，单位为秒。默认30
}

// WechatMiniConfig 微信小程序参数
type WechatMiniConfig struct {
	AppID     string `gorm:"column:wechat_mini_app_id" json:"wechat_mini_app_id"`
	AppSecret string `gorm:"column:wechat_mini_app_secret" json:"wechat_mini_app_secret"`
}

// WechatOaConfig 微信公众号参数
type WechatOaConfig struct {
	AppID          string `gorm:"column:wechat_oa_app_id" json:"wechat_oa_app_id"`
	AppSecret      string `gorm:"column:wechat_oa_app_secret" json:"wechat_oa_app_secret"`
	Token          string `gorm:"column:token" json:"token"`
	EncodingAESKey string `gorm:"column:encoding_aes_key" json:"encoding_aes_key"`
}

// WechatOpenConfig 微信开放平台参数
type WechatOpenConfig struct {
	AppID     string `gorm:"column:wechat_oa_app_id" json:"wechat_oa_app_id"`
	AppSecret string `gorm:"column:wechat_oa_app_secret" json:"wechat_oa_app_secret"`
}

// WechatPayConfig 微信支付参数
type WechatPayConfig struct {
	AppID           string `gorm:"column:wechat_pay_app_id" json:"app_id"`             //【微信支付】appID
	ApiKey          string `gorm:"column:wechat_api_key" json:"api_key"`               //【微信支付】apiKey（v2）
	ApiKeyV3        string `gorm:"column:wechat_api_key_v3" json:"api_key_v3"`         //【微信支付】apiKey,apiV3Key（v3）
	MchID           string `gorm:"column:wechat_pay_mch_id" json:"mch_id"`             //【微信支付】商户ID 或者服务商模式的 sp_mchid
	CertURI         string `gorm:"column:wechat_pay_cert_uri" json:"cert_uri"`         //【微信支付】公钥文件
	KeyURI          string `gorm:"column:wechat_pay_key_uri" json:"key_uri"`           //【微信支付】私钥文件
	CertSerialNo    string `gorm:"column:cert_serial_mo" json:"cert_serial_mo"`        //【微信支付】证书序列号（V3使用）
	CertContent     string `gorm:"column:wechat_pay_cert_content" json:"cert_content"` //【微信支付】私钥文件内容（私钥 apiclient_key.pem 读取后的字符串内容）
	NotifyURL       string `gorm:"column:notify_url" json:"notify_url"`                // 【微信支付】支付回调地址
	RefundNotifyURL string `gorm:"column:refund_notify_url" json:"refund_notify_url"`  // 【微信支付】退款回调地址
	Debug           bool   `gorm:"column:debug" json:"debug"`                          // 【微信支付】是否是调试模式
}

// AliPayConfig 支付宝参数
type AliPayConfig struct {
	AppID      string `gorm:"column:alipay_app_id" json:"alipay_app_id"`           //【支付宝】appID
	PrivateKey string `gorm:"column:alipay_private_key" json:"alipay_private_key"` //【支付宝】密钥（PKCS1）
	NotifyURL  string // 【支付宝】回调地址
	Debug      bool   // 【支付宝】是否是调试模式
}

type ProjectConfig interface {
	Build() error // 构建参数
}

// AliApiConfig 阿里云市场提供的服务的基本配置
type AliApiConfig struct {
	AppCode   string // 一般有这个就够用了
	AppKey    string
	AppSecret string
}

// AliSmsConfig 阿里云短信服务的配置
type AliSmsConfig struct {
	AccessKeyID     string `gorm:"column:oss_access_key_id;type:varchar(128)" json:"oss_access_key_id"`         // 【OSS】密钥ID
	AccessKeySecret string `gorm:"column:oss_access_key_secret;type:varchar(128)" json:"oss_access_key_secret"` // 【OSS】密钥
}
