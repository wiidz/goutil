package appMng

import (
	"fmt"
	"strings"
)

// errorFactory 统一生成所有错误，减少重复文案
type errorFactory struct{}

var errFactory = &errorFactory{}

// ========== 配置构建相关错误 ==========

// databaseEmpty 从数据库加载到空结果
func (f *errorFactory) databaseEmpty(key string) error {
	return fmt.Errorf("❌从数据库加载 %s 配置失败: 数据为空", key)
}

// yamlNotInit 未初始化 YAML
func (f *errorFactory) yamlNotInit(key string) error {
	return fmt.Errorf("❌从 YAML 加载 %s 配置失败: 未初始化 YAML 配置", key)
}

// yamlLoadFailed YAML 解析失败
func (f *errorFactory) yamlLoadFailed(key string, err error) error {
	return fmt.Errorf("❌从 YAML 加载 %s 配置失败: %w", key, err)
}

// databaseLoadFailed 从数据库加载配置失败
func (f *errorFactory) databaseLoadFailed(configName string, err error) error {
	return fmt.Errorf("❌从数据库加载配置 %s 失败: %w", configName, err)
}

// validateFailed 配置验证失败
func (f *errorFactory) validateFailed(configKey string, errMsgs []string) error {
	return fmt.Errorf("❌配置 %s 验证失败: %s", configKey, strings.Join(errMsgs, "; "))
}

// ========== 配置池相关错误 ==========

// yamlFilesEmpty YAML 配置文件列表为空
func (f *errorFactory) yamlFilesEmpty() error {
	return fmt.Errorf("❌YAML 配置文件列表不能为空")
}

// yamlInitFailed 初始化 YAML 配置失败
func (f *errorFactory) yamlInitFailed(err error) error {
	return fmt.Errorf("❌初始化 YAML 配置失败: %w", err)
}

// databaseInitFailed 初始化数据库连接失败
func (f *errorFactory) databaseInitFailed(err error) error {
	return fmt.Errorf("❌初始化数据库连接失败（需要从数据库加载配置，但无法连接数据库）: %w", err)
}

// postgresConnectFailed 连接 PostgreSQL 失败
func (f *errorFactory) postgresConnectFailed(err error) error {
	return fmt.Errorf("❌连接 PostgreSQL 失败: %w", err)
}

// mysqlConnectFailed 连接 MySQL 失败
func (f *errorFactory) mysqlConnectFailed(err error) error {
	return fmt.Errorf("❌连接 MySQL 失败: %w", err)
}

// getDBFailed 获取底层数据库连接失败
func (f *errorFactory) getDBFailed(err error) error {
	return fmt.Errorf("❌获取底层数据库连接失败: %w", err)
}

// yamlConfigNotInit YAML 配置未初始化
func (f *errorFactory) yamlConfigNotInit() error {
	return fmt.Errorf("❌YAML 配置未初始化")
}

// databaseConfigNotFound 无法从 YAML 加载数据库配置
func (f *errorFactory) databaseConfigNotFound() error {
	return fmt.Errorf("❌无法从 YAML 加载数据库配置，不存在 %s 或 %s", ConfigKeys.Postgres.Key, ConfigKeys.Mysql.Key)
}

// dbNotSet 数据库连接未设置
func (f *errorFactory) dbNotSet() error {
	return fmt.Errorf("❌数据库连接未设置")
}

// ========== 基础配置构建器相关错误 ==========

// configPoolNil 配置池不能为 nil
func (f *errorFactory) configPoolNil() error {
	return fmt.Errorf("❌配置池不能为 nil")
}

// strategyNil 策略不能为 nil
func (f *errorFactory) strategyNil() error {
	return fmt.Errorf("❌策略不能为 nil")
}

// configNotInCustom 配置项未在策略的 Custom 中定义
func (f *errorFactory) configNotInCustom(configName string) error {
	return fmt.Errorf("❌配置项 %s 未在策略的 Custom 中定义", configName)
}

// unsupportedSource 不支持的配置来源
func (f *errorFactory) unsupportedSource(source ConfigSource) error {
	return fmt.Errorf("❌不支持的配置来源: %s", source)
}

// loadBaseConfigFailed 加载基础配置失败
func (f *errorFactory) loadBaseConfigFailed(err error) error {
	return fmt.Errorf("❌加载基础配置失败: %w", err)
}

// httpServerConfigNotFound HttpServer 配置不存在
func (f *errorFactory) httpServerConfigNotFound(serverLabel string) error {
	return fmt.Errorf("❌加载 HttpServer 配置失败: 标签 %s 不存在", serverLabel)
}

// ========== 目标参数相关错误 ==========

// targetEmpty target 不能为空
func (f *errorFactory) targetEmpty() error {
	return fmt.Errorf("❌target 不能为空")
}

// targetNotPointer target 必须是指针类型
func (f *errorFactory) targetNotPointer() error {
	return fmt.Errorf("❌target 必须是指针类型")
}

// targetNotStructPointer target 必须是非 nil 指针
func (f *errorFactory) targetNotStructPointer() error {
	return fmt.Errorf("❌target 必须是非 nil 指针")
}

// targetNotStruct target 必须指向结构体
func (f *errorFactory) targetNotStruct() error {
	return fmt.Errorf("❌target 必须指向结构体")
}

// ========== AppMng 相关错误 ==========

// postgresInitFailed 初始化 PostgreSQL 失败
func (f *errorFactory) postgresInitFailed(err error) error {
	return fmt.Errorf("❌appMng: init postgres failed: %w", err)
}

// redisInitFailed 初始化 Redis 失败
func (f *errorFactory) redisInitFailed(err error) error {
	return fmt.Errorf("❌appMng: init redis failed: %w", err)
}

// esInitFailed 初始化 ES 失败
func (f *errorFactory) esInitFailed(err error) error {
	return fmt.Errorf("❌appMng: init es failed: %w", err)
}

// projectBuildFailed 项目配置构建失败
func (f *errorFactory) projectBuildFailed(err error) error {
	return fmt.Errorf("❌appMng: project build failed: %w", err)
}
