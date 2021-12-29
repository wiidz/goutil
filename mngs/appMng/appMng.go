package appMng

import (
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/esMng"
	"github.com/wiidz/goutil/mngs/memoryMng"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"github.com/wiidz/goutil/mngs/rabbitMng"
	"github.com/wiidz/goutil/mngs/redisMng"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
	"time"
)

var cacheM = memoryMng.NewCacheMng()
var mysqlM = mysqlMng.NewMysqlMng()

// GetSingletonAppMng : 获取单例app管理器
func GetSingletonAppMng(appID uint64, mysqlConfig *configStruct.MysqlConfig, configStruct configStruct.ProjectConfig, checkStart *configStruct.CheckStart) (mng *AppMng, err error) {

	//【1】从缓存中提取
	res, isExist := cacheM.Get("app-" + typeHelper.Uint64ToStr(appID) + "-config")
	if isExist && res != nil {
		return res.(*AppMng), nil
	}

	mng = &AppMng{
		ID:            appID,
		ProjectConfig: configStruct,
	}

	//【2】初始化mysql
	if checkStart.Mysql {
		//【2-1】基础配置
		err = mysqlMng.Init(mysqlConfig)
		if err != nil {
			return
		}
		//【2-2】基础配置
		mng.BaseConfig, err = mng.SetBaseConfig(mysqlConfig.DbName, mysqlConfig.SettingTableName)
		if err != nil {
			return
		}
		mng.BaseConfig.MysqlConfig = mysqlConfig
	}

	if mng.BaseConfig.Profile.Debug {
		log.Println("【" + mng.BaseConfig.Profile.Name + "】启动-调试模式")
	} else {
		log.Println("【" + mng.BaseConfig.Profile.Name + "】启动-生产模式")
	}

	//【4】初始化redis
	if checkStart.Redis {
		err = redisMng.Init(mng.BaseConfig.RedisConfig)
		if err != nil {
			return
		}
	}

	//【5】初始化es
	if checkStart.Es {
		err = esMng.Init(mng.BaseConfig.EsConfig)
		if err != nil {
			return
		}
	}

	//【6】初始化mq
	if checkStart.RabbitMQ {
		err = rabbitMng.Init(mng.BaseConfig.RabbitMQConfig)
		if err != nil {
			return
		}
	}

	//【7】项目配置
	err = mng.ProjectConfig.Build(mng.BaseConfig.Profile.Debug)
	if err != nil {
		return
	}

	//【8】写入缓存
	cacheM.Set("app-"+typeHelper.Uint64ToStr(appID)+"-config", mng, time.Minute*30)

	//【9】返回
	return
}

func (mng *AppMng) SetBaseConfig(dbName string, tableName string) (config *configStruct.BaseConfig, err error) {

	config = &configStruct.BaseConfig{}

	//【1】初始化sql
	conn := mysqlM.GetDBConn(dbName)

	//【1】查询数据
	var rows []*DbSettingRow
	err = conn.Table(tableName).Where("belonging = ?", "system").Find(&rows).Error
	if err != nil {
		return
	}

	//【2】初始化配置
	config.Location, _ = getLocationConfig(rows)
	config.Profile = getAppProfile(rows)

	// 数据源
	config.RedisConfig = getRedisConfig(rows, config.Profile.Debug)
	config.EsConfig = getEsConfig(rows, config.Profile.Debug)
	config.RabbitMQConfig = getRabbitMQConfig(rows, config.Profile.Debug)

	// 腾讯系
	config.WechatMiniConfig = getWechatMiniConfig(rows, config.Profile.Debug)
	config.WechatOaConfig = getWechatOaConfig(rows, config.Profile.Debug)
	config.WechatOpenConfig = getWechatOpenConfig(rows, config.Profile.Debug)
	config.WechatPayConfig = getWechatPayConfig(rows, config.Profile.Debug)

	// 阿里系
	config.OssConfig = getOssConfig(rows, config.Profile.Debug)
	config.AliPayConfig = getAliPayConfig(rows, config.Profile.Debug)
	config.AliApiConfig = getAliApiConfig(rows, config.Profile.Debug)
	config.AliSmsConfig = getAliSmsConfig(rows, config.Profile.Debug)

	return
}

// getLocationConfig : 获取时区设置
func getLocationConfig(rows []*DbSettingRow) (location *time.Location, err error) {
	timeZone := GetValueFromRow(rows, "time_zone", "", "", "Asia/Shanghai", false)
	location, err = time.LoadLocation(timeZone)
	if err != nil {
		location = time.FixedZone("CST-8", 8*3600)
	}
	return
}

//func getMysqlConfig(rows []*DbSettingRow) *configStruct.MysqlConfig {
//	return &configStruct.MysqlConfig{
//		Host:     GetValueFromRow(rows, "mysql", "host", "",
//		Username: GetValueFromRow(rows, "mysql", "username", "",
//		Password: GetValueFromRow(rows, "mysql", "password", "",
//		Port:     GetValueFromRow(rows, "mysql", "port", "",
//	}
//}

// getWechatMiniConfig : 获取微信小程序设置
func getWechatMiniConfig(rows []*DbSettingRow, debug bool) *configStruct.WechatMiniConfig {
	return &configStruct.WechatMiniConfig{
		AppID:     GetValueFromRow(rows, "wechat", "mini", "app_id", "", debug),
		AppSecret: GetValueFromRow(rows, "wechat", "mini", "app_secret", "", debug),
	}
}

func getWechatOaConfig(rows []*DbSettingRow, debug bool) *configStruct.WechatOaConfig {
	return &configStruct.WechatOaConfig{
		AppID:          GetValueFromRow(rows, "wechat", "oa", "app_id", "", debug),
		AppSecret:      GetValueFromRow(rows, "wechat", "oa", "app_secret", "", debug),
		Token:          GetValueFromRow(rows, "wechat", "oa", "token", "", debug),
		EncodingAESKey: GetValueFromRow(rows, "wechat", "oa", "encoding_aes_key", "", debug),
	}
}

func getWechatOpenConfig(rows []*DbSettingRow, debug bool) *configStruct.WechatOpenConfig {
	return &configStruct.WechatOpenConfig{
		AppID:     GetValueFromRow(rows, "wechat", "open", "app_id", "", debug),
		AppSecret: GetValueFromRow(rows, "wechat", "open", "app_secret", "", debug),
	}
}
func getWechatPayConfig(rows []*DbSettingRow, debug bool) *configStruct.WechatPayConfig {
	return &configStruct.WechatPayConfig{
		AppID:           GetValueFromRow(rows, "wechat", "pay", "app_id", "", debug),
		ApiKey:          GetValueFromRow(rows, "wechat", "pay", "api_key", "", debug),
		ApiKeyV3:        GetValueFromRow(rows, "wechat", "pay", "api_key_v3", "", debug),
		MchID:           GetValueFromRow(rows, "wechat", "pay", "mch_id", "", debug),
		CertURI:         GetValueFromRow(rows, "wechat", "pay", "cert_uri", "", debug),
		KeyURI:          GetValueFromRow(rows, "wechat", "pay", "key_uri", "", debug),
		CertContent:     GetValueFromRow(rows, "wechat", "pay", "cert_content", "", debug),
		CertSerialNo:    GetValueFromRow(rows, "wechat", "pay", "cert_serial_no", "", debug),
		NotifyURL:       GetValueFromRow(rows, "wechat", "pay", "notify_url", "", debug),
		RefundNotifyURL: GetValueFromRow(rows, "wechat", "pay", "refund_notify_url", "", debug),
		Debug:           GetValueFromRow(rows, "wechat", "pay", "debug", "0", debug) == "1", // 0=生产，1=调试
	}
}

func getAliPayConfig(rows []*DbSettingRow, debug bool) *configStruct.AliPayConfig {
	return &configStruct.AliPayConfig{
		AppID:      GetValueFromRow(rows, "ali", "pay", "app_id", "", debug),
		PrivateKey: GetValueFromRow(rows, "ali", "pay", "private_key", "", debug),
		NotifyURL:  GetValueFromRow(rows, "ali", "pay", "notify_url", "", debug),
		Debug:      GetValueFromRow(rows, "ali", "pay", "debug", "0", debug) == "1", // 0=生产，1=调试
	}
}

func getRedisConfig(rows []*DbSettingRow, debug bool) *configStruct.RedisConfig {
	return &configStruct.RedisConfig{
		Host:        GetValueFromRow(rows, "redis", "host", "", "127.0.0.1", debug),
		Port:        GetValueFromRow(rows, "redis", "port", "", "6379", debug),
		Username:    GetValueFromRow(rows, "redis", "username", "", "", debug),
		Password:    GetValueFromRow(rows, "redis", "password", "", "", debug),
		IdleTimeout: typeHelper.Str2Int(GetValueFromRow(rows, "redis", "idle_timeout", "", "60", debug)),
		Database:    typeHelper.Str2Int(GetValueFromRow(rows, "redis", "database", "", "", debug)),
		MaxActive:   typeHelper.Str2Int(GetValueFromRow(rows, "redis", "max_active", "", "10", debug)),
		MaxIdle:     typeHelper.Str2Int(GetValueFromRow(rows, "redis", "max_idle", "", "10", debug)),
	}
}

func getEsConfig(rows []*DbSettingRow, debug bool) *configStruct.EsConfig {
	return &configStruct.EsConfig{
		Host:     GetValueFromRow(rows, "es", "host", "", "http://127.0.0.1", debug),
		Port:     GetValueFromRow(rows, "es", "port", "", "9200", debug),
		Password: GetValueFromRow(rows, "es", "password", "", "123456", debug),
		Username: GetValueFromRow(rows, "es", "username", "", "es", debug),
	}
}

func getRabbitMQConfig(rows []*DbSettingRow, debug bool) *configStruct.RabbitMQConfig {
	return &configStruct.RabbitMQConfig{
		Host:     GetValueFromRow(rows, "rabbit_mq", "host", "", "http://127.0.0.1", debug),
		Password: GetValueFromRow(rows, "rabbit_mq", "password", "", "123456", debug),
		Username: GetValueFromRow(rows, "rabbit_mq", "username", "", "root", debug),
	}
}

func getOssConfig(rows []*DbSettingRow, debug bool) *configStruct.OssConfig {
	return &configStruct.OssConfig{
		AccessKeyID:     GetValueFromRow(rows, "ali", "oss", "access_key_id", "", debug),
		AccessKeySecret: GetValueFromRow(rows, "ali", "oss", "access_key_secret", "", debug),
		Host:            GetValueFromRow(rows, "ali", "oss", "host", "", debug),
		EndPoint:        GetValueFromRow(rows, "ali", "oss", "end_point", "", debug),
		BucketName:      GetValueFromRow(rows, "ali", "oss", "bucket_name", "", debug),
		ExpireTime:      typeHelper.Str2Int64(GetValueFromRow(rows, "ali", "oss", "expire_time", "30", debug)),
	}
}

// getAppProfile 项目基础信息
func getAppProfile(rows []*DbSettingRow) *configStruct.AppProfile {
	return &configStruct.AppProfile{
		No:      GetValueFromRow(rows, "app", "no", "", "", false),
		Name:    GetValueFromRow(rows, "app", "name", "", "", false),
		Host:    GetValueFromRow(rows, "app", "host", "", "", false),
		Debug:   GetValueFromRow(rows, "app", "debug", "", "", false) == "1", // 0=生产，1=调试
		Version: GetValueFromRow(rows, "app", "version", "", "", false),
	}
}

// getAliApiConfig 阿里云云市场API服务
func getAliApiConfig(rows []*DbSettingRow, debug bool) *configStruct.AliApiConfig {
	return &configStruct.AliApiConfig{
		AppKey:    GetValueFromRow(rows, "ali", "api", "app_key", "", debug),
		AppSecret: GetValueFromRow(rows, "ali", "api", "app_secret", "", debug),
		AppCode:   GetValueFromRow(rows, "ali", "api", "app_code", "", debug),
	}
}

// getAliSmsConfig 阿里云短信服务
func getAliSmsConfig(rows []*DbSettingRow, debug bool) *configStruct.AliSmsConfig {
	return &configStruct.AliSmsConfig{
		AccessKeySecret: GetValueFromRow(rows, "ali", "sms", "access_key_secret", "", debug),
		AccessKeyID:     GetValueFromRow(rows, "ali", "sms", "access_key_id", "", debug),
	}
}

// GetValueFromRow : 从rows中提取row
func GetValueFromRow(rows []*DbSettingRow, name, flag1, flag2 string, defaultValue string, debug bool) (value string) {

	var row = &DbSettingRow{}

	if len(rows) == 0 {
		return
	}

	for k := range rows {
		if rows[k].Name != name {
			continue
		}
		if flag1 != "" && rows[k].Flag1 != flag1 {
			continue
		}
		if flag2 != "" && rows[k].Flag2 != flag2 {
			continue
		}
		row = rows[k]
		break
	}

	// pdc模式下取value_1,debug模式下取value_2
	value = row.Value1
	if debug {
		value = row.Value2
	}

	if value == "" && defaultValue != "" {
		value = defaultValue
	}

	return
}
