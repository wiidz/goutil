package jwtMng

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/kataras/iris/v12"
	"github.com/wiidz/goutil/helpers"
	"strings"
)

type JwtMng struct {
	SaltKey []byte `json:"salt_key"` //盐值
}

func GetJwtMng(saltKey string) *JwtMng {
	return &JwtMng{
		SaltKey: []byte(saltKey),
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

	if err != nil{
		fmt.Println("Couldn't handle this token:", err)
		return err
	}

	if token.Valid {
		fmt.Println("You look nice today")
		return nil
	} else {
		fmt.Println("failed")
		return errors.New("failed")
	}
}

func (mng *JwtMng) Serve(ctx iris.Context) {

	tokenStr, err := mng.FromAuthHeader(ctx.GetHeader("Authorization"))
	if err != nil {
		helpers.ReturnError(ctx,err.Error())
		return
	}

	var tokenData TokenData

	if err := mng.Decrypt(&tokenData, tokenStr); err != nil {
		helpers.ReturnError(ctx,err.Error())
		return
	}
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
