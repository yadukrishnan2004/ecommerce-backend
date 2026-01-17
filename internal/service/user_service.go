package service

import (
	"context"
	"errors"

	"github.com/yadukrishnan2004/ecommerce-backend/helper"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type userService struct {
	repo domain.UserRepositery
}

func NewUserService(repo domain.UserRepositery) domain.UserService{
	return &userService{repo:repo}
}


func (s *userService) Register(ctx context.Context,name,email,password string) error{
	existing,_:=s.repo.GetByEmail(ctx,email)
	if existing != nil && existing.ID != 0 {
        return errors.New("user already exists")
    }

	pass,err:=helper.Hash(password)
	if err != nil {
		return err
	}

	user:=domain.User{
		Name: name,
		Email: email,
		Password: pass,
	}
	return s.repo.Create(ctx,&user)
}