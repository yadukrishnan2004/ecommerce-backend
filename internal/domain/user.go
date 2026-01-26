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

type UserProfile struct {
	UserID  	uint  
    Name    	string `json:"name"`
	Email   	string `json:"email"`
	Role 		string `json:"role"`
	IsBlocked   bool

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


type UserService interface{
	Register(ctx context.Context,name,email,password string)(string,error)
	VerifyOtp(ctx context.Context,email,code string)error
	Login(ctx context.Context,email,password string)(string,error)
	Forgetpassword(ctx context.Context,email string)(string,error)
	Resetpassword(ctx context.Context,email,code,newpassword string)error
	UpdateProfile(ctx context.Context, userID uint, input UserProfile) error
	GetProfile(ctx context.Context, userID uint) (*UserProfile, error)
}


type AdminService interface{
	 
}