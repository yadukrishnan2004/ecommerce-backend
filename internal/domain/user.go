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

	Addresses []Address  `json:"addresses" gorm:"foreignKey:UserID"`
	Cart      []CartItem `json:"cart" gorm:"foreignKey:UserID"`
	Wishlist  []Wishlist `json:"wishlist" gorm:"foreignKey:UserID"`
	Orders    []Order    `json:"orders" gorm:"foreignKey:UserID"`
}

type NotificationClient interface {
	SendOtp(toEmail string, code string) error
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, userID uint) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userId uint) error
	SearchUsers(ctx context.Context, query string) ([]User, error)
	SaveSignup(ctx context.Context, signup *SignupRequest) error
	GetSignup(ctx context.Context, email string) (*SignupRequest, error)
	DeleteSignup(ctx context.Context, email string) error
	GetAllUsers(ctx context.Context) ([]User, error)
}

type SignupRequest struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	Otp       string `json:"otp"`
	OtpExpire int64  `json:"otp_expire"`
}
