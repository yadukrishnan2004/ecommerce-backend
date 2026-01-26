package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/config"
)

type JwtService struct {
	secret []byte
	AccessTTL 	int32
	RefreshTTL  int32
}

func NewJwtService(s *config.JWTConfig) *JwtService {
	return &JwtService{
		secret: []byte(s.Secret),
		AccessTTL: s.AccessTTL,
		RefreshTTL: s.RefreshTTL,
	}
}

func (j *JwtService) GenerateToken(u uint, ttl int32,role string)(string,error){
	claims:=jwt.MapClaims{
		"sub":u,
		"role":role,
		"exp":time.Now().Add(time.Duration(ttl)*time.Second).Unix(),
		"iat":time.Now().Unix(),
	}
	t:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return t.SignedString(j.secret)
}

func (j *JwtService) GenerateAuthToken(r,u string, ttl int32)(string,error){
	claims:=jwt.MapClaims{
		"sub":u,
		"role":r,
		"exp":time.Now().Add(time.Duration(ttl)*time.Second).Unix(),
		"iat":time.Now().Unix(),
	}
	t:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return t.SignedString(j.secret)
}