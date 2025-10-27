package authn

import (
	"context"

	stputil "github.com/click33/sa-token-go/stputil"
)

// Service provides simple authentication operations backed by Sa-Token.
type Service struct{}

func New() *Service { return &Service{} }

func (s *Service) Login(ctx context.Context, loginId, device string) (accessToken, refreshToken string, err error) {
	if device == "" {
		device = "client"
	}
	info, err := stputil.LoginWithRefreshToken(loginId, device)
	if err != nil {
		return "", "", err
	}
	return info.AccessToken, info.RefreshToken, nil
}

func (s *Service) Logout(ctx context.Context) error { return stputil.Logout(nil) }

func (s *Service) CurrentLoginID(ctx context.Context) string {
	id, _ := stputil.GetLoginID("")
	return id
}
