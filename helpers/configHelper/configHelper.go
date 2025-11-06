package configHelper

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/wiidz/goutil/structs/configStruct"
)

// GetViper 获取viper配置对象
func GetViper(data *configStruct.ViperConfig) (viperData *viper.Viper, err error) {
	viperData = viper.New()

	if data == nil {
		data = &configStruct.ViperConfig{
			DirPath:  "./configs",
			FileName: "config",
			FileType: "yaml",
		}
	}

	if data.FileName != "" {
		viperData.SetConfigName(data.FileName)
	}
	if data.FileType != "" {
		viperData.SetConfigType(data.FileType)
	}
	if data.DirPath != "" {
		viperData.AddConfigPath(data.DirPath)
	}

	viperData.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperData.AutomaticEnv()

	if err = viperData.ReadInConfig(); err != nil {
		var nf viper.ConfigFileNotFoundError
		if !errors.As(err, &nf) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	return
}

// LoadConfig 将配置文件内容解码到目标结构体。
// - viperData 允许复用已配置的 viper 实例；传 nil 时会通过 data 创建新的实例。
// - data 指定配置文件信息；当 viperData 非 nil 时可为 nil。
// 返回使用中的 viper 实例，便于后续继续读取其它键。
func LoadConfig(viperData *viper.Viper, target interface{}) error {
	if viperData == nil {
		return errors.New("configHelper: viper instance is required")
	}
	if target == nil {
		return errors.New("configHelper: target struct is required")
	}

	if err := applyDefaultsFromStruct(viperData, reflect.TypeOf(target), ""); err != nil {
		return err
	}

	if err := viperData.ReadInConfig(); err != nil {
		var nf viper.ConfigFileNotFoundError
		if !errors.As(err, &nf) {
			return fmt.Errorf("configHelper: failed to read config file: %w", err)
		}
	}

	if err := viperData.Unmarshal(target); err != nil {
		return fmt.Errorf("configHelper: failed to unmarshal config: %w", err)
	}
	return nil
}

func SimpleLoadConfig(data *configStruct.ViperConfig, target interface{}) error {
	viperData, err := GetViper(data)
	if err != nil {
		return err
	}
	return LoadConfig(viperData, target)
}

func applyDefaultsFromStruct(v *viper.Viper, typ reflect.Type, prefix string) error {
	if typ == nil {
		return nil
	}

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.PkgPath != "" { // unexported
			continue
		}

		// Determine key name
		tag := field.Tag.Get("mapstructure")
		tagName := strings.ToLower(field.Name)
		if tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] == "-" {
				continue
			}
			if parts[0] != "" {
				tagName = parts[0]
			}
		} else {
			tagName = strings.ToLower(field.Name)
		}

		fullKey := tagName
		if prefix != "" {
			if tagName != "" {
				fullKey = prefix + "." + tagName
			} else {
				fullKey = prefix
			}
		}

		// Apply default if present
		if defVal := field.Tag.Get("default"); defVal != "" {
			if !v.IsSet(fullKey) {
				parsed, err := parseDefaultValue(field.Type, defVal)
				if err != nil {
					return fmt.Errorf("configHelper: parse default for %s: %w", fullKey, err)
				}
				v.SetDefault(fullKey, parsed)
			}
		}

		// Recurse into nested structs
		nestedType := field.Type
		for nestedType.Kind() == reflect.Ptr {
			nestedType = nestedType.Elem()
		}

		if nestedType.Kind() == reflect.Struct {
			nextPrefix := fullKey
			if field.Anonymous && (tag == "" || strings.Split(tag, ",")[0] == "") {
				nextPrefix = prefix
			}
			if err := applyDefaultsFromStruct(v, nestedType, nextPrefix); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseDefaultValue(typ reflect.Type, val string) (interface{}, error) {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	switch typ.Kind() {
	case reflect.String:
		return val, nil
	case reflect.Bool:
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		return parsed, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ.PkgPath() == "time" && typ.Name() == "Duration" {
			d, err := time.ParseDuration(val)
			if err != nil {
				return nil, err
			}
			return d, nil
		}
		parsed, err := strconv.ParseInt(val, 10, typ.Bits())
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(parsed).Convert(typ).Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		parsed, err := strconv.ParseUint(val, 10, typ.Bits())
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(parsed).Convert(typ).Interface(), nil
	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(val, typ.Bits())
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(parsed).Convert(typ).Interface(), nil
	}

	return val, nil
}

// SimpleLoadHTTPConfig 简单读取http配置
func SimpleLoadHTTPConfig(viperData *viper.Viper) (*configStruct.HttpConfig, error) {

	if viperData == nil {
		viperData, _ = GetViper(nil)
	}

	viperData.SetDefault("http.ip", "0.0.0.0")
	viperData.SetDefault("http.port", "8080")

	if err := viperData.ReadInConfig(); err != nil {
		var nf viper.ConfigFileNotFoundError
		if !errors.As(err, &nf) {
			return nil, fmt.Errorf("appMng: failed to read config file: %w", err)
		}
	}

	var httpCfg configStruct.HttpConfig
	if err := viperData.UnmarshalKey("http", &httpCfg); err != nil {
		return nil, err
	}
	return &httpCfg, nil
}

// SimpleLoadRepoConfig 数据仓库配置
func SimpleLoadRepoConfig(viperData *viper.Viper, dbType string) (*configStruct.RepoConfig, error) {
	if viperData == nil {
		viperData, _ = GetViper(nil)
	}

	viperData.SetDefault(fmt.Sprintf("%s.dsn", dbType), "")
	viperData.SetDefault(fmt.Sprintf("%s.auto_migrate", dbType), false)

	if err := viperData.ReadInConfig(); err != nil {
		var nf viper.ConfigFileNotFoundError
		if !errors.As(err, &nf) {
			return nil, fmt.Errorf("appMng: failed to read config file: %w", err)
		}
	}

	section := viperData.Sub(dbType)
	if section == nil {
		return nil, fmt.Errorf("appMng: config section %q not found", dbType)
	}

	var repoCfg configStruct.RepoConfig
	if err := section.Unmarshal(&repoCfg); err != nil {
		return nil, fmt.Errorf("appMng: failed to unmarshal %s config: %w", dbType, err)
	}
	return &repoCfg, nil
}
