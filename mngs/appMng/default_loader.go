package appMng

import (
	"context"
	"time"

	"github.com/wiidz/goutil/helpers/configHelper"
	"github.com/wiidz/goutil/structs/configStruct"
)

// DefaultLoader 返回一个 Loader，它会读取当前工作目录下的配置文件
// (读取 ./configs/config.yaml)，并仅初始化 HttpConfig。
// 若找不到配置文件，则使用默认值（IP: 0.0.0.0，Port: 8080）。
func DefaultLoader() Loader {
	return LoaderFunc(func(ctx context.Context) (*LoaderResult, error) {
		viperData, _ := configHelper.GetViper(nil)
		httpCfg, err := configHelper.SimpleLoadHTTPConfig(viperData, nil)
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
