package repositery

import (
	"context"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)


type userRepo struct {
	db *gorm.DB
}

// all the db logic 

func NewUserRepo(db *gorm.DB) domain.UserRepositery{
	return &userRepo{db:db}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    return &user, err
}