package jwtMng

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/kataras/iris/v12"
	"github.com/wiidz/goutil/helpers"
	"github.com/wiidz/goutil/mngs/redisMng"
	"golang.org/x/xerrors"
	"reflect"
	"strings"
	"time"
)

var typeHelper = helpers.TypeHelper{}

type JwtMng struct {
	AppID           int        `json:"app_id"`       // app_id 主要用来区别登陆
	IdentifyKey     string     `json:"identify_key"` // 身份标识键名，这个key必须存在于tokenStruct里
	TokenStruct     jwt.Claims `json:"token_struct"`
	SaltKey         []byte     `json:"salt_key"`          // 盐值
	IsSingletonMode bool       `json:"is_singleton_mode"` // 是否单例登陆模式
}

// GetJwtMng 获取jwt管理器
func GetJwtMng(appID int, isSingletonMode bool, identifyKey, saltKey string, tokenStruct jwt.Claims) *JwtMng {
	return &JwtMng{
		AppID:           appID,
		IsSingletonMode: isSingletonMode,
		IdentifyKey:     identifyKey, // 例如 user_id、staff_id等
		SaltKey:         []byte(saltKey),
		TokenStruct:     tokenStruct,
	}
}

// GetTokenStr： 获取jwt token
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
		helpers.ReturnError(ctx, err.Error())
		return
	}

	//【2】尝试解密
	if err := mng.Decrypt(mng.TokenStruct, tokenStr); err != nil {
		helpers.ReturnError(ctx, err.Error())
		return
	}

	//【3】判断和缓存中的数据
	if mng.IsSingletonMode {
		//【】取出ID
		immutable := reflect.ValueOf(mng.TokenStruct)
		id := immutable.Elem().FieldByName(mng.IdentifyKey).Interface().(int)
		err = mng.CompareJwtCache(mng.AppID, id, tokenStr)
		if err != nil {
			helpers.ReturnError(ctx, err.Error())
			return
		}
	}

	//【4】写入value
	ctx.Values().Set("token_data", mng.TokenStruct)

	//【5】继续下一步处理
	ctx.Next()
}

// FromAuthHeader 从header头中获取jwt
func (mng *JwtMng) FromAuthHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("Authorization header is empty") // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

// RefreshToken 刷新token
func (mng *JwtMng) RefreshToken(ctx iris.Context, validDuration float64) {

	tokenStr, err := mng.FromAuthHeader(ctx.GetHeader("Authorization"))
	if err != nil {
		helpers.ReturnError(ctx, err.Error())
		return
	}

	err = mng.Decrypt(mng.TokenStruct, tokenStr)

	// 判断错误过期
	var expErr *jwt.TokenExpiredError
	if xerrors.As(err, &expErr) || err == nil {

		//【】取出过期时间
		immutable := reflect.ValueOf(err)
		expiredBy := immutable.Elem().FieldByName("ExpiredBy")

		if expiredBy.Interface().(time.Duration) > (time.Duration(validDuration) * time.Second) {
			helpers.ReturnError(ctx, "已超出预定时长")
			return
		}

		newToken, err := mng.GetTokenStr(mng.TokenStruct)
		if err != nil {
			helpers.ReturnError(ctx, err.Error())
			return
		}

		helpers.ReturnResult(ctx, "success", newToken, 200)
		return
	}

	helpers.ReturnError(ctx, err.Error())
	return

}

// StorageJWT 存储kwt至redis中
func (mng *JwtMng) SetCache(appID, userID int, token string) (int64, error) {
	redis := redisMng.NewRedisMng()
	res, err := redis.HSet(typeHelper.Int2Str(appID)+"-jwt", typeHelper.Int2Str(userID), token)

	return res.(int64), err
}

// GetJwtCache 从缓存中读取jwt
func (mng *JwtMng) GetCache(appID, userID int) (string, error) {
	redis := redisMng.NewRedisMng()
	res, err := redis.HGet(typeHelper.Int2Str(appID)+"-jwt", typeHelper.Int2Str(userID))
	return res.(string), err
}

// DeleteCache 从缓存中删除jwt
func (mng *JwtMng) DeleteCache(appID, userID int) (string, error) {
	redis := redisMng.NewRedisMng()
	res, err := redis.HDel(typeHelper.Int2Str(appID)+"-jwt", typeHelper.Int2Str(userID))
	return res.(string), err
}

// CompareJwtCache 判断jwtToken
func (mng *JwtMng) CompareJwtCache(appID, userID int, token string) error {
	//【1】从缓存中读取jwt
	cacheToken, err := mng.GetCache(appID, userID)
	if err != nil {
		return err
	}

	//【2】判断是否相等
	if token != cacheToken {
		return errors.New("该账号正在使用中，请重新登入")
	}

	return nil
}
