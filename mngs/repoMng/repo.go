package repoMng

import (
	"context"

	"gorm.io/gorm"
)

type Repo[T any] struct{ db *gorm.DB }

// RepoOf returns a generic repository for type T bound to db.
func RepoOf[T any](db *gorm.DB) *Repo[T] { return &Repo[T]{db: db} }

// Reads
func (r *Repo[T]) GetByID(ctx context.Context, id any) (*T, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	var m T
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repo[T]) First(ctx context.Context, opts ...Selector) (*T, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	qb := apply(r.db.WithContext(ctx), opts...)
	var m T
	if err := qb.First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repo[T]) List(ctx context.Context, opts ...Selector) ([]*T, int64, error) {
	if r.db == nil {
		return nil, 0, gorm.ErrInvalidDB
	}
	qb := apply(r.db.WithContext(ctx), opts...)
	var total int64
	if err := qb.Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	qb = applyPage(qb, opts...)
	var rows []T
	if err := qb.Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	res := make([]*T, 0, len(rows))
	for i := range rows {
		res = append(res, &rows[i])
	}
	return res, total, nil
}

// Writes
func (r *Repo[T]) Create(ctx context.Context, m *T) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *Repo[T]) Update(ctx context.Context, m *T, cols ...string) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}
	if len(cols) == 0 {
		return r.db.WithContext(ctx).Save(m).Error
	}
	return r.db.WithContext(ctx).Select(cols).Save(m).Error
}

func (r *Repo[T]) Delete(ctx context.Context, opts ...Selector) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}
	qb := apply(r.db.WithContext(ctx), opts...)
	return qb.Delete(new(T)).Error
}
