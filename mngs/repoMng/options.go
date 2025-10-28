package repoMng

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Selector interface{ apply(*gorm.DB) *gorm.DB }

type selFn func(*gorm.DB) *gorm.DB

func (f selFn) apply(db *gorm.DB) *gorm.DB { return f(db) }

func apply(db *gorm.DB, opts ...Selector) *gorm.DB {
	for _, o := range opts {
		if o != nil {
			db = o.apply(db)
		}
	}
	return db
}

// Filters
func WithEq(field string, v any) Selector {
	return selFn(func(db *gorm.DB) *gorm.DB { return db.Where(field+" = ?", v) })
}
func WithIn(field string, vals any) Selector {
	return selFn(func(db *gorm.DB) *gorm.DB { return db.Where(field+" IN ?", vals) })
}
func WithLike(field, kw string) Selector {
	return selFn(func(db *gorm.DB) *gorm.DB { return db.Where(field+" LIKE ?", "%"+kw+"%") })
}

// Projection / preload
func WithSelect(cols ...string) Selector {
	return selFn(func(db *gorm.DB) *gorm.DB { return db.Select(strings.Join(cols, ",")) })
}
func WithPreload(path string, conds ...any) Selector {
	return selFn(func(db *gorm.DB) *gorm.DB { return db.Preload(path, conds...) })
}
func WithScopes(scopes ...func(*gorm.DB) *gorm.DB) Selector {
	return selFn(func(db *gorm.DB) *gorm.DB { return db.Scopes(scopes...) })
}

// Sorting / paging
func WithOrder(order string) Selector {
	return selFn(func(db *gorm.DB) *gorm.DB { return db.Order(order) })
}

// Read route
// WithWriteRoute forces query to use primary connection (for strong-consistent reads outside tx)
func WithWriteRoute() Selector {
	return selFn(func(db *gorm.DB) *gorm.DB { return db.Clauses(dbresolver.Write) })
}

type pager struct{ page, size int }

func WithPage(page, size int) Selector {
	return selFn(func(db *gorm.DB) *gorm.DB { return db.Set("__page__", pager{page: page, size: size}) })
}

func applyPage(db *gorm.DB, opts ...Selector) *gorm.DB {
	val, ok := db.Get("__page__")
	if !ok {
		return db
	}
	p, _ := val.(pager)
	page := p.page
	size := p.size
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}
	return db.Offset((page - 1) * size).Limit(size)
}
