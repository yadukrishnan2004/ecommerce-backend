package repository

import (
	"context"
	"time"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

// User represents the database model
type User struct {
	gorm.Model
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"-"`
	Role      string `json:"role" gorm:"default:'user'"`
	Otp       string `json:"otp"`
	IsActive  bool	 `json:"is_active"`
	IsBlocked bool   `json:"is_block"`
	OtpExpire int64  
}



// ToDomain converts database model to domain entity
func (u *User) ToDomain() *domain.User {
	return &domain.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: func() *time.Time {
			if u.DeletedAt.Valid {
				return &u.DeletedAt.Time
			}
			return nil
		}(),
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		Otp:       u.Otp,
		IsActive:  u.IsActive,
		IsBlocked: u.IsBlocked,
		OtpExpire: u.OtpExpire,
	}
}

// FromDomain converts domain entity to database model
func FromDomain(u *domain.User) *User {
	var deletedAt gorm.DeletedAt
	if u.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *u.DeletedAt, Valid: true}
	}
	return &User{
		Model: gorm.Model{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			DeletedAt: deletedAt,
		},
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		Otp:       u.Otp,
		IsActive:  u.IsActive,
		IsBlocked: u.IsBlocked,
		OtpExpire: u.OtpExpire,
	}
}

func NewUserRepo(db *gorm.DB) domain.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	dbUser := FromDomain(user)
	if err := r.db.WithContext(ctx).Create(dbUser).Error; err != nil {
		return err
	}
	// Update ID back to domain user
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *userRepo) Update(ctx context.Context, user *domain.User) error {
	dbUser := FromDomain(user)
	// We need to save all fields that might have changed
	return r.db.WithContext(ctx).Save(dbUser).Error
}

func (r *userRepo) GetByID(ctx context.Context, userID uint) (*domain.User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *userRepo) Delete(ctx context.Context, userID uint) (error) {
	return r.db.WithContext(ctx).Delete(&domain.User{}, userID).Error
}



func (r *userRepo) SearchUsers(ctx context.Context, query string) ([]domain.User, error) {
    var users []domain.User
    
    searchPattern := "%" + query + "%"

    err := r.db.WithContext(ctx).
        Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern).
        Limit(20).
        Find(&users).Error

    return users, err
}

