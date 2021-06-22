package test

import (
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/wiidz/goutil/mngs/jwtMng"
	"log"
	"time"
)

var mng = jwtMng.GetJwtMng("hujiayilu123")

func GetToken() {

	expiringSeconds := 3600

	time := jwt.Time{time.Now().Add(time.Second * time.Duration(expiringSeconds))}

	claims := jwtMng.TokenData{
		1,
		1,
		2,
		jwt.StandardClaims{
			ExpiresAt: &time,
			Issuer:    "test",
		},
	}

	tokenStr, err := mng.GetTokenStr(claims)
	log.Println("tokenStr", tokenStr)
	log.Println("err", err)
}

func ParsedToken() {

	var tokenData jwtMng.TokenData

	err := mng.Decrypt(&tokenData, "")
	log.Println("tokenData", tokenData)
	log.Println("err", err)
}
