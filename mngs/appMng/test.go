package appMng

import (
	"context"
	"fmt"
	"testing"

	"github.com/wiidz/goutil/structs/configStruct"
)

func TestNewApp(t *testing.T) {

	//【1】定义pool
	configPool, err := NewConfigPool(context.Background(), []*configStruct.ViperConfig{
		{
			DirPath:  "./configs/",
			FileName: "common",
			FileType: "yaml",
		},
	}, "a_setting")
	if err != nil {
		t.Fatalf("NewConfigPool failed: %v", err)
	}

	//【2】构建基础配置构建器
	baseConfigBuilder, err := NewBaseConfigBuilder(configPool, &ConfigSourceStrategy{
		Profile:  SourceDatabase,
		Location: SourceYAML,
		Redis:    SourceDatabase,
		// Es:       SourceDatabase,
		// RabbitMQ: SourceYAML,
		// Postgres: SourceYAML,
		// Mysql:    SourceYAML,
	}, []string{"client", "console"})
	if err != nil {
		t.Fatalf("NewConfigBuilder failed: %v", err)
	}

	//【3】构建appMng
	app, err := NewApp(context.Background(), configPool, baseConfigBuilder, &MyProjectConfig{})
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}
	t.Logf("App: %+v", app)
}

type MyProjectConfig struct {
	TestConfig *TestConfig `mapstructure:"test_config"`
}

type TestConfig struct {
	Test string `mapstructure:"test"`
}

func (c *MyProjectConfig) Build(baseConfig *configStruct.BaseConfig, configPool *ConfigPool) error {
	// 初始化配置
	c.TestConfig = &TestConfig{}

	// 创建配置加载器（自动从 ProjectConfig 获取策略和调试模式）
	loader, err := NewProjectConfigLoader(c, baseConfig, configPool)
	if err != nil {
		return fmt.Errorf("创建配置加载器失败: %w", err)
	}

	// 使用链式调用加载配置（简洁易扩展）
	if err := loader.Load("test_config", c.TestConfig).Error(); err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 如果需要加载更多配置，可以继续链式调用：
	// c.AnotherConfig = &AnotherConfig{}
	// c.ThirdConfig = &ThirdConfig{}
	// if err := loader.
	// 	Load("another_config", c.AnotherConfig).
	// 	Load("third_config", c.ThirdConfig).
	// 	Error(); err != nil {
	// 	return fmt.Errorf("加载配置失败: %w", err)
	// }

	return nil
}

func (c *MyProjectConfig) GetStrategy() *ConfigSourceStrategy {
	return &ConfigSourceStrategy{
		Custom: map[string]ConfigSource{
			"test_config": SourceDatabase,
		},
	}
}
