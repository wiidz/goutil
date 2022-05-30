package mysqlStruct

import (
	"github.com/wiidz/goutil/helpers/timeHelper"
	"gorm.io/gorm"
	"time"
)

type Full struct {
	ID        uint64         `gorm:"primary_key;column:id;type:int(11);not null" json:"id"` // 主键
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp" json:"created_at"`    // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp" json:"updated_at"`    // 修改时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp" json:"-"`             // 删除时间
}

type NoDeletedAt struct {
	ID        uint64    `gorm:"primary_key;column:id;type:int(11);not null" json:"id"` // 主键
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"`    // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updated_at"`    // 修改时间
}

type NoUpdatedAt struct {
	ID        uint64         `gorm:"primary_key;column:id;type:int(11);not null" json:"id"` // 主键
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp" json:"created_at"`    // 创建时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp" json:"-"`             // 删除时间
}

type IDCreatedAt struct {
	ID        uint64    `gorm:"column:id" json:"id"`                 // 主键
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"` // 创建时间
}

type IDDeletedAt struct {
	ID        uint64         `gorm:"column:id" json:"id"`                       // 主键
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp" json:"-"` // 删除时间
}

type FullList struct {
	ID        uint64                 `gorm:"primary_key;column:id;type:int(11);not null" json:"id"` // 主键
	CreatedAt *timeHelper.MyJsonTime `gorm:"column:created_at;type:timestamp" json:"created_at"`    // 创建时间
	UpdatedAt *timeHelper.MyJsonTime `gorm:"column:updated_at;type:timestamp" json:"updated_at"`    // 修改时间
	DeletedAt gorm.DeletedAt         `gorm:"column:deleted_at;type:timestamp" json:"-"`             // 删除时间
}

type NoUpdatedAtList struct {
	ID        uint64                 `gorm:"primary_key;column:id;type:int(11);not null" json:"id"` // 主键
	CreatedAt *timeHelper.MyJsonTime `gorm:"column:created_at;type:timestamp" json:"created_at"`    // 创建时间
	DeletedAt gorm.DeletedAt         `gorm:"column:deleted_at;type:timestamp" json:"-"`             // 删除时间
}

type IDList struct {
	ID        uint64         `gorm:"primary_key;column:id;type:int(11);not null" json:"id"` // 主键
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp" json:"-"`             // 删除时间
}

type IDCreatedAtList struct {
	ID        uint64                 `gorm:"column:id" json:"id"`                 // 主键
	CreatedAt *timeHelper.MyJsonTime `gorm:"column:created_at" json:"created_at"` // 创建时间
}

type Sum struct {
	SumFloat64 float64 `gorm:"column:sum_float64"`
}

type ModelWithTableName interface {
	TableName() string
}
