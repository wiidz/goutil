package configMng

import (
	"fmt"
	"github.com/wiidz/goutil/helpers/osHelper"
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
)

var appConfig = AppConfig{}

type ConfigMng struct{}

func init() {
	buf := osHelper.GetFileBuf(ConfigPath + AppConfigFileName)
	_ = appConfig.UnmarshalJSON(buf)
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

func (*ConfigMng) GetHttpPort() string {
	return appConfig.HttpPort
}

func (*ConfigMng) GetMysql() MysqlConfig {
	buf := getFileBuf(MysqlConfigFileName)
	mysqlConfig := MysqlConfig{}
	_ = mysqlConfig.UnmarshalJSON(buf)
	return mysqlConfig
}

func (*ConfigMng) GetRedis() RedisConfig {
	buf := getFileBuf(RedisConfigFileName)
	redisConfig := RedisConfig{}
	_ = redisConfig.UnmarshalJSON(buf)
	return redisConfig
}

func (*ConfigMng) GetWechat() WechatConfig {
	buf := getFileBuf(WechatConfigFileName)
	wechatConfig := WechatConfig{}
	_ = wechatConfig.UnmarshalJSON(buf)
	return wechatConfig
}

func (*ConfigMng) GetOss() OssConfig {
	buf := getFileBuf(OssConfigFileName)
	ossConfig := OssConfig{}
	_ = ossConfig.UnmarshalJSON(buf)
	fmt.Println("ossConfig", ossConfig)
	return ossConfig
}

func (*ConfigMng) GetEs() EsConfig {
	buf := getFileBuf(EsConfigFileName)
	var esConfig EsConfig
	_ = esConfig.UnmarshalJSON(buf)
	fmt.Println("esConfig", esConfig)
	return esConfig
}
