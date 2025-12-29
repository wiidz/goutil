package appMng

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/wiidz/goutil/helpers/configHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// configErrorFactory 统一生成配置相关错误，减少重复文案
type configErrorFactory struct{}

func newConfigErrorFactory() *configErrorFactory {
	return &configErrorFactory{}
}

// missingField 缺失必填字段
func (f *configErrorFactory) missingField(key, field string) error {
	return fmt.Errorf("%s 配置 %s 为空", GetKeyDisplayName(key), field)
}

// databaseEmpty 从数据库加载到空结果
func (f *configErrorFactory) databaseEmpty(key string) error {
	return fmt.Errorf("从数据库加载 %s 配置失败: 数据为空", GetKeyDisplayName(key))
}

// yamlNotInit 未初始化 YAML
func (f *configErrorFactory) yamlNotInit(key string) error {
	return fmt.Errorf("从 YAML 加载 %s 配置失败: 未初始化 YAML 配置", GetKeyDisplayName(key))
}

// yamlLoadFailed YAML 解析失败
func (f *configErrorFactory) yamlLoadFailed(key string, err error) error {
	return fmt.Errorf("从 YAML 加载 %s 配置失败: %w", GetKeyDisplayName(key), err)
}

var errFactory = newConfigErrorFactory()

// key 对应 BaseConfig 赋值函数
var configAssigners = map[string]func(*configStruct.BaseConfig, interface{}){
	ConfigKeys.Redis.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.RedisConfig); ok {
			cfg.RedisConfig = val
		}
	},
	ConfigKeys.Es.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.EsConfig); ok {
			cfg.EsConfig = val
		}
	},
	ConfigKeys.RabbitMQ.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.RabbitMQConfig); ok {
			cfg.RabbitMQConfig = val
		}
	},
	ConfigKeys.Postgres.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.PostgresConfig); ok {
			cfg.PostgresConfig = val
		}
	},
	ConfigKeys.Mysql.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.MysqlConfig); ok {
			cfg.MysqlConfig = val
		}
	},
	ConfigKeys.WechatMini.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.WechatMiniConfig); ok {
			cfg.WechatMiniConfig = val
		}
	},
	ConfigKeys.WechatOa.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.WechatOaConfig); ok {
			cfg.WechatOaConfig = val
		}
	},
	ConfigKeys.WechatOpen.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.WechatOpenConfig); ok {
			cfg.WechatOpenConfig = val
		}
	},
	ConfigKeys.WechatPayV3.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.WechatPayConfigV3); ok {
			cfg.WechatPayConfigV3 = val
		}
	},
	ConfigKeys.WechatPayV2.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.WechatPayConfigV2); ok {
			cfg.WechatPayConfigV2 = val
		}
	},
	ConfigKeys.AliOss.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.AliOssConfig); ok {
			cfg.AliOssConfig = val
		}
	},
	ConfigKeys.AliPay.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.AliPayConfig); ok {
			cfg.AliPayConfig = val
		}
	},
	ConfigKeys.AliApi.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.AliApiConfig); ok {
			cfg.AliApiConfig = val
		}
	},
	ConfigKeys.AliSms.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.AliSmsConfig); ok {
			cfg.AliSmsConfig = val
		}
	},
	ConfigKeys.AliIot.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.AliIotConfig); ok {
			cfg.AliIotConfig = val
		}
	},
	ConfigKeys.Amap.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.AmapConfig); ok {
			cfg.AmapConfig = val
		}
	},
	ConfigKeys.Yunxin.Key: func(cfg *configStruct.BaseConfig, v interface{}) {
		if val, ok := v.(*configStruct.YunxinConfig); ok {
			cfg.YunxinConfig = val
		}
	},
}

// key 对应的来源选择器
var configSources = map[string]func(*ConfigSourceStrategy) ConfigSource{
	ConfigKeys.Redis.Key:      func(s *ConfigSourceStrategy) ConfigSource { return s.Redis },
	ConfigKeys.Es.Key:         func(s *ConfigSourceStrategy) ConfigSource { return s.Es },
	ConfigKeys.RabbitMQ.Key:   func(s *ConfigSourceStrategy) ConfigSource { return s.RabbitMQ },
	ConfigKeys.Postgres.Key:   func(s *ConfigSourceStrategy) ConfigSource { return s.Postgres },
	ConfigKeys.Mysql.Key:      func(s *ConfigSourceStrategy) ConfigSource { return s.Mysql },
	ConfigKeys.WechatMini.Key: func(s *ConfigSourceStrategy) ConfigSource { return s.WechatMini },
	ConfigKeys.WechatOa.Key:   func(s *ConfigSourceStrategy) ConfigSource { return s.WechatOa },
	ConfigKeys.WechatOpen.Key: func(s *ConfigSourceStrategy) ConfigSource { return s.WechatOpen },
	ConfigKeys.WechatPayV3.Key: func(s *ConfigSourceStrategy) ConfigSource {
		return s.WechatPayV3
	},
	ConfigKeys.WechatPayV2.Key: func(s *ConfigSourceStrategy) ConfigSource {
		return s.WechatPayV2
	},
	ConfigKeys.AliOss.Key: func(s *ConfigSourceStrategy) ConfigSource { return s.AliOss },
	ConfigKeys.AliPay.Key: func(s *ConfigSourceStrategy) ConfigSource { return s.AliPay },
	ConfigKeys.AliApi.Key: func(s *ConfigSourceStrategy) ConfigSource { return s.AliApi },
	ConfigKeys.AliSms.Key: func(s *ConfigSourceStrategy) ConfigSource { return s.AliSms },
	ConfigKeys.AliIot.Key: func(s *ConfigSourceStrategy) ConfigSource { return s.AliIot },
	ConfigKeys.Amap.Key:   func(s *ConfigSourceStrategy) ConfigSource { return s.Amap },
	ConfigKeys.Yunxin.Key: func(s *ConfigSourceStrategy) ConfigSource { return s.Yunxin },
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
		return nil, fmt.Errorf("策略不能为 nil")
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

// 统一加载所有配置项
func (b *ConfigBuilder) loadAllConfigs(cfg *configStruct.BaseConfig, dbRows []*DbSettingRow, debug bool) error {
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
		case SourceYAML:
			if len(b.yamlVipers) == 0 {
				return errFactory.yamlNotInit(key)
			}
			if err := b.yamlVipers[0].UnmarshalKey(key, targetPtr.Interface()); err != nil {
				return errFactory.yamlLoadFailed(key, err)
			}
			applyDefaultsFromTags(targetPtr.Interface())
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

// Build 构建 BaseConfig，根据策略从不同来源加载配置
func (b *ConfigBuilder) Build(ctx context.Context) (*configStruct.BaseConfig, error) {

	cfg := &configStruct.BaseConfig{}

	// 第一步：检查策略中是否有数据库相关的配置，如果有，优先初始化数据库
	// 这样后续的配置才能从数据库中读取
	// 注意：如果策略要求从数据库加载配置，b.db 在 Build 开始时必然是 nil，需要先初始化
	if b.needDatabaseConnection() {
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
		dbRows, err = b.loadAllSettingRows(ctx)
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

	// 初始化HttpServer配置
	for _, serverLabel := range b.initialConfig.HttpServerLabels {
		serverConfig := getServerConfig(dbRows, serverLabel)
		if serverConfig == nil {
			return nil, fmt.Errorf("加载 HttpServer 配置失败: 标签 %s 不存在", serverLabel)
		}
		cfg.HttpServerConfig[serverLabel] = serverConfig
	}

	// 第四步：统一管线加载已注册的配置
	if err := b.loadAllConfigs(cfg, dbRows, debug); err != nil {
		return nil, err
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
			if err := b.yamlVipers[0].UnmarshalKey(ConfigKeys.Postgres.Key, &postgresConfig); err == nil && postgresConfig.DSN != "" {
				return b.initPostgresFromConfig(&postgresConfig)
			}
		}
		// 如果策略要求从数据库加载 MySQL 配置，尝试从 YAML 加载 MySQL 配置来初始化连接
		if b.strategy.Mysql == SourceDatabase {
			var mysqlConfig configStruct.MysqlConfig
			if err := b.yamlVipers[0].UnmarshalKey(ConfigKeys.Mysql.Key, &mysqlConfig); err == nil {
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
			return errFactory.missingField(ConfigKeys.Postgres.Key, "DSN")
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
		return errFactory.missingField(ConfigKeys.Postgres.Key, "DSN")
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
		return fmt.Errorf("%s 配置为空", GetKeyDisplayName(ConfigKeys.Mysql.Key))
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

// loadAllSettingRows 从数据库加载配置行
func (b *ConfigBuilder) loadAllSettingRows(ctx context.Context) ([]*DbSettingRow, error) {
	if b.db == nil {
		return nil, fmt.Errorf("数据库连接未设置")
	}

	var rows []*DbSettingRow
	err := b.db.WithContext(ctx).
		Table(b.getSettingTableName()).
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
		if err := b.yamlVipers[0].UnmarshalKey(ConfigKeys.Profile.Key, &profile); err == nil {
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
		timeZone := b.yamlVipers[0].GetString(ConfigKeys.LocationTZ.Key)
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
		No:   GetValueFromRow(rows, ConfigKeys.App.Key, "", "no", "", false),
		Name: GetValueFromRow(rows, ConfigKeys.App.Key, "", "name", "", false),
		// Host:    GetValueFromRow(rows, ConfigKeys.App.Key, "", "host", "", false),
		// Port:    GetValueFromRow(rows, ConfigKeys.App.Key, "", "port", "127.0.0.1", false),
		// Domain:  GetValueFromRow(rows, ConfigKeys.App.Key, "", "domain", "", false),
		Debug:   GetValueFromRow(rows, ConfigKeys.App.Key, "", "debug", "", false) == "1",
		Version: GetValueFromRow(rows, ConfigKeys.App.Key, "", "version", "", false),
	}
}

func getServerConfig(rows []*DbSettingRow, serverLabel string) *configStruct.HttpServerConfig {
	return &configStruct.HttpServerConfig{
		Label:  GetValueFromRow(rows, ConfigKeys.Server.Key, serverLabel, "label", "", false),
		Host:   GetValueFromRow(rows, ConfigKeys.Server.Key, serverLabel, "host", "", false),
		Port:   GetValueFromRow(rows, ConfigKeys.Server.Key, serverLabel, "port", "", false),
		Domain: GetValueFromRow(rows, ConfigKeys.Server.Key, serverLabel, "domain", "", false),
	}
}
func getLocationConfig(rows []*DbSettingRow) (location *time.Location, err error) {
	timeZone := GetValueFromRow(rows, ConfigKeys.TimeZone.Key, "", "", "Asia/Shanghai", false)
	location, err = time.LoadLocation(timeZone)
	if err != nil {
		location = time.FixedZone("CST-8", 8*3600)
	}
	return
}

// ConfigMap 以键对应结构体定义（默认值从 default tag 读取）
type ConfigMap struct {
	Key  ConfigKey
	Data interface{}
}

// 所有配置的结构映射
var configMaps = map[string]ConfigMap{
	ConfigKeys.Redis.Key:       {Key: ConfigKeys.Redis, Data: configStruct.RedisConfig{}},
	ConfigKeys.Es.Key:          {Key: ConfigKeys.Es, Data: configStruct.EsConfig{}},
	ConfigKeys.RabbitMQ.Key:    {Key: ConfigKeys.RabbitMQ, Data: configStruct.RabbitMQConfig{}},
	ConfigKeys.Postgres.Key:    {Key: ConfigKeys.Postgres, Data: configStruct.PostgresConfig{}},
	ConfigKeys.Mysql.Key:       {Key: ConfigKeys.Mysql, Data: configStruct.MysqlConfig{}},
	ConfigKeys.WechatMini.Key:  {Key: ConfigKeys.WechatMini, Data: configStruct.WechatMiniConfig{}},
	ConfigKeys.WechatOa.Key:    {Key: ConfigKeys.WechatOa, Data: configStruct.WechatOaConfig{}},
	ConfigKeys.WechatOpen.Key:  {Key: ConfigKeys.WechatOpen, Data: configStruct.WechatOpenConfig{}},
	ConfigKeys.WechatPayV3.Key: {Key: ConfigKeys.WechatPayV3, Data: configStruct.WechatPayConfigV3{}},
	ConfigKeys.WechatPayV2.Key: {Key: ConfigKeys.WechatPayV2, Data: configStruct.WechatPayConfigV2{}},
	ConfigKeys.AliOss.Key:      {Key: ConfigKeys.AliOss, Data: configStruct.AliOssConfig{}},
	ConfigKeys.AliPay.Key:      {Key: ConfigKeys.AliPay, Data: configStruct.AliPayConfig{}},
	ConfigKeys.AliApi.Key:      {Key: ConfigKeys.AliApi, Data: configStruct.AliApiConfig{}},
	ConfigKeys.AliSms.Key:      {Key: ConfigKeys.AliSms, Data: configStruct.AliSmsConfig{}},
	ConfigKeys.AliIot.Key:      {Key: ConfigKeys.AliIot, Data: configStruct.AliIotConfig{}},
	ConfigKeys.Amap.Key:        {Key: ConfigKeys.Amap, Data: configStruct.AmapConfig{}},
	ConfigKeys.Yunxin.Key:      {Key: ConfigKeys.Yunxin, Data: configStruct.YunxinConfig{}},
}

// applyDefaultsFromTags 根据 struct 的 default tag 填充零值字段
func applyDefaultsFromTags(target interface{}) {
	if target == nil {
		return
	}
	val := reflect.ValueOf(target)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		ft := typ.Field(i)
		if !field.CanSet() || !field.IsZero() {
			continue
		}
		def := ft.Tag.Get("default")
		if def == "" {
			continue
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(def)
		case reflect.Bool:
			if v, err := strconv.ParseBool(def); err == nil {
				field.SetBool(v)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v, err := strconv.ParseInt(def, 10, 64); err == nil {
				field.SetInt(v)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if v, err := strconv.ParseUint(def, 10, 64); err == nil {
				field.SetUint(v)
			}
		}
	}
}

// fillConfigFromRows 使用 flag1=parentKey、flag2=字段 json/mapstructure 标签，从 rows 填充配置，并按 default tag 设置默认值
func fillConfigFromRows(target interface{}, nameKey, flag1 string, rows []*DbSettingRow, debug bool) error {
	if target == nil {
		return fmt.Errorf("target 不能为空")
	}
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("target 必须是非 nil 指针")
	}
	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("target 必须指向结构体")
	}

	typ := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		ft := typ.Field(i)
		if !field.CanSet() {
			continue
		}

		jsonTag := ft.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = ft.Tag.Get("mapstructure")
		}
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		jsonTag = strings.Split(jsonTag, ",")[0]

		defVal := ft.Tag.Get("default")
		raw := GetValueFromRow(rows, nameKey, flag1, jsonTag, defVal, debug)
		if raw == "" && defVal != "" {
			raw = defVal
		}
		if raw == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(raw)
		case reflect.Bool:
			if v, err := strconv.ParseBool(raw); err == nil {
				field.SetBool(v)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v, err := strconv.ParseInt(raw, 10, 64); err == nil {
				field.SetInt(v)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if v, err := strconv.ParseUint(raw, 10, 64); err == nil {
				field.SetUint(v)
			}
		case reflect.Float32, reflect.Float64:
			if v, err := strconv.ParseFloat(raw, 64); err == nil {
				field.SetFloat(v)
			}
		}
	}

	applyDefaultsFromTags(target)
	return nil
}
