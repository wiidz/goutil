package appMng

import "log"

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
	appM,_ := GetSingletonAppMng(1,&MysqlConfig{
		Host:             "localhost",
		Port:             "3306",
		Username:         "test",
		Password:         "test",
		DbName:           "a_space",
		SettingTableName: "a_setting",
	},&CustomerConfig{})
	log.Println("appM",appM.BaseConfig.Location,appM.ProjectConfig.(*CustomerConfig).TestConfig.Number)
}