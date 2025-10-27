package service

import (
	"context"

	"github.com/wiidz/goutil/mngs/identityMng/subject/dto"
	"github.com/wiidz/goutil/mngs/identityMng/subject/model"
	subrepo "github.com/wiidz/goutil/mngs/identityMng/subject/repo"
)

type Service struct{ r subrepo.SubjectRepository }

func New(r subrepo.SubjectRepository) *Service { return &Service{r: r} }

func (s *Service) GetByID(ctx context.Context, id uint64) (*model.Subject, error) {
	// Adapt gorm repo model if needed via field copy
	m, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.Subject{ID: m.ID, SubjectType: m.SubjectType, LoginID: m.LoginID, ExternalID: m.ExternalID, Status: m.Status}, nil
}

func (s *Service) List(ctx context.Context, req dto.ListSubjectsRequest) ([]*model.Subject, int64, error) {
	rows, total, err := s.r.List(ctx, req.Page, req.PageSize, req.SubjectType, req.Keyword)
	if err != nil {
		return nil, 0, err
	}
	res := make([]*model.Subject, 0, len(rows))
	for _, m := range rows {
		res = append(res, &model.Subject{ID: m.ID, SubjectType: m.SubjectType, LoginID: m.LoginID, ExternalID: m.ExternalID, Status: m.Status})
	}
	return res, total, nil
}
