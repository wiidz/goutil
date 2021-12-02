package appMng

import (
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
)

type CustomerConfig  struct {
	TestConfig *TestConfig
}

// Build : 这个方法完成设置的初始化
func (config *CustomerConfig) Build(){
	config.TestConfig = &TestConfig{
		Number:2,
	}
}

// TestConfig 项目中的设置
type TestConfig struct {
	Number uint64
}

func test(){
	//appM,_ := appMng.GetAppMng(1,"a_space","a_setting",&CustomerConfig{})
	appM,_ := GetSingletonAppMng(1,&configStruct.MysqlConfig{
		Host:             "localhost",
		Port:             "3306",
		Username:         "test",
		Password:         "test",
		DbName:           "a_space",
		Charset:          "utf8mb4",
		SettingTableName: "a_setting",
		Collation:        "utf8mb4_general_ci",
		TimeZone:         "Asia/Shanghai",
		MaxOpenConns:     5,
		MaxIdle:          10,
		MaxLifeTime:      60,
		ParseTime:        true,
	},&CustomerConfig{},&configStruct.CheckStart{
		Mysql: true,
		Redis: false,
		Es:    false,
	})
	log.Println("appM",appM.BaseConfig.Location,appM.ProjectConfig.(*CustomerConfig).TestConfig.Number)
}