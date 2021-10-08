package configMng

import (
	"fmt"
	"github.com/wiidz/goutil/helpers/osHelper"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	//ConfigPath              = "./configs/" // 路径写死
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
	RebbitMQConfigFileName  = "rabbit-mq.json"
)

var appConfig = AppConfig{}

type ConfigMng struct{}


func init() {
	configPath := getAPPRootPath() + "/configs/"
	log.Println(configPath + AppConfigFileName)
	buf := osHelper.GetFileBuf(configPath + AppConfigFileName)
	_ = appConfig.UnmarshalJSON(buf)
	log.Println("appConfig", appConfig)
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
	return getAPPRootPath() + "/configs/" + dir
}

// getFileBuf 根据文件名获取buf
func getFileBuf(fileName string) []byte {
	return osHelper.GetFileBuf(getTargetDir() + fileName)
}

// GetHttpPort 获取本项目占用端口配置
func GetHttpPort() string {
	return appConfig.HttpPort
}

// GetMysql 获取mysql数据库配置
func GetMysql() MysqlConfig {
	buf := getFileBuf(MysqlConfigFileName)
	mysqlConfig := MysqlConfig{}
	_ = mysqlConfig.UnmarshalJSON(buf)
	return mysqlConfig
}

// GetRedis 获取redis服务器配置
func GetRedis() RedisConfig {
	buf := getFileBuf(RedisConfigFileName)
	redisConfig := RedisConfig{}
	_ = redisConfig.UnmarshalJSON(buf)
	return redisConfig
}

// GetWechat 获取微信配置
func GetWechat() WechatConfig {
	buf := getFileBuf(AliPayConfigFileName)
	wechatConfig := WechatConfig{}
	_ = wechatConfig.UnmarshalJSON(buf)
	return wechatConfig
}

// GetWechatPay 获取微信支付配置
func GetWechatPay() *WechatPayConfig {
	buf := getFileBuf(WechatPayConfigFileName)
	wechatConfig := WechatPayConfig{}
	_ = wechatConfig.UnmarshalJSON(buf)
	return &wechatConfig
}

// GetAliPay 获取支付宝配置
func GetAliPay() *AliPayConfig {
	buf := getFileBuf(WechatConfigFileName)
	wechatConfig := AliPayConfig{}
	_ = wechatConfig.UnmarshalJSON(buf)
	return &wechatConfig
}

// GetOss 获取阿里云对象存储配置
func GetOss() OssConfig {
	buf := getFileBuf(OssConfigFileName)
	ossConfig := OssConfig{}
	_ = ossConfig.UnmarshalJSON(buf)
	fmt.Println("ossConfig", ossConfig)
	return ossConfig
}

// GetEs 获取elastic search配置
func GetEs() EsConfig {
	buf := getFileBuf(EsConfigFileName)
	var esConfig EsConfig
	_ = esConfig.UnmarshalJSON(buf)
	fmt.Println("esConfig", esConfig)
	return esConfig
}

// GetRabbitMQ 获取rabbit mq的配置
func GetRabbitMQ() RabbitMQConfig {
	buf := getFileBuf(RebbitMQConfigFileName)
	var rabbitMQConfig RabbitMQConfig
	_ = rabbitMQConfig.UnmarshalJSON(buf)
	fmt.Println("rabbitMQConfig", rabbitMQConfig)
	return rabbitMQConfig
}
