package appMng

import (
	"context"

	"github.com/wiidz/goutil/mngs/esMng"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"github.com/wiidz/goutil/mngs/psqlMng"
	"github.com/wiidz/goutil/mngs/redisMng"
	"github.com/wiidz/goutil/structs/configStruct"
)

// NewApp 直接构建一个 AppMng。
func NewApp(ctx context.Context, configPool *ConfigPool, baseBuilder ConfigBuilder, projectBuilder ProjectConfig) (appMng *AppMng, err error) {

	//【1】构建 base 配置
	var baseCfg *configStruct.BaseConfig
	baseCfg, err = baseBuilder.Build(configPool)
	if err != nil {
		return nil, err
	}

	appMng = &AppMng{
		ID:         baseCfg.Profile.Name,
		BaseConfig: baseCfg,
	}

	// 【3】初始化依赖
	if appMng.BaseConfig.Mysql != nil && appMng.Repos.Mysql == nil {
		appMng.Repos.Mysql, err = mysqlMng.NewMysqlMng(appMng.BaseConfig.Mysql, nil)
		if err != nil {
			return
		}
	}

	if appMng.BaseConfig.Postgres != nil && appMng.Repos.Postgres == nil {
		appMng.Repos.Postgres, err = psqlMng.NewMng(&psqlMng.Config{
			DSN:             appMng.BaseConfig.Postgres.DSN,
			ConnMaxIdle:     appMng.BaseConfig.Postgres.ConnMaxIdle,
			ConnMaxOpen:     appMng.BaseConfig.Postgres.ConnMaxOpen,
			ConnMaxLifetime: appMng.BaseConfig.Postgres.ConnMaxLifetime,
		})
		if err != nil {
			err = errFactory.postgresInitFailed(err)
			return
		}
	}

	if appMng.BaseConfig.Redis != nil {
		if appMng.Repos.Redis, err = redisMng.NewRedisMng(ctx, appMng.BaseConfig.Redis); err != nil {
			err = errFactory.redisInitFailed(err)
			return
		}
	}
	if appMng.BaseConfig.Es != nil {
		if err = esMng.Init(appMng.BaseConfig.Es); err != nil {
			err = errFactory.esInitFailed(err)
			return
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
		appMng.ProjectConfig = projectBuilder
		if err = appMng.ProjectConfig.Build(appMng.BaseConfig, configPool); err != nil {
			err = errFactory.projectBuildFailed(err)
			return
		}
	}

	return
}
