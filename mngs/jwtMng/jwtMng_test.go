package jwtMng

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"testing"
	"time"
)

type MixedTokenData struct {
	PlayerID       uint64 `json:"player_id"`        // v1，角色ID
	UserID         uint64 `json:"user_id"`          // v2，后台用户ID
	ConsoleStaffID uint64 `json:"console_staff_id"` // v4
	Grouping       int8   `json:"grouping"`         // 账号分组

	IsNormalVip bool `json:"is_normal_vip"` // 是否是认证会员
	IsPaidVip   bool `json:"is_paid_vip"`   // 是否是付费会员
	IsRatingVip bool `json:"is_rating_vip"` // 是否是评级会员

	jwt.RegisteredClaims
}

func Test(t *testing.T) {

	//tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6MSwidXNlcl9pZCI6MSwibWluaV91c2VyX2lkIjowLCJjb25zb2xlX3N0YWZmX2lkIjowLCJncm91cGluZyI6MCwiZXhwIjowLCJpc19ub3JtYWxfdmlwIjpmYWxzZSwiaXNfcGFpZF92aXAiOmZhbHNlLCJpc19yYXRpbmdfdmlwIjpmYWxzZSwiaXNzIjoidGVzdCJ9.NNWwjD5EsWHh5ZnWMX_o5XcGjrdDs2btLi4Rb2Uqpk8"
	//mng := GetJwtMng(20, false, "user_id", "hujiayilu123456", &MixedTokenData{})

	//var temp = MixedTokenData{}
	//err := mng.Decrypt(temp, tokenStr)
	//temp, err := jwt.ParseWithClaims(tokenStr, MixedTokenData{}, func(token *jwt.Token) (interface{}, error) {
	//	return "hujiayilu123456", nil
	//})
	//
	//log.Println("123", err)
	//log.Println("temp", temp)
	// Create the Claims
	//claims := &jwt.RegisteredClaims{
	//	ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
	//	Issuer:    "test",
	//}

	tokenStr, err := GetToken()
	log.Println("tokenStr", tokenStr)
	log.Println("err", err)

	token, err := jwt.ParseWithClaims(tokenStr, &MixedTokenData{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		log.Fatal(err)
	} else if claims, ok := token.Claims.(*MixedTokenData); ok {
		fmt.Println(claims, claims.RegisteredClaims.Issuer)
	} else {
		log.Fatal("unknown claims type, cannot proceed")
	}

}

func GetToken() (tokenStr string, err error) {

	claims := &MixedTokenData{
		PlayerID:       1,
		UserID:         2,
		ConsoleStaffID: 3,
		Grouping:       4,
		IsNormalVip:    false,
		IsPaidVip:      true,
		IsRatingVip:    false,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "",
			Subject:  "",
			Audience: nil,
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(-time.Second * time.Duration(7200)),
			},
			NotBefore: nil,
			IssuedAt:  nil,
			ID:        "",
		},
	}

	mySigningKey := []byte("AllYourBase")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString(mySigningKey)

	return
}
