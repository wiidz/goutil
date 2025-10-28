package identityMng

import (
	"context"
	"sync"

	sagin "github.com/click33/sa-token-go/integrations/gin"
	"github.com/click33/sa-token-go/storage/memory"
	stputil "github.com/click33/sa-token-go/stputil"
)

// TokenPair 标准令牌对
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// Claims 标准声明（按需扩展）
type Claims struct {
	LoginID  string
	Device   string
	ExpireAt int64
}

// Config 预留（当前无需配置即可使用）
type Config struct{}

var (
	setupOnce sync.Once
)

// Init 可选初始化；不调用也会在首次使用时自动初始化
func Init(_ ...Config) {
	ensureSetup()
}

func ensureSetup() {
	setupOnce.Do(func() {
		storage := memory.NewStorage()
		cfg := sagin.DefaultConfig()
		manager := sagin.NewManager(storage, cfg)
		sagin.SetManager(manager)
	})
}

// Register 仅工具库占位：不落库，由业务自行实现需要时可在外层调用本函数前完成注册逻辑
func Register(loginID string) error {
	_ = loginID
	return nil
}

// Login 生成令牌对
func Login(_ context.Context, loginID, device string) (TokenPair, error) {
	ensureSetup()
	if device == "" {
		device = "client"
	}
	info, err := stputil.LoginWithRefreshToken(loginID, device)
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{AccessToken: info.AccessToken, RefreshToken: info.RefreshToken}, nil
}

// RefreshByLoginID 通过 loginID 直接刷新（等价于重新登录）
func RefreshByLoginID(ctx context.Context, loginID, device string) (TokenPair, error) {
	return Login(ctx, loginID, device)
}

// Logout 注销当前会话
func Logout(_ context.Context) error {
	ensureSetup()
	return stputil.Logout(nil)
}

// CurrentLoginID 获取当前登录ID
func CurrentLoginID(_ context.Context) string {
	ensureSetup()
	id, _ := stputil.GetLoginID("")
	return id
}
