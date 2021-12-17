package appMng

import (
	"github.com/wiidz/goutil/structs/configStruct"
	"time"
)

type RunningMode int8           // 脚本运行模式
const Singleton RunningMode = 1 // 单例
const Multiton RunningMode = 2  // 多例

type SqlConfigLocation int8           // 配置文件存放位置
const LocalFile SqlConfigLocation = 1 // 本地文件，在/configs/目录下
const SqlRow SqlConfigLocation = 2    // 总库记录，例如center库中存放了以appID为主键的配置记录

type AppMng struct {
	ID uint64 `gorm:"column:id" json:"id"`
	// RunningMode       RunningMode       // 脚本运行模式
	// SqlConfigLocation SqlConfigLocation // sql配置存放位置
	Location      *time.Location `gorm:"-" json:"-"`
	BaseConfig    *configStruct.BaseConfig
	ProjectConfig configStruct.ProjectConfig
}

/******sql******
CREATE TABLE `u_setting` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `kind` tinyint(4) DEFAULT NULL COMMENT '类别，1=一般设定，2=页面配置',
  `belonging` varchar(128) DEFAULT NULL COMMENT '类别',
  `name` varchar(24) DEFAULT NULL COMMENT '名称',
  `flag_1` varchar(128) DEFAULT NULL COMMENT '【属性】补充的一个标识符1',
  `flag_2` varchar(128) DEFAULT NULL COMMENT '【属性】补充的一个标识符2',
  `value` text COMMENT '值',
  `value_1` text COMMENT '值-2',
  `value_2` text COMMENT '值-1',
  `tips` varchar(255) DEFAULT NULL COMMENT '说明',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '【时间】创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '【时间】最后修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '【时间】删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8
******sql******/
// SettingDbRow [...]
type DbSettingRow struct {
	ID        uint64    `gorm:"primary_key;column:id;type:int(11);not null" json:"id"` // sku编号
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"`    // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updated_at"`    // 修改时间
	Kind      int8      `gorm:"column:kind;type:tinyint(4)" json:"kind"`               // 类别，1=一般设定，2=页面配置
	Belonging string    `gorm:"column:belonging;type:varchar(128)" json:"belonging"`   // 类别
	Name      string    `gorm:"column:name;type:varchar(24)" json:"name"`              // 名称
	Flag1     string    `gorm:"column:flag_1;type:varchar(128)" json:"flag_1"`         // 【属性】补充的一个标识符1
	Flag2     string    `gorm:"column:flag_2;type:varchar(128)" json:"flag_2"`         // 【属性】补充的一个标识符2
	Value     string    `gorm:"column:value;type:text" json:"value"`                   // 值
	Value1    string    `gorm:"column:value_1;type:text" json:"value_1"`               // 值-2
	Value2    string    `gorm:"column:value_2;type:text" json:"value_2"`               // 值-1
	Tips      string    `gorm:"column:tips;type:varchar(255)" json:"tips"`             // 说明
}

// SettingPage 页面设置（带json decode）
type SettingPage struct {
	Kind        int8        `gorm:"column:kind;type:tinyint(4)" json:"-"`        // 类别，1=一般设定，2=页面配置
	Belonging   string      `gorm:"column:belonging;type:varchar(128)" json:"-"` // 类别
	Name        string      `gorm:"column:name;type:varchar(24)" json:"name"`    // 名称
	Value       string      `gorm:"column:value;type:text" json:"-"`             // 值
	ValueParsed interface{} `gorm:"-" json:"value"`                              // 值
}
