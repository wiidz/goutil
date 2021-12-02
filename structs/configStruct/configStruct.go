package configStruct

import "time"

// AppProfile App资料
type AppProfile struct {
	Host       string // host地址
	Debug      bool   // 是否调试模式
	No      string
	Name    string
	Version string
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
	SettingTableName string `gorm:"column:setting_table_name" json:"setting_table_name"` // 设置表的表名
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
	IP          string `json:"ip"`
	Port        string `json:"port"`
	Password    string `json:"password"`
	IdleTimeout int    `json:"idle_timeout"`
	Database    int    `json:"datebase"`
	MaxActive   int    `json:"max_active"`
	MaxIdle     int    `json:"max_idle"`
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
}

// WechatMiniConfig 微信小程序参数
type WechatMiniConfig struct {
	AppID     string `gorm:"column:wechat_mini_app_id" json:"wechat_mini_app_id"`
	AppSecret string `gorm:"column:wechat_mini_app_secret" json:"wechat_mini_app_secret"`
}

// WechatOaConfig 微信公众号参数
type WechatOaConfig struct {
	AppID     string `gorm:"column:wechat_oa_app_id" json:"wechat_oa_app_id"`
	AppSecret string `gorm:"column:wechat_oa_app_secret" json:"wechat_oa_app_secret"`
}

// WechatOpenConfig 微信开放平台参数
type WechatOpenConfig struct {
	AppID     string `gorm:"column:wechat_oa_app_id" json:"wechat_oa_app_id"`
	AppSecret string `gorm:"column:wechat_oa_app_secret" json:"wechat_oa_app_secret"`
}

// WechatPayConfig 微信支付参数
type WechatPayConfig struct {
	AppID       string `gorm:"column:wechat_pay_app_id" json:"wechat_pay_app_id"`             //【微信支付】appID
	Secret      string `gorm:"column:wechat_pay_secret" json:"wechat_pay_secret"`             //【微信支付】密钥
	MchID       string `gorm:"column:wechat_pay_mch_id" json:"wechat_pay_mch_id"`             //【微信支付】商户号
	CertURI     string `gorm:"column:wechat_pay_cert_uri" json:"wechat_pay_cert_uri"`         //【微信支付】公钥文件
	KeyURI      string `gorm:"column:wechat_pay_key_uri" json:"wechat_pay_key_uri"`           //【微信支付】私钥文件
	CertContent string `gorm:"column:wechat_pay_cert_content" json:"wechat_pay_cert_content"` //【微信支付】私钥文件内容
}

// AliPayConfig 支付宝参数
type AliPayConfig struct {
	AppID      string `gorm:"column:alipay_app_id" json:"alipay_app_id"`           //【支付宝】appID
	PrivateKey string `gorm:"column:alipay_private_key" json:"alipay_private_key"` //【支付宝】密钥（PKCS1）
}

type ProjectConfig interface {
	Build() // 构建参数
}
