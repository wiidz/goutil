package authMng

import (
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/wiidz/goutil/helpers"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"reflect"
)

var typeHelper = helpers.TypeHelper{}
var sliceHelper = helpers.SliceHelper{}

type AuthMng struct {
	AuthTableName string // 身份表的名称例如user,staff
	OwnerTableName string // auth表的名称
	IdentifyKey string // 使用者在tokenData中的身份标识名称：例如user_id，staff_id
}

// New 获取权限验证器
func New(ownerTableName,authTableName,identifyKey string)*AuthMng{
	return &AuthMng{
		OwnerTableName:ownerTableName,
		AuthTableName:authTableName,
		IdentifyKey:identifyKey,
	}
}
// Serve 注入服务
func (mng *AuthMng) Serve(ctx iris.Context) {

	//【1】获取主键
	tokenData := ctx.Values().Get("token_data")
	immutable := reflect.ValueOf(tokenData)
	id := immutable.Elem().FieldByName(mng.IdentifyKey).Interface().(int)

	//【2】初始化数据库
	mysql := mysqlMng.NewMysqlMng()

	//【3】获取用户资料并判断
	owner,err := mng.getOwnerFromDB(mysql,id)
	if err != nil{
		if mysql.IsNotFound(err){
			helpers.ReturnError(ctx,"找不到您的账户")
			return
		}
		helpers.ReturnError(ctx,err.Error())
		return
	} else if owner.IsActive == 0{
		//判断用户是否被禁用
		helpers.ReturnError(ctx,"账户禁用中")
		return
	}

	//【4】查询当前请求地址的authID
	route := ctx.GetCurrentRoute()
	authID,err := mng.getAuthIDFromDB(mysql,route.Method(), route.Path())
	if err != nil {
		helpers.ReturnError(ctx,err.Error())
		return
	}

	//【4】判断客户的权限集是否包括
	if owner.Grouping != SuperManager{
		authIDs := typeHelper.ExplodeInt(owner.AuthIDs, ",")
		exist := sliceHelper.ExistInt(authID, authIDs)
		if !exist {
			helpers.ReturnError(ctx,"您没有权限操作")
			return
		}
	}

	//【5】继续下一步处理
	ctx.Next()
}


// getAuthFromDB 根据方法和路由地址获取对应权限的主键
func (mng *AuthMng) getAuthIDFromDB(mysql *mysqlMng.MysqlMng,method,route string)(authID int,err error){
	//【2】获取操作方法对应的数字
	numMap := map[string]Method{
		"GET":    Get,
		"POST":   Post,
		"PUT":    Put,
		"DELETE": Delete,
	}
	methodNum := numMap[method]

	var row DBAuthRow
	err = mysql.Conn.Table(mng.AuthTableName).Where("method = ? and route = ?",&methodNum,&route).First(&row).Error

	if mysql.IsNotFound(err){
		err = errors.New("找不到匹配路由的权限")
		return
	}
	authID = row.ID

	return
}

// getOwnerFromDB 从数据库根据主键值获取用户资料
func (mng *AuthMng) getOwnerFromDB(mysql *mysqlMng.MysqlMng,ownerID int)(owner DBAuthOwnerMixed,err error){

	err = mysql.Conn.Table(mng.OwnerTableName).Where("id = ?",&ownerID).First(&owner).Error
	return
}