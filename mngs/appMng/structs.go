package appMng

import (
	"time"
)

// BaseConfig 参数
type BaseConfig struct {
	Location *time.Location `gorm:"-" json:"-"` // 时区

	*MysqlConfig    // 数据库设定
	*RedisConfig    // redis设定
	*OssConfig      // oss对象存储设定
	*EsConfig       // es设定
	*RabbitMQConfig // es设定

	*WechatMiniConfig // 小程序设定
	*WechatOaConfig   // 公众号设定
	*WechatOpenConfig // 开放平台设定
	*WechatPayConfig  // 微信支付设定
	*AliPayConfig     // 支付宝设定
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
	Value     string    `gorm:"column:value;type:text" json:"value"`                   // 值
	Value1    string    `gorm:"column:value_1;type:text" json:"value_1"`               // 值-2
	Value2    string    `gorm:"column:value_2;type:text" json:"value_2"`               // 值-1
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
