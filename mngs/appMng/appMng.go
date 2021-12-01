package appMng

import (
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/memoryMng"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"time"
)

var cacheM = memoryMng.NewCacheMng()
var mysqlM = mysqlMng.NewMysqlMng()

// GetAppConfig 根据appID获取对应的参数设置
func GetAppConfig(appID uint64) (config *AppConfig, err error) {

	//【1】从缓存中提取
	res, isExist := cacheM.Get("app-" + typeHelper.Uint64ToStr(appID) + "-config")
	if isExist && res != nil {
		return res.(*AppConfig), nil
	}

	//【2】从数据库中提取
	centerConn := mysqlM.GetConn()
	err = centerConn.Table("center.c_app").Where("id = ?", &appID).First(&config).Error
	if err != nil {
		return
	}

	//【3】设置时区
	config.Location, _ = time.LoadLocation(config.TimeZone)
	//dbConn := mysqlM.GetDBConn(config.DbName)

	//【4】获取其他设置
	//【4-1】帖子阅读数设定
	//postConfig, err := GetPostConfig(dbConn)
	//if err != nil {
	//	return
	//}
	//config.PostConfig = postConfig

	//【5】写入缓存
	cacheM.Set("app-"+typeHelper.Uint64ToStr(appID)+"-config", config, time.Minute*30)

	//【4】换回
	return
}
