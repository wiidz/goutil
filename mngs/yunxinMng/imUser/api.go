package imUser

import (
	"github.com/wiidz/goutil/mngs/yunxinMng/imClient"
)

type Api struct {
	Client *imClient.Client
}

// 网易云信 IM 账号：以下文档中也称为“用户帐号”，参数名用 “accid” 或 “account” 等表示。
// token：网易云信 IM 账号的密码。创建 IM accid 时可以由开发者 app 的服务端指定。若未指定，则云信会自动生成一个 IM token，并返回给开发者。客户端登录时，需要传参 accid 与 token 给云信服务器鉴权。token 没有过期的概念，除非人为更改。只有最新的token才是唯一有效的。当登录时使用非最新的 token，将会返回的错误码 302。
// 客户端通过网易云信 SDK 连接登录云信服务器时，需要保证网易云信IM账号已经注册， 且确保客户端从自己的服务器已经取得了有效 token；

// https://doc.yunxin.163.com/docs/TM5MzM5Njk/zE2NzA3Mjc?platformId=60353
// #先获取当前时间戳，单位毫秒
// curTime = 1614764611561
// #设置过期时间，单位秒，如600
// ttl = 600
// #生成signature，将appkey、accid、curTime、ttl、appsecret五个字段拼成一个字符串，进行sha1编码
// signature = sha1(appkey + accid + curTime + ttl + appsecret)
// #组装成json
// json = {"signature": "xx", "curTime":1614764611561, "ttl": 600}
// #将json转成字符串后进行base64编码，生成最终的token
// token=base64(json)
