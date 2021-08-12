package mysqlMng

import (
	"errors"
	"fmt"
	"github.com/wiidz/goutil/mngs/configMng"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

var db *gorm.DB

type MysqlMng struct {
	Conn      *gorm.DB // 普通会话
	TransConn *gorm.DB // 事务会话
}

/**
 * @func:   init 获取mysql配置
 * @author: Wiidz
 * @date:   2020-04-15
 */
func init() {

	log.Println("mysql init")

	//【1】获取配置
	config := configMng.GetMysql()

	//【2】构建DSN
	dsn := config.Username + ":" + config.Password +
		"@tcp(" + config.Host + ":" + config.Port + ")/" + config.DbName +
		"?charset=" + config.Charset +
		"&collation=" + config.Collation +
		"&loc=Asia%2FShanghai" +
		"&parseTime=true"

	//【3】构建DB对象
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(5)                   //最大空闲连接数
	sqlDB.SetMaxOpenConns(10)                  //最大连接数
	sqlDB.SetConnMaxLifetime(time.Second * 10) //设置连接空闲超时

	//err := sqlDB.Ping()
	//log.Println("err", err)
}

/**
 * @func:   NewMysqlMng mysql管理器工厂模式
 * @author: Wiidz
 * @date:   2020-04-15
 */
func NewMysqlMng() *MysqlMng {
	mysqlMng := &MysqlMng{}
	return mysqlMng
}

// 获取一个新的会话
func (mysql *MysqlMng) NewCommonConn() {
	mysql.Conn = db.Session(&gorm.Session{
		//WithConditions: true,
		Logger: db.Logger.LogMode(logger.Info),
	})
}

// GetConn 获取一个新的会话
func (mysql *MysqlMng) GetConn() *gorm.DB {
	return db.Session(&gorm.Session{
		//WithConditions: true,
		Logger: db.Logger.LogMode(logger.Info),
	})
}


// GetDBConn 获取一个新的会话，并且切换数据库名
func (mysql *MysqlMng) GetDBConn(dbName string) *gorm.DB {
	return db.Session(&gorm.Session{
		//WithConditions: true,
		Logger: db.Logger.LogMode(logger.Info),
	}).Exec("use "+dbName)
}


// GetDBTransConn 获取一个新的会话，并且切换数据库名
func (mysql *MysqlMng) GetDBTransConn(dbName string) *gorm.DB {
	return db.Begin().Exec("use "+dbName)
}


// 开启一个事务会话
func (mysql *MysqlMng) NewTransConn() {
	mysql.TransConn = db.Session(&gorm.Session{
		//WithConditions: true,
		Logger: db.Logger.LogMode(logger.Info),
	}).Begin()
}

// 回滚事务
func (mysql *MysqlMng) Rollback() {
	mysql.TransConn.Rollback()
}

// 提交事务
func (mysql *MysqlMng) Commit() {
	mysql.TransConn.Commit()
}

// IsNotFound 判断读取结果是否为空错误
func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// WhereBuild 复合condition成为cons、vals的结构
func WhereBuild(condition map[string]interface{}) (whereSQL string, vals []interface{}, err error) {

	for k, v := range condition {
		if whereSQL != "" {
			whereSQL += " AND "
		}
		switch v.(type) {
		case []interface{}:
			switch v.([]interface{})[0] {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
			case "<":
				whereSQL += fmt.Sprint(k, "<?")
			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
			case "<>":
				whereSQL += fmt.Sprint(k, "!=?")
			case "like":
				whereSQL += fmt.Sprint(k, " like ?")
			}

			if v.([]interface{})[0] == "between" {
				whereSQL += fmt.Sprint(k, " between ? AND ?")
				vals = append(vals, v.([]interface{})[1], v.([]interface{})[2])
			} else if v.([]interface{})[0] == "in" {
				whereSQL += fmt.Sprint(k, " IN (")
				if intSlice, ok := v.([]interface{})[1].([]int); ok {
					for k := 0; k < len(intSlice); k++ {
						whereSQL += "?,"
						vals = append(vals, intSlice[k])
					}
					whereSQL = whereSQL[0:len(whereSQL)-1] + ")"
				} else if intSlice, ok := v.([]interface{})[1].(*[]int); ok {
					for k := 0; k < len((*intSlice)); k++ {
						whereSQL += "?,"
						vals = append(vals, (*intSlice)[k])
					}
					whereSQL = whereSQL[0:len(whereSQL)-1] + ")"
				}
			} else {
				vals = append(vals, v.([]interface{})[1])
			}

		default:
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
			}
		}
	}
	return
}

// WhereOrBuild 复合condition成为cons、vals的结构
func WhereOrBuild(condition map[string]interface{}) (whereSQL string, vals []interface{}, err error) {

	for k, v := range condition {
		if whereSQL != "" {
			whereSQL += " OR "
		}
		switch v.(type) {
		case []interface{}:
			switch v.([]interface{})[0] {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
			case "<":
				whereSQL += fmt.Sprint(k, "<?")
			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
			case "<>":
				whereSQL += fmt.Sprint(k, "!=?")
			case "like":
				whereSQL += fmt.Sprint(k, " like ?")
			}

			if v.([]interface{})[0] == "between" {
				whereSQL += fmt.Sprint(k, " between ? AND ?")
				vals = append(vals, v.([]interface{})[1], v.([]interface{})[2])
			} else if v.([]interface{})[0] == "in" {
				whereSQL += fmt.Sprint(k, " IN (")
				if intSlice, ok := v.([]interface{})[1].([]int); ok {
					for k := 0; k < len(intSlice); k++ {
						whereSQL += "?,"
						vals = append(vals, intSlice[k])
					}
					whereSQL = whereSQL[0:len(whereSQL)-1] + ")"
				} else if intSlice, ok := v.([]interface{})[1].(*[]int); ok {
					for k := 0; k < len((*intSlice)); k++ {
						whereSQL += "?,"
						vals = append(vals, (*intSlice)[k])
					}
					whereSQL = whereSQL[0:len(whereSQL)-1] + ")"
				}
			} else {
				vals = append(vals, v.([]interface{})[1])
			}

		default:
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
			}
		}
	}
	return
}

// IsExist 查询是否存在记录
func IsExist(conn *gorm.DB, condition map[string]interface{}, tableName string) (err error) {
	cons, vals, err := WhereBuild(condition)
	if err != nil {
		return err
	}
	var count int64
	err = conn.Table(tableName).Where(cons, vals...).Limit(1).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	return errors.New("记录已存在")
}

// RequireExist 要求存在，否则都报错
func RequireExist(conn *gorm.DB, condition map[string]interface{}, tableName string) (err error) {
	cons, vals, err := WhereBuild(condition)
	if err != nil {
		return err
	}
	var count int64
	err = conn.Table(tableName).Where(cons, vals...).Limit(1).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("记录不存在")
	}
	return nil
}

// WhereFindAll 条件查询全部
func WhereFindAll(conn *gorm.DB, condition map[string]interface{}, rows interface{}) ( err error) {

	cons, vals, err := WhereBuild(condition)
	if err != nil {
		return
	}
	err = conn.Where(cons, vals...).Find(&rows).Error
	if err != nil {
		return
	}
	return
}

// WhereFirst 条件查询一条
func WhereFirst(conn *gorm.DB, condition map[string]interface{}, rows interface{}) (err error) {

	cons, vals, err := WhereBuild(condition)
	if err != nil {
		return
	}
	err = conn.Where(cons, vals...).First(&rows).Error
	if err != nil {
		return
	}
	return
}

// GetOffset 获取偏移量
func GetOffset(pageNow, pageSize int) int {
	var offset int
	if pageNow > 1 {
		offset = (pageNow - 1) * pageSize
	} else {
		offset = 0
	}
	return offset
}



// SimpleUpdate 简单更新操作
func SimpleUpdate(conn *gorm.DB,condition,data map[string]interface{},model interface{},tableName string)(rowsAffected int64,err error) {

	//【1】判断条件是否为空
	if len(condition) == 0 {
		err = errors.New("条件为空")
		return
	}

	//【2】判断值是否为空
	if len(data) == 0 {
		err = errors.New("值为空")
		return
	}

	//【3】拆分条件
	cons,vals, _ := WhereBuild(condition)

	//【4】执行操作
	conn = conn.Debug().Table(tableName).Model(model).Where(cons,vals...).Updates(data)

	//【5】返回
	err = conn.Error
	rowsAffected = conn.RowsAffected

	return
}