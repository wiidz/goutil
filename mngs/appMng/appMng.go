package appMng

import (
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/memoryMng"
	"github.com/wiidz/goutil/mngs/mysqlMng"
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

	//【4】初始化redis、es
	if checkStart.Redis {
		err = redisMng.Init(mng.BaseConfig.RedisConfig)
		if err != nil {
			return
		}
	}

	//【3】项目配置
	mng.ProjectConfig.Build()

	//【5】写入缓存
	cacheM.Set("app-"+typeHelper.Uint64ToStr(appID)+"-config", mng, time.Minute*30)

	//【4】返回
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
	config.WechatMiniConfig = getWechatMiniConfig(rows)
	config.WechatOaConfig = getWechatOaConfig(rows)
	config.WechatOpenConfig = getWechatOpenConfig(rows)
	config.WechatPayConfig = getWechatPayConfig(rows)
	config.AliPayConfig = getAliPayConfig(rows)
	config.OssConfig = getOssConfig(rows)
	config.RedisConfig = getRedisConfig(rows)
	config.Profile = getAppProfile(rows)
	config.AliApiConfig = getAliApiConfig(rows)
	return
}

// getLocationConfig : 获取时区设置
func getLocationConfig(rows []*DbSettingRow) (location *time.Location) {
	timeZone := getRow(rows, "time_zone", "", "", "Asia/Shanghai").Value
	if timeZone == "" {
		timeZone = "Asia/Shanghai" // 默认东八区
	}
	location, _ = time.LoadLocation(timeZone)
	return
}

//func getMysqlConfig(rows []*DbSettingRow) *configStruct.MysqlConfig {
//	return &configStruct.MysqlConfig{
//		Host:     getRow(rows, "mysql", "host", "").Value,
//		Username: getRow(rows, "mysql", "username", "").Value,
//		Password: getRow(rows, "mysql", "password", "").Value,
//		Port:     getRow(rows, "mysql", "port", "").Value,
//	}
//}

// getWechatMiniConfig : 获取微信小程序设置
func getWechatMiniConfig(rows []*DbSettingRow) *configStruct.WechatMiniConfig {
	return &configStruct.WechatMiniConfig{
		AppID:     getRow(rows, "wechat", "mini", "app_id", "").Value,
		AppSecret: getRow(rows, "wechat", "mini", "app_secret", "").Value,
	}
}

func getWechatOaConfig(rows []*DbSettingRow) *configStruct.WechatOaConfig {
	return &configStruct.WechatOaConfig{
		AppID:     getRow(rows, "wechat", "oa", "app_id", "").Value,
		AppSecret: getRow(rows, "wechat", "oa", "app_secret", "").Value,
	}
}

func getWechatOpenConfig(rows []*DbSettingRow) *configStruct.WechatOpenConfig {
	return &configStruct.WechatOpenConfig{
		AppID:     getRow(rows, "wechat", "open", "app_id", "").Value,
		AppSecret: getRow(rows, "wechat", "open", "app_secret", "").Value,
	}
}
func getWechatPayConfig(rows []*DbSettingRow) *configStruct.WechatPayConfig {
	return &configStruct.WechatPayConfig{
		AppID:           getRow(rows, "wechat", "pay", "app_id", "").Value,
		ApiKey:          getRow(rows, "wechat", "pay", "api_key", "").Value,
		MchID:           getRow(rows, "wechat", "pay", "mch_id", "").Value,
		CertURI:         getRow(rows, "wechat", "pay", "cert_uri", "").Value,
		KeyURI:          getRow(rows, "wechat", "pay", "key_uri", "").Value,
		CertContent:     getRow(rows, "wechat", "pay", "cert_content", "").Value,
		NotifyURL:       getRow(rows, "wechat", "pay", "notify_url", "").Value,
		RefundNotifyURL: getRow(rows, "wechat", "pay", "refund_notify_url", "").Value,
		IsProd:          getRow(rows, "wechat", "pay", "is_prod", "false").Value == "1", // 0=调试，1=生产
	}
}

func getAliPayConfig(rows []*DbSettingRow) *configStruct.AliPayConfig {
	return &configStruct.AliPayConfig{
		AppID:      getRow(rows, "ali", "pay", "app_id", "").Value,
		PrivateKey: getRow(rows, "ali", "pay", "private_key", "").Value,
		NotifyURL:  getRow(rows, "ali", "pay", "notify_url", "").Value,
		IsProd:     getRow(rows, "ali", "pay", "is_prod", "false").Value == "1", // 0=调试，1=生产
	}
}

func getRedisConfig(rows []*DbSettingRow) *configStruct.RedisConfig {
	return &configStruct.RedisConfig{
		Host:        getRow(rows, "redis", "host", "", "").Value,
		Port:        getRow(rows, "redis", "port", "", "").Value,
		Password:    getRow(rows, "redis", "password", "", "").Value,
		IdleTimeout: typeHelper.Str2Int(getRow(rows, "redis", "idle_timeout", "", "60").Value),
		Database:    typeHelper.Str2Int(getRow(rows, "redis", "database", "", "").Value),
		MaxActive:   typeHelper.Str2Int(getRow(rows, "redis", "max_active", "", "10").Value),
		MaxIdle:     typeHelper.Str2Int(getRow(rows, "redis", "max_idle", "", "10").Value),
	}
}

func getOssConfig(rows []*DbSettingRow) *configStruct.OssConfig {
	return &configStruct.OssConfig{
		AccessKeyID:     getRow(rows, "ali", "oss", "access_key_id", "").Value,
		AccessKeySecret: getRow(rows, "ali", "oss", "access_key_secret", "").Value,
		Host:            getRow(rows, "ali", "oss", "host", "").Value,
		EndPoint:        getRow(rows, "ali", "oss", "end_point", "").Value,
		BucketName:      getRow(rows, "ali", "oss", "bucket_name", "").Value,
		ExpireTime:      typeHelper.Str2Int64(getRow(rows, "ali", "oss", "expire_time", "30").Value),
	}
}

func getAppProfile(rows []*DbSettingRow) *configStruct.AppProfile {
	return &configStruct.AppProfile{
		No:      getRow(rows, "app", "no", "", "").Value,
		Name:    getRow(rows, "app", "name", "", "").Value,
		Host:    getRow(rows, "app", "host", "", "").Value,
		Debug:   getRow(rows, "app", "debug", "", "").Value == "0", // 0=生产，1=调试
		Version: getRow(rows, "app", "version", "", "").Value,
	}
}

func getAliApiConfig(rows []*DbSettingRow) *configStruct.AliApiConfig {
	return &configStruct.AliApiConfig{
		AppID:     getRow(rows, "ali", "app_id", "", "").Value,
		AppSecret: getRow(rows, "ali", "app_secret", "", "").Value,
		AppCode:   getRow(rows, "ali", "app_code", "", "").Value,
	}
}

// getRow : 从rows中提取row
func getRow(rows []*DbSettingRow, name, flag1, flag2 string, defaultValue string) (row *DbSettingRow) {

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
