# 配置项书写规范

本文档定义了 `appMng` 配置管理系统中配置项的书写规范，包括数据库配置和 YAML 配置的格式要求。

## 目录

- [配置键（ConfigKey）规范](#配置键configkey规范)
- [数据库配置规范](#数据库配置规范)
- [YAML 配置规范](#yaml配置规范)
- [结构体定义规范](#结构体定义规范)
- [示例](#示例)

---

## 配置键（ConfigKey）规范

### 格式要求

- **格式**：首字母大写的驼峰命名法（PascalCase）
- **用途**：作为 `ConfigKey.Key` 的值，直接对应结构体字段名
- **示例**：
  - `Redis` - Redis 配置
  - `Mysql` - MySQL 配置
  - `Postgres` - PostgreSQL 配置
  - `AliApi` - 阿里云 API 配置

### 说明

- `ConfigKey.Key` 的值必须与目标结构体的字段名完全一致（首字母大写）
- 例如：`ConfigKeys.Redis.Key = "Redis"` 对应 `BaseConfig.Redis` 字段

---

## 数据库配置规范

### 表结构

配置存储在 `a_setting` 表中，主要字段：

- `name`: 配置名称（用于标识配置组）
- `flag_1`: 父级配置键（首字母大写）
- `flag_2`: 字段标识（小写下划线）
- `value_1`: 生产环境值
- `value_2`: 调试环境值（debug 模式时使用）

### 字段规范

#### 1. `name` 字段

- **格式**：首字母大写的驼峰命名法（PascalCase）
- **用途**：标识配置组，用于区分不同的配置项
- **示例**：
  - `AliApi` - 阿里云 API 配置组
  - `VolcengineConfig` - 火山引擎配置组
  - `Redis` - Redis 配置组

#### 2. `flag_1` 字段

- **格式**：首字母大写的驼峰命名法（PascalCase）
- **用途**：对应 `ConfigKey.Key` 的值，用于标识父级配置
- **规则**：
  - 必须与 `ConfigKey.Key` 完全一致
  - 对于 `BaseConfig` 中的顶级配置，`flag_1` 等于 `name`
  - 对于嵌套配置（如 `ProjectConfig`），`flag_1` 等于配置名称
- **示例**：
  - `Redis` - 对应 `BaseConfig.Redis`
  - `AliApi` - 对应 `MyProjectConfigData.AliApi`

#### 3. `flag_2` 字段

- **格式**：小写下划线命名法（snake_case）
- **用途**：对应结构体字段的 `json` 或 `mapstructure` tag 值
- **规则**：
  - 必须与结构体字段的 `json` 或 `mapstructure` tag 完全一致
  - 如果字段没有 `json` tag，则使用 `mapstructure` tag
  - 如果两者都没有，该字段不会被加载
  - 格式必须是小写下划线（如 `app_key`，不是 `AppKey`）
- **示例**：
  - `app_key` - 对应 `AliApiConfig.AppKey` 字段的 `mapstructure:"app_key"`
  - `app_secret` - 对应 `AliApiConfig.AppSecret` 字段的 `mapstructure:"app_secret"`
  - `host` - 对应 `HttpServerConfig.Host` 字段的 `mapstructure:"host"`

### 数据库配置示例

```sql
-- 示例 1: BaseConfig 中的 Redis 配置
INSERT INTO "a_setting" ("name", "flag_1", "flag_2", "value_1", "value_2") VALUES
('Redis', 'Redis', 'host', '127.0.0.1', '127.0.0.1'),
('Redis', 'Redis', 'port', '6379', '6379'),
('Redis', 'Redis', 'password', 'your_password', 'dev_password');

-- 示例 2: ProjectConfig 中的 AliApi 配置
INSERT INTO "a_setting" ("name", "flag_1", "flag_2", "value_1", "value_2") VALUES
('AliApi', 'AliApi', 'app_key', '111', '111'),
('AliApi', 'AliApi', 'app_secret', '222', '222'),
('AliApi', 'AliApi', 'app_code', '333', '333');

-- 示例 3: HttpServerConfig 配置
INSERT INTO "a_setting" ("name", "flag_1", "flag_2", "value_1", "value_2") VALUES
('HttpServer', 'client', 'label', 'client', 'client'),
('HttpServer', 'client', 'host', '0.0.0.0', '0.0.0.0'),
('HttpServer', 'client', 'port', '8021', '8021'),
('HttpServer', 'console', 'label', 'console', 'console'),
('HttpServer', 'console', 'host', '0.0.0.0', '0.0.0.0'),
('HttpServer', 'console', 'port', '8022', '8022');
```

---

## YAML 配置规范

### 格式要求

- **格式**：小写下划线命名法（snake_case）
- **用途**：YAML 配置文件中的键名
- **规则**：
  - 顶级配置键使用小写下划线（如 `redis`, `mysql`）
  - 嵌套配置键也使用小写下划线（如 `app_code`, `app_key`）
  - 必须与结构体字段的 `mapstructure` tag 完全一致

### YAML 配置示例

```yaml
# config.yaml

# Profile 配置
Profile:
  no: "001"
  name: "homepage"
  version: "1.0.0"
  debug: true

# Redis 配置
Redis:
  host: "127.0.0.1"
  port: "6379"
  password: "your_password"

# MySQL 配置
Mysql:
  dsn: "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  auto_migrate: "true"

# HTTP Server 配置
HttpServer:
  client:
    label: "client"
    host: "0.0.0.0"
    port: "8021"
    domain: "https://api.example.com"
  console:
    label: "console"
    host: "0.0.0.0"
    port: "8022"
    domain: "https://console.example.com"

# 项目特定配置（ProjectConfig）
AliApi:
  app_key: "203728448"
  app_secret: "123"
  app_code: "456"
```

**注意**：YAML 中的顶级键（如 `Redis`, `Mysql`）使用首字母大写，这是为了与 `BaseConfig` 的字段名保持一致。但嵌套字段（如 `app_key`）使用小写下划线。

---

## 结构体定义规范

### 字段命名

- **格式**：首字母大写的驼峰命名法（PascalCase）
- **示例**：
  - `AppKey` - 应用密钥字段
  - `AppSecret` - 应用密钥字段
  - `HttpServerConfig` - HTTP 服务器配置字段

### Tag 规范

#### 1. `mapstructure` tag

- **格式**：小写下划线命名法（snake_case）
- **用途**：
  - 用于 YAML 配置映射
  - 用于数据库 `flag_2` 字段匹配
- **规则**：
  - 必须使用小写下划线格式
  - 必须与数据库 `flag_2` 字段完全一致
  - 必须与 YAML 配置中的键名完全一致（嵌套字段）
- **示例**：
  ```go
  type AliApiConfig struct {
      AppKey    string `mapstructure:"app_key" validate:"required"`
      AppSecret string `mapstructure:"app_secret" validate:"required"`
      AppCode   string `mapstructure:"app_code" validate:"required"`
  }
  ```

#### 2. `json` tag（可选）

- **格式**：小写下划线命名法（snake_case）
- **用途**：如果存在，优先使用 `json` tag 作为 `flag_2` 的匹配值
- **规则**：如果 `json` tag 为空，则使用 `mapstructure` tag

#### 3. `validate` tag（可选）

- **用途**：配置验证规则
- **示例**：`validate:"required"` - 必填字段

### 结构体定义示例

```go
// BaseConfig 中的顶级配置
type BaseConfig struct {
    Redis    *RedisConfig    `mapstructure:"Redis"`    // 顶级键使用首字母大写
    Mysql    *MysqlConfig    `mapstructure:"Mysql"`
    Postgres *PostgresConfig `mapstructure:"Postgres"`
}

// 嵌套配置结构体
type RedisConfig struct {
    Host     string `mapstructure:"host"`     // 嵌套字段使用小写下划线
    Port     string `mapstructure:"port"`
    Password string `mapstructure:"password"`
}

// ProjectConfig 中的配置
type AliApiConfig struct {
    AppKey    string `mapstructure:"app_key" validate:"required"`
    AppSecret string `mapstructure:"app_secret" validate:"required"`
    AppCode   string `mapstructure:"app_code" validate:"required"`
}
```

---

## 示例

### 完整示例：AliApi 配置

#### 1. 结构体定义

```go
// 在 configStruct/configStruct.go 中
type AliApiConfig struct {
    AppKey    string `mapstructure:"app_key" validate:"required"`
    AppSecret string `mapstructure:"app_secret" validate:"required"`
    AppCode   string `mapstructure:"app_code" validate:"required"`
}

// 在 ProjectConfig 中使用
type MyProjectConfigData struct {
    AliApi *AliApiConfig
}
```

#### 2. ConfigKey 定义

```go
// 在 appMng/structs.go 中
var ConfigKeys = struct {
    // ... 其他配置键
    AliApi ConfigKey
}{
    // ... 其他配置键
    AliApi: NewConfigKey("AliApi", "阿里云API"),
}
```

#### 3. 数据库配置

```sql
INSERT INTO "a_setting" ("name", "flag_1", "flag_2", "value_1", "value_2") VALUES
('AliApi', 'AliApi', 'app_key', '203728448', '203728448'),
('AliApi', 'AliApi', 'app_secret', '123', '123'),
('AliApi', 'AliApi', 'app_code', '456', '456');
```

#### 4. YAML 配置（可选）

```yaml
AliApi:
  app_key: "203728448"
  app_secret: "123"
  app_code: "456"
```

#### 5. ConfigSourceStrategy 配置

```go
// 在 appMng/structs.go 中
type ConfigSourceStrategy struct {
    // ... 其他配置
    AliApi ConfigSource `mapstructure:"AliApi"` // 使用首字母大写
}

// 设置配置来源
strategy := &ConfigSourceStrategy{
    AliApi: SourceDatabase, // 从数据库加载
}
```

---

## 总结

### 命名规范对照表

| 位置 | 格式 | 示例 |
|------|------|------|
| `ConfigKey.Key` | 首字母大写（PascalCase） | `Redis`, `AliApi` |
| 结构体字段名 | 首字母大写（PascalCase） | `AppKey`, `AppSecret` |
| `mapstructure` tag | 小写下划线（snake_case） | `app_key`, `app_secret` |
| 数据库 `name` | 首字母大写（PascalCase） | `AliApi`, `Redis` |
| 数据库 `flag_1` | 首字母大写（PascalCase） | `AliApi`, `Redis` |
| 数据库 `flag_2` | 小写下划线（snake_case） | `app_key`, `app_secret` |
| YAML 顶级键 | 首字母大写（PascalCase） | `Redis`, `Mysql` |
| YAML 嵌套键 | 小写下划线（snake_case） | `app_key`, `host` |

### 关键规则

1. **ConfigKey.Key** = **数据库 flag_1** = **结构体字段名**（首字母大写）
2. **mapstructure tag** = **数据库 flag_2** = **YAML 嵌套键**（小写下划线）
3. **数据库 name** = **配置组名称**（首字母大写）
4. **YAML 顶级键** = **结构体字段名**（首字母大写，与 BaseConfig 字段一致）

---

## 常见问题

### Q1: 为什么 `flag_2` 使用小写下划线而不是首字母大写？

**A**: 这是为了符合 Go 社区的标准惯例。`mapstructure` tag 通常使用小写下划线格式，这样可以：
- 与 YAML 配置保持一致
- 符合 Go 的命名惯例
- 与 JSON tag 保持一致

### Q2: 如果结构体字段名和 `mapstructure` tag 不一致会怎样？

**A**: 系统会使用 `mapstructure` tag 作为 `flag_2` 的匹配值。如果 `mapstructure` tag 为空，则使用 `json` tag。如果两者都为空，该字段不会被加载。

### Q3: YAML 中的顶级键为什么使用首字母大写？

**A**: 因为 `BaseConfig` 的字段名是首字母大写的（如 `Redis`, `Mysql`），为了保持一致性，YAML 中的顶级键也使用首字母大写。但嵌套字段使用小写下划线。

### Q4: 如何添加新的配置项？

**A**: 按照以下步骤：
1. 定义结构体（字段名首字母大写，`mapstructure` tag 小写下划线）
2. 在 `ConfigKeys` 中注册配置键（首字母大写）
3. 在 `ConfigSourceStrategy` 中设置配置来源
4. 在数据库中插入配置（`name` 和 `flag_1` 首字母大写，`flag_2` 小写下划线）
5. （可选）在 YAML 中配置（顶级键首字母大写，嵌套键小写下划线）

---

**最后更新**: 2025-12-30

