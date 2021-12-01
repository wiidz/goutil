package configMng

// AppConfig app设置
//type AppConfig struct {
//	Debug    bool   `json:"debug"`
//	Name     string `json:"name"`
//	Version  string `json:"version"`
//	Host     string `json:"host"`
//	HttpPort string `json:"http_port"`
//}
//
//// MysqlConfig mysql数据库设置
//type MysqlConfig struct {
//	Host      string `json:"host"`
//	Port      string `json:"port"`
//	Username  string `json:"username"`
//	Password  string `json:"password"`
//	Collation string `json:"collation"`
//	DbName    string `json:"db_name"`
//	Charset   string `json:"charset"`
//}
//
//// EsConfig elastic search 设置
//type EsConfig struct {
//	Host     string `json:"host"`
//	Port     string `json:"port"`
//	Username string `json:"username"`
//	Password string `json:"password"`
//}
//
//// OssConfig oss阿里云对象存储设置
//type OssConfig struct {
//	AccessKeyID     string `json:"access_key_id"`
//	AccessKeySecret string `json:"access_key_secret"`
//	EndPoint        string `json:"end_point"`
//	BucketName      string `json:"bucket_name"`
//	UploadPath      string `json:"upload_path"`
//	Host            string `json:"host"`
//	CallBackUrl     string `json:"call_back_url"` //为上传回调服务器的URL，请将下面的IP和Port配置为您自己的真实信息。
//	ExpireTime      int64  `json:"expire_time"`   // 上传策略Policy的失效时间，单位为秒。
//}
//
//// RedisConfig redis服务器设置
//type RedisConfig struct {
//	IP          string `json:"ip"`
//	Port        string `json:"port"`
//	Password    string `json:"password"`
//	IdleTimeout int    `json:"idle_timeout"`
//	Database    int    `json:"datebase"`
//	MaxActive   int    `json:"max_active"`
//	MaxIdle     int    `json:"max_idle"`
//}
//
//// WechatMiniConfig 微信小程序设置
//type WechatMiniConfig struct {
//	AppID     string `json:"app_id"`
//	AppSecret string `json:"app_secret"`
//}
//
//// WechatOaConfig 微信公众号设置
//type WechatOaConfig struct {
//	AppID     string `json:"app_id"`
//	AppSecret string `json:"app_secret"`
//}
//
//// WechatOpenConfig 微信开放平台设置
//type WechatOpenConfig struct {
//	AppID     string `json:"app_id"`
//	AppSecret string `json:"app_secret"`
//}
//
//
//// WechatPayConfig 微信支付配置
//type WechatPayConfig struct {
//	AppID           string `json:"app_id"`            //【微信支付】appID
//	PayKey          string `json:"pay_key"`           //【微信支付】支付密钥
//	MchID           string `json:"mch_id"`            //【微信支付】商户号
//	NotifyURL       string `json:"notify_url"`        //【微信支付】付款回调地址
//	RefundNotifyURL string `json:"refund_notify_url"` //【微信支付】退款回调地址
//	CertPath        string `json:"cert_path"`         //【微信支付】证书路径(cert)
//	CertKeyPath     string `json:"cert_key_path"`     //【微信支付】证书路径(key)
//	CertFileContent string `json:"cert_file_content"` //【微信支付】证书文件(cert)内容
//	IsProd          bool   `json:"is_prod"`           // 是否是正式环境
//}
//
//// AliPayConfig 支付宝参数
//type AliPayConfig struct {
//	AppID      string `json:"wechat_pay_app_id"` //【支付宝】appID
//	PrivateKey string `json:"private_key"`       //【支付宝】应用私钥，支持PKCS1和PKCS8
//	NotifyURL  string `json:"notify_url"`        //【支付宝】回调地址
//	IsProd     bool   `json:"is_prod"`           //【支付宝】是否是正式环境
//}
//
//// RabbitMQConfig rabbit mq配置
//type RabbitMQConfig struct {
//	Host     string `json:"host"`
//	Username string `json:"username"`
//	Password string `json:"password"`
//}
