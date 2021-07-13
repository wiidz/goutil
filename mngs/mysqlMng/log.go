package mysqlMng

import (
	"errors"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"gorm.io/gorm"
	"time"
)

const LogTableName = "a_user_operate_log"

type Log struct {
	ID           int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`   // 编号
	UserID       int       `gorm:"column:user_id;type:int(11);not null" json:"user_id"`    // 职员编号
	AuthID       int       `gorm:"column:auth_id;type:int(11);not null" json:"auth_id"`    // 模块
	Kind         int8      `gorm:"column:kind;type:tinyint(4);not null" json:"kind"`       // 操作类型：1=select，2=insert，3=update，4=delete
	Condition    string    `gorm:"column:condition;type:varchar(256)" json:"condition"`    // 条件
	ROrder       string    `gorm:"column:r_order;type:varchar(256)" json:"r_order"`        // 顺序
	ROffset      int       `gorm:"column:r_offset;type:int(11)" json:"r_offset"`           // offset
	RLimit       int       `gorm:"column:r_limit;type:int(11)" json:"r_limit"`             // limit
	Data         string    `gorm:"column:data;type:varchar(1028)" json:"data"`             // 修改数据
	RowsAffected int       `gorm:"column:affected_rows;type:int(11)" json:"affected_rows"` // 影响行数
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"`     // 创建时间
}

func (x *Log) TableName() string {
	return LogTableName
}

// 读操作记录
type LogReadCreate struct {
	ID           int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`   // 编号
	UserID       int       `gorm:"column:user_id;type:int(11);not null" json:"user_id"`    // 职员编号
	AuthID       int       `gorm:"column:auth_id;type:int(11);not null" json:"auth_id"`    // 模块
	Kind         int8      `gorm:"column:kind;type:tinyint(4);not null" json:"kind"`       // 操作类型：1=select，2=insert，3=update，4=delete
	Condition    string    `gorm:"column:condition;type:varchar(256)" json:"condition"`    // 条件
	ROrder       string    `gorm:"column:r_order;type:varchar(256)" json:"r_order"`        // 顺序
	ROffset      int       `gorm:"column:r_offset;type:int(11)" json:"r_offset"`           // offset
	RLimit       int       `gorm:"column:r_limit;type:int(11)" json:"r_limit"`             // limit
	RowsAffected int       `gorm:"column:affected_rows;type:int(11)" json:"affected_rows"` // 影响行数
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"`     // 创建时间
}

func (x *LogReadCreate) TableName() string {
	return LogTableName
}

// 修改操作记录
type LogUpdateCreate struct {
	ID           int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`   // 编号
	UserID       int       `gorm:"column:user_id;type:int(11);not null" json:"user_id"`    // 职员编号
	AuthID       int       `gorm:"column:auth_id;type:int(11);not null" json:"auth_id"`    // 模块
	Kind         int8      `gorm:"column:kind;type:tinyint(4);not null" json:"kind"`       // 操作类型：1=select，2=insert，3=update，4=delete
	Condition    string    `gorm:"column:condition;type:varchar(256)" json:"condition"`    // 条件
	Data         string    `gorm:"column:data;type:varchar(1028)" json:"data"`             // 修改数据
	RowsAffected int       `gorm:"column:affected_rows;type:int(11)" json:"affected_rows"` // 影响行数
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"`     // 创建时间
}

func (x *LogUpdateCreate) TableName() string {
	return LogTableName
}

// 修改操作记录
type LogDeleteCreate struct {
	ID           int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`   // 编号
	UserID       int       `gorm:"column:user_id;type:int(11);not null" json:"user_id"`    // 职员编号
	AuthID       int       `gorm:"column:auth_id;type:int(11);not null" json:"auth_id"`    // 模块
	Kind         int8      `gorm:"column:kind;type:tinyint(4);not null" json:"kind"`       // 操作类型：1=select，2=insert，3=update，4=delete
	Condition    string    `gorm:"column:condition;type:varchar(256)" json:"condition"`    // 条件
	RowsAffected int       `gorm:"column:affected_rows;type:int(11)" json:"affected_rows"` // 影响行数
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"`     // 创建时间
}

func (x *LogDeleteCreate) TableName() string {
	return LogTableName
}

// 插入改操作记录
type LogInsertCreate struct {
	ID           int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`   // 编号
	UserID       int       `gorm:"column:user_id;type:int(11);not null" json:"user_id"`    // 职员编号
	AuthID       int       `gorm:"column:auth_id;type:int(11);not null" json:"auth_id"`    // 模块
	Kind         int8      `gorm:"column:kind;type:tinyint(4);not null" json:"kind"`       // 操作类型：1=select，2=insert，3=update，4=delete
	Data         string    `gorm:"column:data;type:varchar(1028)" json:"data"`             // 修改数据
	RowsAffected int       `gorm:"column:affected_rows;type:int(11)" json:"affected_rows"` // 影响行数
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"`     // 创建时间
}

func (x *LogInsertCreate) TableName() string {
	return LogTableName
}

/**
 * @func  : 通用方法 获取列表
 * @author: Wiidz
 * @date  : 2020-10-14
 * @params: [mysql] *mysqlMng.MysqlMng 数据库连接
 *			[list] dbStruct.List 查询结构体
 * @return: [err] error 错误
 */
func (mysql *MysqlMng) LogRead(list ReadInterface, userID, authID int) {

	//【1】初始化参数
	offset := list.GetOffset()
	condition := list.GetCondition()
	limit := list.GetLimit()
	order := list.GetOrder()
	preloads := list.GetPreloads()
	rows := list.GetRows()

	thisConn := mysql.Conn

	//【2】拼接
	if len(condition) > 0 {
		cons, vals, _ := WhereBuild(condition)
		thisConn = thisConn.Where(cons, vals...)
	}
	if len(preloads) > 0 {
		for _, v := range preloads {
			thisConn = thisConn.Preload(v, func(db *gorm.DB) *gorm.DB {
				return db.Order("created_at desc")
			})
		}
	}
	if order != "" {
		thisConn = thisConn.Order(order)
	}

	//【3】查询rows
	var err error
	var rowsAffected int64
	if list.GetPageSize() == 1 {
		err = thisConn.First(rows).Error
	} else {
		count := int64(0)
		// rows
		err = thisConn.Offset(offset).Limit(limit).Find(rows).Error
		if err == nil {
			rowsAffected = thisConn.RowsAffected

			// count
			//mysql.NewCommonConn()
			thisConn = mysql.Conn
			if len(condition) > 0 {
				cons, vals, _ := WhereBuild(condition)
				thisConn = thisConn.Where(cons, vals...)
			}
			err = thisConn.Model(rows).Count(&count).Error
			list.SetCount(count)
		}
	}

	// 【4】记录操作
	go func() {
		mysql.NewCommonConn()
		jsonCondition, _ := typeHelper.JsonEncode(condition)
		mysql.Conn.Create(&LogReadCreate{
			UserID:       userID,
			AuthID:       authID,
			Kind:         int8(1),
			Condition:    jsonCondition,
			ROrder:       order,
			ROffset:      offset,
			RLimit:       limit,
			RowsAffected: int(rowsAffected),
		})
	}()

	//【4】返回
	list.SetError(err)
}

/**
 * @func  : 获取新闻列表
 * @author: Wiidz
 * @date  : 2020-04-15
 * @params: [pageNow] int 当前页码
 *			[pageSize] int 页长
 * 			[kind] int 新闻类型
 * @return: [msg] string 消息体
 * 			[data] interface{} 数据
 * 			[statusCode] 状态码
 */
func (mysql *MysqlMng) LogCreateOne(insert InsertInterface, userID, authID int) error {

	//【1】初始化参数
	row := insert.GetRow()

	thisConn := mysql.Conn
	if mysql.TransConn != nil {
		thisConn = mysql.TransConn
	}
	thisConn = thisConn.Create(row)

	//【2】提取结果
	err := thisConn.Error
	rowsAffected := thisConn.RowsAffected
	if err == nil {
		insert.SetRowsAffected(rowsAffected)
	}
	if test, ok := row.(IDInterface); ok {
		insert.SetNewID(test.GetID())
	}

	//【5】记录操作
	go func() {
		mysql.NewCommonConn()
		jsonValue, _ := typeHelper.JsonEncode(row)
		data := LogInsertCreate{
			UserID:       userID,
			AuthID:       authID,
			Kind:         int8(2),
			Data:         jsonValue,
			RowsAffected: int(rowsAffected),
		}
		mysql.Conn.Create(&data)
	}()

	//【5】返回
	return err
}

/**
 * @func  : 通用方法 获取列表
 * @author: Wiidz
 * @date  : 2020-10-14
 * @params: [mysql] *mysqlMng.MysqlMng 数据库连接
 *			[list] dbStruct.List 查询结构体
 * @return: [err] error 错误
 */
func (mysql *MysqlMng) LogUpdate(update UpdateInterface, userID, authID int) error {

	//【1】初始化参数
	condition := update.GetCondition()
	value := update.GetValue()
	tableName := update.GetTableName()
	thisConn := mysql.Conn
	if mysql.TransConn != nil {
		thisConn = mysql.TransConn
	}

	//【2】拼接
	if len(condition) == 0 {
		return errors.New("条件不允许为空")
	}
	if len(value) == 0 {
		return errors.New("值不允许为空")
	}
	//【3】查原始数据 - 暂时不实现

	//【4】修改
	value["updated_at"] = time.Now()
	cons, vals, _ := WhereBuild(condition)
	thisConn = thisConn.Table(tableName).Where(cons, vals...).Updates(value)

	//【4】提取结果
	err := thisConn.Error
	rowsAffected := thisConn.RowsAffected

	//【5】记录操作
	go func() {
		mysql.NewCommonConn()
		jsonCondition, _ := typeHelper.JsonEncode(condition)
		jsonValue, _ := typeHelper.JsonEncode(value)
		data := LogUpdateCreate{
			UserID:       userID,
			AuthID:       authID,
			Kind:         int8(3),
			Condition:    jsonCondition,
			Data:         jsonValue,
			RowsAffected: int(rowsAffected),
		}
		mysql.Conn.Create(&data)
	}()

	//【5】返回
	update.SetError(err)
	if err == nil {
		update.SetRowsAffected(rowsAffected)
	}
	return err
}

/**
 * @func  : 删除某一条新闻
 * @params: [mysql] mysqlMng *MysqlMng 数据库连接
 *          [newsID]  int 新闻的ID
 * @return: [err] error 错误信息
 */
func (mysql *MysqlMng) LogDelete(params DeleteInterface, userID, authID int) error {

	//【1】初始化参数
	condition := params.GetCondition()
	thisConn := mysql.Conn
	if mysql.TransConn != nil {
		thisConn = mysql.TransConn
	}

	row := params.GetRow()

	//【2】拼接
	if len(condition) == 0 {
		return errors.New("条件不允许为空")
	}

	cons, vals, _ := WhereBuild(condition)
	thisConn = thisConn.Where(cons, vals...).Delete(row)

	//【3】提取结果
	err := thisConn.Error
	rowsAffected := thisConn.RowsAffected
	if err == nil {
		params.SetRowsAffected(rowsAffected)
	}

	//【5】记录操作
	go func() {
		mysql.NewCommonConn()
		jsonCondition, _ := typeHelper.JsonEncode(condition)
		data := LogDeleteCreate{
			UserID:       userID,
			AuthID:       authID,
			Kind:         int8(4),
			Condition:    jsonCondition,
			RowsAffected: int(rowsAffected),
		}
		mysql.Conn.Create(&data)
	}()

	//【5】返回
	return err
}
