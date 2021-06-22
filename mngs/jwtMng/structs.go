package jwtMng

import (
	"github.com/dgrijalva/jwt-go/v4"
)

type TokenData struct {
	UserID   int  `json:"user_id"`
	OgID     int  `json:"og_id"`
	Grouping int8 `json:"grouping"`
	jwt.StandardClaims
}
