package configHelper

// AppConfig 例子
type AppConfig struct {
	Env  string `mapstructure:"env" default:"dev"`
	HTTP struct {
		IP   string `mapstructure:"ip" default:"0.0.0.0"`
		Port string `mapstructure:"port" default:"8080"`
	} `mapstructure:"http"`
}
