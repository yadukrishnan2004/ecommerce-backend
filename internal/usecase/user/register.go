package user

import (
	"errors"

	"github.com/yadukrishnan2004/ecommerce-backend/helper"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain/user"
)

type RegisterUserCase struct {
	repo user.Repository
}

func NewRegisterUseCase(repo user.Repository) *RegisterUserCase{
	return &RegisterUserCase{repo:repo}
}

func (uc *RegisterUserCase) Execute(email,password string) error{
	existing,_:=uc.repo.FindByEmail(email)
	if existing !=nil {
		return  errors.New("email already registered")
	}

	hash,errr:=helper.Hash(password)
	if errr != nil {
		return errors.New("please try another password")
	}

	u , err:=user.New(email,hash)
	if err != nil {
		return err
	}

	return  uc.repo.Create(u)
}