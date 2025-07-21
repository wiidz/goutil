package authMng

import (
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/helpers/sliceHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"github.com/wiidz/goutil/structs/configStruct"
	"gorm.io/gorm"
	"reflect"
)

type AuthMng struct {
	mysqlConfig    *configStruct.MysqlConfig
	AuthTableName  string // 身份表的名称例如user,staff
	OwnerTableName string // auth表的名称
	IdentifyKey    string // 使用者在tokenData中的身份标识名称：例如user_id，staff_id
}

// New 获取权限验证器
func New(mysqlConfig *configStruct.MysqlConfig, ownerTableName, authTableName, identifyKey string) *AuthMng {
	return &AuthMng{
		mysqlConfig:    mysqlConfig,
		OwnerTableName: ownerTableName,
		AuthTableName:  authTableName,
		IdentifyKey:    identifyKey,
	}
}

// Serve 注入服务
func (mng *AuthMng) Serve(ctx iris.Context) {

	//【1】获取主键
	tokenData := ctx.Values().Get("token_data")

	immutable := reflect.ValueOf(tokenData)
	if immutable.Elem().FieldByName(mng.IdentifyKey).IsValid() == false {
		networkHelper.ReturnError(ctx, "无效的token")
		return
	}

	id := immutable.Elem().FieldByName(mng.IdentifyKey).Interface().(uint64)
	if id == 0 {
		networkHelper.ReturnError(ctx, "登陆主体为空")
		return
	}

	//【2】初始化数据库
	mysql, _ := mysqlMng.NewMysqlMng(mng.mysqlConfig, nil)
	conn := mysql.GetConn()

	//【3】获取用户资料并判断
	owner, err := mng.getOwnerFromDB(conn, id)
	if err != nil {
		if mysqlMng.IsNotFound(err) {
			networkHelper.ReturnError(ctx, "找不到您的账户")
			return
		}
		networkHelper.ReturnError(ctx, err.Error())
		return
	} else if owner.Status == 0 {
		//判断用户是否被禁用
		networkHelper.ReturnError(ctx, "账户禁用中")
		return
	}

	//【4】查询当前请求地址的authID
	route := ctx.GetCurrentRoute()
	authID, err := mng.getAuthIDFromDB(conn, route.Method(), route.Path())
	if err != nil {
		networkHelper.ReturnError(ctx, err.Error())
		return
	}

	//【4】判断客户的权限集是否包括
	if owner.Grouping != SuperManager {
		authIDs := typeHelper.ExplodeUint64(owner.AuthIDs, ",")
		exist := sliceHelper.ExistUint64(authID, authIDs)
		if !exist {
			networkHelper.ReturnError(ctx, "您没有权限操作")
			return
		}
	}

	//【5】继续下一步处理
	ctx.Next()
}

// getAuthFromDB 根据方法和路由地址获取对应权限的主键
func (mng *AuthMng) getAuthIDFromDB(conn *gorm.DB, method, route string) (authID uint64, err error) {
	//【2】获取操作方法对应的数字
	numMap := map[string]Method{
		"GET":    Get,
		"POST":   Post,
		"PUT":    Put,
		"DELETE": Delete,
	}
	methodNum := numMap[method]

	var row DBAuthRow
	err = conn.Table(mng.AuthTableName).Where("method = ? and route = ?", &methodNum, &route).First(&row).Error

	if mysqlMng.IsNotFound(err) {
		err = errors.New("找不到匹配路由的权限")
		return
	}
	authID = row.ID

	return
}

// getOwnerFromDB 从数据库根据主键值获取用户资料
func (mng *AuthMng) getOwnerFromDB(conn *gorm.DB, ownerID uint64) (owner DBAuthOwnerMixed, err error) {

	err = conn.Table(mng.OwnerTableName).Where("id = ?", &ownerID).First(&owner).Error
	return
}
