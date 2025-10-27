package gorm

import (
	"context"

	gormlib "gorm.io/gorm"

	"github.com/wiidz/goutil/mngs/identityMng/subject/model"
)

// SubjectEntity is embedded here to keep goutil independent from gin_template's entity path.
type SubjectEntity struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	SubjectType string `gorm:"index;size:16;not null"`
	LoginID     string `gorm:"index;size:128;not null"`
	ExternalID  string `gorm:"size:128"`
	Status      int    `gorm:"not null;default:1"`
	CreatedAt   int64  `gorm:"autoCreateTime"`
	UpdatedAt   int64  `gorm:"autoUpdateTime"`
}

func EntitiesForMigrate() []interface{} { return []interface{}{&SubjectEntity{}} }

type subjectRepository struct{ db *gormlib.DB }

func NewSubjectRepository(db *gormlib.DB) *subjectRepository { return &subjectRepository{db: db} }

// Model mirrors minimal view used by repo.
func (r *subjectRepository) GetByID(ctx context.Context, id uint64) (*model.Subject, error) {
	if r.db == nil {
		return nil, gormlib.ErrInvalidDB
	}
	var se SubjectEntity
	if err := r.db.WithContext(ctx).First(&se, id).Error; err != nil {
		return nil, err
	}
	return &model.Subject{ID: se.ID, SubjectType: se.SubjectType, LoginID: se.LoginID, ExternalID: se.ExternalID, Status: se.Status}, nil
}

func (r *subjectRepository) List(ctx context.Context, page, pageSize int, subjectType, keyword string) ([]*model.Subject, int64, error) {
	if r.db == nil {
		return nil, 0, gormlib.ErrInvalidDB
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	qb := r.db.WithContext(ctx).Model(&SubjectEntity{})
	if subjectType != "" {
		qb = qb.Where("subject_type = ?", subjectType)
	}
	if keyword != "" {
		qb = qb.Where("login_id LIKE ?", "%"+keyword+"%")
	}
	var total int64
	if err := qb.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []SubjectEntity
	if err := qb.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	res := make([]*model.Subject, 0, len(rows))
	for _, se := range rows {
		res = append(res, &model.Subject{ID: se.ID, SubjectType: se.SubjectType, LoginID: se.LoginID, ExternalID: se.ExternalID, Status: se.Status})
	}
	return res, total, nil
}

func (r *subjectRepository) FindByTypeAndLoginID(ctx context.Context, subjectType, loginID string) (*model.Subject, error) {
	if r.db == nil {
		return nil, gormlib.ErrInvalidDB
	}
	var se SubjectEntity
	if err := r.db.WithContext(ctx).Where("subject_type = ? AND login_id = ?", subjectType, loginID).First(&se).Error; err != nil {
		return nil, err
	}
	return &model.Subject{ID: se.ID, SubjectType: se.SubjectType, LoginID: se.LoginID, ExternalID: se.ExternalID, Status: se.Status}, nil
}

func (r *subjectRepository) Create(ctx context.Context, s *model.Subject) error {
	if r.db == nil {
		return gormlib.ErrInvalidDB
	}
	se := SubjectEntity{SubjectType: s.SubjectType, LoginID: s.LoginID, ExternalID: s.ExternalID, Status: s.Status}
	if err := r.db.WithContext(ctx).Create(&se).Error; err != nil {
		return err
	}
	s.ID = se.ID
	return nil
}
