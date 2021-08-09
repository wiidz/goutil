package configMng

import (
	"fmt"
	"github.com/wiidz/goutil/helpers/osHelper"
	"log"
)

const (
	ConfigPath           = "./configs/" // 路径写死
	DevPath = "dev/"
	PdcPath = "pdc/"
	AppConfigFileName    = "app.json"
	WechatConfigFileName = "wechat.json"
	MysqlConfigFileName  = "mysql.json"
	EsConfigFileName     = "es.json"
	RedisConfigFileName  = "redis.json"
	OssConfigFileName    = "oss.json"
	RebbitMQConfigFileName     = "rabbit-mq.json"
)

var appConfig = AppConfig{}

type ConfigMng struct{}

func init() {
	log.Println(ConfigPath + AppConfigFileName)
	buf := osHelper.GetFileBuf(ConfigPath + AppConfigFileName)
	_ = appConfig.UnmarshalJSON(buf)
	log.Println("appConfig",appConfig)
}

// 获取指定目录里的配置文件
func getTargetDir() string {
	var dir string
	if appConfig.Debug == false {
		dir = PdcPath
	} else {
		dir = DevPath
	}
	return ConfigPath + dir
}

func getFileBuf(fileName string) []byte {
	return osHelper.GetFileBuf(getTargetDir() + fileName)
}

func  GetHttpPort() string {
	return appConfig.HttpPort
}

func GetMysql() MysqlConfig {
	buf := getFileBuf(MysqlConfigFileName)
	mysqlConfig := MysqlConfig{}
	_ = mysqlConfig.UnmarshalJSON(buf)
	return mysqlConfig
}

func  GetRedis() RedisConfig {
	buf := getFileBuf(RedisConfigFileName)
	redisConfig := RedisConfig{}
	_ = redisConfig.UnmarshalJSON(buf)
	return redisConfig
}

func  GetWechat() WechatConfig {
	buf := getFileBuf(WechatConfigFileName)
	wechatConfig := WechatConfig{}
	_ = wechatConfig.UnmarshalJSON(buf)
	return wechatConfig
}

func  GetOss() OssConfig {
	buf := getFileBuf(OssConfigFileName)
	ossConfig := OssConfig{}
	_ = ossConfig.UnmarshalJSON(buf)
	fmt.Println("ossConfig", ossConfig)
	return ossConfig
}

func  GetEs() EsConfig {
	buf := getFileBuf(EsConfigFileName)
	var esConfig EsConfig
	_ = esConfig.UnmarshalJSON(buf)
	fmt.Println("esConfig", esConfig)
	return esConfig
}

// GetRabbitMQ 获取rabbit mq的配置
func  GetRabbitMQ() RabbitMQConfig {
	buf := getFileBuf(RebbitMQConfigFileName)
	var rabbitMQConfig RabbitMQConfig
	_ = rabbitMQConfig.UnmarshalJSON(buf)
	fmt.Println("rabbitMQConfig", rabbitMQConfig)
	return rabbitMQConfig
}