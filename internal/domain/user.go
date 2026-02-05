package domain

import (
	"context"
	"time"
)

type User struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Role      string     `json:"role" gorm:"default:'user'"`
	Otp       string     `json:"-"`
	IsActive  bool       `json:"is_active"`
	IsBlocked bool       `json:"is_blocked"`
	OtpExpire int64      `json:"otp_expire"`
}


type NotificationClient interface {
	SendOtp(toEmail string, code string) error
}


type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, userID uint) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context,userId uint) error
	SearchUsers(ctx context.Context, query string) ([]User, error)
}







