package appMng

import (
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/memoryMng"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"github.com/wiidz/goutil/mngs/redisMng"
	"time"
)

var cacheM = memoryMng.NewCacheMng()
var mysqlM = mysqlMng.NewMysqlMng()

type RunningMode int8           // 脚本运行模式
const Singleton RunningMode = 1 // 单例
const Multiton RunningMode = 2  // 多例

type SqlConfigLocation int8           // 配置文件存放位置
const LocalFile SqlConfigLocation = 1 // 本地文件，在/configs/目录下
const SqlRow SqlConfigLocation = 2    // 总库记录，例如center库中存放了以appID为主键的配置记录

type AppMng struct {
	ID                uint64            `gorm:"column:id" json:"id"`
	// RunningMode       RunningMode       // 脚本运行模式
	// SqlConfigLocation SqlConfigLocation // sql配置存放位置
	AppNo             string            `gorm:"column:app_no" json:"app_no"`
	AppName           string            `gorm:"column:app_name" json:"app_name"`
	Location          *time.Location    `gorm:"-" json:"-"`
	BaseConfig        *BaseConfig
	ProjectConfig     ProjectConfig
}

// GetSingletonAppMng : 获取单例app管理器
func GetSingletonAppMng(appID uint64, mysqlConfig *MysqlConfig, configStruct ProjectConfig) (mng *AppMng, err error) {

	//【1】从缓存中提取
	res, isExist := cacheM.Get("app-" + typeHelper.Uint64ToStr(appID) + "-config")
	if isExist && res != nil {
		return res.(*AppMng), nil
	}

	mng = &AppMng{
		ID:            appID,
		ProjectConfig: configStruct,
	}

	//【2】初始化mysql
	mysqlMng.Init(mysqlConfig)

	//【3】基础配置
	mng.BaseConfig, err = mng.SetBaseConfig(mysqlConfig.DbName, mysqlConfig.SettingTableName)
	if err != nil {
		return
	}
	mng.BaseConfig.MysqlConfig = mysqlConfig

	//【4】初始化redis、es
	redisMng.Init(mng.BaseConfig.RedisConfig)

	//【3】项目配置
	mng.ProjectConfig.Build()

	//【5】写入缓存
	cacheM.Set("app-"+typeHelper.Uint64ToStr(appID)+"-config", mng, time.Minute*30)

	//【4】返回
	return
}

func (mng *AppMng) SetBaseConfig(dbName string, tableName string) (config *BaseConfig, err error) {

	config = &BaseConfig{}

	//【1】初始化sql
	conn := mysqlM.GetDBConn(dbName)

	//【1】查询数据
	var rows []*DbSettingRow
	err = conn.Table(tableName).Where("belonging = ?", "system").Find(&rows).Error
	if err != nil {
		return
	}

	//【2】初始化配置
	config.Location = getLocationConfig(rows)
	config.WechatMiniConfig = getWechatMiniConfig(rows)
	config.WechatOaConfig = getWechatOaConfig(rows)
	config.WechatOpenConfig = getWechatOpenConfig(rows)
	config.AliPayConfig = getAliPayConfig(rows)
	config.OssConfig = getOssConfig(rows)
	return
}

func getLocationConfig(rows []*DbSettingRow) (location *time.Location) {
	timeZone := getRow(rows, "time_zone", "", "").Value
	if timeZone == "" {
		timeZone = "Asia/Shanghai" // 默认东八区
	}
	location, _ = time.LoadLocation(timeZone)
	return
}

func getMysqlConfig(rows []*DbSettingRow) *MysqlConfig {
	return &MysqlConfig{
		Host:     getRow(rows, "mysql", "host", "").Value,
		Username: getRow(rows, "mysql", "username", "").Value,
		Password: getRow(rows, "mysql", "password", "").Value,
		Port:     getRow(rows, "mysql", "port", "").Value,
	}
}

func getWechatMiniConfig(rows []*DbSettingRow) *WechatMiniConfig {
	return &WechatMiniConfig{
		AppID:     getRow(rows, "wechat", "mini", "app_id").Value,
		AppSecret: getRow(rows, "wechat", "mini", "app_secret").Value,
	}
}

func getWechatOaConfig(rows []*DbSettingRow) *WechatOaConfig {
	return &WechatOaConfig{
		AppID:     getRow(rows, "wechat", "oa", "app_id").Value,
		AppSecret: getRow(rows, "wechat", "oa", "app_secret").Value,
	}
}

func getWechatOpenConfig(rows []*DbSettingRow) *WechatOpenConfig {
	return &WechatOpenConfig{
		AppID:     getRow(rows, "wechat", "open", "app_id").Value,
		AppSecret: getRow(rows, "wechat", "open", "app_secret").Value,
	}
}

func getAliPayConfig(rows []*DbSettingRow) *AliPayConfig {
	return &AliPayConfig{
		AppID:      getRow(rows, "ali", "pay", "app_id").Value,
		PrivateKey: getRow(rows, "ali", "pay", "private_key").Value,
	}
}

func getOssConfig(rows []*DbSettingRow) *OssConfig {
	return &OssConfig{
		AccessKeyID:     getRow(rows, "ali", "oss", "access_key_id").Value,
		AccessKeySecret: getRow(rows, "ali", "oss", "access_key_secret").Value,
		Host:            getRow(rows, "ali", "oss", "host").Value,
		EndPoint:        getRow(rows, "ali", "oss", "end_point").Value,
		BucketName:      getRow(rows, "ali", "oss", "bucket_name").Value,
	}
}

// getRow : 从rows中提取row
func getRow(rows []*DbSettingRow, name, flag1, flag2 string) (row *DbSettingRow) {

	row = &DbSettingRow{}

	if len(rows) == 0 {
		return nil
	}

	for k := range rows {
		if rows[k].Name != name {
			continue
		}
		if flag1 != "" && rows[k].Flag1 != flag1 {
			continue
		}
		if flag2 != "" && rows[k].Flag2 != flag2 {
			continue
		}
		row = rows[k]
		break
	}

	return
}
