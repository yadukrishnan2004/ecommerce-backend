package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secret []byte
}

func NewJwtService(s string) *JwtService {
	return &JwtService{
		secret: []byte(s),
	}
}

func (j *JwtService) GenerateToken(u string, ttl time.Duration)(string,error){
	claims:=jwt.MapClaims{
		"sub":u,
		"exp":time.Now().Add(ttl).Unix(),
		"iat":time.Now().Unix(),
	}
	t:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return t.SignedString(j.secret)
}