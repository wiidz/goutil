package configMng

//easyjson:json
type AppConfig struct {
	Debug    bool   `json:"debug"`
	HttpPort string `json:"http_port"`
}

//easyjson:json
type MysqlConfig struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Collation string `json:"collation"`
	DbName    string `json:"db_name"`
	Charset   string `json:"charset"`
}

//easyjson:json
type EsConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//easyjson:json
type OssConfig struct {
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	EndPoint        string `json:"end_point"`
	BucketName      string `json:"bucket_name"`
	UploadPath      string `json:"upload_path"`
	Host            string `json:"host"`
	CallBackUrl     string `json:"call_back_url"` //为上传回调服务器的URL，请将下面的IP和Port配置为您自己的真实信息。
	ExpireTime      int64  `json:"expire_time"`   // 上传策略Policy的失效时间，单位为秒。
}

//easyjson:json
type RedisConfig struct {
	IP          string `json:"ip"`
	Port        string `json:"port"`
	Password    string `json:"password"`
	IdleTimeout int    `json:"idle_timeout"`
	Database    int    `json:"datebase"`
	MaxActive   int    `json:"max_active"`
	MaxIdle     int    `json:"max_idle"`
}

//easyjson:json
type WechatConfig struct {
	AppID       string `json:"app_id"`
	AppSecret   string `json:"app_secret"`
	GrantType   string `json:"grant_type"`
	PayKey      string `json:"pay_key"`
	MechID      string `json:"mech_id"`
	NotifyUrl   string `json:"notify_url"`
	RefundUrl   string `json:"refund_url"`
	CertPath    string `json:"cert_path"`
	CertKeyPath string `json:"cert_key_path"`
}

// RabbitMQConfig rabbit mq配置
type RabbitMQConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}