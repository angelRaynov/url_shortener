package model

import "github.com/golang-jwt/jwt"

type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UID      string `json:"uid"`
	jwt.StandardClaims
}
