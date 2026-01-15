package user

import (
	"errors"

	"github.com/yadukrishnan2004/ecommerce-backend/helper"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain/user"
)

type LoginUserCase struct {
	repo user.Repository
}

func NewLoginUseCase(repo user.Repository) *LoginUserCase{
	return &LoginUserCase{repo: repo}
}

func (uc *LoginUserCase) Execute(email,password string)error{
	u,err:=uc.repo.FindByEmail(email)
	if err !=nil || u == nil {
		return errors.New("invalid email or password")
	}

	if err:=helper.VerifyHash(u.Password,password); !err{
		return errors.New("invalid email or password")
	}
	return nil
}