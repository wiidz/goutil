package appMng

import (
	"fmt"
	"reflect"

	"github.com/wiidz/goutil/structs/configStruct"
)

// ProjectConfigLoader 项目配置加载器，提供更简洁的配置加载接口
type ProjectConfigLoader struct {
	configPool *ConfigPool
	strategy   *ConfigSourceStrategy
	debug      bool
	err        error // 用于链式调用时保存错误
}

// NewProjectConfigLoader 从 ProjectConfig 创建配置加载器
// 自动从 ProjectConfig 获取策略和调试模式
func NewProjectConfigLoader(projectConfig ProjectConfig, baseConfig *configStruct.BaseConfig, configPool *ConfigPool) (*ProjectConfigLoader, error) {
	if projectConfig == nil {
		return nil, fmt.Errorf("ProjectConfig 不能为 nil")
	}
	strategy := projectConfig.GetStrategy()
	if strategy == nil {
		return nil, fmt.Errorf("策略不能为 nil")
	}
	debug := baseConfig.Profile != nil && baseConfig.Profile.Debug
	return &ProjectConfigLoader{
		configPool: configPool,
		strategy:   strategy,
		debug:      debug,
	}, nil
}

// Load 加载单个配置项（支持链式调用）
// configName: 配置项名称（数据库中的 name 字段，或 YAML 中的配置键）
// target: 目标配置结构体的指针
// 如果之前有错误，会跳过本次加载
func (l *ProjectConfigLoader) Load(configName string, target interface{}) *ProjectConfigLoader {
	if l.err != nil {
		return l // 如果已有错误，跳过
	}
	l.err = LoadProjectConfig(configName, configName, target, l.configPool, l.strategy, l.debug)
	return l
}

// LoadWithKey 加载配置项（可以指定不同的 YAML 配置键，支持链式调用）
// configName: 配置项名称（数据库中的 name 字段）
// configKey: YAML 中的配置键（如果从 YAML 加载）
// target: 目标配置结构体的指针
// 如果之前有错误，会跳过本次加载
func (l *ProjectConfigLoader) LoadWithKey(configName string, configKey string, target interface{}) *ProjectConfigLoader {
	if l.err != nil {
		return l // 如果已有错误，跳过
	}
	l.err = LoadProjectConfig(configName, configKey, target, l.configPool, l.strategy, l.debug)
	return l
}

// Error 返回加载过程中的错误（用于链式调用结束时检查）
func (l *ProjectConfigLoader) Error() error {
	return l.err
}

// LoadProjectConfig 为 ProjectConfig 提供通用的配置加载功能
// configName: 配置项名称（用于从数据库或 YAML 中查找）
// configKey: YAML 中的配置键（如果从 YAML 加载）
// target: 目标配置结构体的指针
// configPool: 配置池
// strategy: 配置来源策略（从 GetStrategy() 获取）
// debug: 是否开启调试模式
func LoadProjectConfig(
	configName string,
	configKey string,
	target interface{},
	configPool *ConfigPool,
	strategy *ConfigSourceStrategy,
	debug bool,
) error {
	if strategy == nil {
		return fmt.Errorf("配置策略不能为 nil")
	}

	// 获取配置来源
	var source ConfigSource
	if strategy.Custom != nil {
		if src, ok := strategy.Custom[configName]; ok {
			source = src
		}
	}

	// 如果 Custom 中没有找到，返回错误（要求明确指定）
	if source == "" {
		return fmt.Errorf("配置项 %s 未在策略的 Custom 中定义", configName)
	}

	// 创建目标指针
	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		return fmt.Errorf("target 必须是指针类型")
	}
	targetPtr := reflect.New(targetType.Elem())

	// 根据来源加载配置
	switch source {
	case SourceDatabase:
		dbRows := configPool.GetDBRows()
		if len(dbRows) == 0 {
			return errFactory.databaseEmpty(configName)
		}
		// 使用 configName 作为数据库中的 name 字段，configKey 作为 flag1
		if err := fillConfigFromRows(targetPtr.Interface(), configName, configKey, dbRows, debug); err != nil {
			return fmt.Errorf("从数据库加载配置 %s 失败: %w", configName, err)
		}
		// 应用默认值
		applyDefaultsFromTags(targetPtr.Interface())
		// 验证配置
		if err := validateConfig(targetPtr.Interface(), configName); err != nil {
			return err
		}

	case SourceYAML:
		if configPool == nil || len(configPool.GetYAML()) == 0 {
			return errFactory.yamlNotInit(configName)
		}
		// 如果 configKey 为空，使用 configName
		yamlKey := configKey
		if yamlKey == "" {
			yamlKey = configName
		}
		if err := configPool.GetYAML()[0].UnmarshalKey(yamlKey, targetPtr.Interface()); err != nil {
			return errFactory.yamlLoadFailed(configName, err)
		}
		// 应用默认值
		applyDefaultsFromTags(targetPtr.Interface())
		// 验证配置
		if err := validateConfig(targetPtr.Interface(), configName); err != nil {
			return err
		}

	default:
		return fmt.Errorf("不支持的配置来源: %s", source)
	}

	// 将加载的配置赋值给 target
	reflect.ValueOf(target).Elem().Set(targetPtr.Elem())

	return nil
}
