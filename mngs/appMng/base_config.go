package appMng

import (
	"context"
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

// GenericProjectConfig 泛型项目配置（简化扩展，类型安全）
// T: 项目配置结构体类型
// 使用示例：
//
//	type MyConfig struct {
//	    ServiceA *ServiceAConfig
//	    ServiceB *ServiceBConfig
//	}
//	cfg := appMng.NewGenericProjectConfig[MyConfig](strategy)
//	// 在 Build 方法中链式加载：
//	func (c *MyConfig) Build(baseConfig *configStruct.BaseConfig, configPool *appMng.ConfigPool) error {
//	    return c.GenericProjectConfig.Build(baseConfig, configPool).
//	        Load("service_a", &c.Data.ServiceA).
//	        Load("service_b", &c.Data.ServiceB).
//	        Error()
//	}
type GenericProjectConfig[T any] struct {
	Data       T                     // 配置数据
	strategy   *ConfigSourceStrategy // 配置策略
	configPool *ConfigPool           // 配置池
	debug      bool                  // 调试模式
	err        error                 // 用于链式调用时保存错误
}

// NewGenericProjectConfig 创建泛型项目配置
// strategy: 配置来源策略（必须包含 Custom 字段定义配置项来源）
func NewGenericProjectConfig[T any](strategy *ConfigSourceStrategy) *GenericProjectConfig[T] {
	return &GenericProjectConfig[T]{
		strategy: strategy,
	}
}

// GetStrategy 返回配置策略（实现 ConfigStrategyProvider 接口）
func (g *GenericProjectConfig[T]) GetStrategy() *ConfigSourceStrategy {
	return g.strategy
}

// Build 构建项目配置（实现 ProjectConfig 接口）
// 初始化配置池和调试模式，并自动加载所有在 Custom 策略中定义的配置项
// 应用层通常不需要重写此方法，除非有特殊需求
func (g *GenericProjectConfig[T]) Build(baseConfig *configStruct.BaseConfig, configPool *ConfigPool) error {
	if configPool == nil {
		return errFactory.configPoolNil()
	}
	g.configPool = configPool
	g.debug = baseConfig.Profile != nil && baseConfig.Profile.Debug

	// 自动加载所有在 Custom 策略中定义的配置项
	return g.AutoLoad()
}

// AutoLoad 自动加载所有在 Custom 策略中定义的配置项
// 通过反射自动初始化指针并加载配置，应用层无需手动处理每个字段
// 使用示例：
//
//	func (c *MyProjectConfig) Build(baseConfig *configStruct.BaseConfig, configPool *appMng.ConfigPool) error {
//	    return c.GenericProjectConfig.Build(baseConfig, configPool).AutoLoad()
//	}
func (g *GenericProjectConfig[T]) AutoLoad() error {
	if g.strategy == nil {
		return errFactory.strategyNil()
	}

	// 获取策略中的 Custom 配置
	customStrategy := g.strategy.Custom
	if customStrategy == nil || len(customStrategy) == 0 {
		return nil // 如果没有自定义配置，直接返回
	}

	// 通过反射自动初始化指针并加载配置
	dataVal := reflect.ValueOf(&g.Data).Elem()
	dataType := dataVal.Type()

	for i := 0; i < dataVal.NumField(); i++ {
		field := dataVal.Field(i)
		fieldType := field.Type()
		fieldName := dataType.Field(i).Name

		// 只处理指针类型字段
		if fieldType.Kind() != reflect.Ptr {
			continue
		}

		// 检查策略中是否有该配置项
		if _, exists := customStrategy[fieldName]; !exists {
			continue // 如果策略中没有定义，跳过
		}

		// 自动初始化指针（如果为 nil）
		if field.IsNil() {
			newValue := reflect.New(fieldType.Elem())
			field.Set(newValue)
		}

		// 自动加载配置
		if err := g.loadConfig(fieldName, fieldName, field.Interface()); err != nil {
			return err
		}
	}

	return nil
}

// Load 加载配置项（链式调用，简洁易用）
// configName: 配置项名称（必须在 strategy.Custom 中定义）
// target: 目标字段的指针（通常是 g.Data 的某个字段）
func (g *GenericProjectConfig[T]) Load(configName string, target interface{}) *GenericProjectConfig[T] {
	if g.err != nil {
		return g // 如果已有错误，跳过
	}
	g.err = g.loadConfig(configName, configName, target)
	return g
}

// LoadWithKey 加载配置项（可指定 YAML 键名）
func (g *GenericProjectConfig[T]) LoadWithKey(configName, configKey string, target interface{}) *GenericProjectConfig[T] {
	if g.err != nil {
		return g // 如果已有错误，跳过
	}
	g.err = g.loadConfig(configName, configKey, target)
	return g
}

// Error 返回加载过程中的错误
func (g *GenericProjectConfig[T]) Error() error {
	return g.err
}

// loadConfig 内部方法：加载配置的核心逻辑
func (g *GenericProjectConfig[T]) loadConfig(configName, configKey string, target interface{}) error {
	if g.strategy == nil {
		return errFactory.strategyNil()
	}

	// 获取配置来源
	var source ConfigSource
	if g.strategy.Custom != nil {
		if src, ok := g.strategy.Custom[configName]; ok {
			source = src
		}
	}

	// 如果 Custom 中没有找到，返回错误（要求明确指定）
	if source == "" {
		return errFactory.configNotInCustom(configName)
	}

	// 创建目标指针
	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		return errFactory.targetNotPointer()
	}
	targetPtr := reflect.New(targetType.Elem())

	// 根据来源加载配置（使用公共加载逻辑）
	if err := loadConfigFromSource(source, configName, configKey, targetPtr.Interface(), g.configPool, g.debug); err != nil {
		return err
	}

	// 将加载的配置赋值给 target
	reflect.ValueOf(target).Elem().Set(targetPtr.Elem())

	return nil
}

// NewBaseConfigBuilder 创建基础构建器
// configPool: 配置池（包含数据库连接和 YAML 配置）
// strategy: 配置来源策略，如果为 nil 则使用默认策略（所有配置从数据库加载）
// httpServerLabels: HttpServer标签列表，用于区分不同的HttpServer，例如(client和console)
// 注意：策略中指定的配置项必须成功加载，如果加载失败会报错
func NewBaseConfigBuilder(configPool *ConfigPool, strategy *ConfigSourceStrategy, httpServerLabels []string) (*BaseConfigBuilder, error) {

	// 如果策略为 nil，使用默认策略
	if strategy == nil {
		return nil, errFactory.strategyNil()
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
	for key, configType := range configTypes {
		src := getConfigSource(b.strategy, key)
		if src == "" {
			continue
		}

		targetPtr := reflect.New(configType)

		// 使用公共加载逻辑
		if err := loadConfigFromSource(src, key, key, targetPtr.Interface(), configPool, debug); err != nil {
			return err
		}

		assignConfigToBaseConfig(cfg, key, targetPtr.Interface())
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
			err = errFactory.configPoolNil()
			return
		}
		if err = configPool.InitDatabaseFromYAML(); err != nil {
			err = errFactory.databaseInitFailed(err)
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
		err = errFactory.loadBaseConfigFailed(err)
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
			err = errFactory.httpServerConfigNotFound(serverLabel)
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
