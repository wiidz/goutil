package jwtMng

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/kataras/iris/v12"
	"github.com/wiidz/goutil/helpers"
	"golang.org/x/xerrors"
	"log"
	"reflect"
	"strings"
	"time"
)

var structHelper = helpers.StructHelper{}
var typeHelper = helpers.TypeHelper{}

type JwtMng struct {
	TokenStruct jwt.Claims `json:"token_struct"`
	SaltKey     []byte     `json:"salt_key"` //盐值
}

func GetJwtMng(saltKey string, tokenStruct jwt.Claims) *JwtMng {
	return &JwtMng{
		SaltKey:     []byte(saltKey),
		TokenStruct: tokenStruct,
	}
}

/**
 * GetTokenStr ： 获取jwt token
 **/
func (mng *JwtMng) GetTokenStr(claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mng.SaltKey)
	return ss, err
}

/**
 * Decrypt ： 解码
 **/
func (mng *JwtMng) Decrypt(claims jwt.Claims, tokenStr string) error {

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return mng.SaltKey, nil
	})

	if token != nil && token.Valid {
		return nil
	}
	return err
}

/**
 * Decrypt ： 解码
 **/
func (mng *JwtMng) DecryptWithoutValidation(claims jwt.Claims, tokenStr string) error {

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return mng.SaltKey, nil
	}, jwt.WithoutClaimsValidation())

	if !token.Valid {
		return errors.New("token验证失败")
	}

	return err
}

func (mng *JwtMng) Serve(ctx iris.Context) {

	log.Println("serve")

	tokenStr, err := mng.FromAuthHeader(ctx.GetHeader("Authorization"))
	if err != nil {
		helpers.ReturnError(ctx, err.Error())
		return
	}
	log.Println("tokenStr",tokenStr)

	if err := mng.Decrypt(mng.TokenStruct, tokenStr); err != nil {
		helpers.ReturnError(ctx, err.Error())
		return
	}

	log.Println(" mng.TokenStruct", mng.TokenStruct)

	ctx.Values().Set("token_data", mng.TokenStruct)

	log.Println("saved")

	// If everything ok then call next.
	ctx.Next()
}

// FromAuthHeader is a "TokenExtractor" that takes a give context and extracts
// the JWT token from the Authorization header.
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

		//log.Println("aa", expiredBy.Interface().(time.Duration))
		//log.Println("bb", (time.Duration(validDuration) * time.Second))

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