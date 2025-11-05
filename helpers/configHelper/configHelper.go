package configHelper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/wiidz/goutil/structs/configStruct"
)

// GetViper 获取viper配置对象
func GetViper(data *configStruct.ViperConfig) (viperData *viper.Viper, err error) {
	viperData = viper.New()

	viperData.SetConfigName(data.FileName)
	viperData.SetConfigType(data.FileType)
	viperData.AddConfigPath(data.DirPath)

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

// SimpleLoadHTTPConfig 简单读取http配置
func SimpleLoadHTTPConfig(viperData *viper.Viper) (*configStruct.HttpConfig, error) {
	v := viper.New()

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

	var cfg struct {
		HTTP configStruct.HttpConfig `mapstructure:"http"`
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("appMng: failed to unmarshal config: %w", err)
	}

	httpCfg := cfg.HTTP
	return &httpCfg, nil
}

// SimpleLoadRepoConfig 数据仓库配置
func SimpleLoadRepoConfig(viperData *viper.Viper, dbType string) (*configStruct.RepoConfig, error) {

	if viperData == nil {
		viperData, _ = GetViper(nil)
	}

	viperData.SetDefault(dbType+".dsn", "")
	viperData.SetDefault(dbType+"auto_migrate", "false")

	if err := viperData.ReadInConfig(); err != nil {
		var nf viper.ConfigFileNotFoundError
		if !errors.As(err, &nf) {
			return nil, fmt.Errorf("appMng: failed to read config file: %w", err)
		}
	}

	var cfg struct {
		Repo configStruct.RepoConfig `mapstructure:"repo"`
	}
	if err := viperData.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("appMng: failed to unmarshal config: %w", err)
	}

	repoCfg := cfg.Repo
	return &repoCfg, nil
}
