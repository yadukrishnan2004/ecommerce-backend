package domain

import (
	"context"


	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name 	  string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Otp		  string `json:"-"`
	IsActive  bool `json:"is_active"` 
	OtpExpire int64 `json:"otp_expire"` 
}


type NotificationClint interface{
	SendOtp(toEmail string,code string) error
}

type UserRepositery interface {
	Create(ctx context.Context,user *User) error
	GetByEmail(ctx context.Context,email string) (*User,error)
	Update(ctx context.Context,user *User)error
}

type UserService interface{
	Register(ctx context.Context,name,email,password string)error
	VerifyOtp(ctx context.Context,email,code string)error
}