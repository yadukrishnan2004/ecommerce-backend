package user

import (
	"errors"
	"time"

	"github.com/yadukrishnan2004/ecommerce-backend/helper"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain/user"
)

type LoginUserCase struct {
	repo 		user.Repository
	tokens 		TokenGenerator
	accessTTL 	time.Duration
	refreshTTL 	time.Duration
}

func NewLoginUseCase(
	repo 		user.Repository,
	tokens 		TokenGenerator,
	accessTTL 	time.Duration,
	refreshTTL  time.Duration,
	) *LoginUserCase{
	return &LoginUserCase{
		repo: 		repo,
		tokens: 	tokens,
		accessTTL: 	accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (uc *LoginUserCase) Execute(email,password string)(string, string ,error){
	u,err:=uc.repo.FindByEmail(email)
	if err !=nil || u == nil {
		return "","",errors.New("invalid email or password")
	}
	if err:=helper.VerifyHash(u.Password,password); !err{
		return "","",errors.New("invalid email or password")
	}
	accesstoken,err:=uc.tokens.GenerateToken(u.ID,uc.accessTTL)
	if err != nil {
		return "","",err
	}
	refreshToken,err:=uc.tokens.GenerateToken(u.ID,uc.refreshTTL)
	if err != nil {
		return "","",err
	}
	return accesstoken,refreshToken,nil 
}

// ------------------------------------------------------------------

type TokenGenerator interface{
	GenerateToken(userID string,ttl time.Duration)(string,error)
}



