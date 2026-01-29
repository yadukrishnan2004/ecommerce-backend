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
	Role 	  string `json:"role"`
	Otp		  string `json:"-"`
	IsActive  bool   `json:"is_active"`
	IsBlocked bool   `json:"is_blocked"` 
	OtpExpire int64  `json:"otp_expire"` 
}


type NotificationClint interface{
	SendOtp(toEmail string,code string) error
}


type UserRepositery interface {
	Create(ctx context.Context,user *User) error
	GetByEmail(ctx context.Context,email string) (*User,error)
	GetByID(ctx context.Context,userID uint) (*User,error)
	Update(ctx context.Context,user *User)error
}





