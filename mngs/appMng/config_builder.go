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
	"net/url"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/wiidz/goutil/helpers/configHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConfigSource 配置来源类型
type ConfigSource string

const (
	SourceYAML     ConfigSource = "yaml"     // 从 YAML 文件加载
	SourceDatabase ConfigSource = "database" // 从数据库加载
)

// 配置键名常量（用于 YAML 和数据库配置）
const (
	// YAML 配置键（用于 UnmarshalKey）
	ConfigKeyProfile     = "profile"           // Profile 配置键
	ConfigKeyLocation    = "location"          // Location 配置键
	ConfigKeyRedis       = "redis"             // Redis 配置键
	ConfigKeyEs          = "es"                // Elasticsearch 配置键
	ConfigKeyRabbitMQ    = "rabbit_mq"         // RabbitMQ 配置键
	ConfigKeyPostgres    = "postgres"          // PostgreSQL 配置键
	ConfigKeyMysql       = "mysql"             // MySQL 配置键
	ConfigKeyWechatMini  = "wechat_mini"       // 微信小程序配置键
	ConfigKeyWechatOa    = "wechat_oa"         // 微信公众号配置键
	ConfigKeyWechatOpen  = "wechat_open"       // 微信开放平台配置键
	ConfigKeyWechatPayV3 = "wechat_pay_v3"     // 微信支付 V3 配置键
	ConfigKeyWechatPayV2 = "wechat_pay_v2"     // 微信支付 V2 配置键
	ConfigKeyAliOss      = "ali_oss"           // 阿里云 OSS 配置键
	ConfigKeyAliPay      = "ali_pay"           // 支付宝配置键
	ConfigKeyAliApi      = "ali_api"           // 阿里云 API 配置键
	ConfigKeyAliSms      = "ali_sms"           // 阿里云短信配置键
	ConfigKeyAliIot      = "ali_iot"           // 阿里云 IoT 配置键
	ConfigKeyAmap        = "amap"              // 高德地图配置键
	ConfigKeyYunxin      = "yunxin"            // 网易云信配置键
	ConfigKeyLocationTZ  = "location.timezone" // Location 时区配置键

	// 数据库配置键（用于 GetValueFromRow 的 name 参数）
	ConfigKeyApp      = "app"       // 应用配置键
	ConfigKeyTimeZone = "time_zone" // 时区配置键
	ConfigKeyWechat   = "wechat"    // 微信配置键
	ConfigKeyAli      = "ali"       // 阿里配置键
	ConfigKeyNetease  = "netease"   // 网易配置键

	// 数据库配置子键（用于 GetValueFromRow 的 flag1 参数）
	ConfigKeyWechatMiniFlag  = "mini"   // 微信小程序子键
	ConfigKeyWechatOaFlag    = "oa"     // 微信公众号子键
	ConfigKeyWechatOpenFlag  = "open"   // 微信开放平台子键
	ConfigKeyWechatPayV3Flag = "pay_v3" // 微信支付 V3 子键
	ConfigKeyWechatPayV2Flag = "pay_v2" // 微信支付 V2 子键
	ConfigKeyAliPayFlag      = "pay"    // 支付宝子键
	ConfigKeyAliOssFlag      = "oss"    // 阿里云 OSS 子键
	ConfigKeyAliApiFlag      = "api"    // 阿里云 API 子键
	ConfigKeyAliSmsFlag      = "sms"    // 阿里云短信子键
	ConfigKeyAliIotFlag      = "iot"    // 阿里云 IoT 子键
	ConfigKeyAliAmapFlag     = "amap"   // 高德地图子键
	ConfigKeyYunxinFlag      = "yunxin" // 网易云信子键
)

// errConfigFromDatabaseEmpty 生成从数据库加载配置失败的错误信息（数据库配置行为空）
func errConfigFromDatabaseEmpty(configName string) string {
	return fmt.Sprintf("策略要求从数据库加载%s配置，但数据库配置行为空", configName)
}

// errConfigFromYAMLNotInit 生成从 YAML 加载配置失败的错误信息（YAML 配置未初始化）
func errConfigFromYAMLNotInit(configName string) string {
	return fmt.Sprintf("策略要求从 YAML 加载%s配置，但 YAML 配置未初始化", configName)
}

// errConfigLoadFailed 生成从 YAML 加载配置失败的错误信息（YAML 配置未初始化）
func errConfigLoadFailed(configName string, err error) string {
	return fmt.Sprintf("加载 %s 配置失败: %w", configName, err)
}

// configLoadSuccess 加载成功
func configLoadSuccess(configName string) string {
	return fmt.Sprintf("加载 %s 配置成功", configName)
}

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

// needDatabaseConnection 判断是否需要数据库连接
// 返回 true 如果：
// 1. strategy 中有任何配置项需要从数据库读取（SourceDatabase）
// 2. 或者 initialConfig.DB 不为 nil
func (b *ConfigBuilder) needDatabaseConnection() bool {
	// 检查 initialConfig.DB 是否不为 nil
	if b.initialConfig != nil && b.initialConfig.DB != nil {
		return true
	}

	// 检查 strategy 中是否有任何配置项需要从数据库读取
	if b.strategy == nil {
		return false
	}

	// 检查所有配置项是否从数据库加载
	return b.strategy.Profile == SourceDatabase ||
		b.strategy.Location == SourceDatabase ||
		b.strategy.Redis == SourceDatabase ||
		b.strategy.Es == SourceDatabase ||
		b.strategy.RabbitMQ == SourceDatabase ||
		b.strategy.Postgres == SourceDatabase ||
		b.strategy.Mysql == SourceDatabase ||
		b.strategy.WechatMini == SourceDatabase ||
		b.strategy.WechatOa == SourceDatabase ||
		b.strategy.WechatOpen == SourceDatabase ||
		b.strategy.WechatPayV3 == SourceDatabase ||
		b.strategy.WechatPayV2 == SourceDatabase ||
		b.strategy.AliOss == SourceDatabase ||
		b.strategy.AliPay == SourceDatabase ||
		b.strategy.AliApi == SourceDatabase ||
		b.strategy.AliSms == SourceDatabase ||
		b.strategy.AliIot == SourceDatabase ||
		b.strategy.Amap == SourceDatabase ||
		b.strategy.Yunxin == SourceDatabase
}

// Build 构建 BaseConfig，根据策略从不同来源加载配置
func (b *ConfigBuilder) Build(ctx context.Context) (*configStruct.BaseConfig, error) {
	cfg := &configStruct.BaseConfig{}

	// 第一步：检查策略中是否有数据库相关的配置，如果有，优先初始化数据库
	// 这样后续的配置才能从数据库中读取
	// 注意：如果策略要求从数据库加载配置，b.db 在 Build 开始时必然是 nil，需要先初始化
	needDatabase := b.needDatabaseConnection()
	if needDatabase {
		// 需要数据库连接，从 YAML 或 InitialConfig.DB 初始化
		// b.db 在 Build 开始时必然是 nil，因为数据库连接是在这里初始化的
		if err := b.initDatabaseFromConfig(); err != nil {
			return nil, fmt.Errorf("初始化数据库连接失败: %w", err)
		}
	}

	// 第二步：加载数据库配置行（如果需要）
	var dbRows []*DbSettingRow
	if b.db != nil {
		var err error
		dbRows, err = b.loadDatabaseRows(ctx)
		if err != nil {
			log.Printf("警告: 无法从数据库加载配置: %v", err)
		} else {
			log.Printf("成功: 从数据库加载了 %d 条配置", len(dbRows))
		}
	}

	// 第三步：加载 Profile 和 Location（基础配置）
	if err := b.loadProfileAndLocation(cfg, dbRows); err != nil {
		return nil, fmt.Errorf("加载基础配置失败: %w", err)
	}

	debug := cfg.Profile != nil && cfg.Profile.Debug

	// 第四步：根据策略加载各个配置项（如果策略中指定了配置来源，则必须成功加载）
	// 注意：数据库配置（Postgres/Mysql）如果策略要求从数据库加载，此时应该已经可以从 dbRows 中读取了
	if b.strategy.Redis != "" {
		if err := b.loadRedisConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载 Redis 配置失败: %w", err)
		}
		log.Printf("成功: Redis 配置已加载")
	}
	if b.strategy.Es != "" {
		if err := b.loadEsConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载 Elasticsearch 配置失败: %w", err)
		}
		log.Printf("成功: Elasticsearch 配置已加载")
	}
	if b.strategy.RabbitMQ != "" {
		if err := b.loadRabbitMQConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载 RabbitMQ 配置失败: %w", err)
		}
		log.Printf("成功: RabbitMQ 配置已加载")
	}
	if b.strategy.Postgres != "" {
		if err := b.loadPostgresConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载 PostgreSQL 配置失败: %w", err)
		}
		log.Printf("成功: PostgreSQL 配置已加载")
	}
	if b.strategy.Mysql != "" {
		if err := b.loadMysqlConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载 MySQL 配置失败: %w", err)
		}
		log.Printf("成功: MySQL 配置已加载")
	}
	if b.strategy.WechatMini != "" {
		if err := b.loadWechatMiniConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载微信小程序配置失败: %w", err)
		}
		log.Printf("成功: 微信小程序配置已加载")
	}
	if b.strategy.WechatOa != "" {
		if err := b.loadWechatOaConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载微信公众号配置失败: %w", err)
		}
		log.Printf("成功: 微信公众号配置已加载")
	}
	if b.strategy.WechatOpen != "" {
		if err := b.loadWechatOpenConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载微信开放平台配置失败: %w", err)
		}
		log.Printf("成功: 微信开放平台配置已加载")
	}
	if b.strategy.WechatPayV3 != "" {
		if err := b.loadWechatPayV3Config(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载微信支付 V3 配置失败: %w", err)
		}
		log.Printf("成功: 微信支付 V3 配置已加载")
	}
	if b.strategy.WechatPayV2 != "" {
		if err := b.loadWechatPayV2Config(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载微信支付 V2 配置失败: %w", err)
		}
		log.Printf("成功: 微信支付 V2 配置已加载")
	}
	if b.strategy.AliOss != "" {
		if err := b.loadAliOssConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载阿里云 OSS 配置失败: %w", err)
		}
		log.Printf("成功: 阿里云 OSS 配置已加载")
	}
	if b.strategy.AliPay != "" {
		if err := b.loadAliPayConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载支付宝配置失败: %w", err)
		}
		log.Printf("成功: 支付宝配置已加载")
	}
	if b.strategy.AliApi != "" {
		if err := b.loadAliApiConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载阿里云 API 配置失败: %w", err)
		}
		log.Printf("成功: 阿里云 API 配置已加载")
	}
	if b.strategy.AliSms != "" {
		if err := b.loadAliSmsConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载阿里云短信配置失败: %w", err)
		}
		log.Printf("成功: 阿里云短信配置已加载")
	}
	if b.strategy.AliIot != "" {
		if err := b.loadAliIotConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载阿里云 IoT 配置失败: %w", err)
		}
		log.Printf("成功: 阿里云 IoT 配置已加载")
	}
	if b.strategy.Amap != "" {
		if err := b.loadAmapConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载高德地图配置失败: %w", err)
		}
		log.Printf("成功: 高德地图配置已加载")
	}
	if b.strategy.Yunxin != "" {
		if err := b.loadYunxinConfig(cfg, dbRows, debug); err != nil {
			return nil, fmt.Errorf("加载网易云信配置失败: %w", err)
		}
		log.Printf("成功: 网易云信配置已加载")
	}

	return cfg, nil
}

// initDatabaseFromConfig 从 YAML 或 InitialConfig.DB 初始化数据库连接
// 当策略要求从数据库加载配置时，需要先初始化数据库连接
func (b *ConfigBuilder) initDatabaseFromConfig() error {
	// 优先使用 InitialConfig.DB 中的配置
	if b.initialConfig.DB != nil {
		return b.initDatabaseFromDBConfig(b.initialConfig.DB)
	}

	// 如果没有 InitialConfig.DB，尝试从 YAML 加载数据库配置
	// 根据策略中需要的数据库类型，尝试从 YAML 加载对应的配置
	if len(b.yamlVipers) > 0 {
		// 如果策略要求从数据库加载 PostgreSQL 配置，尝试从 YAML 加载 PostgreSQL 配置来初始化连接
		if b.strategy.Postgres == SourceDatabase {
			var postgresConfig configStruct.PostgresConfig
			if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyPostgres, &postgresConfig); err == nil && postgresConfig.DSN != "" {
				return b.initPostgresFromConfig(&postgresConfig)
			}
		}
		// 如果策略要求从数据库加载 MySQL 配置，尝试从 YAML 加载 MySQL 配置来初始化连接
		if b.strategy.Mysql == SourceDatabase {
			var mysqlConfig configStruct.MysqlConfig
			if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyMysql, &mysqlConfig); err == nil {
				return b.initMysqlFromConfig(&mysqlConfig)
			}
		}
	}

	return fmt.Errorf("无法初始化数据库连接：未找到数据库配置（需要 InitialConfig.DB 或 YAML 中的数据库配置）")
}

// initDatabaseFromDBConfig 从 DBConfig 初始化数据库连接
func (b *ConfigBuilder) initDatabaseFromDBConfig(dbConfig *configStruct.DBConfig) error {
	if dbConfig == nil {
		return fmt.Errorf("DBConfig 为空")
	}

	switch dbConfig.Type {
	case configStruct.DBTypePostgres:
		// 构建 PostgreSQL 配置
		dsn := dbConfig.DSN
		if dsn == "" {
			return fmt.Errorf("PostgreSQL DSN 为空")
		}
		return b.initPostgresFromDSN(dsn, dbConfig.ConnMaxIdle, dbConfig.ConnMaxOpen, dbConfig.ConnMaxLifetime, dbConfig.Logger)

	case configStruct.DBTypeMysql:
		// 构建 MySQL 配置
		dsn := dbConfig.DSN
		if dsn == "" {
			// 从 Host/Port 等字段构建 DSN
			if dbConfig.Host == "" || dbConfig.Port == "" || dbConfig.Username == "" || dbConfig.DbName == "" {
				return fmt.Errorf("MySQL 配置不完整：需要 Host, Port, Username, DbName")
			}
			charset := dbConfig.Charset
			if charset == "" {
				charset = "utf8mb4"
			}
			collation := dbConfig.Collation
			if collation == "" {
				collation = "utf8mb4_unicode_ci"
			}
			timeZone := dbConfig.TimeZone
			if timeZone == "" {
				timeZone = "Asia/Shanghai"
			}
			parseTime := dbConfig.ParseTime
			if !parseTime {
				parseTime = true
			}
			dsn = dbConfig.Username + ":" + dbConfig.Password +
				"@tcp(" + dbConfig.Host + ":" + dbConfig.Port + ")/" + dbConfig.DbName +
				"?charset=" + charset +
				"&collation=" + collation +
				"&loc=" + url.QueryEscape(timeZone) +
				"&parseTime=" + strconv.FormatBool(parseTime)
		}
		return b.initMysqlFromDSN(dsn, dbConfig.ConnMaxIdle, dbConfig.ConnMaxOpen, dbConfig.ConnMaxLifetime, dbConfig.Logger)

	default:
		return fmt.Errorf("不支持的数据库类型: %s", dbConfig.Type)
	}
}

// initPostgresFromDSN 从 DSN 初始化 PostgreSQL 连接
func (b *ConfigBuilder) initPostgresFromDSN(dsn string, maxIdle, maxOpen int, maxLifetime time.Duration, loggerInterface logger.Interface) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: loggerInterface,
	})
	if err != nil {
		return fmt.Errorf("连接 PostgreSQL 失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	if maxIdle > 0 {
		sqlDB.SetMaxIdleConns(maxIdle)
	}
	if maxOpen > 0 {
		sqlDB.SetMaxOpenConns(maxOpen)
	}
	if maxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(maxLifetime)
	}

	b.db = db
	log.Printf("成功: PostgreSQL 数据库连接已初始化")
	return nil
}

// initPostgresFromConfig 从 PostgresConfig 初始化 PostgreSQL 连接
func (b *ConfigBuilder) initPostgresFromConfig(postgresConfig *configStruct.PostgresConfig) error {
	if postgresConfig == nil || postgresConfig.DSN == "" {
		return fmt.Errorf("PostgreSQL 配置 DSN 为空")
	}
	return b.initPostgresFromDSN(postgresConfig.DSN, postgresConfig.ConnMaxIdle, postgresConfig.ConnMaxOpen, postgresConfig.ConnMaxLifetime, postgresConfig.Logger)
}

// initMysqlFromDSN 从 DSN 初始化 MySQL 连接
func (b *ConfigBuilder) initMysqlFromDSN(dsn string, maxIdle, maxOpen int, maxLifetime time.Duration, loggerInterface logger.Interface) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: loggerInterface,
	})
	if err != nil {
		return fmt.Errorf("连接 MySQL 失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	if maxIdle > 0 {
		sqlDB.SetMaxIdleConns(maxIdle)
	}
	if maxOpen > 0 {
		sqlDB.SetMaxOpenConns(maxOpen)
	}
	if maxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(maxLifetime)
	}

	b.db = db
	log.Printf("成功: MySQL 数据库连接已初始化")
	return nil
}

// initMysqlFromConfig 从 MysqlConfig 初始化 MySQL 连接
func (b *ConfigBuilder) initMysqlFromConfig(mysqlConfig *configStruct.MysqlConfig) error {
	if mysqlConfig == nil {
		return fmt.Errorf("MySQL 配置为空")
	}

	// 构建 DSN
	charset := mysqlConfig.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	collation := mysqlConfig.Collation
	if collation == "" {
		collation = "utf8mb4_unicode_ci"
	}
	timeZone := mysqlConfig.TimeZone
	if timeZone == "" {
		timeZone = "Asia/Shanghai"
	}
	parseTime := mysqlConfig.ParseTime
	if !parseTime {
		parseTime = true
	}

	dsn := mysqlConfig.Username + ":" + mysqlConfig.Password +
		"@tcp(" + mysqlConfig.Host + ":" + mysqlConfig.Port + ")/" + mysqlConfig.DbName +
		"?charset=" + charset +
		"&collation=" + collation +
		"&loc=" + url.QueryEscape(timeZone) +
		"&parseTime=" + strconv.FormatBool(parseTime)

	maxIdle := mysqlConfig.MaxIdle
	maxOpen := mysqlConfig.MaxOpenConns
	maxLifetime := time.Duration(mysqlConfig.MaxLifeTime) * time.Second

	return b.initMysqlFromDSN(dsn, maxIdle, maxOpen, maxLifetime, mysqlConfig.Logger)
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
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyProfile, &profile); err == nil {
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
		timeZone := b.yamlVipers[0].GetString(ConfigKeyLocationTZ)
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
	switch b.strategy.Redis {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("Redis"))
		}
		redisConfig, err := getRedisConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.RedisConfig = redisConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("Redis"))
		}
		var redisConfig configStruct.RedisConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyRedis, &redisConfig); err != nil {
			return fmt.Errorf("从 YAML 加载 Redis 配置失败: %w", err)
		}
		cfg.RedisConfig = &redisConfig
	}
	return nil
}

func (b *ConfigBuilder) loadEsConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.Es {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("Elasticsearch"))
		}
		esConfig, err := getEsConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.EsConfig = esConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("Elasticsearch"))
		}
		var esConfig configStruct.EsConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyEs, &esConfig); err != nil {
			return fmt.Errorf("从 YAML 加载 Elasticsearch 配置失败: %w", err)
		}
		cfg.EsConfig = &esConfig
	}
	return nil
}

func (b *ConfigBuilder) loadRabbitMQConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.RabbitMQ {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("RabbitMQ"))
		}
		rabbitMQConfig, err := getRabbitMQConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.RabbitMQConfig = rabbitMQConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("RabbitMQ"))
		}
		var rabbitMQConfig configStruct.RabbitMQConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyRabbitMQ, &rabbitMQConfig); err != nil {
			return fmt.Errorf("从 YAML 加载 RabbitMQ 配置失败: %w", err)
		}
		cfg.RabbitMQConfig = &rabbitMQConfig
	}
	return nil
}

func (b *ConfigBuilder) loadPostgresConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.Postgres {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("PostgreSQL"))
		}
		postgresConfig, err := getPostgresConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.PostgresConfig = postgresConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("PostgreSQL"))
		}
		var postgresConfig configStruct.PostgresConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyPostgres, &postgresConfig); err != nil {
			return fmt.Errorf("从 YAML 加载 PostgreSQL 配置失败: %w", err)
		}
		cfg.PostgresConfig = &postgresConfig
	}
	return nil
}

func (b *ConfigBuilder) loadMysqlConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.Mysql {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("MySQL"))
		}
		mysqlConfig, err := getMysqlConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.MysqlConfig = mysqlConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("MySQL"))
		}
		var mysqlConfig configStruct.MysqlConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyMysql, &mysqlConfig); err != nil {
			return fmt.Errorf("从 YAML 加载 MySQL 配置失败: %w", err)
		}
		cfg.MysqlConfig = &mysqlConfig
	}
	return nil
}

func (b *ConfigBuilder) loadWechatMiniConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.WechatMini {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("微信小程序"))
		}
		wechatMiniConfig, err := getWechatMiniConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.WechatMiniConfig = wechatMiniConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("微信小程序"))
		}
		var wechatMiniConfig configStruct.WechatMiniConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyWechatMini, &wechatMiniConfig); err != nil {
			return fmt.Errorf("从 YAML 加载微信小程序配置失败: %w", err)
		}
		cfg.WechatMiniConfig = &wechatMiniConfig
	}
	return nil
}

func (b *ConfigBuilder) loadWechatOaConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.WechatOa {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("微信公众号"))
		}
		wechatOaConfig, err := getWechatOaConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.WechatOaConfig = wechatOaConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("微信公众号"))
		}
		var wechatOaConfig configStruct.WechatOaConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyWechatOa, &wechatOaConfig); err != nil {
			return fmt.Errorf("从 YAML 加载微信公众号配置失败: %w", err)
		}
		cfg.WechatOaConfig = &wechatOaConfig
	}
	return nil
}

func (b *ConfigBuilder) loadWechatOpenConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.WechatOpen {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("微信开放平台"))
		}
		wechatOpenConfig, err := getWechatOpenConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.WechatOpenConfig = wechatOpenConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("微信开放平台"))
		}
		var wechatOpenConfig configStruct.WechatOpenConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyWechatOpen, &wechatOpenConfig); err != nil {
			return fmt.Errorf("从 YAML 加载微信开放平台配置失败: %w", err)
		}
		cfg.WechatOpenConfig = &wechatOpenConfig
	}
	return nil
}

func (b *ConfigBuilder) loadWechatPayV3Config(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.WechatPayV3 {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("微信支付 V3"))
		}
		wechatPayV3Config, err := getWechatPayConfigV3(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.WechatPayConfigV3 = wechatPayV3Config
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("微信支付 V3"))
		}
		var wechatPayV3Config configStruct.WechatPayConfigV3
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyWechatPayV3, &wechatPayV3Config); err != nil {
			return fmt.Errorf("从 YAML 加载微信支付 V3 配置失败: %w", err)
		}
		cfg.WechatPayConfigV3 = &wechatPayV3Config
	}
	return nil
}

func (b *ConfigBuilder) loadWechatPayV2Config(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.WechatPayV2 {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("微信支付 V2"))
		}
		wechatPayV2Config, err := getWechatPayConfigV2(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.WechatPayConfigV2 = wechatPayV2Config
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("微信支付 V2"))
		}
		var wechatPayV2Config configStruct.WechatPayConfigV2
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyWechatPayV2, &wechatPayV2Config); err != nil {
			return fmt.Errorf("从 YAML 加载微信支付 V2 配置失败: %w", err)
		}
		cfg.WechatPayConfigV2 = &wechatPayV2Config
	}
	return nil
}

func (b *ConfigBuilder) loadAliOssConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.AliOss {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("阿里云 OSS"))
		}
		aliOssConfig, err := getAliOssConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.AliOssConfig = aliOssConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("阿里云 OSS"))
		}
		var aliOssConfig configStruct.AliOssConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyAliOss, &aliOssConfig); err != nil {
			return fmt.Errorf("从 YAML 加载阿里云 OSS 配置失败: %w", err)
		}
		cfg.AliOssConfig = &aliOssConfig
	}
	return nil
}

func (b *ConfigBuilder) loadAliPayConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.AliPay {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("支付宝"))
		}
		aliPayConfig, err := getAliPayConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.AliPayConfig = aliPayConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("支付宝"))
		}
		var aliPayConfig configStruct.AliPayConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyAliPay, &aliPayConfig); err != nil {
			return fmt.Errorf("从 YAML 加载支付宝配置失败: %w", err)
		}
		cfg.AliPayConfig = &aliPayConfig
	}
	return nil
}

func (b *ConfigBuilder) loadAliApiConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.AliApi {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("阿里云 API"))
		}
		aliApiConfig, err := getAliApiConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.AliApiConfig = aliApiConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("阿里云 API"))
		}
		var aliApiConfig configStruct.AliApiConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyAliApi, &aliApiConfig); err != nil {
			return fmt.Errorf("从 YAML 加载阿里云 API 配置失败: %w", err)
		}
		cfg.AliApiConfig = &aliApiConfig
	}
	return nil
}

func (b *ConfigBuilder) loadAliSmsConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.AliSms {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("阿里云短信"))
		}
		aliSmsConfig, err := getAliSmsConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.AliSmsConfig = aliSmsConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("阿里云短信"))
		}
		var aliSmsConfig configStruct.AliSmsConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyAliSms, &aliSmsConfig); err != nil {
			return fmt.Errorf("从 YAML 加载阿里云短信配置失败: %w", err)
		}
		cfg.AliSmsConfig = &aliSmsConfig
	}
	return nil
}

func (b *ConfigBuilder) loadAliIotConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.AliIot {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("阿里云 IoT"))
		}
		aliIotConfig, err := getAliIotConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.AliIotConfig = aliIotConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("阿里云 IoT"))
		}
		var aliIotConfig configStruct.AliIotConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyAliIot, &aliIotConfig); err != nil {
			return fmt.Errorf("从 YAML 加载阿里云 IoT 配置失败: %w", err)
		}
		cfg.AliIotConfig = &aliIotConfig
	}
	return nil
}

func (b *ConfigBuilder) loadAmapConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.Amap {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("高德地图"))
		}
		amapConfig, err := getAmapConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.AmapConfig = amapConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("高德地图"))
		}
		var amapConfig configStruct.AmapConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyAmap, &amapConfig); err != nil {
			return fmt.Errorf("从 YAML 加载高德地图配置失败: %w", err)
		}
		cfg.AmapConfig = &amapConfig
	}
	return nil
}

func (b *ConfigBuilder) loadYunxinConfig(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
	switch b.strategy.Yunxin {
	case SourceDatabase:
		if len(dbRows) == 0 {
			return fmt.Errorf(errConfigFromDatabaseEmpty("网易云信"))
		}
		yunxinConfig, err := getYunXinConfig(dbRows, debug)
		if err != nil {
			return err
		}
		cfg.YunxinConfig = yunxinConfig
	case SourceYAML:
		if len(b.yamlVipers) == 0 {
			return fmt.Errorf(errConfigFromYAMLNotInit("网易云信"))
		}
		var yunxinConfig configStruct.YunxinConfig
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeyYunxin, &yunxinConfig); err != nil {
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
		No:      GetValueFromRow(rows, ConfigKeyApp, "", "no", "", false),
		Name:    GetValueFromRow(rows, ConfigKeyApp, "", "name", "", false),
		Host:    GetValueFromRow(rows, ConfigKeyApp, "", "host", "", false),
		Port:    GetValueFromRow(rows, ConfigKeyApp, "", "port", "127.0.0.1", false),
		Domain:  GetValueFromRow(rows, ConfigKeyApp, "", "domain", "", false),
		Debug:   GetValueFromRow(rows, ConfigKeyApp, "", "debug", "", false) == "1",
		Version: GetValueFromRow(rows, ConfigKeyApp, "", "version", "", false),
	}
}

func getLocationConfig(rows []*DbSettingRow) (location *time.Location, err error) {
	timeZone := GetValueFromRow(rows, ConfigKeyTimeZone, "", "", "Asia/Shanghai", false)
	location, err = time.LoadLocation(timeZone)
	if err != nil {
		location = time.FixedZone("CST-8", 8*3600)
	}
	return
}

func getWechatMiniConfig(rows []*DbSettingRow, debug bool) (*configStruct.WechatMiniConfig, error) {
	cfg := &configStruct.WechatMiniConfig{
		AppID:     GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatMiniFlag, "app_id", "", debug),
		AppSecret: GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatMiniFlag, "app_secret", "", debug),
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
		AppID:          GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatOaFlag, "app_id", "", debug),
		AppSecret:      GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatOaFlag, "app_secret", "", debug),
		Token:          GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatOaFlag, "token", "", debug),
		EncodingAESKey: GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatOaFlag, "encoding_aes_key", "", debug),
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
		AppID:     GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatOpenFlag, "app_id", "", debug),
		AppSecret: GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatOpenFlag, "app_secret", "", debug),
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
		AppID:                     GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV3Flag, "app_id", "", debug),
		ApiKeyV3:                  GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV3Flag, "api_key", "", debug),
		MchID:                     GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV3Flag, "mch_id", "", debug),
		CertURI:                   GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV3Flag, "cert_uri", "", debug),
		KeyURI:                    GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV3Flag, "key_uri", "", debug),
		PEMPrivateKeyContent:      GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV3Flag, "pem_private_key_content", "", debug),
		PEMCertContent:            GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV3Flag, "pem_cert_content", "", debug),
		CertSerialNo:              GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV3Flag, "cert_serial_no", "", debug),
		NotifyURL:                 GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV3Flag, "notify_url", "", debug),
		RefundNotifyURL:           GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV3Flag, "refund_notify_url", "", debug),
		MerchantTransferNotifyURL: GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV3Flag, "merchant_transfer_notify_url", "", debug),
		Debug:                     GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV3Flag, "debug", "0", debug) == "1",
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
		AppID:           GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV2Flag, "app_id", "", debug),
		ApiKey:          GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV2Flag, "api_key", "", debug),
		MchID:           GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV2Flag, "mch_id", "", debug),
		CertURI:         GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV2Flag, "cert_uri", "", debug),
		KeyURI:          GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV2Flag, "key_uri", "", debug),
		P12CertFilePath: GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV2Flag, "p12_cert_file_path", "", debug),
		CertSerialNo:    GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV2Flag, "cert_serial_no", "", debug),
		NotifyURL:       GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV2Flag, "notify_url", "", debug),
		RefundNotifyURL: GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV2Flag, "refund_notify_url", "", debug),
		Debug:           GetValueFromRow(rows, ConfigKeyWechat, ConfigKeyWechatPayV2Flag, "debug", "0", debug) == "1",
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
		AppID:            GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliPayFlag, "app_id", "", debug),
		PrivateKey:       GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliPayFlag, "private_key", "", debug),
		NotifyURL:        GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliPayFlag, "notify_url", "", debug),
		Debug:            GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliPayFlag, "debug", "0", debug) == "1",
		AppCertPublicKey: GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliPayFlag, "app_cert_public_key", "", debug),
		CertPublicKey:    GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliPayFlag, "cert_public_key", "", debug),
		RootCert:         GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliPayFlag, "root_cert", "", debug),
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
		Host:        GetValueFromRow(rows, ConfigKeyRedis, "host", "", "127.0.0.1", debug),
		Port:        GetValueFromRow(rows, ConfigKeyRedis, "port", "", "6379", debug),
		Username:    GetValueFromRow(rows, ConfigKeyRedis, "username", "", "", debug),
		Password:    GetValueFromRow(rows, ConfigKeyRedis, "password", "", "", debug),
		IdleTimeout: typeHelper.Str2Int(GetValueFromRow(rows, ConfigKeyRedis, "idle_timeout", "", "60", debug)),
		Database:    typeHelper.Str2Int(GetValueFromRow(rows, ConfigKeyRedis, "database", "", "", debug)),
		MaxActive:   typeHelper.Str2Int(GetValueFromRow(rows, ConfigKeyRedis, "max_active", "", "10", debug)),
		MaxIdle:     typeHelper.Str2Int(GetValueFromRow(rows, ConfigKeyRedis, "max_idle", "", "10", debug)),
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
		Host:     GetValueFromRow(rows, ConfigKeyEs, "host", "", "http://127.0.0.1", debug),
		Port:     GetValueFromRow(rows, ConfigKeyEs, "port", "", "9200", debug),
		Password: GetValueFromRow(rows, ConfigKeyEs, "password", "", "123456", debug),
		Username: GetValueFromRow(rows, ConfigKeyEs, "username", "", "es", debug),
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
		Host:     GetValueFromRow(rows, ConfigKeyRabbitMQ, "host", "", "http://127.0.0.1", debug),
		Password: GetValueFromRow(rows, ConfigKeyRabbitMQ, "password", "", "123456", debug),
		Username: GetValueFromRow(rows, ConfigKeyRabbitMQ, "username", "", "root", debug),
	}
	if cfg.Host == "" {
		return nil, fmt.Errorf("RabbitMQ 配置 Host 为空")
	}
	return cfg, nil
}

func getPostgresConfig(rows []*DbSettingRow, debug bool) (*configStruct.PostgresConfig, error) {
	dsn := GetValueFromRow(rows, ConfigKeyPostgres, "", "dsn", "", debug)
	if dsn == "" {
		dsn = GetValueFromRow(rows, ConfigKeyPostgres, "", "", "", debug)
	}
	if dsn == "" {
		return nil, nil // DSN 为空时返回 nil，这是允许的
	}

	cfg := &configStruct.PostgresConfig{DSN: dsn}
	if cfg.DSN == "" {
		return nil, fmt.Errorf("PostgreSQL 配置 DSN 为空")
	}
	if v := GetValueFromRow(rows, ConfigKeyPostgres, "", "conn_max_idle", "", debug); v != "" {
		cfg.ConnMaxIdle = typeHelper.Str2Int(v)
	}
	if v := GetValueFromRow(rows, ConfigKeyPostgres, "", "conn_max_open", "", debug); v != "" {
		cfg.ConnMaxOpen = typeHelper.Str2Int(v)
	}
	if v := GetValueFromRow(rows, ConfigKeyPostgres, "", "conn_max_lifetime", "", debug); v != "" {
		cfg.ConnMaxLifetime = time.Duration(typeHelper.Str2Int64(v)) * time.Second
	}
	return cfg, nil
}

func getMysqlConfig(rows []*DbSettingRow, debug bool) (*configStruct.MysqlConfig, error) {
	cfg := &configStruct.MysqlConfig{
		Host:             GetValueFromRow(rows, ConfigKeyMysql, "", "host", "127.0.0.1", debug),
		Port:             GetValueFromRow(rows, ConfigKeyMysql, "", "port", "3306", debug),
		Username:         GetValueFromRow(rows, ConfigKeyMysql, "", "username", "", debug),
		Password:         GetValueFromRow(rows, ConfigKeyMysql, "", "password", "", debug),
		DbName:           GetValueFromRow(rows, ConfigKeyMysql, "", "db_name", "", debug),
		Charset:          GetValueFromRow(rows, ConfigKeyMysql, "", "charset", "utf8mb4", debug),
		Collation:        GetValueFromRow(rows, ConfigKeyMysql, "", "collation", "utf8mb4_unicode_ci", debug),
		SettingTableName: GetValueFromRow(rows, ConfigKeyMysql, "", "setting_table_name", "u_setting", debug),
		TimeZone:         GetValueFromRow(rows, ConfigKeyMysql, "", "time_zone", "Asia/Shanghai", debug),
		ParseTime:        GetValueFromRow(rows, ConfigKeyMysql, "", "parse_time", "true", debug) == "true",
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
	if v := GetValueFromRow(rows, ConfigKeyMysql, "", "max_open_conns", "", debug); v != "" {
		cfg.MaxOpenConns = typeHelper.Str2Int(v)
	}
	if v := GetValueFromRow(rows, ConfigKeyMysql, "", "max_idle", "", debug); v != "" {
		cfg.MaxIdle = typeHelper.Str2Int(v)
	}
	if v := GetValueFromRow(rows, ConfigKeyMysql, "", "max_life_time", "", debug); v != "" {
		cfg.MaxLifeTime = typeHelper.Str2Int(v)
	}
	return cfg, nil
}

func getAliOssConfig(rows []*DbSettingRow, debug bool) (*configStruct.AliOssConfig, error) {
	cfg := &configStruct.AliOssConfig{
		AccessKeyID:     GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliOssFlag, "access_key_id", "", debug),
		AccessKeySecret: GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliOssFlag, "access_key_secret", "", debug),
		Host:            GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliOssFlag, "host", "", debug),
		EndPoint:        GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliOssFlag, "end_point", "", debug),
		BucketName:      GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliOssFlag, "bucket_name", "", debug),
		ExpireTime:      typeHelper.Str2Int64(GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliOssFlag, "expire_time", "30", debug)),
		ARN:             GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliOssFlag, "arn", "", debug),
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
		AppKey:    GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliApiFlag, "app_key", "", debug),
		AppSecret: GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliApiFlag, "app_secret", "", debug),
		AppCode:   GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliApiFlag, "app_code", "", debug),
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
		AccessKeySecret: GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliSmsFlag, "access_key_secret", "", debug),
		AccessKeyID:     GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliSmsFlag, "access_key_id", "", debug),
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
		AccessKeySecret: GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliIotFlag, "access_key_secret", "", debug),
		AccessKeyID:     GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliIotFlag, "access_key_id", "", debug),
		EndPoint:        GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliIotFlag, "end_point", "", debug),
		RegionID:        GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliIotFlag, "region_id", "", debug),
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
	cfg := &configStruct.AmapConfig{Key: GetValueFromRow(rows, ConfigKeyAli, ConfigKeyAliAmapFlag, "key", "", debug)}
	if cfg.Key == "" {
		return nil, fmt.Errorf("高德地图配置 Key 为空")
	}
	return cfg, nil
}

func getYunXinConfig(rows []*DbSettingRow, debug bool) (*configStruct.YunxinConfig, error) {
	cfg := &configStruct.YunxinConfig{
		AppKey:    GetValueFromRow(rows, ConfigKeyNetease, ConfigKeyYunxinFlag, "app_key", "", debug),
		AppSecret: GetValueFromRow(rows, ConfigKeyNetease, ConfigKeyYunxinFlag, "app_secret", "", debug),
	}
	if cfg.AppKey == "" {
		return nil, fmt.Errorf("网易云信配置 AppKey 为空")
	}
	if cfg.AppSecret == "" {
		return nil, fmt.Errorf("网易云信配置 AppSecret 为空")
	}
	return cfg, nil
}
