package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/yadukrishnan2004/ecommerce-backend/helper"

	auth "github.com/yadukrishnan2004/ecommerce-backend/internal/Auth"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type UpdateUserInput struct {
	Name *string
}

type UserProfileOutput struct {
	Name  string
	Email string
	Role  string
}

type UserUseCase interface {
	SignUp(ctx context.Context, name, email, password string) (string, error)
	VerifyOtp(ctx context.Context, email, code string) error
	SignIn(ctx context.Context, email, password string) (string, error)
	ForgotPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, email, code, newPassword string) error
	UpdateProfile(ctx context.Context, userID uint, input UpdateUserInput) error
	GetProfile(ctx context.Context, userID uint) (*UserProfileOutput, error)
	GetOrderDetail(ctx context.Context, orderID, userID uint) (*domain.Order, error)
	CancelOrder(ctx context.Context, orderID, userID uint) error
}

type userUseCase struct {
	repo domain.UserRepository
	otp  domain.NotificationClient
	jwt  auth.JwtService
	orders domain.OrderRepository
}

func NewUserUseCase(repo domain.UserRepository, otp domain.NotificationClient, jwt auth.JwtService,oreders domain.OrderRepository) UserUseCase {
	return &userUseCase{
		repo: repo,
		otp:  otp,
		jwt:  jwt,
		orders: oreders,
	}
}

func (s *userUseCase) SignUp(ctx context.Context, name, email, password string) (string, error) {
	
	user, err := s.repo.GetByEmail(ctx, email)

	if err == nil && user.ID != 0 {

	
		if user.IsActive {
			return "", errors.New("user already exists")
		}
		hashedPass, err := helper.Hash(password)
		if err != nil {
			return "", err
		}

		otp := helper.GenerateOtp()

		user.Name = name 
		user.Password = hashedPass
		user.Otp = otp
		user.OtpExpire = time.Now().Add(10 * time.Minute).Unix()
		user.Role = "user"

	
		if err := s.repo.Update(ctx, user); err != nil {
			return "", err
		}

		token, erro := s.jwt.GenerateAuthToken(user.Role, user.Email, 10*60)
		if erro != nil {
			return "", errors.New("forgot pass is not generated")
		}
	
		s.otp.SendOtp(user.Email, user.Otp)

		return token, nil
	}



	hashedPass, err := helper.Hash(password)
	if err != nil {
		return "", err
	}

	otp := helper.GenerateOtp()

	newUser := &domain.User{
		Name:      name,
		Email:     email,
		Password:  hashedPass,
		IsActive:  false,
		Otp:       otp,
		OtpExpire: time.Now().Add(10 * time.Minute).Unix(),
	}

	
	if err := s.repo.Create(ctx, newUser); err != nil {
		return "", err
	}

	token, erro := s.jwt.GenerateToken(newUser.ID, 10*60, newUser.Role)
	if erro != nil {
		return "", errors.New("forgot pass is not generated")
	}
	s.otp.SendOtp(newUser.Email, newUser.Otp)

	return token, nil
}

func (s *userUseCase) VerifyOtp(ctx context.Context, email, code string) error {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}


	if user.IsActive {
		return errors.New("user already active")
	}
	if user.Otp != code {
		return errors.New("invalid code")
	}
	if time.Now().Unix() > user.OtpExpire {
		return errors.New("code expired")
	}


	user.IsActive = true
	user.Otp = "" 
	return s.repo.Update(ctx, user)
}

func (s *userUseCase) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}


	if !user.IsActive {
		return "", errors.New("account not verified")
	}


	if err := helper.VerifyHash(user.Password, password); !err {
		return "", errors.New("invalid email or password")
	}


	acc, erro := s.jwt.GenerateToken(user.ID, s.jwt.AccessTTL, user.Role)
	if erro != nil {
		return "", erro
	}
	return acc, nil

}

func (s *userUseCase) ForgotPassword(ctx context.Context, email string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	otp := helper.GenerateOtp()
	s.otp.SendOtp(user.Email, otp)
	user.Otp = otp
	user.OtpExpire = time.Now().Add(10 * time.Minute).Unix()
	if err := s.repo.Update(ctx, user); err != nil {
		return "", errors.New("something went wrong please try again later")
	}
	token, erro := s.jwt.GenerateAuthToken(user.Role, user.Email, 10*60)
	if erro != nil {
		return "", errors.New("forgot pass is not generated")
	}
	return token, nil
}

func (s *userUseCase) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}
	if user.OtpExpire < time.Now().Unix() {
		return errors.New("time Expired")
	}
	if user.Otp != code || code == "" {
		return errors.New("code not match")
	}
	if newPassword == "" {
		return errors.New("please enter an valid password")
	}
	hash, hasherr := helper.Hash(newPassword)
	if hasherr != nil {
		return errors.New("something went wrong")
	}
	user.Password = hash
	user.Otp = ""
	if erro := s.repo.Update(ctx, user); erro != nil {
		return errors.New("sorry password not updated please try again")
	}
	return nil
}

func (s *userUseCase) UpdateProfile(ctx context.Context, userID uint, input UpdateUserInput) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	if input.Name != nil {
		user.Name = *input.Name
	}
	return s.repo.Update(ctx, user)
}

func (s *userUseCase) GetProfile(ctx context.Context, userID uint) (*UserProfileOutput, error) {
	users, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user := UserProfileOutput{
		Name:  users.Name,
		Email: users.Email,
		Role:  users.Role,
	}
	return &user, nil
}

func (s *userUseCase) GetOrderDetail(ctx context.Context, orderID, userID uint) (*domain.Order, error) {
    return s.orders.GetByIDAndUser(ctx, orderID, userID)
}


func (s *userUseCase) CancelOrder(ctx context.Context, orderID, userID uint) error {

    order, err := s.orders.GetByIDAndUser(ctx, orderID, userID)
    if err != nil {
        return errors.New("order not found")
    }

    if order.Status != "Pending" {
        return errors.New("only pending orders can be cancelled")
    }

    return s.orders.CancelOrder(ctx, orderID, userID)
}



