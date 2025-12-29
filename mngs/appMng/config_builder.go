package appMng

import (
	"fmt"
	"log"
	"time"

	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
)

// ConfigLoadOptions 配置加载选项，控制哪些配置需要加载
type ConfigLoadOptions struct {
	// 基础配置（总是加载）
	// Profile 和 Location 总是加载，不需要选项

	// 存储相关配置
	Redis    bool // Redis 配置
	Es       bool // Elasticsearch 配置
	RabbitMQ bool // RabbitMQ 配置
	Postgres bool // PostgreSQL 配置

	// 微信相关配置
	WechatMini  bool // 微信小程序配置
	WechatOa    bool // 微信公众号配置
	WechatOpen  bool // 微信开放平台配置
	WechatPayV3 bool // 微信支付 V3 配置
	WechatPayV2 bool // 微信支付 V2 配置

	// 阿里相关配置
	AliOss bool // 阿里云 OSS 配置
	AliPay bool // 支付宝配置
	AliApi bool // 阿里云 API 配置
	AliSms bool // 阿里云短信配置
	AliIot bool // 阿里云 IoT 配置
	Amap   bool // 高德地图配置

	// 其他配置
	Yunxin bool // 网易云信配置
}

// DefaultConfigLoadOptions 返回默认的配置加载选项（所有配置都加载）
func DefaultConfigLoadOptions() *ConfigLoadOptions {
	return &ConfigLoadOptions{
		Redis:       true,
		Es:          true,
		RabbitMQ:    true,
		Postgres:    true,
		WechatMini:  true,
		WechatOa:    true,
		WechatOpen:  true,
		WechatPayV3: true,
		WechatPayV2: true,
		AliOss:      true,
		AliPay:      true,
		AliApi:      true,
		AliSms:      true,
		AliIot:      true,
		Amap:        true,
		Yunxin:      true,
	}
}

// BuildBaseConfig 从数据库配置行构建 BaseConfig
// rows: 数据库配置行
// options: 配置加载选项，如果为 nil 则使用默认选项（加载所有配置）
func BuildBaseConfig(rows []*DbSettingRow, options *ConfigLoadOptions) (*configStruct.BaseConfig, error) {
	cfg := &configStruct.BaseConfig{}
	if len(rows) == 0 {
		cfg.Profile = &configStruct.AppProfile{}
		return cfg, nil
	}

	location, err := getLocationConfig(rows)
	if err == nil {
		cfg.Location = location
	}
	cfg.Profile = getAppProfile(rows)

	debug := cfg.Profile != nil && cfg.Profile.Debug

	// 如果 options 为 nil，使用默认选项（加载所有配置）
	if options == nil {
		log.Println("options 为 nil，使用默认选项")
		options = DefaultConfigLoadOptions()
	}

	// 根据选项决定是否加载各个配置
	// 注意：只有在 options 中明确指定要加载的配置才会被加载和验证
	// 如果某个配置在 options 中未指定，即使数据库中有相关配置也不会加载
	// 每个 get*Config 函数已经负责验证配置的有效性（包括关键字段），这里直接赋值即可
	if options.Redis {
		redisConfig, err := getRedisConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载 Redis 配置失败: %w", err)
		}
		cfg.RedisConfig = redisConfig
	}
	if options.Es {
		esConfig, err := getEsConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载 Elasticsearch 配置失败: %w", err)
		}
		cfg.EsConfig = esConfig
	}
	if options.RabbitMQ {
		rabbitMQConfig, err := getRabbitMQConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载 RabbitMQ 配置失败: %w", err)
		}
		cfg.RabbitMQConfig = rabbitMQConfig
	}
	if options.Postgres {
		postgresConfig, err := getPostgresConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载 PostgreSQL 配置失败: %w", err)
		}
		// PostgresConfig 可以为 nil（如果 DSN 为空），这是允许的
		cfg.PostgresConfig = postgresConfig
	}

	if options.WechatMini {
		wechatMiniConfig, err := getWechatMiniConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载微信小程序配置失败: %w", err)
		}
		cfg.WechatMiniConfig = wechatMiniConfig
	}
	if options.WechatOa {
		wechatOaConfig, err := getWechatOaConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载微信公众号配置失败: %w", err)
		}
		cfg.WechatOaConfig = wechatOaConfig
	}
	if options.WechatOpen {
		wechatOpenConfig, err := getWechatOpenConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载微信开放平台配置失败: %w", err)
		}
		cfg.WechatOpenConfig = wechatOpenConfig
	}
	if options.WechatPayV3 {
		wechatPayV3Config, err := getWechatPayConfigV3(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载微信支付 V3 配置失败: %w", err)
		}
		cfg.WechatPayConfigV3 = wechatPayV3Config
	}
	if options.WechatPayV2 {
		wechatPayV2Config, err := getWechatPayConfigV2(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载微信支付 V2 配置失败: %w", err)
		}
		cfg.WechatPayConfigV2 = wechatPayV2Config
	}

	if options.AliOss {
		aliOssConfig, err := getAliOssConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载阿里云 OSS 配置失败: %w", err)
		}
		cfg.AliOssConfig = aliOssConfig
	}
	if options.AliPay {
		aliPayConfig, err := getAliPayConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载支付宝配置失败: %w", err)
		}
		cfg.AliPayConfig = aliPayConfig
	}
	if options.AliApi {
		aliApiConfig, err := getAliApiConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载阿里云 API 配置失败: %w", err)
		}
		cfg.AliApiConfig = aliApiConfig
	}
	if options.AliSms {
		aliSmsConfig, err := getAliSmsConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载阿里云短信配置失败: %w", err)
		}
		cfg.AliSmsConfig = aliSmsConfig
	}
	if options.AliIot {
		aliIotConfig, err := getAliIotConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载阿里云 IoT 配置失败: %w", err)
		}
		cfg.AliIotConfig = aliIotConfig
	}
	if options.Amap {
		amapConfig, err := getAmapConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载高德地图配置失败: %w", err)
		}
		cfg.AmapConfig = amapConfig
	}

	if options.Yunxin {
		yunxinConfig, err := getYunXinConfig(rows, debug)
		if err != nil {
			return nil, fmt.Errorf("加载网易云信配置失败: %w", err)
		}
		cfg.YunxinConfig = yunxinConfig
	}

	return cfg, nil
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

func getLocationConfig(rows []*DbSettingRow) (location *time.Location, err error) {
	timeZone := GetValueFromRow(rows, "time_zone", "", "", "Asia/Shanghai", false)
	location, err = time.LoadLocation(timeZone)
	if err != nil {
		location = time.FixedZone("CST-8", 8*3600)
	}
	return
}

func getWechatMiniConfig(rows []*DbSettingRow, debug bool) (*configStruct.WechatMiniConfig, error) {
	cfg := &configStruct.WechatMiniConfig{
		AppID:     GetValueFromRow(rows, "wechat", "mini", "app_id", "", debug),
		AppSecret: GetValueFromRow(rows, "wechat", "mini", "app_secret", "", debug),
	}
	if cfg == nil {
		return nil, fmt.Errorf("微信小程序配置为 nil")
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("微信小程序配置 AppID 为空")
	}
	if cfg.AppSecret == "" {
		return nil, fmt.Errorf("微信小程序配置 AppSecret 为空")
	}
	return cfg, nil
}

func getWechatOaConfig(rows []*DbSettingRow, debug bool) (*configStruct.WechatOaConfig, error) {
	cfg := &configStruct.WechatOaConfig{
		AppID:          GetValueFromRow(rows, "wechat", "oa", "app_id", "", debug),
		AppSecret:      GetValueFromRow(rows, "wechat", "oa", "app_secret", "", debug),
		Token:          GetValueFromRow(rows, "wechat", "oa", "token", "", debug),
		EncodingAESKey: GetValueFromRow(rows, "wechat", "oa", "encoding_aes_key", "", debug),
	}
	if cfg == nil {
		return nil, fmt.Errorf("微信公众号配置为 nil")
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("微信公众号配置 AppID 为空")
	}
	if cfg.AppSecret == "" {
		return nil, fmt.Errorf("微信公众号配置 AppSecret 为空")
	}
	return cfg, nil
}

func getWechatOpenConfig(rows []*DbSettingRow, debug bool) (*configStruct.WechatOpenConfig, error) {
	cfg := &configStruct.WechatOpenConfig{
		AppID:     GetValueFromRow(rows, "wechat", "open", "app_id", "", debug),
		AppSecret: GetValueFromRow(rows, "wechat", "open", "app_secret", "", debug),
	}
	if cfg == nil {
		return nil, fmt.Errorf("微信开放平台配置为 nil")
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("微信开放平台配置 AppID 为空")
	}
	if cfg.AppSecret == "" {
		return nil, fmt.Errorf("微信开放平台配置 AppSecret 为空")
	}
	return cfg, nil
}

func getWechatPayConfigV3(rows []*DbSettingRow, debug bool) (*configStruct.WechatPayConfigV3, error) {
	cfg := &configStruct.WechatPayConfigV3{
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
	if cfg == nil {
		return nil, fmt.Errorf("微信支付 V3 配置为 nil")
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("微信支付 V3 配置 AppID 为空")
	}
	if cfg.MchID == "" {
		return nil, fmt.Errorf("微信支付 V3 配置 MchID 为空")
	}
	return cfg, nil
}

func getWechatPayConfigV2(rows []*DbSettingRow, debug bool) (*configStruct.WechatPayConfigV2, error) {
	cfg := &configStruct.WechatPayConfigV2{
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
	if cfg == nil {
		return nil, fmt.Errorf("微信支付 V2 配置为 nil")
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("微信支付 V2 配置 AppID 为空")
	}
	if cfg.MchID == "" {
		return nil, fmt.Errorf("微信支付 V2 配置 MchID 为空")
	}
	return cfg, nil
}

func getAliPayConfig(rows []*DbSettingRow, debug bool) (*configStruct.AliPayConfig, error) {
	cfg := &configStruct.AliPayConfig{
		AppID:            GetValueFromRow(rows, "ali", "pay", "app_id", "", debug),
		PrivateKey:       GetValueFromRow(rows, "ali", "pay", "private_key", "", debug),
		NotifyURL:        GetValueFromRow(rows, "ali", "pay", "notify_url", "", debug),
		Debug:            GetValueFromRow(rows, "ali", "pay", "debug", "0", debug) == "1",
		AppCertPublicKey: GetValueFromRow(rows, "ali", "pay", "app_cert_public_key", "", debug),
		CertPublicKey:    GetValueFromRow(rows, "ali", "pay", "cert_public_key", "", debug),
		RootCert:         GetValueFromRow(rows, "ali", "pay", "root_cert", "", debug),
	}
	if cfg == nil {
		return nil, fmt.Errorf("支付宝配置为 nil")
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("支付宝配置 AppID 为空")
	}
	if cfg.PrivateKey == "" {
		return nil, fmt.Errorf("支付宝配置 PrivateKey 为空")
	}
	return cfg, nil
}

func getRedisConfig(rows []*DbSettingRow, debug bool) (*configStruct.RedisConfig, error) {
	cfg := &configStruct.RedisConfig{
		Host:        GetValueFromRow(rows, "redis", "host", "", "127.0.0.1", debug),
		Port:        GetValueFromRow(rows, "redis", "port", "", "6379", debug),
		Username:    GetValueFromRow(rows, "redis", "username", "", "", debug),
		Password:    GetValueFromRow(rows, "redis", "password", "", "", debug),
		IdleTimeout: typeHelper.Str2Int(GetValueFromRow(rows, "redis", "idle_timeout", "", "60", debug)),
		Database:    typeHelper.Str2Int(GetValueFromRow(rows, "redis", "database", "", "", debug)),
		MaxActive:   typeHelper.Str2Int(GetValueFromRow(rows, "redis", "max_active", "", "10", debug)),
		MaxIdle:     typeHelper.Str2Int(GetValueFromRow(rows, "redis", "max_idle", "", "10", debug)),
	}
	if cfg == nil {
		return nil, fmt.Errorf("Redis 配置为 nil")
	}
	if cfg.Host == "" {
		return nil, fmt.Errorf("Redis 配置 Host 为空")
	}
	if cfg.Port == "" {
		return nil, fmt.Errorf("Redis 配置 Port 为空")
	}
	return cfg, nil
}

func getEsConfig(rows []*DbSettingRow, debug bool) (*configStruct.EsConfig, error) {
	cfg := &configStruct.EsConfig{
		Host:     GetValueFromRow(rows, "es", "host", "", "http://127.0.0.1", debug),
		Port:     GetValueFromRow(rows, "es", "port", "", "9200", debug),
		Password: GetValueFromRow(rows, "es", "password", "", "123456", debug),
		Username: GetValueFromRow(rows, "es", "username", "", "es", debug),
	}
	if cfg == nil {
		return nil, fmt.Errorf("Elasticsearch 配置为 nil")
	}
	if cfg.Host == "" {
		return nil, fmt.Errorf("Elasticsearch 配置 Host 为空")
	}
	if cfg.Port == "" {
		return nil, fmt.Errorf("Elasticsearch 配置 Port 为空")
	}
	return cfg, nil
}

func getRabbitMQConfig(rows []*DbSettingRow, debug bool) (*configStruct.RabbitMQConfig, error) {
	cfg := &configStruct.RabbitMQConfig{
		Host:     GetValueFromRow(rows, "rabbit_mq", "host", "", "http://127.0.0.1", debug),
		Password: GetValueFromRow(rows, "rabbit_mq", "password", "", "123456", debug),
		Username: GetValueFromRow(rows, "rabbit_mq", "username", "", "root", debug),
	}
	if cfg == nil {
		return nil, fmt.Errorf("RabbitMQ 配置为 nil")
	}
	if cfg.Host == "" {
		return nil, fmt.Errorf("RabbitMQ 配置 Host 为空")
	}
	return cfg, nil
}

func getPostgresConfig(rows []*DbSettingRow, debug bool) (*configStruct.PostgresConfig, error) {
	dsn := GetValueFromRow(rows, "postgres", "", "dsn", "", debug)
	if dsn == "" {
		dsn = GetValueFromRow(rows, "postgres", "", "", "", debug)
	}
	if dsn == "" {
		return nil, nil // DSN 为空时返回 nil，这是允许的
	}

	cfg := &configStruct.PostgresConfig{DSN: dsn}
	if cfg == nil {
		return nil, fmt.Errorf("PostgreSQL 配置为 nil")
	}
	if cfg.DSN == "" {
		return nil, fmt.Errorf("PostgreSQL 配置 DSN 为空")
	}
	if v := GetValueFromRow(rows, "postgres", "", "conn_max_idle", "", debug); v != "" {
		cfg.ConnMaxIdle = typeHelper.Str2Int(v)
	}
	if v := GetValueFromRow(rows, "postgres", "", "conn_max_open", "", debug); v != "" {
		cfg.ConnMaxOpen = typeHelper.Str2Int(v)
	}
	if v := GetValueFromRow(rows, "postgres", "", "conn_max_lifetime", "", debug); v != "" {
		cfg.ConnMaxLifetime = time.Duration(typeHelper.Str2Int64(v)) * time.Second
	}
	return cfg, nil
}

func getAliOssConfig(rows []*DbSettingRow, debug bool) (*configStruct.AliOssConfig, error) {
	cfg := &configStruct.AliOssConfig{
		AccessKeyID:     GetValueFromRow(rows, "ali", "oss", "access_key_id", "", debug),
		AccessKeySecret: GetValueFromRow(rows, "ali", "oss", "access_key_secret", "", debug),
		Host:            GetValueFromRow(rows, "ali", "oss", "host", "", debug),
		EndPoint:        GetValueFromRow(rows, "ali", "oss", "end_point", "", debug),
		BucketName:      GetValueFromRow(rows, "ali", "oss", "bucket_name", "", debug),
		ExpireTime:      typeHelper.Str2Int64(GetValueFromRow(rows, "ali", "oss", "expire_time", "30", debug)),
		ARN:             GetValueFromRow(rows, "ali", "oss", "arn", "", debug),
	}
	if cfg == nil {
		return nil, fmt.Errorf("阿里云 OSS 配置为 nil")
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("阿里云 OSS 配置 AccessKeyID 为空")
	}
	if cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("阿里云 OSS 配置 AccessKeySecret 为空")
	}
	return cfg, nil
}

func getAliApiConfig(rows []*DbSettingRow, debug bool) (*configStruct.AliApiConfig, error) {
	cfg := &configStruct.AliApiConfig{
		AppKey:    GetValueFromRow(rows, "ali", "api", "app_key", "", debug),
		AppSecret: GetValueFromRow(rows, "ali", "api", "app_secret", "", debug),
		AppCode:   GetValueFromRow(rows, "ali", "api", "app_code", "", debug),
	}
	if cfg == nil {
		return nil, fmt.Errorf("阿里云 API 配置为 nil")
	}
	if cfg.AppKey == "" {
		return nil, fmt.Errorf("阿里云 API 配置 AppKey 为空")
	}
	if cfg.AppSecret == "" {
		return nil, fmt.Errorf("阿里云 API 配置 AppSecret 为空")
	}
	if cfg.AppCode == "" {
		return nil, fmt.Errorf("阿里云 API 配置 AppCode 为空")
	}
	return cfg, nil
}

func getAliSmsConfig(rows []*DbSettingRow, debug bool) (*configStruct.AliSmsConfig, error) {
	cfg := &configStruct.AliSmsConfig{
		AccessKeySecret: GetValueFromRow(rows, "ali", "sms", "access_key_secret", "", debug),
		AccessKeyID:     GetValueFromRow(rows, "ali", "sms", "access_key_id", "", debug),
	}
	if cfg == nil {
		return nil, fmt.Errorf("阿里云短信配置为 nil")
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("阿里云短信配置 AccessKeyID 为空")
	}
	if cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("阿里云短信配置 AccessKeySecret 为空")
	}
	return cfg, nil
}

func getAliIotConfig(rows []*DbSettingRow, debug bool) (*configStruct.AliIotConfig, error) {
	cfg := &configStruct.AliIotConfig{
		AccessKeySecret: GetValueFromRow(rows, "ali", "iot", "access_key_secret", "", debug),
		AccessKeyID:     GetValueFromRow(rows, "ali", "iot", "access_key_id", "", debug),
		EndPoint:        GetValueFromRow(rows, "ali", "iot", "end_point", "", debug),
		RegionID:        GetValueFromRow(rows, "ali", "iot", "region_id", "", debug),
	}
	if cfg == nil {
		return nil, fmt.Errorf("阿里云 IoT 配置为 nil")
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("阿里云 IoT 配置 AccessKeyID 为空")
	}
	if cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("阿里云 IoT 配置 AccessKeySecret 为空")
	}
	return cfg, nil
}

func getAmapConfig(rows []*DbSettingRow, debug bool) (*configStruct.AmapConfig, error) {
	cfg := &configStruct.AmapConfig{Key: GetValueFromRow(rows, "ali", "amap", "key", "", debug)}
	if cfg == nil {
		return nil, fmt.Errorf("高德地图配置为 nil")
	}
	if cfg.Key == "" {
		return nil, fmt.Errorf("高德地图配置 Key 为空")
	}
	return cfg, nil
}

func getYunXinConfig(rows []*DbSettingRow, debug bool) (*configStruct.YunxinConfig, error) {
	cfg := &configStruct.YunxinConfig{
		AppKey:    GetValueFromRow(rows, "netease", "yunxin", "app_key", "", debug),
		AppSecret: GetValueFromRow(rows, "netease", "yunxin", "app_secret", "", debug),
	}
	if cfg == nil {
		return nil, fmt.Errorf("网易云信配置为 nil")
	}
	if cfg.AppKey == "" {
		return nil, fmt.Errorf("网易云信配置 AppKey 为空")
	}
	if cfg.AppSecret == "" {
		return nil, fmt.Errorf("网易云信配置 AppSecret 为空")
	}
	return cfg, nil
}
