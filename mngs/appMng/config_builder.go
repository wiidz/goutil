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

// assignConfigToBaseConfig 使用反射将配置值赋值到 BaseConfig 的对应字段
// key 直接作为字段名使用（首字母大写）
func assignConfigToBaseConfig(cfg *configStruct.BaseConfig, key string, value interface{}) {
	if key == "" {
		return
	}
	cfgVal := reflect.ValueOf(cfg).Elem()
	field := cfgVal.FieldByName(key) // key 直接作为字段名
	if !field.IsValid() || !field.CanSet() {
		return
	}
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		field.Set(val)
	} else {
		field.Set(val.Addr())
	}
}

// getConfigSource 使用反射从 ConfigSourceStrategy 获取配置来源
// key 直接作为字段名使用（首字母大写）
func getConfigSource(strategy *ConfigSourceStrategy, key string) ConfigSource {
	if key == "" {
		return ""
	}
	strategyVal := reflect.ValueOf(strategy).Elem()
	field := strategyVal.FieldByName(key) // key 直接作为字段名
	if !field.IsValid() {
		return ""
	}
	return field.Interface().(ConfigSource)
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
		ID:      GetValueFromRow(rows, ConfigKeys.App.Key, "", "id", "", false),
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

// 所有配置的类型映射（用于反射创建实例，仅包含基础配置）
// configTypes 所有配置的类型映射（通过 ConfigKey 自动生成）
var configTypes = initConfigTypes()

// initConfigTypes 从 ConfigKeys 中提取类型映射
func initConfigTypes() map[string]reflect.Type {
	result := make(map[string]reflect.Type)
	configKeys := []ConfigKey{
		ConfigKeys.Redis,
		ConfigKeys.Es,
		ConfigKeys.Rabbitmq,
		ConfigKeys.Postgres,
		ConfigKeys.Mysql,
	}

	// 通过反射从 BaseConfig 中获取字段类型
	baseConfigType := reflect.TypeOf(configStruct.BaseConfig{})
	fieldTypeMap := make(map[string]reflect.Type)
	for i := 0; i < baseConfigType.NumField(); i++ {
		field := baseConfigType.Field(i)
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		fieldTypeMap[field.Name] = fieldType
	}

	// 根据 ConfigKey 的 Key（直接作为字段名）查找类型
	for _, ck := range configKeys {
		if ck.Key != "" {
			if fieldType, ok := fieldTypeMap[ck.Key]; ok {
				result[ck.Key] = fieldType
			}
		}
	}

	return result
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
		return errFactory.validateFailed(configKey, errMsgs)
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

// loadConfigFromSource 从指定来源加载配置（公共逻辑）
// source: 配置来源（SourceDatabase 或 SourceYAML）
// nameKey: 配置名称（用于数据库 name 字段和错误信息）
// configKey: 配置键（用于 YAML 键名和数据库 flag1）
// targetPtr: 目标配置指针
// configPool: 配置池
// debug: 调试模式
func loadConfigFromSource(source ConfigSource, nameKey, configKey string, targetPtr interface{}, configPool *ConfigPool, debug bool) error {
	switch source {
	case SourceDatabase:
		dbRows := configPool.GetDBRows()
		if len(dbRows) == 0 {
			return errFactory.databaseEmpty(nameKey)
		}
		if err := fillConfigFromRows(targetPtr, nameKey, configKey, dbRows, debug); err != nil {
			return errFactory.databaseLoadFailed(nameKey, err)
		}
		// 应用默认值
		applyDefaultsFromTags(targetPtr)
		// 验证配置
		if err := validateConfig(targetPtr, nameKey); err != nil {
			return err
		}

	case SourceYAML:
		if configPool == nil || len(configPool.GetYAML()) == 0 {
			return errFactory.yamlNotInit(nameKey)
		}
		// 如果 configKey 为空，使用 nameKey
		yamlKey := configKey
		if yamlKey == "" {
			yamlKey = nameKey
		}
		if err := configPool.GetYAML()[0].UnmarshalKey(yamlKey, targetPtr); err != nil {
			return errFactory.yamlLoadFailed(nameKey, err)
		}
		// 应用默认值
		applyDefaultsFromTags(targetPtr)
		// 验证配置
		if err := validateConfig(targetPtr, nameKey); err != nil {
			return err
		}

	default:
		return errFactory.unsupportedSource(source)
	}

	return nil
}

// fillConfigFromRows 使用 flag1=parentKey、flag2=字段 json/mapstructure 标签，从 rows 填充配置，并按 default tag 设置默认值
func fillConfigFromRows(target interface{}, nameKey, flag1 string, rows []*DbSettingRow, debug bool) error {
	if target == nil {
		return errFactory.targetEmpty()
	}
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errFactory.targetNotStructPointer()
	}
	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return errFactory.targetNotStruct()
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
