package appMng

import (
	"context"
	"fmt"

	"github.com/wiidz/goutil/mngs/esMng"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"github.com/wiidz/goutil/mngs/psqlMng"
	"github.com/wiidz/goutil/mngs/redisMng"
)

// NewApp 直接构建一个 AppMng。
func NewApp(ctx context.Context, configPool *ConfigPool, baseBuilder ConfigBuilder, projectBuilder ProjectConfig) (*AppMng, error) {

	//【1】构建 base 配置
	baseCfg, err := baseBuilder.Build(configPool)
	if err != nil {
		return nil, err
	}

	app := &AppMng{
		ID:         baseCfg.Profile.Name,
		BaseConfig: baseCfg,
	}

	// 【3】初始化依赖
	if app.BaseConfig.MysqlConfig != nil && app.Repos.Mysql == nil {
		app.Repos.Mysql, err = mysqlMng.NewMysqlMng(app.BaseConfig.MysqlConfig, nil)
	}

	if app.BaseConfig.PostgresConfig != nil && app.Repos.Postgres == nil {
		app.Repos.Postgres, err = psqlMng.NewMng(&psqlMng.Config{
			DSN:             app.BaseConfig.PostgresConfig.DSN,
			ConnMaxIdle:     app.BaseConfig.PostgresConfig.ConnMaxIdle,
			ConnMaxOpen:     app.BaseConfig.PostgresConfig.ConnMaxOpen,
			ConnMaxLifetime: app.BaseConfig.PostgresConfig.ConnMaxLifetime,
		})
		if err != nil {
			return nil, fmt.Errorf("appMng: init postgres failed: %w", err)
		}
	}

	if app.BaseConfig.RedisConfig != nil {
		if app.Repos.Redis, err = redisMng.NewRedisMng(ctx, app.BaseConfig.RedisConfig); err != nil {
			return nil, fmt.Errorf("appMng: init redis failed: %w", err)
		}
	}
	if app.BaseConfig.EsConfig != nil {
		if err = esMng.Init(app.BaseConfig.EsConfig); err != nil {
			return nil, fmt.Errorf("appMng: init es failed: %w", err)
		}
	}

	// if app.BaseConfig.RabbitMQConfig != nil {
	// 	if err = amqpMng.Init(app.BaseConfig.RabbitMQConfig); err != nil {
	// 		return nil, fmt.Errorf("appMng: init rabbitmq failed: %w", err)
	// 	}
	// 	app.Repos.RabbitMQ, err = amqpMng.NewRabbitMQ(app.BaseConfig.RabbitMQConfig)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("appMng: init rabbitmq failed: %w", err)
	// 	}
	// }

	// 【4】项目级配置构建

	// 设置 ProjectConfig（如果提供了）
	if projectBuilder != nil {
		app.ProjectConfig = projectBuilder
		if err = app.ProjectConfig.Build(app.BaseConfig, configPool); err != nil {
			return nil, fmt.Errorf("appMng: project build failed: %w", err)
		}
	}

	return app, nil
}
