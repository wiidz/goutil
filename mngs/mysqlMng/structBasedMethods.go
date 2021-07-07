package mysqlMng

import (
	"errors"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"log"
)

var typeH = typeHelper.TypeHelper{}

/**
 * @func  : 通用方法 获取列表
 * @author: Wiidz
 * @date  : 2020-10-14
 * @params: [mysql] *mysqlMng.MysqlMng 数据库连接
 *			[list] dbStruct.List 查询结构体
 * @return: [err] error 错误
 */
func (mysql *MysqlMng) Read(list ReadInterface) {

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
		cons, vals, _ := mysql.WhereBuild(condition)
		thisConn = thisConn.Where(cons, vals...)
	}
	if len(preloads) > 0 {
		for _, v := range preloads {
			thisConn = thisConn.Preload(v)
		}
	}
	if order != "" {
		thisConn = thisConn.Order(order)
	}

	//【3】查询rows
	var err error
	if list.GetPageSize() == 1 {
		err = thisConn.First(rows).Error
	} else {
		count := int64(0)
		// rows
		err = thisConn.Offset(offset).Limit(limit).Find(rows).Error
		if err == nil {
			// count
			thisConn = mysql.Conn
			if len(condition) > 0 {
				cons, vals, _ := mysql.WhereBuild(condition)
				thisConn = thisConn.Where(cons, vals...)
			}
			err = thisConn.Model(rows).Count(&count).Error
			list.SetCount(count)
		}
	}

	//【4】返回
	list.SetError(err)
}

/**
 * @func  : 通用方法 获取列表
 * @author: Wiidz
 * @date  : 2020-10-14
 * @params: [mysql] *mysqlMng.MysqlMng 数据库连接
 *			[list] dbStruct.List 查询结构体
 * @return: [err] error 错误
 */
func (mysql *MysqlMng) Update(update UpdateInterface) error {

	//【1】初始化参数
	condition := update.GetCondition()
	value := update.GetValue()
	log.Println("【value】", typeH.GetType(value), value)
	tableName := update.GetTableName()
	thisConn := mysql.Conn

	//【2】拼接
	if len(condition) == 0 {
		return errors.New("条件不允许为空")
	}
	if len(value) == 0 {
		return errors.New("值不允许为空")
	}

	//【3】修改
	cons, vals, _ := mysql.WhereBuild(condition)
	thisConn = thisConn.Table(tableName).Where(cons, vals...).Updates(value)

	//【4】提取结果
	err := thisConn.Error
	if err == nil {
		update.SetRowsAffected(thisConn.RowsAffected)
	}

	//【5】返回
	return err
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
func (mysql *MysqlMng) CreateOne(insert InsertInterface) {

	//【1】初始化参数
	row := insert.GetRow()
	thisConn := mysql.Conn
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

	//【5】返回
	insert.SetError(err)
}

/**
 * @func  : 删除某一条新闻
 * @params: [mysql] mysqlMng *MysqlMng 数据库连接
 *          [newsID]  int 新闻的ID
 * @return: [err] error 错误信息
 */
func (mysql *MysqlMng) Delete(params DeleteInterface) error {

	//【1】初始化参数
	condition := params.GetCondition()
	row := params.GetRow()
	thisConn := mysql.Conn

	//【2】拼接
	if len(condition) == 0 {
		return errors.New("条件不允许为空")
	}

	cons, vals, _ := mysql.WhereBuild(condition)
	thisConn = thisConn.Where(cons, vals...).Delete(row)

	//【2】提取结果
	err := thisConn.Error
	if err == nil {
		params.SetRowsAffected(thisConn.RowsAffected)
	}

	//【5】返回
	return err
}
