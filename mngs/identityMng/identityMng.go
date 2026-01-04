package identityMng

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/click33/sa-token-go/core/config"
	sagin "github.com/click33/sa-token-go/integrations/gin"
	"github.com/click33/sa-token-go/storage/memory"
	"github.com/click33/sa-token-go/stputil"
)

var setOnce sync.Once

// NewMng 创建实例（可多实例，首次创建时初始化 Sa-Token 全局管理器）
func NewMng(config *Config) (mng *IdentityMng, err error) {

	if config == nil {
		err = errors.New("config is required")
		return
	}

	saCfg, err := generateSaConfig(config)
	if err != nil {
		return
	}

	mng = &IdentityMng{config33: saCfg, defaultDevice: config.DefaultDevice, debug: config.Debug}

	if config.StorageType == "" {
		mng.Storage = memory.NewStorage()
	} else {
		switch config.StorageType {
		case Memory, "":
			mng.Storage = memory.NewStorage()
		case Redis:
			if config.RedisClient == nil {
				err = errors.New("redis client is required")
				return
			}
			mng.Storage = NewRedisStorage(config.RedisClient, config.Debug)
		}
	}
	mng.dbg("init storage=%T storageType=%s redisNil=%v", mng.Storage, config.StorageType, config.RedisClient == nil)

	setOnce.Do(func() {
		manager := sagin.NewManager(mng.Storage, saCfg)
		sagin.SetManager(manager)
		mng.dbg("sa-token manager initialized (once)")
	})
	return
}

// generateSaConfig 构建配置
func generateSaConfig(cfg *Config) (*config.Config, error) {
	if cfg.SaConfig != nil {
		return cfg.SaConfig, nil
	}
	sa := sagin.DefaultConfig()
	sa.TokenStyle = cfg.TokenStyle
	if sa.TokenStyle == config.TokenStyleJWT {
		if cfg.Salt == "" {
			return nil, errors.New("token salt is required for JWT")
		}
		sa.JwtSecretKey = cfg.Salt
	}
	if cfg.Timeout > 0 {
		sa.Timeout = int64(cfg.Timeout / time.Second)
	} else if sa.Timeout <= 0 {
		sa.Timeout = 3600
	}
	return sa, nil
}

// Login 生成令牌对
func (mng *IdentityMng) Login(_ context.Context, loginID string, device string) (TokenPair, error) {
	if device == "" {
		if mng.defaultDevice != "" {
			device = mng.defaultDevice
		} else {
			device = "client"
		}
	}
	info, err := stputil.LoginWithRefreshToken(loginID, device)
	if err != nil {
		return TokenPair{}, fmt.Errorf("login failed: %w", err)
	}
	pair := TokenPair{AccessToken: info.AccessToken, RefreshToken: info.RefreshToken}
	mng.dbg("login ok loginID=%s device=%s access=%s refresh=%s", loginID, device, pair.AccessToken, pair.RefreshToken)
	return pair, nil
}

// RefreshByLoginID 通过 loginID 直接刷新（等价于重新登录）
func (mng *IdentityMng) RefreshByLoginID(ctx context.Context, loginID string, device string) (TokenPair, error) {
	pair, err := mng.Login(ctx, loginID, device)
	if err == nil {
		mng.dbg("refresh ok loginID=%s device=%s access=%s", loginID, device, pair.AccessToken)
	}
	return pair, err
}

// Logout 注销当前会话
// LogoutCurrent 注销当前会话
func (mng *IdentityMng) LogoutCurrent(_ context.Context) error { return stputil.Logout(nil) }

// LogoutByLoginID 注销指定主体
func (mng *IdentityMng) LogoutByLoginID(_ context.Context, loginID string) error {
	err := stputil.Logout(loginID)
	if err == nil {
		mng.dbg("logout by loginID ok loginID=%s", loginID)
	}
	return err
}

// 为兼容旧签名，可保留一个 Logout 代理到当前会话
func (mng *IdentityMng) Logout(_ context.Context) error { return mng.LogoutCurrent(nil) }

// CurrentLoginID 获取当前登录ID
func (mng *IdentityMng) CurrentLoginID(_ context.Context) string {
	id, _ := stputil.GetLoginID("")
	mng.dbg("current login id=%s", id)
	return id
}

func (mng *IdentityMng) dbg(format string, args ...interface{}) {
	if mng == nil || !mng.debug {
		return
	}
	log.Printf("[identityMng] "+format, args...)
}

func (mng *IdentityMng) CheckLogin(token string) error {
	return stputil.CheckLogin(token)
}

func (mng *IdentityMng) GetLoginID(token string) (string, error) {
	return stputil.GetLoginID(token)
}
