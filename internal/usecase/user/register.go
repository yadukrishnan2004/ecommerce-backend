package user

import (
	"errors"

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

	u , err:=user.New(email,password)
	if err != nil {
		return err
	}

	return  uc.repo.Create(u)
}