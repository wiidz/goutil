package appMng

import (
	"time"

	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
)

func buildBaseConfig(rows []*DbSettingRow) *configStruct.BaseConfig {
	cfg := &configStruct.BaseConfig{}
	if len(rows) == 0 {
		cfg.Profile = &configStruct.AppProfile{}
		return cfg
	}

	location, err := getLocationConfig(rows)
	if err == nil {
		cfg.Location = location
	}
	cfg.Profile = getAppProfile(rows)

	debug := cfg.Profile != nil && cfg.Profile.Debug

	cfg.RedisConfig = getRedisConfig(rows, debug)
	cfg.EsConfig = getEsConfig(rows, debug)
	cfg.RabbitMQConfig = getRabbitMQConfig(rows, debug)
	cfg.PostgresConfig = getPostgresConfig(rows, debug)

	cfg.WechatMiniConfig = getWechatMiniConfig(rows, debug)
	cfg.WechatOaConfig = getWechatOaConfig(rows, debug)
	cfg.WechatOpenConfig = getWechatOpenConfig(rows, debug)
	cfg.WechatPayConfigV3 = getWechatPayConfigV3(rows, debug)
	cfg.WechatPayConfigV2 = getWechatPayConfigV2(rows, debug)

	cfg.AliOssConfig = getAliOssConfig(rows, debug)
	cfg.AliPayConfig = getAliPayConfig(rows, debug)
	cfg.AliApiConfig = getAliApiConfig(rows, debug)
	cfg.AliSmsConfig = getAliSmsConfig(rows, debug)
	cfg.AliIotConfig = getAliIotConfig(rows, debug)
	cfg.AmapConfig = getAmapConfig(rows, debug)

	cfg.YunxinConfig = getYunXinConfig(rows, debug)

	return cfg
}

func getLocationConfig(rows []*DbSettingRow) (location *time.Location, err error) {
	timeZone := GetValueFromRow(rows, "time_zone", "", "", "Asia/Shanghai", false)
	location, err = time.LoadLocation(timeZone)
	if err != nil {
		location = time.FixedZone("CST-8", 8*3600)
	}
	return
}

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

func getWechatPayConfigV3(rows []*DbSettingRow, debug bool) *configStruct.WechatPayConfigV3 {
	return &configStruct.WechatPayConfigV3{
		AppID:                     GetValueFromRow(rows, "wechat", "pay_v3", "app_id", "", debug),
		ApiKeyV3:                  GetValueFromRow(rows, "wechat", "pay_v3", "api_key", "", debug),
		MchID:                     GetValueFromRow(rows, "wechat", "pay_v3", "mch_id", "", debug),
		CertURI:                   GetValueFromRow(rows, "wechat", "pay_v3", "cert_uri", "", debug),
		KeyURI:                    GetValueFromRow(rows, "wechat", "pay_v3", "key_uri", "", debug),
		PEMPrivateKeyContent:      GetValueFromRow(rows, "wechat", "pay_v3", "pem_private_key_content", "", debug),
		PEMCertContent:            GetValueFromRow(rows, "wechat", "pay_v3", "pem_cert_content", "", debug),
		CertSerialNo:              GetValueFromRow(rows, "wechat", "pay_v3", "cert_serial_no", "", debug),
		NotifyURL:                 GetValueFromRow(rows, "wechat", "pay_v3", "notify_url", "", debug),
		RefundNotifyURL:           GetValueFromRow(rows, "wechat", "pay_v3", "refund_notify_url", "", debug),
		MerchantTransferNotifyURL: GetValueFromRow(rows, "wechat", "pay_v3", "merchant_transfer_notify_url", "", debug),
		Debug:                     GetValueFromRow(rows, "wechat", "pay_v3", "debug", "0", debug) == "1",
	}
}

func getWechatPayConfigV2(rows []*DbSettingRow, debug bool) *configStruct.WechatPayConfigV2 {
	return &configStruct.WechatPayConfigV2{
		AppID:           GetValueFromRow(rows, "wechat", "pay_v2", "app_id", "", debug),
		ApiKey:          GetValueFromRow(rows, "wechat", "pay_v2", "api_key", "", debug),
		MchID:           GetValueFromRow(rows, "wechat", "pay_v2", "mch_id", "", debug),
		CertURI:         GetValueFromRow(rows, "wechat", "pay_v2", "cert_uri", "", debug),
		KeyURI:          GetValueFromRow(rows, "wechat", "pay_v2", "key_uri", "", debug),
		P12CertFilePath: GetValueFromRow(rows, "wechat", "pay_v2", "p12_cert_file_path", "", debug),
		CertSerialNo:    GetValueFromRow(rows, "wechat", "pay_v2", "cert_serial_no", "", debug),
		NotifyURL:       GetValueFromRow(rows, "wechat", "pay_v2", "notify_url", "", debug),
		RefundNotifyURL: GetValueFromRow(rows, "wechat", "pay_v2", "refund_notify_url", "", debug),
		Debug:           GetValueFromRow(rows, "wechat", "pay_v2", "debug", "0", debug) == "1",
	}
}

func getAliPayConfig(rows []*DbSettingRow, debug bool) *configStruct.AliPayConfig {
	return &configStruct.AliPayConfig{
		AppID:            GetValueFromRow(rows, "ali", "pay", "app_id", "", debug),
		PrivateKey:       GetValueFromRow(rows, "ali", "pay", "private_key", "", debug),
		NotifyURL:        GetValueFromRow(rows, "ali", "pay", "notify_url", "", debug),
		Debug:            GetValueFromRow(rows, "ali", "pay", "debug", "0", debug) == "1",
		AppCertPublicKey: GetValueFromRow(rows, "ali", "pay", "app_cert_public_key", "", debug),
		CertPublicKey:    GetValueFromRow(rows, "ali", "pay", "cert_public_key", "", debug),
		RootCert:         GetValueFromRow(rows, "ali", "pay", "root_cert", "", debug),
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

func getPostgresConfig(rows []*DbSettingRow, debug bool) *configStruct.PostgresConfig {
	dsn := GetValueFromRow(rows, "postgres", "", "dsn", "", debug)
	if dsn == "" {
		dsn = GetValueFromRow(rows, "postgres", "", "", "", debug)
	}
	if dsn == "" {
		return nil
	}

	cfg := &configStruct.PostgresConfig{DSN: dsn}
	if v := GetValueFromRow(rows, "postgres", "", "conn_max_idle", "", debug); v != "" {
		cfg.ConnMaxIdle = typeHelper.Str2Int(v)
	}
	if v := GetValueFromRow(rows, "postgres", "", "conn_max_open", "", debug); v != "" {
		cfg.ConnMaxOpen = typeHelper.Str2Int(v)
	}
	if v := GetValueFromRow(rows, "postgres", "", "conn_max_lifetime", "", debug); v != "" {
		cfg.ConnMaxLifetime = time.Duration(typeHelper.Str2Int64(v)) * time.Second
	}
	return cfg
}

func getAliOssConfig(rows []*DbSettingRow, debug bool) *configStruct.AliOssConfig {
	return &configStruct.AliOssConfig{
		AccessKeyID:     GetValueFromRow(rows, "ali", "oss", "access_key_id", "", debug),
		AccessKeySecret: GetValueFromRow(rows, "ali", "oss", "access_key_secret", "", debug),
		Host:            GetValueFromRow(rows, "ali", "oss", "host", "", debug),
		EndPoint:        GetValueFromRow(rows, "ali", "oss", "end_point", "", debug),
		BucketName:      GetValueFromRow(rows, "ali", "oss", "bucket_name", "", debug),
		ExpireTime:      typeHelper.Str2Int64(GetValueFromRow(rows, "ali", "oss", "expire_time", "30", debug)),
		ARN:             GetValueFromRow(rows, "ali", "oss", "arn", "", debug),
	}
}

func getAppProfile(rows []*DbSettingRow) *configStruct.AppProfile {
	return &configStruct.AppProfile{
		No:      GetValueFromRow(rows, "app", "", "no", "", false),
		Name:    GetValueFromRow(rows, "app", "", "name", "", false),
		Host:    GetValueFromRow(rows, "app", "", "host", "", false),
		Port:    GetValueFromRow(rows, "app", "", "port", "127.0.0.1", false),
		Domain:  GetValueFromRow(rows, "app", "", "domain", "", false),
		Debug:   GetValueFromRow(rows, "app", "", "debug", "", false) == "1",
		Version: GetValueFromRow(rows, "app", "", "version", "", false),
	}
}

func getAliApiConfig(rows []*DbSettingRow, debug bool) *configStruct.AliApiConfig {
	return &configStruct.AliApiConfig{
		AppKey:    GetValueFromRow(rows, "ali", "api", "app_key", "", debug),
		AppSecret: GetValueFromRow(rows, "ali", "api", "app_secret", "", debug),
		AppCode:   GetValueFromRow(rows, "ali", "api", "app_code", "", debug),
	}
}

func getAliSmsConfig(rows []*DbSettingRow, debug bool) *configStruct.AliSmsConfig {
	return &configStruct.AliSmsConfig{
		AccessKeySecret: GetValueFromRow(rows, "ali", "sms", "access_key_secret", "", debug),
		AccessKeyID:     GetValueFromRow(rows, "ali", "sms", "access_key_id", "", debug),
	}
}

func getAliIotConfig(rows []*DbSettingRow, debug bool) *configStruct.AliIotConfig {
	return &configStruct.AliIotConfig{
		AccessKeySecret: GetValueFromRow(rows, "ali", "iot", "access_key_secret", "", debug),
		AccessKeyID:     GetValueFromRow(rows, "ali", "iot", "access_key_id", "", debug),
		EndPoint:        GetValueFromRow(rows, "ali", "iot", "end_point", "", debug),
		RegionID:        GetValueFromRow(rows, "ali", "iot", "region_id", "", debug),
	}
}

func getAmapConfig(rows []*DbSettingRow, debug bool) *configStruct.AmapConfig {
	return &configStruct.AmapConfig{Key: GetValueFromRow(rows, "ali", "amap", "key", "", debug)}
}

func getYunXinConfig(rows []*DbSettingRow, debug bool) *configStruct.YunxinConfig {
	return &configStruct.YunxinConfig{
		AppKey:    GetValueFromRow(rows, "netease", "yunxin", "app_key", "", debug),
		AppSecret: GetValueFromRow(rows, "netease", "yunxin", "app_secret", "", debug),
	}
}

// GetValueFromRow 从 rows 中检索符合条件的数据。
func GetValueFromRow(rows []*DbSettingRow, name, flag1, flag2, defaultValue string, debug bool) (value string) {
	if len(rows) == 0 {
		return
	}

	var row *DbSettingRow
	for i := range rows {
		item := rows[i]
		if item.Name != name {
			continue
		}
		if flag1 != "" && item.Flag1 != flag1 {
			continue
		}
		if flag2 != "" && item.Flag2 != flag2 {
			continue
		}
		row = item
		break
	}

	if row == nil {
		value = defaultValue
		return
	}

	value = row.Value1
	if debug {
		value = row.Value2
	}
	if value == "" && defaultValue != "" {
		value = defaultValue
	}
	return
}

