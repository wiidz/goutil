package mysqlMng

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/wiidz/goutil/structs/configStruct"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type MysqlMng struct {
	db        *gorm.DB
	config    *configStruct.MysqlConfig   // 配置项
	Conn      *gorm.DB                    // 普通会话
	TransConn *gorm.DB                    // 事务会话
	Slaves    []*configStruct.MysqlConfig // 从库配置项
}

// NewMysqlMng 获取一个mysql实例
func NewMysqlMng(master *configStruct.MysqlConfig, slaves []*configStruct.MysqlConfig) (mysqlMng *MysqlMng, err error) {
	mysqlMng = &MysqlMng{
		config: master,
		Slaves: slaves,
	}

	err = mysqlMng.Init()
	return
}

// getDsn 从配置项获取dsn
func getDsn(config *configStruct.MysqlConfig) string {
	return config.Username + ":" + config.Password +
		"@tcp(" + config.Host + ":" + config.Port + ")/" + config.DbName +
		"?charset=" + config.Charset +
		"&collation=" + config.Collation +
		"&loc=" + url.QueryEscape(config.TimeZone) +
		"&parseTime=" + strconv.FormatBool(config.ParseTime)
}

func (mng *MysqlMng) Init() (err error) {

	//【1】构建DSN
	dsn := getDsn(mng.config)

	//【3】构建主库DB对象
	log.Println("【mysql-dsn】", dsn)
	mng.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: mng.config.Logger,
	})
	if err != nil {
		log.Println("【mysql-init-err】", err)
		return
	}
	sqlDB, _ := mng.db.DB()
	sqlDB.SetMaxIdleConns(mng.config.MaxIdle)                                     //最大空闲连接数
	sqlDB.SetMaxOpenConns(mng.config.MaxOpenConns)                                //最大连接数
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(mng.config.MaxLifeTime)) //设置连接空闲超时

	//【4】构建从库
	if mng.Slaves != nil && len(mng.Slaves) > 0 {
		log.Println("【mysql-slaves】", len(mng.Slaves))
		dialectors := []gorm.Dialector{}
		for _, v := range mng.Slaves {
			dialectors = append(dialectors, mysql.Open(getDsn(v)))
		}

		// 配置主从（读写分离）
		err = mng.db.Use(
			dbresolver.Register(
				// 指定源（主库）和 replicas（从库）
				dbresolver.Config{
					Replicas: dialectors,
				},
				// 指定作用于哪些表，可选
				// &User{},
			).
				// 配置读写连接池参数（可选）
				SetMaxIdleConns(mng.Slaves[0].MaxIdle).
				SetMaxOpenConns(mng.Slaves[0].MaxOpenConns),
		)
		log.Println("【mysql-slaves-init】", err)
	}

	return
}

// NewCommonConn 获取一个新的会话
func (mng *MysqlMng) NewCommonConn() {
	mng.Conn = mng.db.Session(&gorm.Session{
		//WithConditions: true,
		Logger: mng.config.Logger,
	})
}

// SetLogger 设置日志
func (mng *MysqlMng) SetLogger(logger logger.Interface) {
	mng.config.Logger = logger
	mng.db.Logger = logger
}

// GetConn 获取一个新的会话
func (mng *MysqlMng) GetConn() *gorm.DB {
	return mng.db.Session(&gorm.Session{
		//WithConditions: true,
		Logger: mng.config.Logger,
	})
}

// GetDBConn 获取一个新的会话，并且切换数据库名
func (mng *MysqlMng) GetDBConn(dbName string) *gorm.DB {
	return mng.db.Session(&gorm.Session{
		//WithConditions: true,
		Logger: mng.config.Logger,
	}).Exec("use " + dbName)
}

// GetDBTransConn 获取一个新的会话，并且切换数据库名
func (mng *MysqlMng) GetDBTransConn(dbName string) *gorm.DB {
	return mng.db.Begin().Exec("use " + dbName)
}

// NewTransConn 开启一个事务会话
func (mng *MysqlMng) NewTransConn() {
	mng.TransConn = mng.db.Session(&gorm.Session{
		//WithConditions: true,
		Logger: mng.config.Logger,
	}).Begin()
}

// Rollback 回滚事务
func (mng *MysqlMng) Rollback() {
	mng.TransConn.Rollback()
}

// Commit 提交事务
func (mng *MysqlMng) Commit() {
	mng.TransConn.Commit()
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
				//whereSQL += fmt.Sprint("`"+k+"`", "=?")
				whereSQL += fmt.Sprint(k, "=?")
			case ">":
				//whereSQL += fmt.Sprint("`"+k+"`", ">?")
				whereSQL += fmt.Sprint(k, ">?")
			case ">=":
				//whereSQL += fmt.Sprint("`"+k+"`", ">=?")
				whereSQL += fmt.Sprint(k, ">=?")
			case "<":
				//whereSQL += fmt.Sprint("`"+k+"`", "<?")
				whereSQL += fmt.Sprint(k, "<?")
			case "<=":
				//whereSQL += fmt.Sprint("`"+k+"`", "<=?")
				whereSQL += fmt.Sprint(k, "<=?")
			case "!=":
				//whereSQL += fmt.Sprint("`"+k+"`", "!=?")
				whereSQL += fmt.Sprint(k, "!=?")
			case "<>":
				//whereSQL += fmt.Sprint("`"+k+"`", "!=?")
				whereSQL += fmt.Sprint(k, "!=?")
			case "like":
				//whereSQL += fmt.Sprint("`"+k+"`", " like ?")
				whereSQL += fmt.Sprint(k, " like ?")
			}

			if v.([]interface{})[0] == "between" {
				whereSQL += fmt.Sprint(k, " between ? AND ?")
				vals = append(vals, v.([]interface{})[1], v.([]interface{})[2])
			} else if v.([]interface{})[0] == "in" || v.([]interface{})[0] == "not in" {
				whereSQL += fmt.Sprint(k, " "+v.([]interface{})[0].(string)+" (")
				if intSlice, ok := v.([]interface{})[1].([]int); ok {
					for k := 0; k < len(intSlice); k++ {
						whereSQL += "?,"
						vals = append(vals, intSlice[k])
					}
					whereSQL = whereSQL[0:len(whereSQL)-1] + ")"
				} else if intSlice, ok := v.([]interface{})[1].(*[]int); ok {
					for k := 0; k < len(*intSlice); k++ {
						whereSQL += "?,"
						vals = append(vals, (*intSlice)[k])
					}
					whereSQL = whereSQL[0:len(whereSQL)-1] + ")"
				} else if int8Slice, ok := v.([]interface{})[1].([]int8); ok {
					for k := 0; k < len(int8Slice); k++ {
						whereSQL += "?,"
						vals = append(vals, (int8Slice)[k])
					}
					whereSQL = whereSQL[0:len(whereSQL)-1] + ")"
				} else if uint64Slice, ok := v.([]interface{})[1].([]uint64); ok {
					for k := 0; k < len(uint64Slice); k++ {
						whereSQL += "?,"
						vals = append(vals, (uint64Slice)[k])
					}
					whereSQL = whereSQL[0:len(whereSQL)-1] + ")"
				} else {
					err = errors.New("无法匹配的类型")
				}
			} else {
				vals = append(vals, v.([]interface{})[1])
			}

		default:
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					//whereSQL += fmt.Sprint("`"+k+"`", " IS NOT NULL")
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					//whereSQL += fmt.Sprint("`"+k+"`", " IS NULL")
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				//whereSQL += fmt.Sprint("`"+k+"`", "=?")
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
				//whereSQL += fmt.Sprint("`"+k+"`", "=?")
				whereSQL += fmt.Sprint(k, "=?")
			case ">":
				//whereSQL += fmt.Sprint("`"+k+"`", ">?")
				whereSQL += fmt.Sprint(k, ">?")
			case ">=":
				//whereSQL += fmt.Sprint("`"+k+"`", ">=?")
				whereSQL += fmt.Sprint(k, ">=?")
			case "<":
				//whereSQL += fmt.Sprint("`"+k+"`", "<?")
				whereSQL += fmt.Sprint(k, "<?")
			case "<=":
				//whereSQL += fmt.Sprint("`"+k+"`", "<=?")
				whereSQL += fmt.Sprint(k, "<=?")
			case "!=":
				//whereSQL += fmt.Sprint("`"+k+"`", "!=?")
				whereSQL += fmt.Sprint(k, "!=?")
			case "<>":
				//whereSQL += fmt.Sprint("`"+k+"`", "!=?")
				whereSQL += fmt.Sprint(k, "!=?")
			case "like":
				//whereSQL += fmt.Sprint("`"+k+"`", " like ?")
				whereSQL += fmt.Sprint(k, " like ?")
			}

			if v.([]interface{})[0] == "between" {
				//whereSQL += fmt.Sprint("`"+k+"`", " between ? AND ?")
				whereSQL += fmt.Sprint(k, " between ? AND ?")
				vals = append(vals, v.([]interface{})[1], v.([]interface{})[2])
			} else if v.([]interface{})[0] == "in" {
				//whereSQL += fmt.Sprint("`"+k+"`", " IN (")
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
					//whereSQL += fmt.Sprint("`"+k+"`", " IS NOT NULL")
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					//whereSQL += fmt.Sprint("`"+k+"`", " IS NULL")
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				//whereSQL += fmt.Sprint("`"+k+"`", "=?")
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
func WhereFindAll(conn *gorm.DB, condition map[string]interface{}, rows interface{}) (err error) {

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

//var user dbStruct.User
//err := mysqlM.GetConn().Model(dbStruct.User{
//	ID: 2,
//}).First(&user).Error

// ok
//var user = dbStruct.User{ID: 10}
//err := mysqlM.GetConn().First(&user).Error

// ok Where中放map
//var users []dbStruct.User
//err := mysqlM.GetConn().Where(map[string]interface{}{
//	dbStruct.UserColumns.ID:           2,
//	dbStruct.UserColumns.PointBalance: 0, // 这一条 会作为条件
//}).Find(&users).Error

// ok Where中放model
//var users []dbStruct.User
//err := mysqlM.GetConn().Where(&dbStruct.User{
//	ID:           2,
//	PointBalance: 0, // 这一条 并不会作为条件 注意
//}).Find(&users).Error
