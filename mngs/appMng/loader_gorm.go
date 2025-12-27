package appMng

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// NewGormSettingLoader 基于 gorm Dialector 的装载器，多用于 MySQL/Postgres。
// enrich 允许对生成的 BaseConfig 做额外加工（可选）。
func NewGormSettingLoader(dialector gorm.Dialector, tableName string, enrich func(*LoaderResult) error) Loader {
	return LoaderFunc(func(ctx context.Context) (*LoaderResult, error) {
		db, err := gorm.Open(dialector, &gorm.Config{})
		if err != nil {
			return nil, err
		}
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.SetMaxIdleConns(2)
			sqlDB.SetMaxOpenConns(4)
			sqlDB.SetConnMaxLifetime(5 * time.Minute)
			defer sqlDB.Close()
		}

		var rows []*DbSettingRow
		if err := db.WithContext(ctx).Table(tableName).Where("belonging = ?", "system").Find(&rows).Error; err != nil {
			return nil, err
		}

		baseConfig, err := buildBaseConfig(rows, nil) // nil 表示使用默认选项（加载所有配置）
		if err != nil {
			return nil, err
		}
		res := &LoaderResult{BaseConfig: baseConfig}
		if enrich != nil {
			if err := enrich(res); err != nil {
				return nil, err
			}
		}
		return res, nil
	})
}
