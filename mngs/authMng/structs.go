package authMng

import "book/structs/baseStruct"


const (
	SuperManager int8 = 1
)

type Method int8
const (
	Get Method = 1
	Post Method = 2
	Put Method = 3
	Delete Method = 4
)

// DBAuthOwnerMixed 权限ID集合作为string存放在拥有者表里的情况
type DBAuthOwnerMixed struct {
	ID int `gorm:"column:id;type:int(11);not null" json:"id"`
	AuthIDs string `gorm:"column:auth_ids;type:text;not null" json:"auth_ids"`
	IsActive int8 `gorm:"column:is_active;type:tinyint" json:"is_active"` // 用来判断账户是否被禁用
	Grouping int8 `gorm:"column:grouping;type:tinyint" json:"grouping"` // 客户分组，1=超级管理员，不需要判断权限
}

// DBAuthRow 权限表的结构
type DBAuthRow struct {
	ID        int                    `gorm:"primary_key;column:id;type:int(11);not null" json:"id"`  // 编号
	Target    string                 `gorm:"column:target;type:varchar(128);not null" json:"target"` // 权限名称
	Name      string                 `gorm:"column:name;type:varchar(128);not null" json:"name"`     // 权限名称
	Method    int8                   `gorm:"column:method;type:tinyint(4);not null" json:"method"`   // 请求方法，1=get，2=post，3=put，4=delete
	Route     string                 `gorm:"column:route;type:varchar(128);not null" json:"route"`   // 路由地址
	Tips      string                 `gorm:"column:tips;type:varchar(128)" json:"tips"`              // 备注
	CreatedAt *baseStruct.MyJsonTime `gorm:"column:created_at;type:timestamp" json:"created_at"`     // 创建时间
	UpdatedAt *baseStruct.MyJsonTime `gorm:"column:updated_at;type:timestamp" json:"updated_at"`     // 最后修改时间
	DeletedAt *baseStruct.MyJsonTime `gorm:"column:deleted_at;type:timestamp" json:"-"`              // 最后修改时间
}

// DBAuthCreate 增加权限
type DBAuthCreate struct {
	ID        int                    `gorm:"primary_key;column:id;type:int(11);not null" json:"id"`  // 编号
	Target    string                 `gorm:"column:target;type:varchar(128);not null" json:"target"` // 权限名称
	Name      string                 `gorm:"column:name;type:varchar(128);not null" json:"name"`     // 权限名称
	Method    int8                   `gorm:"column:method;type:tinyint(4);not null" json:"method"`   // 请求方法，1=get，2=post，3=put，4=delete
	Route     string                 `gorm:"column:route;type:varchar(128);not null" json:"route"`   // 路由地址
	Tips      string                 `gorm:"column:tips;type:varchar(128)" json:"tips"`              // 备注
	CreatedAt *baseStruct.MyJsonTime `gorm:"column:created_at;type:timestamp" json:"created_at"`     // 创建时间
	UpdatedAt *baseStruct.MyJsonTime `gorm:"column:updated_at;type:timestamp" json:"updated_at"`     // 最后修改时间
}

// DBAuthOption 权限下拉菜单
type DBAuthOption struct {
	ID        int                    `gorm:"primary_key;column:id;type:int(11);not null" json:"id"`  // 编号
	Target    string                 `gorm:"column:target;type:varchar(128);not null" json:"target"` // 权限名称
	Name      string                 `gorm:"column:name;type:varchar(128);not null" json:"label"`    // 权限名称
	DeletedAt *baseStruct.MyJsonTime `gorm:"column:deleted_at;type:timestamp" json:"-"`              // 最后修改时间
}

// DBAuthPreview 权限预览数据
type DBAuthPreview struct {
	ID        int                    `gorm:"primary_key;column:id;type:int(11);not null" json:"id"`  // 编号
	Target    string                 `gorm:"column:target;type:varchar(128);not null" json:"target"` // 权限名称
	Name      string                 `gorm:"column:name;type:varchar(128);not null" json:"name"`     // 权限名称
	Method    int8                   `gorm:"column:method;type:tinyint(4);not null" json:"method"`   // 请求方法，1=get，2=post，3=put，4=delete
	Route     string                 `gorm:"column:route;type:varchar(128);not null" json:"route"` // 路由地址
	WebRouter string 				`gorm:"column:web_route;varchar(550)" json:"web_route"` // web项目的地址
	DeletedAt *baseStruct.MyJsonTime `gorm:"column:deleted_at;type:timestamp" json:"-"`              // 最后修改时间
}