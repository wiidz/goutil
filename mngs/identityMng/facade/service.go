package facade

import (
	"context"

	"github.com/wiidz/goutil/mngs/identityMng/authn"
	"github.com/wiidz/goutil/mngs/identityMng/subject/model"
	subrepo "github.com/wiidz/goutil/mngs/identityMng/subject/repo"
)

type Service struct {
	auth    *authn.Service
	subject subrepo.SubjectRepository
	// default subject type for this port (e.g., USER for client, STAFF for console)
	defaultSubjectType string
}

func New(auth *authn.Service, subject subrepo.SubjectRepository, defaultSubjectType string) *Service {
	return &Service{auth: auth, subject: subject, defaultSubjectType: defaultSubjectType}
}

func (s *Service) Login(ctx context.Context, loginId, device string) (accessToken, refreshToken string, err error) {
	if s.subject != nil {
		if _, err := s.subject.FindByTypeAndLoginID(ctx, s.defaultSubjectType, loginId); err != nil {
			_ = s.subject.Create(ctx, &model.Subject{SubjectType: s.defaultSubjectType, LoginID: loginId, Status: 1})
		}
	}
	return s.auth.Login(ctx, loginId, device)
}

func (s *Service) Logout(ctx context.Context) error { return s.auth.Logout(ctx) }

func (s *Service) CurrentLoginID(ctx context.Context) string { return s.auth.CurrentLoginID(ctx) }
