package mysqlMng

import (
	"errors"
	"gorm.io/gorm"
)

/**
 * @func  : 通用方法 获取列表
 * @author: Wiidz
 * @date  : 2020-10-14
 * @params: [mysql] *mysqlMng.MysqlMng 数据库连接
 *			[list] dbStruct.List 查询结构体
 * @return: [err] error 错误
 */
func (mysql *MysqlMng) Read(list ReadInterface,isSingle,doCount bool) {

	//【1】初始化参数
	offset := list.GetOffset()
	condition := list.GetCondition()
	limit := list.GetLimit()
	order := list.GetOrder()
	preloads := list.GetPreloads()

	var model = list.GetRows()
	if isSingle {
		model = list.GetRow()
	}

	thisConn := mysql.GetConn()

	//【2】拼接
	if len(condition) > 0 {
		cons, vals, _ := WhereBuild(condition)
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
	if isSingle {
		err = thisConn.First(model).Error 		// 查单条
	} else {
		err = thisConn.Offset(offset).Limit(limit).Find(model).Error 		// 查若干
	}
	if err == nil {
		list.SetError(err)
		return
	}

	//【4】查count
	if doCount {
		var count int64
		thisConn = thisConn.Session(&gorm.Session{NewDB: true})
		if len(condition) > 0 {
			cons, vals, _ := WhereBuild(condition)
			thisConn = thisConn.Where(cons, vals...)
		}
		err = thisConn.Model(model).Count(&count).Error
		list.SetCount(count)
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
	//tableName := update.GetTableName()
	thisConn := mysql.GetConn()
	model := update.GetRow()
	if model == nil {
		return errors.New("")
	}

	//【2】拼接
	if len(condition) == 0 {
		return errors.New("条件不允许为空")
	}
	if len(value) == 0 {
		return errors.New("值不允许为空")
	}

	//【3】修改
	cons, vals, _ := WhereBuild(condition)
	thisConn = thisConn.Model(&model).Where(cons, vals...).Updates(value)

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
	thisConn := mysql.GetConn()
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
	thisConn := mysql.GetConn()

	//【2】拼接
	if len(condition) == 0 {
		return errors.New("条件不允许为空")
	}

	cons, vals, _ := WhereBuild(condition)
	thisConn = thisConn.Where(cons, vals...).Delete(row)

	//【2】提取结果
	err := thisConn.Error
	if err == nil {
		params.SetRowsAffected(thisConn.RowsAffected)
	}

	//【5】返回
	return err
}

/**
 * @func  : 获取新闻列表
 * @author: Wiidz
 * @date  : 2020-10-15
 * @params: [pageNow] int 当前页码
 *			[pageSize] int 页长
 * @return: [msg] string 消息体
 * 			[data] interface{} 数据
 * 			[statusCode] 状态码
 */
func (mysql *MysqlMng) SimpleGetListWithLog(read ReadInterface, userID, authID int) (msg string, data interface{}, statusCode int) {

	//【3】查询
	mysql.LogRead(read, userID, authID)
	if read.GetError() != nil {
		return read.GetError().Error(), nil, 400
	}

	//【4】返回
	return "ok", map[string]interface{}{
		"rows":  read.GetRows(),
		"count": read.GetCount(),
	}, 200
}

// SimpleGetDetailWithLog 简单获取记录
func (mysql *MysqlMng) SimpleGetDetailWithLog(params ReadInterface, userID, authID int) (msg string, data interface{}, statusCode int) {

	//【2】查询
	mysql.LogRead(params, userID, authID)
	if params.GetError() != nil {
		return params.GetError().Error(), nil, 400
	}

	//【3】返回
	return "ok", params.GetRows(), 200
}

// SimpleUpdate 简单修改
func (mysql *MysqlMng) SimpleUpdate(params UpdateInterface) (msg string, data interface{}, statusCode int) {

	//【1】修改
	err := mysql.Update(params)
	if err != nil {
		return err.Error(), nil, 400
	}

	//【2】返回
	return "ok", params.GetRowsAffected(), 201
}

// SimpleUpdateMany 简单修改多条
func (mysql *MysqlMng) SimpleUpdateMany(params UpdateInterface) (msg string, data interface{}, statusCode int) {

	//【2】修改
	err := mysql.Update(params)
	if err != nil {
		return err.Error(), nil, 400
	}

	return "ok", params.GetRowsAffected(), 201
}

// SimpleCreateOne 简单插入
func (mysql *MysqlMng) SimpleCreateOne(params InsertInterface) (msg string, data interface{}, statusCode int) {

	//【1】写入数据库
	mysql.CreateOne(params)
	if err := params.GetError(); err != nil {
		return err.Error(), nil, 400
	}

	//【2】返回
	return "ok", params.GetNewID(), 201
}

// SimpleDelete 简单删除
func (mysql *MysqlMng) SimpleDelete(params DeleteInterface) (msg string, data interface{}, statusCode int) {

	//【2】写入数据库
	_ = mysql.Delete(params)
	if err := params.GetError(); err != nil {
		return err.Error(), nil, 400
	}

	//【3】返回
	return "ok", params.GetRowsAffected(), 200
}

// SimpleGetList 简单获取列表
func (mysql *MysqlMng) SimpleGetList(read ReadInterface,isSingle,doCount bool) (msg string, data interface{}, statusCode int) {

	//【1】查询
	mysql.Read(read,isSingle,doCount)
	if read.GetError() != nil {
		return read.GetError().Error(), nil, 400
	}

	//【2】组装数据
	if doCount {
		data = map[string]interface{}{
			"rows":  read.GetRows(),
			"count": read.GetCount(),
		}
	} else {
		data = read.GetRows()
	}

	//【3】返回
	return "ok", data, 200
}

// SimpleGetDetail 简单获取详情
func (mysql *MysqlMng) SimpleGetDetail(params ReadInterface) (msg string, data interface{}, statusCode int) {

	//【1】查询
	mysql.Read(params,true,false)
	if params.GetError() != nil {
		return params.GetError().Error(), nil, 400
	}

	//【2】返回
	return "ok", params.GetRow(), 200
}
