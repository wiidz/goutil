package configMng

import (
	"github.com/wiidz/goutil/helpers/osHelper"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	ConfigPath              = "/configs/" // 路径写死
	DevPath                 = "dev/"
	PdcPath                 = "pdc/"
	AppConfigFileName       = "app.json"
	WechatConfigFileName    = "wechat.json"
	WechatPayConfigFileName = "wechatPay.json"
	AliPayConfigFileName    = "alipay.json"
	MysqlConfigFileName     = "mysql.json"
	EsConfigFileName        = "es.json"
	RedisConfigFileName     = "redis.json"
	OssConfigFileName       = "oss.json"
	RabbitMQConfigFileName  = "rabbit-mq.json"
)

var appConfig AppConfig

type ConfigMng struct{}


func init() {
	configPath := getAPPRootPath() + ConfigPath
	buf := osHelper.GetFileBuf(configPath + AppConfigFileName)
	_ = appConfig.UnmarshalJSON(buf)
	log.Println("【app-config】", appConfig)
}

// getAPPRootPath 获取项目运行的根目录
func getAPPRootPath() string {

	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}

	p, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	return filepath.Dir(p)
}

// 获取指定目录里的配置文件
func getTargetDir() string {
	var dir string
	if appConfig.Debug == false {
		dir = PdcPath
	} else {
		dir = DevPath
	}
	return getAPPRootPath() + ConfigPath + dir
}

// getFileBuf 根据文件名获取buf
func getFileBuf(fileName string) []byte {
	return osHelper.GetFileBuf(getTargetDir() + fileName)
}

// GetAppConfig 获取本项目配置
func GetAppConfig() *AppConfig {
	return &appConfig
}

// GetMysql 获取mysql数据库配置
func GetMysql() *MysqlConfig {
	buf := getFileBuf(MysqlConfigFileName)
	var mysqlConfig MysqlConfig
	_ = mysqlConfig.UnmarshalJSON(buf)
	log.Println("【mysql-config】", mysqlConfig)
	return &mysqlConfig
}

// GetRedis 获取redis服务器配置
func GetRedis() *RedisConfig {
	buf := getFileBuf(RedisConfigFileName)
	var redisConfig RedisConfig
	_ = redisConfig.UnmarshalJSON(buf)
	log.Println("【redis-config】", redisConfig)
	return &redisConfig
}

// GetWechatMini 获取微信小程序配置
func GetWechatMini() *WechatMiniConfig {
	buf := getFileBuf(AliPayConfigFileName)
	var miniConfig WechatMiniConfig
	_ = miniConfig.UnmarshalJSON(buf)
	log.Println("【wechat-mini-config】", miniConfig)
	return &miniConfig
}

// GetWechatOa 获取微信公众号配置
func GetWechatOa() *WechatOaConfig {
	buf := getFileBuf(AliPayConfigFileName)
	var oaConfig WechatOaConfig
	_ = oaConfig.UnmarshalJSON(buf)
	log.Println("【wechat-oa-config】", oaConfig)
	return &oaConfig
}

// GetWechatOpen 获取微信开放平台
func GetWechatOpen() *WechatOpenConfig {
	buf := getFileBuf(AliPayConfigFileName)
	var openConfig WechatOpenConfig
	_ = openConfig.UnmarshalJSON(buf)
	log.Println("【wechat-open-config】", openConfig)
	return &openConfig
}


// GetWechatPay 获取微信支付配置
func GetWechatPay() *WechatPayConfig {
	buf := getFileBuf(WechatPayConfigFileName)
	var wechatPayConfig WechatPayConfig
	_ = wechatPayConfig.UnmarshalJSON(buf)
	log.Println("【wechat-pay-config】", wechatPayConfig)
	return &wechatPayConfig
}

// GetAliPay 获取支付宝配置
func GetAliPay() *AliPayConfig {
	buf := getFileBuf(WechatConfigFileName)
	var aliConfig AliPayConfig
	_ = aliConfig.UnmarshalJSON(buf)
	log.Println("【ali-config】", aliConfig)
	return &aliConfig
}

// GetOss 获取阿里云对象存储配置
func GetOss() *OssConfig {
	buf := getFileBuf(OssConfigFileName)
	var ossConfig OssConfig
	_ = ossConfig.UnmarshalJSON(buf)
	log.Println("【oss-config】", ossConfig)
	return &ossConfig
}

// GetEs 获取elastic search配置
func GetEs() *EsConfig {
	buf := getFileBuf(EsConfigFileName)
	var esConfig EsConfig
	_ = esConfig.UnmarshalJSON(buf)
	log.Println("【es-config】", esConfig)
	return &esConfig
}

// GetRabbitMQ 获取rabbit mq的配置
func GetRabbitMQ() *RabbitMQConfig {
	buf := getFileBuf(RabbitMQConfigFileName)
	var rabbitMQConfig RabbitMQConfig
	_ = rabbitMQConfig.UnmarshalJSON(buf)
	log.Println("【rabbitMQ-config】", rabbitMQConfig)
	return &rabbitMQConfig
}


//// GetWechatMiniFromDB 从数据库获取微信小程序配置
//func GetWechatMiniFromDB(conn *gorm.DB) *WechatMiniConfig {
//
//	var miniConfig WechatMiniConfig
//	return &miniConfig
//}