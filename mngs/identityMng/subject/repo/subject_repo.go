package repo

import (
	"context"

	"github.com/wiidz/goutil/mngs/identityMng/subject/model"
)

type SubjectRepository interface {
	GetByID(ctx context.Context, id uint64) (*model.Subject, error)
	List(ctx context.Context, page, pageSize int, subjectType, keyword string) ([]*model.Subject, int64, error)
	FindByTypeAndLoginID(ctx context.Context, subjectType, loginID string) (*model.Subject, error)
	Create(ctx context.Context, s *model.Subject) error
}
