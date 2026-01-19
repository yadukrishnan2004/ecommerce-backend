package service

import (
	"context"
	"errors"
	"time"

	"github.com/yadukrishnan2004/ecommerce-backend/helper"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type userService struct {
	repo domain.UserRepositery
	otp  domain.NotificationClint
}

func NewUserService(repo domain.UserRepositery,otp domain.NotificationClint) domain.UserService{
	return &userService{
		repo:repo,
		otp :otp,
	}
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

	otp:=helper.GenerateOtp()

	user:=&domain.User{
		Name: 	  name,
		Email:	  email,
		Password: pass,
		IsActive: false,
		Otp: otp,
		OtpExpire: time.Now().Add(10*time.Minute).Unix(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
        return err
    }
	return s.otp.SendOtp(user.Email,user.Otp)
}

func (s *userService) VerifyOtp(ctx context.Context,email,code string)error{
	user, err := s.repo.GetByEmail(ctx, email)
    if err != nil {
        return errors.New("user not found")
    }

    //  Validate Logic
    if user.IsActive {
        return errors.New("user already active")
    }
    if user.Otp != code {
        return errors.New("invalid code")
    }
    if time.Now().Unix() > user.OtpExpire {
        return errors.New("code expired")
    }

    //  Activate User
    user.IsActive = true
    user.Otp = "" // Clear the code
    return s.repo.Update(ctx, user)
}