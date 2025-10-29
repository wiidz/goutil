package appMng

import (
	"context"
	"errors"

	"github.com/wiidz/goutil/structs/configStruct"
)

// StaticLoader 返回一个使用固定 BaseConfig 的 Loader。
func StaticLoader(cfg *configStruct.BaseConfig) Loader {
	return LoaderFunc(func(ctx context.Context) (*Result, error) {
		if cfg == nil {
			return nil, errors.New("appMng: static loader requires base config")
		}
		return &Result{BaseConfig: cfg}, nil
	})
}
