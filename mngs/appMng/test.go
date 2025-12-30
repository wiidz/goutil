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
	app, err := NewApp(context.Background(), configPool, baseConfigBuilder, NewMyProjectConfig())
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}
	t.Logf("App: %+v", app)
}

// MyProjectConfigData 测试配置数据
type MyProjectConfigData struct {
	TestConfig *TestConfig
}

type TestConfig struct {
	Test string `mapstructure:"test"`
}

// MyProjectConfig 测试项目配置（使用泛型）
type MyProjectConfig struct {
	*GenericProjectConfig[MyProjectConfigData]
}

func NewMyProjectConfig() *MyProjectConfig {
	return &MyProjectConfig{
		GenericProjectConfig: NewGenericProjectConfig[MyProjectConfigData](&ConfigSourceStrategy{
			Custom: map[string]ConfigSource{
				"test_config": SourceDatabase,
			},
		}),
	}
}

func (c *MyProjectConfig) Build(baseConfig *configStruct.BaseConfig, configPool *ConfigPool) error {
	// 初始化
	if err := c.GenericProjectConfig.Build(baseConfig, configPool); err != nil {
		return err
	}

	// 初始化配置
	c.Data.TestConfig = &TestConfig{}

	// 使用链式调用加载配置（简洁易扩展）
	if err := c.GenericProjectConfig.
		Load("test_config", c.Data.TestConfig).
		Error(); err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 如果需要加载更多配置，可以继续链式调用：
	// c.Data.AnotherConfig = &AnotherConfig{}
	// c.Data.ThirdConfig = &ThirdConfig{}
	// if err := c.GenericProjectConfig.
	// 	Load("another_config", c.Data.AnotherConfig).
	// 	Load("third_config", c.Data.ThirdConfig).
	// 	Error(); err != nil {
	// 	return fmt.Errorf("加载配置失败: %w", err)
	// }

	return nil
}
