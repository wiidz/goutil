package appMng

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/wiidz/goutil/structs/configStruct"
)

// configErrorFactory 统一生成配置相关错误，减少重复文案
type configErrorFactory struct{}

func newConfigErrorFactory() *configErrorFactory {
	return &configErrorFactory{}
}

// missingField 缺失必填字段
func (f *configErrorFactory) missingField(key, field string) error {
	return fmt.Errorf("%s 配置 %s 为空", GetKeyDisplayName(key), field)
}

// databaseEmpty 从数据库加载到空结果
func (f *configErrorFactory) databaseEmpty(key string) error {
	return fmt.Errorf("从数据库加载 %s 配置失败: 数据为空", GetKeyDisplayName(key))
}

// yamlNotInit 未初始化 YAML
func (f *configErrorFactory) yamlNotInit(key string) error {
	return fmt.Errorf("从 YAML 加载 %s 配置失败: 未初始化 YAML 配置", GetKeyDisplayName(key))
}

// yamlLoadFailed YAML 解析失败
func (f *configErrorFactory) yamlLoadFailed(key string, err error) error {
	return fmt.Errorf("从 YAML 加载 %s 配置失败: %w", GetKeyDisplayName(key), err)
}

var errFactory = newConfigErrorFactory()

// key 对应 BaseConfig 赋值函数
var configAssigners = map[string]func(*configStruct.BaseConfig, interface{}){
	ConfigKeys.Redis.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.RedisConfig); ok {
			cfg.RedisConfig = val
		}
	},
	ConfigKeys.Es.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.EsConfig); ok {
			cfg.EsConfig = val
		}
	},
	ConfigKeys.RabbitMQ.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.RabbitMQConfig); ok {
			cfg.RabbitMQConfig = val
		}
	},
	ConfigKeys.Postgres.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.PostgresConfig); ok {
			cfg.PostgresConfig = val
		}
	},
	ConfigKeys.Mysql.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.MysqlConfig); ok {
			cfg.MysqlConfig = val
		}
	},
	ConfigKeys.WechatMini.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.WechatMiniConfig); ok {
			cfg.WechatMiniConfig = val
		}
	},
	ConfigKeys.WechatOa.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.WechatOaConfig); ok {
			cfg.WechatOaConfig = val
		}
	},
	ConfigKeys.WechatOpen.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.WechatOpenConfig); ok {
			cfg.WechatOpenConfig = val
		}
	},
	ConfigKeys.WechatPayV3.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.WechatPayConfigV3); ok {
			cfg.WechatPayConfigV3 = val
		}
	},
	ConfigKeys.WechatPayV2.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.WechatPayConfigV2); ok {
			cfg.WechatPayConfigV2 = val
		}
	},
	ConfigKeys.AliOss.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.AliOssConfig); ok {
			cfg.AliOssConfig = val
		}
	},
	ConfigKeys.AliPay.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.AliPayConfig); ok {
			cfg.AliPayConfig = val
		}
	},
	ConfigKeys.AliApi.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.AliApiConfig); ok {
			cfg.AliApiConfig = val
		}
	},
	ConfigKeys.AliSms.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.AliSmsConfig); ok {
			cfg.AliSmsConfig = val
		}
	},
	ConfigKeys.AliIot.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.AliIotConfig); ok {
			cfg.AliIotConfig = val
		}
	},
	ConfigKeys.Amap.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.AmapConfig); ok {
			cfg.AmapConfig = val
		}
	},
	ConfigKeys.Yunxin.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.YunxinConfig); ok {
			cfg.YunxinConfig = val
		}
	},
	ConfigKeys.Volcengine.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.VolcengineConfig); ok {
			cfg.VolcengineConfig = val
		}
	},
}

// key 对应的来源选择器
var configSources = map[string]func(*ConfigSourceStrategy) ConfigSource{
	ConfigKeys.Redis.Key:      func(s *ConfigSourceStrategy) ConfigSource { return s.Redis },
	ConfigKeys.Es.Key:         func(s *ConfigSourceStrategy) ConfigSource { return s.Es },
	ConfigKeys.RabbitMQ.Key:   func(s *ConfigSourceStrategy) ConfigSource { return s.RabbitMQ },
	ConfigKeys.Postgres.Key:   func(s *ConfigSourceStrategy) ConfigSource { return s.Postgres },
	ConfigKeys.Mysql.Key:      func(s *ConfigSourceStrategy) ConfigSource { return s.Mysql },
	ConfigKeys.WechatMini.Key: func(s *ConfigSourceStrategy) ConfigSource { return s.WechatMini },
	ConfigKeys.WechatOa.Key:   func(s *ConfigSourceStrategy) ConfigSource { return s.WechatOa },
	ConfigKeys.WechatOpen.Key: func(s *ConfigSourceStrategy) ConfigSource { return s.WechatOpen },
	ConfigKeys.WechatPayV3.Key: func(s *ConfigSourceStrategy) ConfigSource {
		return s.WechatPayV3
	},
	ConfigKeys.WechatPayV2.Key: func(s *ConfigSourceStrategy) ConfigSource {
		return s.WechatPayV2
	},
	ConfigKeys.AliOss.Key:     func(s *ConfigSourceStrategy) ConfigSource { return s.AliOss },
	ConfigKeys.AliPay.Key:     func(s *ConfigSourceStrategy) ConfigSource { return s.AliPay },
	ConfigKeys.AliApi.Key:     func(s *ConfigSourceStrategy) ConfigSource { return s.AliApi },
	ConfigKeys.AliSms.Key:     func(s *ConfigSourceStrategy) ConfigSource { return s.AliSms },
	ConfigKeys.AliIot.Key:     func(s *ConfigSourceStrategy) ConfigSource { return s.AliIot },
	ConfigKeys.Amap.Key:       func(s *ConfigSourceStrategy) ConfigSource { return s.Amap },
	ConfigKeys.Yunxin.Key:     func(s *ConfigSourceStrategy) ConfigSource { return s.Yunxin },
	ConfigKeys.Volcengine.Key: func(s *ConfigSourceStrategy) ConfigSource { return s.Volcengine },
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
		No:      GetValueFromRow(rows, ConfigKeys.App.Key, "", "no", "", false),
		Name:    GetValueFromRow(rows, ConfigKeys.App.Key, "", "name", "", false),
		Debug:   GetValueFromRow(rows, ConfigKeys.App.Key, "", "debug", "", false) == "1",
		Version: GetValueFromRow(rows, ConfigKeys.App.Key, "", "version", "", false),
	}
}

func getHttpServerConfig(rows []*DbSettingRow, serverLabel string) *configStruct.HttpServerConfig {
	return &configStruct.HttpServerConfig{
		Label:  GetValueFromRow(rows, ConfigKeys.HttpServer.Key, serverLabel, "label", "", false),
		Host:   GetValueFromRow(rows, ConfigKeys.HttpServer.Key, serverLabel, "host", "", false),
		Port:   GetValueFromRow(rows, ConfigKeys.HttpServer.Key, serverLabel, "port", "", false),
		Domain: GetValueFromRow(rows, ConfigKeys.HttpServer.Key, serverLabel, "domain", "", false),
	}
}
func getLocationConfig(rows []*DbSettingRow) (location *time.Location, err error) {
	timeZone := GetValueFromRow(rows, ConfigKeys.TimeZone.Key, "", "", "Asia/Shanghai", false)
	location, err = time.LoadLocation(timeZone)
	if err != nil {
		location = time.FixedZone("CST-8", 8*3600)
	}
	return
}

// ConfigMap 以键对应结构体定义（默认值从 default tag 读取）
type ConfigMap struct {
	Key  ConfigKey
	Data interface{}
}

// 所有配置的结构映射
var configMaps = map[string]ConfigMap{
	ConfigKeys.Redis.Key:       {Key: ConfigKeys.Redis, Data: configStruct.RedisConfig{}},
	ConfigKeys.Es.Key:          {Key: ConfigKeys.Es, Data: configStruct.EsConfig{}},
	ConfigKeys.RabbitMQ.Key:    {Key: ConfigKeys.RabbitMQ, Data: configStruct.RabbitMQConfig{}},
	ConfigKeys.Postgres.Key:    {Key: ConfigKeys.Postgres, Data: configStruct.PostgresConfig{}},
	ConfigKeys.Mysql.Key:       {Key: ConfigKeys.Mysql, Data: configStruct.MysqlConfig{}},
	ConfigKeys.WechatMini.Key:  {Key: ConfigKeys.WechatMini, Data: configStruct.WechatMiniConfig{}},
	ConfigKeys.WechatOa.Key:    {Key: ConfigKeys.WechatOa, Data: configStruct.WechatOaConfig{}},
	ConfigKeys.WechatOpen.Key:  {Key: ConfigKeys.WechatOpen, Data: configStruct.WechatOpenConfig{}},
	ConfigKeys.WechatPayV3.Key: {Key: ConfigKeys.WechatPayV3, Data: configStruct.WechatPayConfigV3{}},
	ConfigKeys.WechatPayV2.Key: {Key: ConfigKeys.WechatPayV2, Data: configStruct.WechatPayConfigV2{}},
	ConfigKeys.AliOss.Key:      {Key: ConfigKeys.AliOss, Data: configStruct.AliOssConfig{}},
	ConfigKeys.AliPay.Key:      {Key: ConfigKeys.AliPay, Data: configStruct.AliPayConfig{}},
	ConfigKeys.AliApi.Key:      {Key: ConfigKeys.AliApi, Data: configStruct.AliApiConfig{}},
	ConfigKeys.AliSms.Key:      {Key: ConfigKeys.AliSms, Data: configStruct.AliSmsConfig{}},
	ConfigKeys.AliIot.Key:      {Key: ConfigKeys.AliIot, Data: configStruct.AliIotConfig{}},
	ConfigKeys.Amap.Key:        {Key: ConfigKeys.Amap, Data: configStruct.AmapConfig{}},
	ConfigKeys.Yunxin.Key:      {Key: ConfigKeys.Yunxin, Data: configStruct.YunxinConfig{}},
	ConfigKeys.Volcengine.Key:  {Key: ConfigKeys.Volcengine, Data: configStruct.VolcengineConfig{}},
}

// validateConfig 使用 validator 库验证配置结构体
func validateConfig(target interface{}, configKey string) error {
	if target == nil {
		return nil
	}

	// 创建 validator 实例
	validate := validator.New()

	// 验证结构体
	if err := validate.Struct(target); err != nil {
		// 格式化验证错误信息
		var errMsgs []string
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				errMsgs = append(errMsgs, fmt.Sprintf("字段 %s: %s", fieldError.Field(), getValidationErrorMsg(fieldError)))
			}
		} else {
			errMsgs = append(errMsgs, err.Error())
		}
		return fmt.Errorf("配置 %s 验证失败: %s", configKey, strings.Join(errMsgs, "; "))
	}

	return nil
}

// getValidationErrorMsg 获取验证错误的友好提示信息
func getValidationErrorMsg(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "必填字段不能为空"
	case "min":
		return fmt.Sprintf("值不能小于 %s", fieldError.Param())
	case "max":
		return fmt.Sprintf("值不能大于 %s", fieldError.Param())
	case "len":
		return fmt.Sprintf("长度必须等于 %s", fieldError.Param())
	case "email":
		return "必须是有效的邮箱地址"
	case "url":
		return "必须是有效的URL地址"
	case "ip":
		return "必须是有效的IP地址"
	case "oneof":
		return fmt.Sprintf("值必须是以下之一: %s", fieldError.Param())
	default:
		return fmt.Sprintf("验证失败 (tag: %s, param: %s)", fieldError.Tag(), fieldError.Param())
	}
}

// applyDefaultsFromTags 根据 struct 的 default tag 填充零值字段
func applyDefaultsFromTags(target interface{}) {
	if target == nil {
		return
	}
	val := reflect.ValueOf(target)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		ft := typ.Field(i)
		if !field.CanSet() || !field.IsZero() {
			continue
		}
		def := ft.Tag.Get("default")
		if def == "" {
			continue
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(def)
		case reflect.Bool:
			if v, err := strconv.ParseBool(def); err == nil {
				field.SetBool(v)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v, err := strconv.ParseInt(def, 10, 64); err == nil {
				field.SetInt(v)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if v, err := strconv.ParseUint(def, 10, 64); err == nil {
				field.SetUint(v)
			}
		}
	}
}

// fillConfigFromRows 使用 flag1=parentKey、flag2=字段 json/mapstructure 标签，从 rows 填充配置，并按 default tag 设置默认值
func fillConfigFromRows(target interface{}, nameKey, flag1 string, rows []*DbSettingRow, debug bool) error {
	if target == nil {
		return fmt.Errorf("target 不能为空")
	}
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("target 必须是非 nil 指针")
	}
	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("target 必须指向结构体")
	}

	typ := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		ft := typ.Field(i)
		if !field.CanSet() {
			continue
		}

		jsonTag := ft.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = ft.Tag.Get("mapstructure")
		}
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		jsonTag = strings.Split(jsonTag, ",")[0]

		defVal := ft.Tag.Get("default")
		raw := GetValueFromRow(rows, nameKey, flag1, jsonTag, defVal, debug)
		if raw == "" && defVal != "" {
			raw = defVal
		}
		if raw == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(raw)
		case reflect.Bool:
			if v, err := strconv.ParseBool(raw); err == nil {
				field.SetBool(v)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v, err := strconv.ParseInt(raw, 10, 64); err == nil {
				field.SetInt(v)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if v, err := strconv.ParseUint(raw, 10, 64); err == nil {
				field.SetUint(v)
			}
		case reflect.Float32, reflect.Float64:
			if v, err := strconv.ParseFloat(raw, 64); err == nil {
				field.SetFloat(v)
			}
		}
	}

	applyDefaultsFromTags(target)
	return nil
}
