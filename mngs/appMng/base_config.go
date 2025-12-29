package appMng

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/wiidz/goutil/structs/configStruct"
)

// ConfigStrategyProvider 配置策略提供者接口
// 提供获取配置来源策略的方法
type ConfigStrategyProvider interface {
	GetStrategy() *ConfigSourceStrategy
}

// ConfigBuilder 基础配置构建器接口
// 负责构建 BaseConfig
type ConfigBuilder interface {
	ConfigStrategyProvider
	Build(configPool *ConfigPool) (config *configStruct.BaseConfig, err error) // 构建参数
}

// ProjectConfig 项目配置接口
// 负责构建项目特定的配置
type ProjectConfig interface {
	ConfigStrategyProvider
	Build(baseConfig *configStruct.BaseConfig, configPool *ConfigPool) error // 构建参数
}

// NewBaseConfigBuilder 创建基础构建器
// configPool: 配置池（包含数据库连接和 YAML 配置）
// strategy: 配置来源策略，如果为 nil 则使用默认策略（所有配置从数据库加载）
// httpServerLabels: HttpServer标签列表，用于区分不同的HttpServer，例如(client和console)
// 注意：策略中指定的配置项必须成功加载，如果加载失败会报错
func NewBaseConfigBuilder(configPool *ConfigPool, strategy *ConfigSourceStrategy, httpServerLabels []string) (*BaseConfigBuilder, error) {

	// 如果策略为 nil，使用默认策略
	if strategy == nil {
		return nil, fmt.Errorf("策略不能为 nil")
	}

	builder := &BaseConfigBuilder{
		strategy:         strategy,
		HttpServerLabels: httpServerLabels,
	}

	return builder, nil
}

// BaseConfigBuilder 配置构建器，支持从多个来源加载配置
type BaseConfigBuilder struct {
	strategy *ConfigSourceStrategy

	// HttpServer标签列表，用于区分不同的HttpServer，例如(client和console)
	HttpServerLabels []string `mapstructure:"http_server_labels"`
}

// loadProfileAndLocation 加载 Profile 和 Location 配置
func (b *BaseConfigBuilder) loadProfileAndLocation(cfg *configStruct.BaseConfig, configPool *ConfigPool) error {
	// 使用传入的 configPool（如果为 nil，使用 builder 自己的 configPool）

	// 从配置池中获取 dbRows
	dbRows := configPool.GetDBRows()
	// 加载 Profile
	if b.strategy.Profile == SourceDatabase && len(dbRows) > 0 {
		cfg.Profile = getAppProfile(dbRows)
	} else if b.strategy.Profile == SourceYAML && configPool != nil && len(configPool.GetYAML()) > 0 {
		// 从第一个 YAML 文件加载 Profile
		var profile configStruct.AppProfile
		if err := configPool.GetYAML()[0].UnmarshalKey(ConfigKeys.Profile.Key, &profile); err == nil {
			cfg.Profile = &profile
		}
	}

	// 如果 Profile 仍为 nil，创建默认值
	if cfg.Profile == nil {
		cfg.Profile = &configStruct.AppProfile{}
	}

	// 加载 Location
	if b.strategy.Location == SourceDatabase && len(dbRows) > 0 {
		location, err := getLocationConfig(dbRows)
		if err == nil {
			cfg.Location = location
		} else {
			cfg.Location = time.Local
		}
	} else if b.strategy.Location == SourceYAML && configPool != nil && len(configPool.GetYAML()) > 0 {
		// 从第一个 YAML 文件加载 Location
		timeZone := configPool.GetYAML()[0].GetString(ConfigKeys.LocationTZ.Key)
		if timeZone == "" {
			timeZone = "Asia/Shanghai"
		}
		location, err := time.LoadLocation(timeZone)
		if err != nil {
			location = time.FixedZone("CST-8", 8*3600)
		}
		cfg.Location = location
	} else {
		cfg.Location = time.Local
	}

	return nil
}

// 统一加载所有配置项
func (b *BaseConfigBuilder) loadAllConfigs(cfg *configStruct.BaseConfig, configPool *ConfigPool, debug bool) error {
	// 使用传入的 configPool（如果为 nil，使用 builder 自己的 configPool）

	// 从配置池获取 dbRows
	var dbRows []*DbSettingRow
	if configPool != nil {
		dbRows = configPool.GetDBRows()
	}

	for key, cm := range configMaps {
		sourceSelector, ok := configSources[key]
		if !ok {
			continue
		}
		src := sourceSelector(b.strategy)
		if src == "" {
			continue
		}

		targetPtr := reflect.New(reflect.TypeOf(cm.Data))

		switch src {
		case SourceDatabase:
			if len(dbRows) == 0 {
				return errFactory.databaseEmpty(key)
			}
			if err := fillConfigFromRows(targetPtr.Interface(), key, key, dbRows, debug); err != nil {
				return err
			}
			// 应用默认值
			applyDefaultsFromTags(targetPtr.Interface())
			// 验证配置
			if err := validateConfig(targetPtr.Interface(), key); err != nil {
				return err
			}
		case SourceYAML:
			if configPool == nil || len(configPool.GetYAML()) == 0 {
				return errFactory.yamlNotInit(key)
			}
			if err := configPool.GetYAML()[0].UnmarshalKey(key, targetPtr.Interface()); err != nil {
				return errFactory.yamlLoadFailed(key, err)
			}
			applyDefaultsFromTags(targetPtr.Interface())
			// 验证配置
			if err := validateConfig(targetPtr.Interface(), key); err != nil {
				return err
			}
		default:
			continue
		}

		assigner, ok := configAssigners[key]
		if !ok {
			continue
		}
		assigner(cfg, targetPtr.Interface())
	}
	return nil
}

// GetStrategy 返回配置来源策略
func (b *BaseConfigBuilder) GetStrategy() *ConfigSourceStrategy {
	return b.strategy
}

// Build 构建 BaseConfig，根据策略从不同来源加载配置
func (b *BaseConfigBuilder) Build(configPool *ConfigPool) (config *configStruct.BaseConfig, err error) {
	config = &configStruct.BaseConfig{}

	// 第一步：检查策略中是否有数据库相关的配置，如果有，优先初始化数据库
	// 这样后续的配置才能从数据库中读取
	// 注意：如果策略要求从数据库加载配置，configPool 中可能还没有数据库连接，需要先初始化
	if configPool == nil || configPool.GetDB() == nil {
		// 需要数据库连接，从 YAML 初始化
		if configPool == nil {
			err = fmt.Errorf("配置池未初始化")
			return
		}
		if err = configPool.InitDatabaseFromYAML(); err != nil {
			err = fmt.Errorf("初始化数据库连接失败: %w", err)
			return
		}
		// 数据库初始化后，重新加载配置行并更新到配置池
		if configPool.GetDB() != nil {
			newDbRows, loadErr := configPool.LoadSettingRows(context.Background())
			if loadErr == nil {
				configPool.dbRows = newDbRows
			}
		}
	}

	// 第二步：加载 Profile 和 Location（基础配置）
	if err = b.loadProfileAndLocation(config, configPool); err != nil {
		err = fmt.Errorf("加载基础配置失败: %w", err)
		return
	}

	debug := config.Profile != nil && config.Profile.Debug

	// 第三步：根据策略加载各个配置项（如果策略中指定了配置来源，则必须成功加载）

	// 初始化HttpServer配置
	config.HttpServerConfig = map[string]*configStruct.HttpServerConfig{}
	// 从配置池获取 dbRows
	dbRows := configPool.GetDBRows()
	for _, serverLabel := range b.HttpServerLabels {
		serverConfig := getHttpServerConfig(dbRows, serverLabel)
		if serverConfig == nil {
			err = fmt.Errorf("加载 HttpServer 配置失败: 标签 %s 不存在", serverLabel)
			return
		}
		config.HttpServerConfig[serverLabel] = serverConfig
	}

	// 第四步：统一管线加载已注册的配置
	if err = b.loadAllConfigs(config, configPool, debug); err != nil {
		err = err
		return
	}

	return
}
