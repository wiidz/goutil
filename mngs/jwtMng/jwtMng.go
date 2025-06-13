package jwtMng

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12"
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/redisMng"
	"golang.org/x/xerrors"
	"reflect"
	"strings"
	"time"
)

const TokenKeyName = "token_data"

type JwtMng struct {
	AppID uint64 `json:"app_id"` // app_id 主要用来区别登陆
	//RouterKey       int8       `json:"router_keys"`    // 表示进入哪个端（一般用于前后端api合体在一个进程里的情况，自行定义，例如0=client,1=console）
	IdentifyKey     string     `json:"identify_key"`  // 身份标识键名，这个key必须存在于tokenStruct里
	IdentifyKeys    []string   `json:"identify_keys"` // 身份标识键名，这个key必须存在于tokenStruct里，和router_keys对应排列
	TokenStruct     jwt.Claims `json:"token_struct"`
	SaltKey         []byte     `json:"salt_key"`          // 盐值
	IsSingletonMode bool       `json:"is_singleton_mode"` // 是否单例登陆模式
}

// GetJwtMng 获取jwt管理器（单体）
func GetJwtMng(appID uint64, isSingletonMode bool, identifyKey, saltKey string, tokenStruct jwt.Claims) *JwtMng {
	return &JwtMng{
		AppID:           appID,
		IsSingletonMode: isSingletonMode,
		IdentifyKey:     identifyKey, // 例如 user_id、staff_id等
		SaltKey:         []byte(saltKey),
		TokenStruct:     tokenStruct,
	}
}

// GetJwtMngMixed 获取jwt管理器（混合体，例如client和console的api共用一个进程的项目）
func GetJwtMngMixed(appID uint64, isSingletonMode bool, identifyKeys []string, saltKey string, tokenStruct jwt.Claims) *JwtMng {
	return &JwtMng{
		AppID:           appID,
		IsSingletonMode: isSingletonMode,
		//RouterKey:       0,
		IdentifyKey:  "", // 例如 user_id、staff_id等
		IdentifyKeys: identifyKeys,
		SaltKey:      []byte(saltKey),
		TokenStruct:  tokenStruct,
	}
}

// GetTokenStr  获取jwt token
func (mng *JwtMng) GetTokenStr(claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mng.SaltKey)
	return ss, err
}

// Decrypt 解码
func (mng *JwtMng) Decrypt(claims jwt.Claims, tokenStr string) error {

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return mng.SaltKey, nil
	})

	if err != nil {
		return err
	}

	if token != nil && token.Valid {
		return nil
	}
	return err
}

// DecryptWithoutValidation 解码但不验证时间
func (mng *JwtMng) DecryptWithoutValidation(claims jwt.Claims, tokenStr string) error {

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return mng.SaltKey, nil
	}, jwt.WithoutClaimsValidation())

	//if err != nil {
	//	return err
	//}

	if !token.Valid {
		return errors.New("token验证失败")
	}

	return err
}

// Serve 注入服务
func (mng *JwtMng) Serve(ctx iris.Context) {

	//【1】从头部获取jwt
	tokenStr, err := mng.FromAuthHeader(ctx.GetHeader("Authorization"))
	if err != nil {
		// 全都视作未登录
		networkHelper.ReturnResult(ctx, err.Error(), nil, 401)
		return
	}

	//【2】尝试解密
	if err = mng.Decrypt(mng.TokenStruct, tokenStr); err != nil {
		networkHelper.ReturnError(ctx, err.Error())
		return
	}

	//【】取出ID
	immutable := reflect.ValueOf(mng.TokenStruct)
	id := immutable.Elem().FieldByName(mng.IdentifyKey).Interface().(uint64)
	if id == 0 {
		networkHelper.ReturnError(ctx, "登陆主体为空")
		return
	}

	//【3】判断和缓存中的数据
	if mng.IsSingletonMode {

		err = mng.CompareJwtCache(ctx, mng.AppID, typeHelper.Uint64ToStr(id), tokenStr)
		if err != nil {
			networkHelper.ReturnError(ctx, err.Error())
			return
		}
	}

	//【4】写入value
	ctx.Values().Set(TokenKeyName, mng.TokenStruct)

	//【5】继续下一步处理
	ctx.Next()
}

// ServeMixed 注入服务（混合体）
// SetRouterKey 也在这一步，SetRouterFlag请在networkHelper中另外调用
// 注意这里的key是比数组下标大1的数值
func (mng *JwtMng) ServeMixed(ctx iris.Context) {

	//【1】从头部获取jwt
	tokenStr, err := mng.FromAuthHeader(ctx.GetHeader("Authorization"))
	if err != nil {
		// 全都视作未登录
		networkHelper.ReturnResult(ctx, err.Error(), nil, 401)
		return
	}

	//【2】尝试解密
	if err = mng.Decrypt(mng.TokenStruct, tokenStr); err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			networkHelper.ReturnResult(ctx, "登陆信息已过期", nil, 401)
			return
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			networkHelper.ReturnResult(ctx, "身份令牌格式错误", nil, 401)
			return
		}
		networkHelper.ReturnResult(ctx, err.Error(), nil, 401)
		return
	}

	//【3】取出ID
	immutable := reflect.ValueOf(mng.TokenStruct)

	// 注意有些项目是混合在一起，有多个router_keys的情况
	// 例如客户端包含了推介员功能
	// 那么userID和promoterID都不为空
	// 所以这里的router_keys采用数组的形式
	routerKeySlice := []int{}
	for k, v := range mng.IdentifyKeys {
		id := immutable.Elem().FieldByName(v).Interface().(uint64)
		if id != 0 {
			//mng.RouterKey = int8(k) // 这一步没有意义
			routerKeySlice = append(routerKeySlice, k+1)
		}
	}
	ctx.Values().Set("router_keys", routerKeySlice)

	if len(routerKeySlice) == 0 {
		networkHelper.ReturnError(ctx, "登陆主体为空")
		return
	}

	//【3】判断和缓存中的数据
	if mng.IsSingletonMode {
		cacheKeyName := typeHelper.ImplodeInt(routerKeySlice, "-")
		err = mng.CompareJwtCache(ctx, mng.AppID, cacheKeyName, tokenStr)
		if err != nil {
			networkHelper.ReturnError(ctx, err.Error())
			return
		}
	}

	//【4】写入value
	ctx.Values().Set(TokenKeyName, mng.TokenStruct)

	//【5】继续下一步处理
	ctx.Next()
}

// FromAuthHeader 从header头中获取jwt
func (mng *JwtMng) FromAuthHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("您尚未登陆") // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		//return "", fmt.Errorf("Authorization header format must be Bearer {token}")
		return "", fmt.Errorf("密钥格式错误")
	}

	return authHeaderParts[1], nil
}

// RefreshToken 刷新token
func (mng *JwtMng) RefreshToken(ctx iris.Context, validDuration float64) {

	tokenStr, err := mng.FromAuthHeader(ctx.GetHeader("Authorization"))
	if err != nil {
		networkHelper.ReturnError(ctx, err.Error())
		return
	}

	err = mng.Decrypt(mng.TokenStruct, tokenStr)

	// 判断错误过期
	//var expErr *jwt.ErrTokenExpired
	if xerrors.As(err, jwt.ErrTokenExpired) || err == nil {

		//【】取出过期时间
		immutable := reflect.ValueOf(err)
		expiredBy := immutable.Elem().FieldByName("ExpiredBy")

		if expiredBy.Interface().(time.Duration) > (time.Duration(validDuration) * time.Second) {
			networkHelper.ReturnError(ctx, "已超出预定时长")
			return
		}

		newToken, err := mng.GetTokenStr(mng.TokenStruct)
		if err != nil {
			networkHelper.ReturnError(ctx, err.Error())
			return
		}

		networkHelper.ReturnResult(ctx, "success", newToken, 200)
		return
	}

	networkHelper.ReturnError(ctx, err.Error())
	return
}

// SetCache StorageJWT 存储kwt至redis中
func (mng *JwtMng) SetCache(ctx context.Context, appID uint64, routerKeyName, token string) (bool, error) {
	redis := redisMng.NewRedisMng()
	res, err := redis.HSetNX(ctx, typeHelper.Uint64ToStr(appID)+"-jwt", routerKeyName, token)
	return res, err
}

// GetCache GetJwtCache 从缓存中读取jwt
func (mng *JwtMng) GetCache(ctx context.Context, appID uint64, routerKeyName string) (string, error) {
	redis := redisMng.NewRedisMng()
	res, err := redis.HGetString(ctx, typeHelper.Uint64ToStr(appID)+"-jwt", routerKeyName)
	return res, err
}

// DeleteCache 从缓存中删除jwt
func (mng *JwtMng) DeleteCache(ctx context.Context, appID, userID uint64) (int64, error) {
	redis := redisMng.NewRedisMng()
	res, err := redis.HDel(ctx, typeHelper.Uint64ToStr(appID)+"-jwt", []string{typeHelper.Uint64ToStr(userID)})
	return res, err
}

// CompareJwtCache 判断jwtToken
func (mng *JwtMng) CompareJwtCache(ctx context.Context, appID uint64, routerKeyName, token string) error {
	//【1】从缓存中读取jwt
	cacheToken, err := mng.GetCache(ctx, appID, routerKeyName)
	if err != nil {
		return err
	}

	//【2】判断是否相等
	if token != cacheToken {
		return errors.New("该账号正在使用中，请重新登入")
	}

	return nil
}

// IsPkSet 主要用来判断是否是前端请求
// Tips：此方法试用于前后端非同表的项目，判断是否是前端（tokenData里是否有jwtMng约定的主键）
func (mng *JwtMng) IsPkSet(tokenData jwt.Claims) bool {
	immutable := reflect.ValueOf(tokenData)
	if immutable.IsValid() == false {
		return false
	}

	temp := immutable.Elem().FieldByName(mng.IdentifyKey)

	if temp.IsValid() == false {
		return false
	}

	id := temp.Interface().(uint64)
	if id == uint64(0) {
		return false
	} else {
		return true
	}
}

// GetTokenData 获取token
func (mng *JwtMng) GetTokenData(ctx iris.Context) (data jwt.Claims, err error) {
	tempData := ctx.Values().Get(TokenKeyName)
	if tempData == nil {
		err = errors.New("token数据为空")
		return
	}

	var ok bool
	data, ok = tempData.(jwt.Claims)
	if !ok {
		err = errors.New("token解析失败")
		return
	}

	if mng.IsPkSet(data) == false {
		err = errors.New("登陆主体为空")
	}

	return
}
