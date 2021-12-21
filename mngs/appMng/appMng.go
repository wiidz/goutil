package appMng

import (
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/esMng"
	"github.com/wiidz/goutil/mngs/memoryMng"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"github.com/wiidz/goutil/mngs/rabbitMng"
	"github.com/wiidz/goutil/mngs/redisMng"
	"github.com/wiidz/goutil/structs/configStruct"
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
	err = mng.ProjectConfig.Build()
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
	config.Location = getLocationConfig(rows)
	config.Profile = getAppProfile(rows)

	// 数据源
	config.RedisConfig = getRedisConfig(rows)
	config.EsConfig = getEsConfig(rows)
	config.RabbitMQConfig = getRabbitMQConfig(rows)

	// 腾讯系
	config.WechatMiniConfig = getWechatMiniConfig(rows)
	config.WechatOaConfig = getWechatOaConfig(rows)
	config.WechatOpenConfig = getWechatOpenConfig(rows)
	config.WechatPayConfig = getWechatPayConfig(rows)

	// 阿里系
	config.OssConfig = getOssConfig(rows)
	config.AliPayConfig = getAliPayConfig(rows)
	config.AliApiConfig = getAliApiConfig(rows)
	config.AliSmsConfig = getAliSmsConfig(rows)

	return
}

// getLocationConfig : 获取时区设置
func getLocationConfig(rows []*DbSettingRow) (location *time.Location) {
	timeZone := GetValueFromRow(rows, "time_zone", "", "", "Asia/Shanghai").Value
	if timeZone == "" {
		timeZone = "Asia/Shanghai" // 默认东八区
	}
	location, _ = time.LoadLocation(timeZone)
	return
}

//func getMysqlConfig(rows []*DbSettingRow) *configStruct.MysqlConfig {
//	return &configStruct.MysqlConfig{
//		Host:     GetValueFromRow(rows, "mysql", "host", "").Value,
//		Username: GetValueFromRow(rows, "mysql", "username", "").Value,
//		Password: GetValueFromRow(rows, "mysql", "password", "").Value,
//		Port:     GetValueFromRow(rows, "mysql", "port", "").Value,
//	}
//}

// getWechatMiniConfig : 获取微信小程序设置
func getWechatMiniConfig(rows []*DbSettingRow) *configStruct.WechatMiniConfig {
	return &configStruct.WechatMiniConfig{
		AppID:     GetValueFromRow(rows, "wechat", "mini", "app_id", "").Value,
		AppSecret: GetValueFromRow(rows, "wechat", "mini", "app_secret", "").Value,
	}
}

func getWechatOaConfig(rows []*DbSettingRow) *configStruct.WechatOaConfig {
	return &configStruct.WechatOaConfig{
		AppID:     GetValueFromRow(rows, "wechat", "oa", "app_id", "").Value,
		AppSecret: GetValueFromRow(rows, "wechat", "oa", "app_secret", "").Value,
	}
}

func getWechatOpenConfig(rows []*DbSettingRow) *configStruct.WechatOpenConfig {
	return &configStruct.WechatOpenConfig{
		AppID:     GetValueFromRow(rows, "wechat", "open", "app_id", "").Value,
		AppSecret: GetValueFromRow(rows, "wechat", "open", "app_secret", "").Value,
	}
}
func getWechatPayConfig(rows []*DbSettingRow) *configStruct.WechatPayConfig {
	return &configStruct.WechatPayConfig{
		AppID:           GetValueFromRow(rows, "wechat", "pay", "app_id", "").Value,
		ApiKey:          GetValueFromRow(rows, "wechat", "pay", "api_key", "").Value,
		ApiKeyV3:        GetValueFromRow(rows, "wechat", "pay", "api_key_v3", "").Value,
		MchID:           GetValueFromRow(rows, "wechat", "pay", "mch_id", "").Value,
		CertURI:         GetValueFromRow(rows, "wechat", "pay", "cert_uri", "").Value,
		KeyURI:          GetValueFromRow(rows, "wechat", "pay", "key_uri", "").Value,
		CertContent:     GetValueFromRow(rows, "wechat", "pay", "cert_content", "").Value,
		CertSerialNo: GetValueFromRow(rows, "wechat", "pay", "cert_serial_no", "").Value,
		NotifyURL:       GetValueFromRow(rows, "wechat", "pay", "notify_url", "").Value,
		RefundNotifyURL: GetValueFromRow(rows, "wechat", "pay", "refund_notify_url", "").Value,
		IsProd:          GetValueFromRow(rows, "wechat", "pay", "is_prod", "0").Value == "1", // 0=调试，1=生产
	}
}

func getAliPayConfig(rows []*DbSettingRow) *configStruct.AliPayConfig {
	return &configStruct.AliPayConfig{
		AppID:      GetValueFromRow(rows, "ali", "pay", "app_id", "").Value,
		PrivateKey: GetValueFromRow(rows, "ali", "pay", "private_key", "").Value,
		NotifyURL:  GetValueFromRow(rows, "ali", "pay", "notify_url", "").Value,
		IsProd:     GetValueFromRow(rows, "ali", "pay", "is_prod", "false").Value == "1", // 0=调试，1=生产
	}
}

func getRedisConfig(rows []*DbSettingRow) *configStruct.RedisConfig {
	return &configStruct.RedisConfig{
		Host:        GetValueFromRow(rows, "redis", "host", "", "127.0.0.1").Value,
		Port:        GetValueFromRow(rows, "redis", "port", "", "6379").Value,
		Password:    GetValueFromRow(rows, "redis", "password", "", "").Value,
		IdleTimeout: typeHelper.Str2Int(GetValueFromRow(rows, "redis", "idle_timeout", "", "60").Value),
		Database:    typeHelper.Str2Int(GetValueFromRow(rows, "redis", "database", "", "").Value),
		MaxActive:   typeHelper.Str2Int(GetValueFromRow(rows, "redis", "max_active", "", "10").Value),
		MaxIdle:     typeHelper.Str2Int(GetValueFromRow(rows, "redis", "max_idle", "", "10").Value),
	}
}

func getEsConfig(rows []*DbSettingRow) *configStruct.EsConfig {
	return &configStruct.EsConfig{
		Host:     GetValueFromRow(rows, "es", "host", "", "http://127.0.0.1").Value,
		Port:     GetValueFromRow(rows, "es", "port", "", "9200").Value,
		Password: GetValueFromRow(rows, "es", "password", "", "123456").Value,
		Username: GetValueFromRow(rows, "es", "username", "", "es").Value,
	}
}


func getRabbitMQConfig(rows []*DbSettingRow) *configStruct.RabbitMQConfig {
	return &configStruct.RabbitMQConfig{
		Host:     GetValueFromRow(rows, "rabbit_mq", "host", "", "http://127.0.0.1").Value,
		Password: GetValueFromRow(rows, "rabbit_mq", "password", "", "123456").Value,
		Username: GetValueFromRow(rows, "rabbit_mq", "username", "", "root").Value,
	}
}

func getOssConfig(rows []*DbSettingRow) *configStruct.OssConfig {
	return &configStruct.OssConfig{
		AccessKeyID:     GetValueFromRow(rows, "ali", "oss", "access_key_id", "").Value,
		AccessKeySecret: GetValueFromRow(rows, "ali", "oss", "access_key_secret", "").Value,
		Host:            GetValueFromRow(rows, "ali", "oss", "host", "").Value,
		EndPoint:        GetValueFromRow(rows, "ali", "oss", "end_point", "").Value,
		BucketName:      GetValueFromRow(rows, "ali", "oss", "bucket_name", "").Value,
		ExpireTime:      typeHelper.Str2Int64(GetValueFromRow(rows, "ali", "oss", "expire_time", "30").Value),
	}
}

// getAppProfile 项目基础信息
func getAppProfile(rows []*DbSettingRow) *configStruct.AppProfile {
	return &configStruct.AppProfile{
		No:      GetValueFromRow(rows, "app", "no", "", "").Value,
		Name:    GetValueFromRow(rows, "app", "name", "", "").Value,
		Host:    GetValueFromRow(rows, "app", "host", "", "").Value,
		Debug:   GetValueFromRow(rows, "app", "debug", "", "").Value == "0", // 0=生产，1=调试
		Version: GetValueFromRow(rows, "app", "version", "", "").Value,
	}
}

// getAliApiConfig 阿里云云市场API服务
func getAliApiConfig(rows []*DbSettingRow) *configStruct.AliApiConfig {
	return &configStruct.AliApiConfig{
		AppID:     GetValueFromRow(rows, "ali", "api", "app_id", "").Value,
		AppSecret: GetValueFromRow(rows, "ali", "api", "app_secret", "").Value,
		AppCode:   GetValueFromRow(rows, "ali", "api", "app_code", "").Value,
	}
}

// getAliSmsConfig 阿里云短信服务
func getAliSmsConfig(rows []*DbSettingRow) *configStruct.AliSmsConfig {
	return &configStruct.AliSmsConfig{
		AccessKeySecret: GetValueFromRow(rows, "ali", "sms", "access_key_secret", "").Value,
		AccessKeyID:     GetValueFromRow(rows, "ali", "sms", "access_key_id", "").Value,
	}
}

// GetValueFromRow : 从rows中提取row
func GetValueFromRow(rows []*DbSettingRow, name, flag1, flag2 string, defaultValue string) (row *DbSettingRow) {

	row = &DbSettingRow{}

	if len(rows) == 0 {
		return nil
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

	if row.Value == "" && defaultValue != "" {
		row.Value = defaultValue
	}

	return
}
