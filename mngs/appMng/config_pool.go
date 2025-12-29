package appMng

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/wiidz/goutil/helpers/configHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConfigPool 配置池，管理数据库连接和 YAML 配置
type ConfigPool struct {
	yamlVipers []*viper.Viper // YAML 配置实例列表

	dbType           configStruct.DBType // 数据库类型（postgres 或 mysql）
	settingTableName string              // 配置表名
	db               *gorm.DB            // 数据库连接（可能是 MySQL 或 PostgreSQL）
	dbRows           []*DbSettingRow     // 数据库配置行（从数据库加载后缓存）
}

// NewConfigPool 创建配置池
// yamlFiles: YAML 配置文件列表（必须，数据库配置也从 YAML 中读取）
// settingTableName: 配置表名（可选，默认 "a_setting"）
func NewConfigPool(ctx context.Context, yamlFiles []*configStruct.ViperConfig, settingTableName string) (*ConfigPool, error) {
	pool := &ConfigPool{
		settingTableName: settingTableName,
		yamlVipers:       make([]*viper.Viper, 0),
	}

	// 第一步：初始化 YAML 配置（必须先初始化 YAML，才能从中读取数据库配置）
	if len(yamlFiles) == 0 {
		return nil, fmt.Errorf("YAML 配置文件列表不能为空")
	}
	if err := pool.initYAML(yamlFiles); err != nil {
		return nil, fmt.Errorf("初始化 YAML 配置失败: %w", err)
	}

	// 第二步：如果传入了 settingTableName，说明需要从数据库加载配置，需要初始化数据库连接
	if settingTableName != "" {
		// 从 YAML 初始化数据库连接
		if err := pool.InitDatabaseFromYAML(); err != nil {
			return nil, fmt.Errorf("初始化数据库连接失败（需要从数据库加载配置，但无法连接数据库）: %w", err)
		}

		// 第三步：数据库初始化后，立即加载配置行并保存
		if pool.db != nil {
			dbRows, err := pool.LoadSettingRows(ctx)
			if err != nil {
				log.Printf("警告: 无法从数据库加载配置行: %v", err)
			} else {
				pool.dbRows = dbRows
				log.Printf("成功: 从数据库加载了 %d 条配置行并缓存到配置池", len(dbRows))
			}
		}
	}

	return pool, nil
}

// initPostgresFromDSN 从 DSN 初始化 PostgreSQL 连接
func (p *ConfigPool) initPostgresFromDSN(dsn string, maxIdle, maxOpen int, maxLifetime time.Duration, loggerInterface logger.Interface) error {
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

	p.db = db
	log.Printf("成功: PostgreSQL 数据库连接已初始化")
	return nil
}

// initMysqlFromDSN 从 DSN 初始化 MySQL 连接
func (p *ConfigPool) initMysqlFromDSN(dsn string, maxIdle, maxOpen int, maxLifetime time.Duration, loggerInterface logger.Interface) error {
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

	p.db = db
	log.Printf("成功: MySQL 数据库连接已初始化")
	return nil
}

// InitDatabaseFromYAML 从 YAML 配置初始化数据库连接
func (p *ConfigPool) InitDatabaseFromYAML() error {
	if len(p.yamlVipers) == 0 {
		return fmt.Errorf("YAML 配置未初始化")
	}

	// 优先尝试从 YAML 加载 PostgreSQL 配置来初始化连接
	var postgresConfig configStruct.PostgresConfig
	if err := p.yamlVipers[0].UnmarshalKey(ConfigKeys.Postgres.Key, &postgresConfig); err == nil && postgresConfig.DSN != "" {
		p.dbType = configStruct.DBTypePostgres
		return p.initPostgresFromDSN(postgresConfig.DSN, postgresConfig.ConnMaxIdle, postgresConfig.ConnMaxOpen, postgresConfig.ConnMaxLifetime, postgresConfig.Logger)
	}

	// 如果 PostgreSQL 配置不存在，尝试从 YAML 加载 MySQL 配置来初始化连接
	var mysqlConfig configStruct.MysqlConfig
	if err := p.yamlVipers[0].UnmarshalKey(ConfigKeys.Mysql.Key, &mysqlConfig); err == nil {
		if mysqlConfig.Host != "" && mysqlConfig.DbName != "" {
			p.dbType = configStruct.DBTypeMysql
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

			return p.initMysqlFromDSN(dsn, maxIdle, maxOpen, maxLifetime, mysqlConfig.Logger)
		}
	}

	return fmt.Errorf("无法从 YAML 加载数据库配置，不存在" + ConfigKeys.Postgres.Key + "或" + ConfigKeys.Mysql.Key)
}

// initYAML 初始化 YAML 配置
func (p *ConfigPool) initYAML(yamlFiles []*configStruct.ViperConfig) error {
	for _, yamlConfig := range yamlFiles {
		v, err := configHelper.GetViper(yamlConfig)
		if err != nil {
			log.Printf("警告: 无法加载 YAML 配置文件 %s/%s.%s: %v", yamlConfig.DirPath, yamlConfig.FileName, yamlConfig.FileType, err)
			continue
		}
		p.yamlVipers = append(p.yamlVipers, v)
	}
	return nil
}

// GetDB 获取数据库连接
func (p *ConfigPool) GetDB() *gorm.DB {
	return p.db
}

// GetYAML 获取 YAML 配置实例列表
func (p *ConfigPool) GetYAML() []*viper.Viper {
	return p.yamlVipers
}

// GetSettingTableName 获取配置表名
func (p *ConfigPool) GetSettingTableName() string {
	if p.settingTableName != "" {
		return p.settingTableName
	}
	return "a_setting" // 默认值
}

// LoadSettingRows 从数据库加载配置行
func (p *ConfigPool) LoadSettingRows(ctx context.Context) ([]*DbSettingRow, error) {
	if p.db == nil {
		return nil, fmt.Errorf("数据库连接未设置")
	}

	var rows []*DbSettingRow
	err := p.db.WithContext(ctx).
		Table(p.GetSettingTableName()).
		Where("kind = ? AND deleted_at IS NULL", 1).
		Find(&rows).Error

	return rows, err
}

// GetDBRows 获取缓存的数据库配置行
func (p *ConfigPool) GetDBRows() []*DbSettingRow {
	return p.dbRows
}

// GetDBType 获取数据库类型（postgres 或 mysql）
func (p *ConfigPool) GetDBType() configStruct.DBType {
	return p.dbType
}
