package identityMng

import (
	"sync"
	"time"

	"github.com/click33/sa-token-go/core"
	"github.com/click33/sa-token-go/core/config"
	"github.com/go-redis/redis/v9"
)

type tokenStyle = struct {
	UUID      config.TokenStyle
	Simple    config.TokenStyle
	Random32  config.TokenStyle
	Random64  config.TokenStyle
	Random128 config.TokenStyle
	JWT       config.TokenStyle
	Hash      config.TokenStyle
	Timestamp config.TokenStyle
	Tik       config.TokenStyle
}

var TokenStyle = tokenStyle{
	UUID:      config.TokenStyleUUID,
	Simple:    config.TokenStyleSimple,
	Random32:  config.TokenStyleRandom32,
	Random64:  config.TokenStyleRandom64,
	Random128: config.TokenStyleRandom128,
	JWT:       config.TokenStyleJWT,
	Hash:      config.TokenStyleHash,
	Timestamp: config.TokenStyleTimestamp,
	Tik:       config.TokenStyleTik,
}

type IdentityMng struct {
	config        *config.Config
	defaultDevice string
	once          sync.Once
	Storage       core.Storage // 存储实现（默认 memory）
}

type StorageType string

const (
	Redis  StorageType = "redis"
	Memory StorageType = "memory"
)

// Config
type Config struct {
	TokenStyle    config.TokenStyle // token风格
	Salt          string            // 盐值（当 TokenStyle=JWT 时必填）
	Timeout       time.Duration     // token有效期（单位：Duration，将转换为秒）
	DefaultDevice string            // 设备默认值（device 为空时使用）

	StorageType StorageType   // 存储类型（memory/redis），为空则走默认
	RedisClient *redis.Client // 仅当StorageType=redis时需要

	SaConfig *config.Config // 直接传入底层配置（可选）
}

// TokenPair 标准令牌对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Claims 标准声明（按需扩展）
type Claims struct {
	LoginID  string
	Device   string
	ExpireAt int64
}
