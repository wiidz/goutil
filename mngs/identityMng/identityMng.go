package identityMng

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/click33/sa-token-go/core/config"
	sagin "github.com/click33/sa-token-go/integrations/gin"
	"github.com/click33/sa-token-go/storage/memory"
	"github.com/click33/sa-token-go/stputil"
)

var setOnce sync.Once

// NewMng 创建实例（可多实例，首次创建时初始化 Sa-Token 全局管理器）
func NewMng(in *Config) (*IdentityMng, error) {
	cfg := in
	if cfg == nil {
		cfg = &Config{}
	}

	saCfg, err := generateSaConfig(cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Storage == nil {
		cfg.Storage = memory.NewStorage()
	}

	mng := &IdentityMng{config: saCfg, defaultDevice: cfg.DefaultDevice}
	setOnce.Do(func() {
		manager := sagin.NewManager(cfg.Storage, saCfg)
		sagin.SetManager(manager)
	})
	return mng, nil
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
	return TokenPair{AccessToken: info.AccessToken, RefreshToken: info.RefreshToken}, nil
}

// RefreshByLoginID 通过 loginID 直接刷新（等价于重新登录）
func (mng *IdentityMng) RefreshByLoginID(ctx context.Context, loginID string, device string) (TokenPair, error) {
	return mng.Login(ctx, loginID, device)
}

// Logout 注销当前会话
// LogoutCurrent 注销当前会话
func (mng *IdentityMng) LogoutCurrent(_ context.Context) error { return stputil.Logout(nil) }

// LogoutByLoginID 注销指定主体
func (mng *IdentityMng) LogoutByLoginID(_ context.Context, loginID string) error {
	return stputil.Logout(loginID)
}

// 为兼容旧签名，可保留一个 Logout 代理到当前会话
func (mng *IdentityMng) Logout(_ context.Context) error { return mng.LogoutCurrent(nil) }

// CurrentLoginID 获取当前登录ID
func (mng *IdentityMng) CurrentLoginID(_ context.Context) string {
	id, _ := stputil.GetLoginID("")
	return id
}
