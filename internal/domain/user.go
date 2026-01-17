package domain

import (
	"context"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name 	 string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type UserRepositery interface {
	Create(ctx context.Context,user *User) error
	GetByEmail(ctx context.Context,email string) (*User,error)
}



type UserService interface{
	Register(ctx context.Context,name,email,password string)error

}