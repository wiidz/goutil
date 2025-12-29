package appMng

// 配置加载机制使用示例：
//
// 示例1：所有配置从数据库加载（需要数据库连接）
//   initialConfig := &appMng.InitialConfig{
//       DB: &configStruct.DBConfig{
//           Type: configStruct.DBTypePostgres,
//           DSN: "postgres://user:password@host:port/dbname",
//           SettingTableName: "a_setting", // 可选，默认值为 "a_setting"
//       },
//   }
//   builder, _ := appMng.NewConfigBuilder(initialConfig, nil)
//   builder.SetDatabase(db) // 设置数据库连接
//   baseConfig, _ := builder.Build(ctx)
//
// 示例2：所有配置从 YAML 加载（不需要数据库）
//   initialConfig := &appMng.InitialConfig{
//       YAMLFiles: []*configStruct.ViperConfig{
//           {DirPath: "./configs", FileName: "common", FileType: "yaml"},
//       },
//   }
//   strategy := &appMng.ConfigSourceStrategy{
//       Redis: appMng.SourceYAML,
//       Es:    appMng.SourceYAML,
//       // ... 其他配置项都设置为 SourceYAML
//   }
//   builder, _ := appMng.NewConfigBuilder(initialConfig, strategy)
//   baseConfig, _ := builder.Build(ctx)
//
// 示例3：混合加载（部分从数据库，部分从 YAML）
//   initialConfig := &appMng.InitialConfig{
//       DB: &configStruct.DBConfig{
//           Type: configStruct.DBTypePostgres,
//           DSN: "postgres://user:password@host:port/dbname",
//           SettingTableName: "custom_setting", // 可选，指定自定义配置表名
//       },
//       YAMLFiles: []*configStruct.ViperConfig{
//           {DirPath: "./configs", FileName: "common", FileType: "yaml"},
//       },
//   }
//   strategy := &appMng.ConfigSourceStrategy{
//       Redis:      appMng.SourceDatabase, // Redis 从数据库加载（必须成功）
//       AliApi:     appMng.SourceYAML,      // 阿里云 API 从 YAML 加载（必须成功）
//       WechatMini: appMng.SourceDatabase, // 微信小程序从数据库加载（必须成功）
//       // ... 其他配置项根据需要设置，如果设置了就必须成功加载
//   }
//   builder, _ := appMng.NewConfigBuilder(initialConfig, strategy)
//   builder.SetDatabase(db)
//   baseConfig, _ := builder.Build(ctx) // 如果策略中指定的配置加载失败会报错

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/wiidz/goutil/helpers/configHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"gorm.io/gorm"
)

// ConfigSource 配置来源类型
type ConfigSource string

const (
	SourceYAML     ConfigSource = "yaml"     // 从 YAML 文件加载
	SourceDatabase ConfigSource = "database" // 从数据库加载
)

// InitialConfig 初始配置，在应用构建之初传入
// 如果项目涉及数据库连接，必须包含数据库连接信息
type InitialConfig struct {
	// 数据库连接配置（统一的数据库配置，支持 PostgreSQL 和 MySQL）
	DB *configStruct.DBConfig `mapstructure:"db"`

	// 配置表名（从数据库加载配置时使用的表名）
	// 优先级：InitialConfig.SettingTableName > DB.SettingTableName > 默认值 "a_setting"
	SettingTableName string `mapstructure:"setting_table_name"`

	// YAML 配置文件列表（支持多个 YAML 文件）
	YAMLFiles []*configStruct.ViperConfig `mapstructure:"yaml_files"`
}

// ConfigSourceStrategy 配置来源策略，定义每个配置项应该从哪个来源加载
type ConfigSourceStrategy struct {

	// Profile 和 Location 配置来源（通常从数据库或第一个 YAML 文件）
	Profile  ConfigSource `mapstructure:"profile"`  // Profile 配置来源
	Location ConfigSource `mapstructure:"location"` // Location 配置来源

	// 存储相关配置
	Redis    ConfigSource `mapstructure:"redis"`    // Redis 配置来源
	Es       ConfigSource `mapstructure:"es"`       // Elasticsearch 配置来源
	RabbitMQ ConfigSource `mapstructure:"rabbitmq"` // RabbitMQ 配置来源
	Postgres ConfigSource `mapstructure:"postgres"` // PostgreSQL 配置来源
	Mysql    ConfigSource `mapstructure:"mysql"`    // MySQL 配置来源

	// 微信相关配置
	WechatMini  ConfigSource `mapstructure:"wechat_mini"`   // 微信小程序配置来源
	WechatOa    ConfigSource `mapstructure:"wechat_oa"`     // 微信公众号配置来源
	WechatOpen  ConfigSource `mapstructure:"wechat_open"`   // 微信开放平台配置来源
	WechatPayV3 ConfigSource `mapstructure:"wechat_pay_v3"` // 微信支付 V3 配置来源
	WechatPayV2 ConfigSource `mapstructure:"wechat_pay_v2"` // 微信支付 V2 配置来源

	// 阿里相关配置
	AliOss ConfigSource `mapstructure:"ali_oss"` // 阿里云 OSS 配置来源
	AliPay ConfigSource `mapstructure:"ali_pay"` // 支付宝配置来源
	AliApi ConfigSource `mapstructure:"ali_api"` // 阿里云 API 配置来源
	AliSms ConfigSource `mapstructure:"ali_sms"` // 阿里云短信配置来源
	AliIot ConfigSource `mapstructure:"ali_iot"` // 阿里云 IoT 配置来源
	Amap   ConfigSource `mapstructure:"amap"`    // 高德地图配置来源

	// 其他配置
	Yunxin ConfigSource `mapstructure:"yunxin"` // 网易云信配置来源

}

// DefaultConfigSourceStrategy 返回默认的配置来源策略（所有配置都从数据库加载）
func DefaultConfigSourceStrategy() *ConfigSourceStrategy {
	return &ConfigSourceStrategy{
		Redis:       SourceDatabase,
		Es:          SourceDatabase,
		RabbitMQ:    SourceDatabase,
		Postgres:    SourceDatabase,
		Mysql:       SourceDatabase,
		WechatMini:  SourceDatabase,
		WechatOa:    SourceDatabase,
		WechatOpen:  SourceDatabase,
		WechatPayV3: SourceDatabase,
		WechatPayV2: SourceDatabase,
		AliOss:      SourceDatabase,
		AliPay:      SourceDatabase,
		AliApi:      SourceDatabase,
		AliSms:      SourceDatabase,
		AliIot:      SourceDatabase,
		Amap:        SourceDatabase,
		Yunxin:      SourceDatabase,
		Profile:     SourceDatabase,
		Location:    SourceDatabase,
	}
}

// ConfigBuilder 配置构建器，支持从多个来源加载配置
type ConfigBuilder struct {
	initialConfig *InitialConfig
	strategy      *ConfigSourceStrategy
	db            *gorm.DB       // 数据库连接（如果使用数据库）
	yamlVipers    []*viper.Viper // YAML 配置实例列表
}

// NewConfigBuilder 创建配置构建器
// initialConfig: 初始配置，如果使用数据库，必须包含数据库连接信息
// strategy: 配置来源策略，如果为 nil 则使用默认策略（所有配置从数据库加载）
// 注意：策略中指定的配置项必须成功加载，如果加载失败会报错
func NewConfigBuilder(initialConfig *InitialConfig, strategy *ConfigSourceStrategy) (*ConfigBuilder, error) {
	if initialConfig == nil {
		return nil, fmt.Errorf("初始配置不能为 nil")
	}

	// 如果策略为 nil，使用默认策略
	if strategy == nil {
		strategy = DefaultConfigSourceStrategy()
	}

	builder := &ConfigBuilder{
		initialConfig: initialConfig,
		strategy:      strategy,
		yamlVipers:    make([]*viper.Viper, 0),
	}

	// 初始化 YAML 配置
	if len(initialConfig.YAMLFiles) > 0 {
		for _, yamlConfig := range initialConfig.YAMLFiles {
			v, err := configHelper.GetViper(yamlConfig)
			if err != nil {
				// YAML 文件不存在时记录警告但不报错（允许部分配置从数据库加载）
				log.Printf("警告: 无法加载 YAML 配置文件 %s/%s.%s: %v", yamlConfig.DirPath, yamlConfig.FileName, yamlConfig.FileType, err)
				continue
			}
			builder.yamlVipers = append(builder.yamlVipers, v)
		}
	}

	return builder, nil
}

// SetDatabase 设置数据库连接（用于从数据库加载配置）
func (b *ConfigBuilder) SetDatabase(db *gorm.DB) {
	b.db = db
}

// Build 构建 BaseConfig，根据策略从不同来源加载配置
func (b *ConfigBuilder) Build(ctx context.Context) (*configStruct.BaseConfig, error) {
	cfg := &configStruct.BaseConfig{}

	// 加载数据库配置行（如果需要）
	var dbRows []*DbSettingRow
	if b.db != nil {
		var err error
		dbRows, err = b.loadDatabaseRows(ctx)
		if err != nil {
			log.Printf("警告: 无法从数据库加载配置: %v", err)
		}
	}

	// 加载 Profile 和 Location（基础配置）
	if err := b.loadProfileAndLocation(cfg, dbRows); err != nil {
		return nil, fmt.Errorf("加载基础配置失败: %w", err)
	}

	debug := cfg.Profile != nil && cfg.Profile.Debug

	// 根据策略加载各个配置项（如果策略中指定了配置来源，则必须成功加载）
	if b.strategy.Redis != "" {
		if err := b.loadRedisConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载 Redis 配置失败: %w", err)
		}
	}
	if b.strategy.Es != "" {
		if err := b.loadEsConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载 Elasticsearch 配置失败: %w", err)
		}
	}
	if b.strategy.RabbitMQ != "" {
		if err := b.loadRabbitMQConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载 RabbitMQ 配置失败: %w", err)
		}
	}
	if b.strategy.Postgres != "" {
		if err := b.loadPostgresConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载 PostgreSQL 配置失败: %w", err)
		}
	}
	if b.strategy.Mysql != "" {
		if err := b.loadMysqlConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载 MySQL 配置失败: %w", err)
		}
	}
	if b.strategy.WechatMini != "" {
		if err := b.loadWechatMiniConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载微信小程序配置失败: %w", err)
		}
	}
	if b.strategy.WechatOa != "" {
		if err := b.loadWechatOaConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载微信公众号配置失败: %w", err)
		}
	}
	if b.strategy.WechatOpen != "" {
		if err := b.loadWechatOpenConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载微信开放平台配置失败: %w", err)
		}
	}
	if b.strategy.WechatPayV3 != "" {
		if err := b.loadWechatPayV3Config(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载微信支付 V3 配置失败: %w", err)
		}
	}
	if b.strategy.WechatPayV2 != "" {
		if err := b.loadWechatPayV2Config(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载微信支付 V2 配置失败: %w", err)
		}
	}
	if b.strategy.AliOss != "" {
		if err := b.loadAliOssConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载阿里云 OSS 配置失败: %w", err)
		}
	}
	if b.strategy.AliPay != "" {
		if err := b.loadAliPayConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载支付宝配置失败: %w", err)
		}
	}
	if b.strategy.AliApi != "" {
		if err := b.loadAliApiConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载阿里云 API 配置失败: %w", err)
		}
	}
	if b.strategy.AliSms != "" {
		if err := b.loadAliSmsConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载阿里云短信配置失败: %w", err)
		}
	}
	if b.strategy.AliIot != "" {
		if err := b.loadAliIotConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载阿里云 IoT 配置失败: %w", err)
		}
	}
	if b.strategy.Amap != "" {
		if err := b.loadAmapConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载高德地图配置失败: %w", err)
		}
	}
	if b.strategy.Yunxin != "" {
		if err := b.loadYunxinConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载网易云信配置失败: %w", err)
		}
	}

	return cfg, nil
}

// loadDatabaseRows 从数据库加载配置行
func (b *ConfigBuilder) loadDatabaseRows(ctx context.Context) ([]*DbSettingRow, error) {
	if b.db == nil {
		return nil, fmt.Errorf("数据库连接未设置")
	}

	// 确定使用的表名
	tableName := b.getSettingTableName()

	var rows []*DbSettingRow
	err := b.db.WithContext(ctx).
		Table(tableName).
		Where("kind = ? AND deleted_at IS NULL", 1).
		Find(&rows).Error

	return rows, err
}

// getSettingTableName 获取配置表名
// 优先级：InitialConfig.SettingTableName > DB.SettingTableName > 默认值 "a_setting"
func (b *ConfigBuilder) getSettingTableName() string {
	// 如果 InitialConfig 中设置了表名，优先使用
	if b.initialConfig.SettingTableName != "" {
		return b.initialConfig.SettingTableName
	}

	// 如果统一的 DB 配置中设置了表名，使用 DB 配置中的值
	if b.initialConfig.DB != nil && b.initialConfig.DB.SettingTableName != "" {
		return b.initialConfig.DB.SettingTableName
	}

	// 默认值
	return "a_setting"
}

// loadProfileAndLocation 加载 Profile 和 Location 配置
func (b *ConfigBuilder) loadProfileAndLocation(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow) error {
	// 加载 Profile
	if b.strategy.Profile == SourceDatabase && len(dbRows) > 0 {
		cfg.Profile = getAppProfile(dbRows)
	} else if b.strategy.Profile == SourceYAML && len(b.yamlVipers) > 0 {
		// 从第一个 YAML 文件加载 Profile
		var profile configStruct.AppProfile
		if err := b.yamlVipers[0].UnmarshalKey("profile", &profile); err == nil {
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
	} else if b.strategy.Location == SourceYAML && len(b.yamlVipers) > 0 {
		// 从第一个 YAML 文件加载 Location
		timeZone := b.yamlVipers[0].GetString("location.timezone")
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

// 各个配置项的加载函数（根据策略从数据库或 YAML 加载）

func (b *ConfigBuilder) loadRedisConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.Redis == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载 Redis 配置，但数据库配置行为空")
		}
		redisConfig, err := getRedisConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.RedisConfig = redisConfig
	} else if b.strategy.Redis == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载 Redis 配置，但 YAML 配置未初始化")
		}
		var redisConfig configStruct.RedisConfig
		if err := b.yamlVipers[0].UnmarshalKey("redis_config", &redisConfig); err != nil {
			return fmt.Errorf("从 YAML 加载 Redis 配置失败: %w", err)
		}
		cfg.RedisConfig = &redisConfig
	}
	return nil
}

func (b *ConfigBuilder) loadEsConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.Es == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载 Elasticsearch 配置，但数据库配置行为空")
		}
		esConfig, err := getEsConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.EsConfig = esConfig
	} else if b.strategy.Es == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载 Elasticsearch 配置，但 YAML 配置未初始化")
		}
		var esConfig configStruct.EsConfig
		if err := b.yamlVipers[0].UnmarshalKey("es_config", &esConfig); err != nil {
			return fmt.Errorf("从 YAML 加载 Elasticsearch 配置失败: %w", err)
		}
		cfg.EsConfig = &esConfig
	}
	return nil
}

func (b *ConfigBuilder) loadRabbitMQConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.RabbitMQ == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载 RabbitMQ 配置，但数据库配置行为空")
		}
		rabbitMQConfig, err := getRabbitMQConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.RabbitMQConfig = rabbitMQConfig
	} else if b.strategy.RabbitMQ == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载 RabbitMQ 配置，但 YAML 配置未初始化")
		}
		var rabbitMQConfig configStruct.RabbitMQConfig
		if err := b.yamlVipers[0].UnmarshalKey("rabbitmq_config", &rabbitMQConfig); err != nil {
			return fmt.Errorf("从 YAML 加载 RabbitMQ 配置失败: %w", err)
		}
		cfg.RabbitMQConfig = &rabbitMQConfig
	}
	return nil
}

func (b *ConfigBuilder) loadPostgresConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.Postgres == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载 PostgreSQL 配置，但数据库配置行为空")
		}
		postgresConfig, err := getPostgresConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.PostgresConfig = postgresConfig
	} else if b.strategy.Postgres == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载 PostgreSQL 配置，但 YAML 配置未初始化")
		}
		var postgresConfig configStruct.PostgresConfig
		if err := b.yamlVipers[0].UnmarshalKey("postgres_config", &postgresConfig); err != nil {
			return fmt.Errorf("从 YAML 加载 PostgreSQL 配置失败: %w", err)
		}
		cfg.PostgresConfig = &postgresConfig
	}
	return nil
}

func (b *ConfigBuilder) loadMysqlConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.Mysql == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载 MySQL 配置，但数据库配置行为空")
		}
		mysqlConfig, err := getMysqlConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.MysqlConfig = mysqlConfig
	} else if b.strategy.Mysql == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载 MySQL 配置，但 YAML 配置未初始化")
		}
		var mysqlConfig configStruct.MysqlConfig
		if err := b.yamlVipers[0].UnmarshalKey("mysql_config", &mysqlConfig); err != nil {
			return fmt.Errorf("从 YAML 加载 MySQL 配置失败: %w", err)
		}
		cfg.MysqlConfig = &mysqlConfig
	}
	return nil
}

func (b *ConfigBuilder) loadWechatMiniConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.WechatMini == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载微信小程序配置，但数据库配置行为空")
		}
		wechatMiniConfig, err := getWechatMiniConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.WechatMiniConfig = wechatMiniConfig
	} else if b.strategy.WechatMini == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载微信小程序配置，但 YAML 配置未初始化")
		}
		var wechatMiniConfig configStruct.WechatMiniConfig
		if err := b.yamlVipers[0].UnmarshalKey("wechat_mini_config", &wechatMiniConfig); err != nil {
			return fmt.Errorf("从 YAML 加载微信小程序配置失败: %w", err)
		}
		cfg.WechatMiniConfig = &wechatMiniConfig
	}
	return nil
}

func (b *ConfigBuilder) loadWechatOaConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.WechatOa == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载微信公众号配置，但数据库配置行为空")
		}
		wechatOaConfig, err := getWechatOaConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.WechatOaConfig = wechatOaConfig
	} else if b.strategy.WechatOa == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载微信公众号配置，但 YAML 配置未初始化")
		}
		var wechatOaConfig configStruct.WechatOaConfig
		if err := b.yamlVipers[0].UnmarshalKey("wechat_oa_config", &wechatOaConfig); err != nil {
			return fmt.Errorf("从 YAML 加载微信公众号配置失败: %w", err)
		}
		cfg.WechatOaConfig = &wechatOaConfig
	}
	return nil
}

func (b *ConfigBuilder) loadWechatOpenConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.WechatOpen == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载微信开放平台配置，但数据库配置行为空")
		}
		wechatOpenConfig, err := getWechatOpenConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.WechatOpenConfig = wechatOpenConfig
	} else if b.strategy.WechatOpen == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载微信开放平台配置，但 YAML 配置未初始化")
		}
		var wechatOpenConfig configStruct.WechatOpenConfig
		if err := b.yamlVipers[0].UnmarshalKey("wechat_open_config", &wechatOpenConfig); err != nil {
			return fmt.Errorf("从 YAML 加载微信开放平台配置失败: %w", err)
		}
		cfg.WechatOpenConfig = &wechatOpenConfig
	}
	return nil
}

func (b *ConfigBuilder) loadWechatPayV3Config(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.WechatPayV3 == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载微信支付 V3 配置，但数据库配置行为空")
		}
		wechatPayV3Config, err := getWechatPayConfigV3(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.WechatPayConfigV3 = wechatPayV3Config
	} else if b.strategy.WechatPayV3 == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载微信支付 V3 配置，但 YAML 配置未初始化")
		}
		var wechatPayV3Config configStruct.WechatPayConfigV3
		if err := b.yamlVipers[0].UnmarshalKey("wechat_pay_config_v3", &wechatPayV3Config); err != nil {
			return fmt.Errorf("从 YAML 加载微信支付 V3 配置失败: %w", err)
		}
		cfg.WechatPayConfigV3 = &wechatPayV3Config
	}
	return nil
}

func (b *ConfigBuilder) loadWechatPayV2Config(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.WechatPayV2 == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载微信支付 V2 配置，但数据库配置行为空")
		}
		wechatPayV2Config, err := getWechatPayConfigV2(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.WechatPayConfigV2 = wechatPayV2Config
	} else if b.strategy.WechatPayV2 == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载微信支付 V2 配置，但 YAML 配置未初始化")
		}
		var wechatPayV2Config configStruct.WechatPayConfigV2
		if err := b.yamlVipers[0].UnmarshalKey("wechat_pay_config_v2", &wechatPayV2Config); err != nil {
			return fmt.Errorf("从 YAML 加载微信支付 V2 配置失败: %w", err)
		}
		cfg.WechatPayConfigV2 = &wechatPayV2Config
	}
	return nil
}

func (b *ConfigBuilder) loadAliOssConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.AliOss == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载阿里云 OSS 配置，但数据库配置行为空")
		}
		aliOssConfig, err := getAliOssConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.AliOssConfig = aliOssConfig
	} else if b.strategy.AliOss == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载阿里云 OSS 配置，但 YAML 配置未初始化")
		}
		var aliOssConfig configStruct.AliOssConfig
		if err := b.yamlVipers[0].UnmarshalKey("ali_oss_config", &aliOssConfig); err != nil {
			return fmt.Errorf("从 YAML 加载阿里云 OSS 配置失败: %w", err)
		}
		cfg.AliOssConfig = &aliOssConfig
	}
	return nil
}

func (b *ConfigBuilder) loadAliPayConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.AliPay == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载支付宝配置，但数据库配置行为空")
		}
		aliPayConfig, err := getAliPayConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.AliPayConfig = aliPayConfig
	} else if b.strategy.AliPay == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载支付宝配置，但 YAML 配置未初始化")
		}
		var aliPayConfig configStruct.AliPayConfig
		if err := b.yamlVipers[0].UnmarshalKey("ali_pay_config", &aliPayConfig); err != nil {
			return fmt.Errorf("从 YAML 加载支付宝配置失败: %w", err)
		}
		cfg.AliPayConfig = &aliPayConfig
	}
	return nil
}

func (b *ConfigBuilder) loadAliApiConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.AliApi == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载阿里云 API 配置，但数据库配置行为空")
		}
		aliApiConfig, err := getAliApiConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.AliApiConfig = aliApiConfig
	} else if b.strategy.AliApi == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载阿里云 API 配置，但 YAML 配置未初始化")
		}
		var aliApiConfig configStruct.AliApiConfig
		if err := b.yamlVipers[0].UnmarshalKey("ali_api_config", &aliApiConfig); err != nil {
			return fmt.Errorf("从 YAML 加载阿里云 API 配置失败: %w", err)
		}
		cfg.AliApiConfig = &aliApiConfig
	}
	return nil
}

func (b *ConfigBuilder) loadAliSmsConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.AliSms == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载阿里云短信配置，但数据库配置行为空")
		}
		aliSmsConfig, err := getAliSmsConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.AliSmsConfig = aliSmsConfig
	} else if b.strategy.AliSms == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载阿里云短信配置，但 YAML 配置未初始化")
		}
		var aliSmsConfig configStruct.AliSmsConfig
		if err := b.yamlVipers[0].UnmarshalKey("ali_sms_config", &aliSmsConfig); err != nil {
			return fmt.Errorf("从 YAML 加载阿里云短信配置失败: %w", err)
		}
		cfg.AliSmsConfig = &aliSmsConfig
	}
	return nil
}

func (b *ConfigBuilder) loadAliIotConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.AliIot == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载阿里云 IoT 配置，但数据库配置行为空")
		}
		aliIotConfig, err := getAliIotConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.AliIotConfig = aliIotConfig
	} else if b.strategy.AliIot == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载阿里云 IoT 配置，但 YAML 配置未初始化")
		}
		var aliIotConfig configStruct.AliIotConfig
		if err := b.yamlVipers[0].UnmarshalKey("ali_iot_config", &aliIotConfig); err != nil {
			return fmt.Errorf("从 YAML 加载阿里云 IoT 配置失败: %w", err)
		}
		cfg.AliIotConfig = &aliIotConfig
	}
	return nil
}

func (b *ConfigBuilder) loadAmapConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.Amap == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载高德地图配置，但数据库配置行为空")
		}
		amapConfig, err := getAmapConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.AmapConfig = amapConfig
	} else if b.strategy.Amap == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载高德地图配置，但 YAML 配置未初始化")
		}
		var amapConfig configStruct.AmapConfig
		if err := b.yamlVipers[0].UnmarshalKey("amap_config", &amapConfig); err != nil {
			return fmt.Errorf("从 YAML 加载高德地图配置失败: %w", err)
		}
		cfg.AmapConfig = &amapConfig
	}
	return nil
}

func (b *ConfigBuilder) loadYunxinConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	if b.strategy.Yunxin == SourceDatabase {
		if len(dbRows) == 0 {
			return fmt.Errorf("策略要求从数据库加载网易云信配置，但数据库配置行为空")
		}
		yunxinConfig, err := getYunXinConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.YunxinConfig = yunxinConfig
	} else if b.strategy.Yunxin == SourceYAML {
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf("策略要求从 YAML 加载网易云信配置，但 YAML 配置未初始化")
		}
		var yunxinConfig configStruct.YunxinConfig
		if err := b.yamlVipers[0].UnmarshalKey("yunxin_config", &yunxinConfig); err != nil {
			return fmt.Errorf("从 YAML 加载网易云信配置失败: %w", err)
		}
		cfg.YunxinConfig = &yunxinConfig
	}
	return nil
}

// GetValueFromRow 从 rows 中检索符合条件的数据。
func GetValueFromRow(rows []*DbSettingRow, name, flag1, flag2, defaultValue string, debug bool) (value string) {
	if len(rows) == 0 {
		return
	}

	var row *DbSettingRow
	for i := range rows {
		item := rows[i]
		if item.Name != name {
			continue
		}
		if flag1 != "" && item.Flag1 != flag1 {
			continue
		}
		if flag2 != "" && item.Flag2 != flag2 {
			continue
		}
		row = item
		break
	}

	if row == nil {
		value = defaultValue
		return
	}

	value = row.Value1
	if debug {
		value = row.Value2
	}
	if value == "" && defaultValue != "" {
		value = defaultValue
	}
	return
}

func getAppProfile(rows []*DbSettingRow) *configStruct.AppProfile {
	return &configStruct.AppProfile{
		No:      GetValueFromRow(rows, "app", "", "no", "", false),
		Name:    GetValueFromRow(rows, "app", "", "name", "", false),
		Host:    GetValueFromRow(rows, "app", "", "host", "", false),
		Port:    GetValueFromRow(rows, "app", "", "port", "127.0.0.1", false),
		Domain:  GetValueFromRow(rows, "app", "", "domain", "", false),
		Debug:   GetValueFromRow(rows, "app", "", "debug", "", false) == "1",
		Version: GetValueFromRow(rows, "app", "", "version", "", false),
	}
}

func getLocationConfig(rows []*DbSettingRow) (location *time.Location, err error) {
	timeZone := GetValueFromRow(rows, "time_zone", "", "", "Asia/Shanghai", false)
	location, err = time.LoadLocation(timeZone)
	if err != nil {
		location = time.FixedZone("CST-8", 8*3600)
	}
	return
}

func getWechatMiniConfig(rows []*DbSettingRow, debug bool) (*configStruct.WechatMiniConfig, error) {
	cfg := &configStruct.WechatMiniConfig{
		AppID:     GetValueFromRow(rows, "wechat", "mini", "app_id", "", debug),
		AppSecret: GetValueFromRow(rows, "wechat", "mini", "app_secret", "", debug),
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("微信小程序配置 AppID 为空")
	}
	if cfg.AppSecret == "" {
		return nil, fmt.Errorf("微信小程序配置 AppSecret 为空")
	}
	return cfg, nil
}

func getWechatOaConfig(rows []*DbSettingRow, debug bool) (*configStruct.WechatOaConfig, error) {
	cfg := &configStruct.WechatOaConfig{
		AppID:          GetValueFromRow(rows, "wechat", "oa", "app_id", "", debug),
		AppSecret:      GetValueFromRow(rows, "wechat", "oa", "app_secret", "", debug),
		Token:          GetValueFromRow(rows, "wechat", "oa", "token", "", debug),
		EncodingAESKey: GetValueFromRow(rows, "wechat", "oa", "encoding_aes_key", "", debug),
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("微信公众号配置 AppID 为空")
	}
	if cfg.AppSecret == "" {
		return nil, fmt.Errorf("微信公众号配置 AppSecret 为空")
	}
	return cfg, nil
}

func getWechatOpenConfig(rows []*DbSettingRow, debug bool) (*configStruct.WechatOpenConfig, error) {
	cfg := &configStruct.WechatOpenConfig{
		AppID:     GetValueFromRow(rows, "wechat", "open", "app_id", "", debug),
		AppSecret: GetValueFromRow(rows, "wechat", "open", "app_secret", "", debug),
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("微信开放平台配置 AppID 为空")
	}
	if cfg.AppSecret == "" {
		return nil, fmt.Errorf("微信开放平台配置 AppSecret 为空")
	}
	return cfg, nil
}

func getWechatPayConfigV3(rows []*DbSettingRow, debug bool) (*configStruct.WechatPayConfigV3, error) {
	cfg := &configStruct.WechatPayConfigV3{
		AppID:                     GetValueFromRow(rows, "wechat", "pay_v3", "app_id", "", debug),
		ApiKeyV3:                  GetValueFromRow(rows, "wechat", "pay_v3", "api_key", "", debug),
		MchID:                     GetValueFromRow(rows, "wechat", "pay_v3", "mch_id", "", debug),
		CertURI:                   GetValueFromRow(rows, "wechat", "pay_v3", "cert_uri", "", debug),
		KeyURI:                    GetValueFromRow(rows, "wechat", "pay_v3", "key_uri", "", debug),
		PEMPrivateKeyContent:      GetValueFromRow(rows, "wechat", "pay_v3", "pem_private_key_content", "", debug),
		PEMCertContent:            GetValueFromRow(rows, "wechat", "pay_v3", "pem_cert_content", "", debug),
		CertSerialNo:              GetValueFromRow(rows, "wechat", "pay_v3", "cert_serial_no", "", debug),
		NotifyURL:                 GetValueFromRow(rows, "wechat", "pay_v3", "notify_url", "", debug),
		RefundNotifyURL:           GetValueFromRow(rows, "wechat", "pay_v3", "refund_notify_url", "", debug),
		MerchantTransferNotifyURL: GetValueFromRow(rows, "wechat", "pay_v3", "merchant_transfer_notify_url", "", debug),
		Debug:                     GetValueFromRow(rows, "wechat", "pay_v3", "debug", "0", debug) == "1",
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("微信支付 V3 配置 AppID 为空")
	}
	if cfg.MchID == "" {
		return nil, fmt.Errorf("微信支付 V3 配置 MchID 为空")
	}
	return cfg, nil
}

func getWechatPayConfigV2(rows []*DbSettingRow, debug bool) (*configStruct.WechatPayConfigV2, error) {
	cfg := &configStruct.WechatPayConfigV2{
		AppID:           GetValueFromRow(rows, "wechat", "pay_v2", "app_id", "", debug),
		ApiKey:          GetValueFromRow(rows, "wechat", "pay_v2", "api_key", "", debug),
		MchID:           GetValueFromRow(rows, "wechat", "pay_v2", "mch_id", "", debug),
		CertURI:         GetValueFromRow(rows, "wechat", "pay_v2", "cert_uri", "", debug),
		KeyURI:          GetValueFromRow(rows, "wechat", "pay_v2", "key_uri", "", debug),
		P12CertFilePath: GetValueFromRow(rows, "wechat", "pay_v2", "p12_cert_file_path", "", debug),
		CertSerialNo:    GetValueFromRow(rows, "wechat", "pay_v2", "cert_serial_no", "", debug),
		NotifyURL:       GetValueFromRow(rows, "wechat", "pay_v2", "notify_url", "", debug),
		RefundNotifyURL: GetValueFromRow(rows, "wechat", "pay_v2", "refund_notify_url", "", debug),
		Debug:           GetValueFromRow(rows, "wechat", "pay_v2", "debug", "0", debug) == "1",
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("微信支付 V2 配置 AppID 为空")
	}
	if cfg.MchID == "" {
		return nil, fmt.Errorf("微信支付 V2 配置 MchID 为空")
	}
	return cfg, nil
}

func getAliPayConfig(rows []*DbSettingRow, debug bool) (*configStruct.AliPayConfig, error) {
	cfg := &configStruct.AliPayConfig{
		AppID:            GetValueFromRow(rows, "ali", "pay", "app_id", "", debug),
		PrivateKey:       GetValueFromRow(rows, "ali", "pay", "private_key", "", debug),
		NotifyURL:        GetValueFromRow(rows, "ali", "pay", "notify_url", "", debug),
		Debug:            GetValueFromRow(rows, "ali", "pay", "debug", "0", debug) == "1",
		AppCertPublicKey: GetValueFromRow(rows, "ali", "pay", "app_cert_public_key", "", debug),
		CertPublicKey:    GetValueFromRow(rows, "ali", "pay", "cert_public_key", "", debug),
		RootCert:         GetValueFromRow(rows, "ali", "pay", "root_cert", "", debug),
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("支付宝配置 AppID 为空")
	}
	if cfg.PrivateKey == "" {
		return nil, fmt.Errorf("支付宝配置 PrivateKey 为空")
	}
	return cfg, nil
}

func getRedisConfig(rows []*DbSettingRow, debug bool) (*configStruct.RedisConfig, error) {
	cfg := &configStruct.RedisConfig{
		Host:        GetValueFromRow(rows, "redis", "host", "", "127.0.0.1", debug),
		Port:        GetValueFromRow(rows, "redis", "port", "", "6379", debug),
		Username:    GetValueFromRow(rows, "redis", "username", "", "", debug),
		Password:    GetValueFromRow(rows, "redis", "password", "", "", debug),
		IdleTimeout: typeHelper.Str2Int(GetValueFromRow(rows, "redis", "idle_timeout", "", "60", debug)),
		Database:    typeHelper.Str2Int(GetValueFromRow(rows, "redis", "database", "", "", debug)),
		MaxActive:   typeHelper.Str2Int(GetValueFromRow(rows, "redis", "max_active", "", "10", debug)),
		MaxIdle:     typeHelper.Str2Int(GetValueFromRow(rows, "redis", "max_idle", "", "10", debug)),
	}
	if cfg.Host == "" {
		return nil, fmt.Errorf("Redis 配置 Host 为空")
	}
	if cfg.Port == "" {
		return nil, fmt.Errorf("Redis 配置 Port 为空")
	}
	return cfg, nil
}

func getEsConfig(rows []*DbSettingRow, debug bool) (*configStruct.EsConfig, error) {
	cfg := &configStruct.EsConfig{
		Host:     GetValueFromRow(rows, "es", "host", "", "http://127.0.0.1", debug),
		Port:     GetValueFromRow(rows, "es", "port", "", "9200", debug),
		Password: GetValueFromRow(rows, "es", "password", "", "123456", debug),
		Username: GetValueFromRow(rows, "es", "username", "", "es", debug),
	}
	if cfg.Host == "" {
		return nil, fmt.Errorf("Elasticsearch 配置 Host 为空")
	}
	if cfg.Port == "" {
		return nil, fmt.Errorf("Elasticsearch 配置 Port 为空")
	}
	return cfg, nil
}

func getRabbitMQConfig(rows []*DbSettingRow, debug bool) (*configStruct.RabbitMQConfig, error) {
	cfg := &configStruct.RabbitMQConfig{
		Host:     GetValueFromRow(rows, "rabbit_mq", "host", "", "http://127.0.0.1", debug),
		Password: GetValueFromRow(rows, "rabbit_mq", "password", "", "123456", debug),
		Username: GetValueFromRow(rows, "rabbit_mq", "username", "", "root", debug),
	}
	if cfg.Host == "" {
		return nil, fmt.Errorf("RabbitMQ 配置 Host 为空")
	}
	return cfg, nil
}

func getPostgresConfig(rows []*DbSettingRow, debug bool) (*configStruct.PostgresConfig, error) {
	dsn := GetValueFromRow(rows, "postgres", "", "dsn", "", debug)
	if dsn == "" {
		dsn = GetValueFromRow(rows, "postgres", "", "", "", debug)
	}
	if dsn == "" {
		return nil, nil // DSN 为空时返回 nil，这是允许的
	}

	cfg := &configStruct.PostgresConfig{DSN: dsn}
	if cfg.DSN == "" {
		return nil, fmt.Errorf("PostgreSQL 配置 DSN 为空")
	}
	if v := GetValueFromRow(rows, "postgres", "", "conn_max_idle", "", debug); v != "" {
		cfg.ConnMaxIdle = typeHelper.Str2Int(v)
	}
	if v := GetValueFromRow(rows, "postgres", "", "conn_max_open", "", debug); v != "" {
		cfg.ConnMaxOpen = typeHelper.Str2Int(v)
	}
	if v := GetValueFromRow(rows, "postgres", "", "conn_max_lifetime", "", debug); v != "" {
		cfg.ConnMaxLifetime = time.Duration(typeHelper.Str2Int64(v)) * time.Second
	}
	return cfg, nil
}

func getMysqlConfig(rows []*DbSettingRow, debug bool) (*configStruct.MysqlConfig, error) {
	cfg := &configStruct.MysqlConfig{
		Host:             GetValueFromRow(rows, "mysql", "", "host", "127.0.0.1", debug),
		Port:             GetValueFromRow(rows, "mysql", "", "port", "3306", debug),
		Username:         GetValueFromRow(rows, "mysql", "", "username", "", debug),
		Password:         GetValueFromRow(rows, "mysql", "", "password", "", debug),
		DbName:           GetValueFromRow(rows, "mysql", "", "db_name", "", debug),
		Charset:          GetValueFromRow(rows, "mysql", "", "charset", "utf8mb4", debug),
		Collation:        GetValueFromRow(rows, "mysql", "", "collation", "utf8mb4_unicode_ci", debug),
		SettingTableName: GetValueFromRow(rows, "mysql", "", "setting_table_name", "u_setting", debug),
		TimeZone:         GetValueFromRow(rows, "mysql", "", "time_zone", "Asia/Shanghai", debug),
		ParseTime:        GetValueFromRow(rows, "mysql", "", "parse_time", "true", debug) == "true",
	}
	if cfg.Host == "" {
		return nil, fmt.Errorf("MySQL 配置 Host 为空")
	}
	if cfg.Port == "" {
		return nil, fmt.Errorf("MySQL 配置 Port 为空")
	}
	if cfg.Username == "" {
		return nil, fmt.Errorf("MySQL 配置 Username 为空")
	}
	if cfg.DbName == "" {
		return nil, fmt.Errorf("MySQL 配置 DbName 为空")
	}
	if v := GetValueFromRow(rows, "mysql", "", "max_open_conns", "", debug); v != "" {
		cfg.MaxOpenConns = typeHelper.Str2Int(v)
	}
	if v := GetValueFromRow(rows, "mysql", "", "max_idle", "", debug); v != "" {
		cfg.MaxIdle = typeHelper.Str2Int(v)
	}
	if v := GetValueFromRow(rows, "mysql", "", "max_life_time", "", debug); v != "" {
		cfg.MaxLifeTime = typeHelper.Str2Int(v)
	}
	return cfg, nil
}

func getAliOssConfig(rows []*DbSettingRow, debug bool) (*configStruct.AliOssConfig, error) {
	cfg := &configStruct.AliOssConfig{
		AccessKeyID:     GetValueFromRow(rows, "ali", "oss", "access_key_id", "", debug),
		AccessKeySecret: GetValueFromRow(rows, "ali", "oss", "access_key_secret", "", debug),
		Host:            GetValueFromRow(rows, "ali", "oss", "host", "", debug),
		EndPoint:        GetValueFromRow(rows, "ali", "oss", "end_point", "", debug),
		BucketName:      GetValueFromRow(rows, "ali", "oss", "bucket_name", "", debug),
		ExpireTime:      typeHelper.Str2Int64(GetValueFromRow(rows, "ali", "oss", "expire_time", "30", debug)),
		ARN:             GetValueFromRow(rows, "ali", "oss", "arn", "", debug),
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("阿里云 OSS 配置 AccessKeyID 为空")
	}
	if cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("阿里云 OSS 配置 AccessKeySecret 为空")
	}
	return cfg, nil
}

func getAliApiConfig(rows []*DbSettingRow, debug bool) (*configStruct.AliApiConfig, error) {
	cfg := &configStruct.AliApiConfig{
		AppKey:    GetValueFromRow(rows, "ali", "api", "app_key", "", debug),
		AppSecret: GetValueFromRow(rows, "ali", "api", "app_secret", "", debug),
		AppCode:   GetValueFromRow(rows, "ali", "api", "app_code", "", debug),
	}
	if cfg.AppKey == "" {
		return nil, fmt.Errorf("阿里云 API 配置 AppKey 为空")
	}
	if cfg.AppSecret == "" {
		return nil, fmt.Errorf("阿里云 API 配置 AppSecret 为空")
	}
	if cfg.AppCode == "" {
		return nil, fmt.Errorf("阿里云 API 配置 AppCode 为空")
	}
	return cfg, nil
}

func getAliSmsConfig(rows []*DbSettingRow, debug bool) (*configStruct.AliSmsConfig, error) {
	cfg := &configStruct.AliSmsConfig{
		AccessKeySecret: GetValueFromRow(rows, "ali", "sms", "access_key_secret", "", debug),
		AccessKeyID:     GetValueFromRow(rows, "ali", "sms", "access_key_id", "", debug),
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("阿里云短信配置 AccessKeyID 为空")
	}
	if cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("阿里云短信配置 AccessKeySecret 为空")
	}
	return cfg, nil
}

func getAliIotConfig(rows []*DbSettingRow, debug bool) (*configStruct.AliIotConfig, error) {
	cfg := &configStruct.AliIotConfig{
		AccessKeySecret: GetValueFromRow(rows, "ali", "iot", "access_key_secret", "", debug),
		AccessKeyID:     GetValueFromRow(rows, "ali", "iot", "access_key_id", "", debug),
		EndPoint:        GetValueFromRow(rows, "ali", "iot", "end_point", "", debug),
		RegionID:        GetValueFromRow(rows, "ali", "iot", "region_id", "", debug),
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("阿里云 IoT 配置 AccessKeyID 为空")
	}
	if cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("阿里云 IoT 配置 AccessKeySecret 为空")
	}
	return cfg, nil
}

func getAmapConfig(rows []*DbSettingRow, debug bool) (*configStruct.AmapConfig, error) {
	cfg := &configStruct.AmapConfig{Key: GetValueFromRow(rows, "ali", "amap", "key", "", debug)}
	if cfg.Key == "" {
		return nil, fmt.Errorf("高德地图配置 Key 为空")
	}
	return cfg, nil
}

func getYunXinConfig(rows []*DbSettingRow, debug bool) (*configStruct.YunxinConfig, error) {
	cfg := &configStruct.YunxinConfig{
		AppKey:    GetValueFromRow(rows, "netease", "yunxin", "app_key", "", debug),
		AppSecret: GetValueFromRow(rows, "netease", "yunxin", "app_secret", "", debug),
	}
	if cfg.AppKey == "" {
		return nil, fmt.Errorf("网易云信配置 AppKey 为空")
	}
	if cfg.AppSecret == "" {
		return nil, fmt.Errorf("网易云信配置 AppSecret 为空")
	}
	return cfg, nil
}
