package appMng

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/wiidz/goutil/structs/configStruct"
)

// DefaultLoader 返回一个 Loader，它会读取当前工作目录下的配置文件
// (读取 ./configs/config.yaml)，并仅初始化 HttpConfig。
// 若找不到配置文件，则使用默认值（IP: 0.0.0.0，Port: 8080）。
func DefaultLoader() Loader {
	return LoaderFunc(func(ctx context.Context) (*LoaderResult, error) {
		httpCfg, err := SimpleLoadHTTPConfig(nil)
		if err != nil {
			return nil, err
		}

		profile := &configStruct.AppProfile{}
		if httpCfg != nil {
			profile.Host = httpCfg.IP
			profile.Port = httpCfg.Port
		}

		return &LoaderResult{
			BaseConfig: &configStruct.BaseConfig{
				Profile:  profile,
				Location: time.Local,
			},
		}, nil
	})
}

// SimpleLoadHTTPConfig 简单读取http配置
func SimpleLoadHTTPConfig(data *configStruct.ViperConfig) (*configStruct.HttpConfig, error) {
	v := viper.New()

	if data == nil {
		data = &configStruct.ViperConfig{
			DirPath:  "./configs",
			FileName: "config",
			FileType: "yaml",
		}
	}

	v.SetConfigName(data.FileName)
	v.SetConfigType(data.FileType)
	v.AddConfigPath(data.DirPath)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("http.ip", "0.0.0.0")
	v.SetDefault("http.port", "8080")

	if err := v.ReadInConfig(); err != nil {
		var nf viper.ConfigFileNotFoundError
		if !errors.As(err, &nf) {
			return nil, fmt.Errorf("appMng: failed to read config file: %w", err)
		}
	}

	var cfg struct {
		HTTP configStruct.HttpConfig `mapstructure:"http"`
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("appMng: failed to unmarshal config: %w", err)
	}

	httpCfg := cfg.HTTP
	return &httpCfg, nil
}
