package appMng

import (
	"gorm.io/gorm"
	"time"
)

type AppDBRow struct {
	ID        int             `gorm:"column:id"`
	DBName    string          `gorm:"column:db_name"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at"`
}

// AppConfig 参数
type AppConfig struct {
	ID               int            `gorm:"column:id" json:"id"`
	AppNo            string         `gorm:"column:app_no" json:"app_no"`
	AppName          string         `gorm:"column:app_name" json:"app_name"`
	TimeZone         string         `gorm:"column:time_zone" json:"time_zone"` // 时区，默认 Asia/Shanghai
	Location         *time.Location `gorm:"-" json:"-"`
	MysqlConfig                     // 数据库设定
	WechatMiniConfig                // 小程序设定
	WechatOaConfig                  // 公众号设定
	WechatPayConfig                 // 微信支付设定
	OSSConfig                       // oss对象存储设定
	AliPayConfig                    // 支付宝设定
	//PromoConfig         *PromoConfig         `gorm:"-" json:"-"` // 促销推广设定
}

// MysqlConfig mysql数据库参数
type MysqlConfig struct {
	Host     string `gorm:"column:db_host" json:"host"`
	Port     string `gorm:"column:db_host" json:"port"`
	Username string `gorm:"column:db_account" json:"username"`
	Password string `gorm:"column:db_password" json:"password"`
	DbName   string `gorm:"column:db_name" json:"db_name"`
}

// WechatMiniConfig 微信小程序参数
type WechatMiniConfig struct {
	WechatMiniAppID     string `gorm:"column:wechat_mini_app_id" json:"wechat_mini_app_id"`
	WechatMiniAppSecret string `gorm:"column:wechat_mini_app_secret" json:"wechat_mini_app_secret"`
}

// WechatOaConfig 微信公众号参数
type WechatOaConfig struct {
	WechatOaAppID     string `gorm:"column:wechat_oa_app_id" json:"wechat_oa_app_id"`
	WechatOaAppSecret string `gorm:"column:wechat_oa_app_secret" json:"wechat_oa_app_secret"`
}

// WechatPayConfig 微信支付参数
type WechatPayConfig struct {
	WechatPayAppID       string `gorm:"column:wechat_pay_app_id" json:"wechat_pay_app_id"`             //【微信支付】appID
	WechatPaySecret      string `gorm:"column:wechat_pay_secret" json:"wechat_pay_secret"`             //【微信支付】密钥
	WechatPayMchID       string `gorm:"column:wechat_pay_mch_id" json:"wechat_pay_mch_id"`             //【微信支付】商户号
	WechatPayCertURI     string `gorm:"column:wechat_pay_cert_uri" json:"wechat_pay_cert_uri"`         //【微信支付】公钥文件
	WechatPayKeyURI      string `gorm:"column:wechat_pay_key_uri" json:"wechat_pay_key_uri"`           //【微信支付】私钥文件
	WechatPayCertContent string `gorm:"column:wechat_pay_cert_content" json:"wechat_pay_cert_content"` //【微信支付】私钥文件内容
}

// AliPayConfig 支付宝参数
type AliPayConfig struct {
	AliPayAppID   string `gorm:"column:alipay_app_id" json:"alipay_app_id"`           //【支付宝】appID
	AliPrivateKey string `gorm:"column:alipay_private_key" json:"alipay_private_key"` //【支付宝】密钥（PKCS1）
}

// OSSConfig oss参数
type OSSConfig struct {
	OssAccessKeyID     string `gorm:"column:oss_access_key_id;type:varchar(128)" json:"oss_access_key_id"`         // 【OSS】密钥ID
	OssAccessKeySecret string `gorm:"column:oss_access_key_secret;type:varchar(128)" json:"oss_access_key_secret"` // 【OSS】密钥
	OssHost            string `gorm:"column:oss_host;type:varchar(128)" json:"oss_host"`                           // 【OSS】域名
	OssEndPoint        string `gorm:"column:oss_end_point;type:varchar(128)" json:"oss_end_point"`                 // 【OSS】端
	OssBucketName      string `gorm:"column:oss_bucket_name;type:varchar(128)" json:"oss_bucket_name"`             // 【OSS】bucket名称
}
