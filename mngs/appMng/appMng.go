package appMng

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/wiidz/goutil/mngs/amqpMng"
	"github.com/wiidz/goutil/mngs/esMng"
	"github.com/wiidz/goutil/mngs/identityMng"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"github.com/wiidz/goutil/mngs/psqlMng"
	"github.com/wiidz/goutil/mngs/redisMng"
	"github.com/wiidz/goutil/structs/configStruct"
)

// NewApp 直接构建一个 AppMng。
func NewApp(ctx context.Context, loggerBuilder AppLogger, configPool *ConfigPool, baseBuilder ConfigBuilder, projectBuilder ProjectConfig) (appMng *AppMng, err error) {
	// 记录启动开始时间
	startTime := time.Now()

	//【0】构建日志记录器
	appMng.Log = loggerBuilder
	err = appMng.Log.Build()
	if err != nil {
		return
	}
	log.Printf("✅成功: 日志记录器已构建完成")

	//【1】构建 base 配置
	var baseCfg *configStruct.BaseConfig
	baseCfg, err = baseBuilder.Build(configPool)
	if err != nil {
		return
	}
	log.Printf("✅成功: 基础配置已构建完成")

	appMng = &AppMng{
		ID:         baseCfg.Profile.ID,
		BaseConfig: baseCfg,
	}

	// 【3】初始化依赖
	if appMng.BaseConfig.Mysql != nil && appMng.Repos.Mysql == nil {
		appMng.Repos.Mysql, err = mysqlMng.NewMysqlMng(appMng.BaseConfig.Mysql, nil)
		if err != nil {
			return
		}
		log.Printf("✅成功: MySQL 数据库连接已初始化")
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
		log.Printf("✅成功: PostgreSQL 数据库连接已初始化")
	}

	if appMng.BaseConfig.Redis != nil {
		if appMng.Repos.Redis, err = redisMng.NewRedisMng(ctx, appMng.BaseConfig.Redis); err != nil {
			err = errFactory.redisInitFailed(err)
			return
		}
		log.Printf("✅成功: Redis 连接已初始化")
	}

	if appMng.BaseConfig.Es != nil {
		if err = esMng.Init(appMng.BaseConfig.Es); err != nil {
			err = errFactory.esInitFailed(err)
			return
		}
		log.Printf("✅成功: Elasticsearch 连接已初始化")
	}

	if appMng.BaseConfig.Rabbitmq != nil && appMng.Repos.RabbitMQ == nil {
		// 初始化 RabbitMQ 连接
		if err = amqpMng.Init(appMng.BaseConfig.Rabbitmq); err != nil {
			err = errFactory.rabbitmqInitFailed(err)
			return
		}
		// 创建最小化的 Config（只包含连接信息，其他配置在使用时再设置）
		amqpConfig := &amqpMng.Config{
			Host:     appMng.BaseConfig.Rabbitmq.Host,
			Username: appMng.BaseConfig.Rabbitmq.Username,
			Password: appMng.BaseConfig.Rabbitmq.Password,
		}
		appMng.Repos.RabbitMQ, err = amqpMng.NewRabbitMQ(amqpConfig)
		if err != nil {
			err = errFactory.rabbitmqInitFailed(err)
			return
		}
		log.Printf("✅成功: RabbitMQ 连接已初始化")
	}

	// 【4】项目级配置构建

	// 设置 ProjectConfig（如果提供了）
	if projectBuilder != nil {
		appMng.ProjectConfig = projectBuilder
		if err = appMng.ProjectConfig.Build(appMng.BaseConfig, configPool); err != nil {
			err = errFactory.projectBuildFailed(err)
			return
		}
		log.Printf("✅成功: 项目配置已构建完成")
	}

	// 计算启动耗时
	elapsed := time.Since(startTime)
	log.Printf("✅成功: 应用初始化完成 (ID: %s, 名称: %s, 版本: %s, 耗时: %v)", appMng.ID, appMng.BaseConfig.Profile.Name, appMng.BaseConfig.Profile.Version, elapsed)
	return
}

func (mng *AppMng) InitIdentityMng(config *identityMng.Config) (err error) {

	if config.StorageType == "redis" {
		if mng.Repos.Redis.Client == nil {
			err = errors.New("redis client is nil")
			return
		}
		config.RedisClient = mng.Repos.Redis.Client
	}

	mng.IdMng, err = identityMng.NewMng(config)

	if err != nil {
		log.Printf("❌错误: 身份管理器构建失败")
	} else {
		log.Printf("✅成功: 身份管理器已构建完成")
	}

	return
}
