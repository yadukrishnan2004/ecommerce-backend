package usecase

import (
	"context"
	"errors"

	// "fmt"
	"strings"
	"time"

	"github.com/yadukrishnan2004/ecommerce-backend/helper"

	auth "github.com/yadukrishnan2004/ecommerce-backend/internal/Auth"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type UpdateUserInput struct {
	Name *string
}

type UserProfileOutput struct {
	Name      string                    `json:"name"`
	Email     string                    `json:"email"`
	Role      string                    `json:"role"`
	CreatedAt time.Time                 `json:"created_at"`
	Address   []domain.Address          `json:"address"`
	Wishlist  []domain.WishlistItemView `json:"wishlist"`
	Cart      []domain.CartItemView     `json:"cart"`
	Orders    []domain.OrderItem        `json:"orders"`
}

type UserUseCase interface {
	SignUp(ctx context.Context, name, email, password string) (string, error)
	VerifyOtp(ctx context.Context, email, code string) (string, error)
	SignIn(ctx context.Context, email, password string) (*UserProfileOutput, string, error)
	ForgotPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, email, code, newPassword string) error
	UpdateProfile(ctx context.Context, userID uint, input UpdateUserInput) error
	GetProfile(ctx context.Context, userID uint) (*UserProfileOutput, error)
	GetProduct(ctx context.Context, id uint) (*domain.Product, error)
	GetOrderDetail(ctx context.Context, orderID, userID uint) ([]domain.Order, error)
	CancelOrder(ctx context.Context, orderID, userID uint) error
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
	SearchProducts(ctx context.Context, query string) ([]domain.Product, error)
	FilterProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error)
	GetOrderItemDetails(ctx context.Context, orderID uint) ([]domain.OrderItem, error)
}

type userUseCase struct {
	repo        domain.UserRepository
	otp         domain.NotificationClient
	jwt         auth.JwtService
	orders      domain.OrderRepository
	productrepo domain.ProductRepository
	cart        domain.CartRepository
	wishlist    domain.WishlistRepository
	address     domain.AddressRepository
}

func NewUserUseCase(repo domain.UserRepository,
	otp domain.NotificationClient,
	jwt auth.JwtService,
	oreders domain.OrderRepository,
	productrepo domain.ProductRepository,
	cart domain.CartRepository,
	wishlist domain.WishlistRepository,
	address domain.AddressRepository,
) UserUseCase {
	return &userUseCase{
		repo:        repo,
		otp:         otp,
		jwt:         jwt,
		orders:      oreders,
		productrepo: productrepo,
		cart:        cart,
		wishlist:    wishlist,
		address:     address,
	}
}

func (s *userUseCase) SignUp(ctx context.Context, name, email, password string) (string, error) {

	user, err := s.repo.GetByEmail(ctx, email)

	if err == nil && user.ID != 0 {
		return "", errors.New("user already exists")
	}

	hashedPass, err := helper.Hash(password)
	if err != nil {
		return "", err
	}

	otp := helper.GenerateOtp()

	signupRequest := &domain.SignupRequest{
		Name:      name,
		Email:     email,
		Password:  hashedPass,
		Role:      "user",
		Otp:       otp,
		OtpExpire: time.Now().Add(10 * time.Minute).Unix(),
	}

	if err := s.repo.SaveSignup(ctx, signupRequest); err != nil {
		return "", err
	}

	if err := s.otp.SendOtp(email, otp); err != nil {
		return "", errors.New("failed to send otp")
	}

	return "OTP sent successfully", nil
}

func (s *userUseCase) VerifyOtp(ctx context.Context, email, code string) (string, error) {
	signup, err := s.repo.GetSignup(ctx, email)
	if err != nil {
		return "", errors.New("invalid request or user not found")
	}

	if signup.Otp != code {
		return "", errors.New("invalid code")
	}
	if time.Now().Unix() > signup.OtpExpire {
		return "", errors.New("code expired")
	}

	// Create User
	newUser := &domain.User{
		Name:     signup.Name,
		Email:    signup.Email,
		Password: signup.Password,
		Role:     signup.Role,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return "", err
	}

	// Delete Signup Request
	_ = s.repo.DeleteSignup(ctx, email)

	// Generate Token
	token, err := s.jwt.GenerateToken(newUser.ID, s.jwt.AccessTTL, newUser.Role)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *userUseCase) SignIn(ctx context.Context, email, password string) (*UserProfileOutput, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, "", errors.New("account not verified")
	}

	if err := helper.VerifyHash(user.Password, password); !err {
		return nil, "", errors.New("invalid email or password")
	}
	addre, _ := s.address.GetByUserID(ctx, user.ID)
	cart, _ := s.cart.GetCart(ctx, user.ID)
	wish, _ := s.wishlist.GetAll(ctx, user.ID)
	order, _ := s.orders.GetAllOrdersByUserID(ctx, user.ID)

	var orderitems []domain.OrderItem

	for _, ord := range order {
		o, _ := s.orders.GetOrdersByOrderID(ctx, ord.ID)
		orderitems = append(orderitems, o...)
	}
	userdata := UserProfileOutput{
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		Cart:     cart,
		Wishlist: wish,
		Orders:   orderitems,
		Address:  addre,
	}

	acc, erro := s.jwt.GenerateToken(user.ID, s.jwt.AccessTTL, user.Role)
	if erro != nil {
		return nil, "", erro
	}
	// fmt.Println("user",userdata,"acc",acc)
	return &userdata, acc, nil
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
	addre, _ := s.address.GetByUserID(ctx, users.ID)
	cart, _ := s.cart.GetCart(ctx, users.ID)
	wish, _ := s.wishlist.GetAll(ctx, users.ID)
	order, _ := s.orders.GetAllOrdersByUserID(ctx, users.ID)

	var orderitems []domain.OrderItem

	for _, ord := range order {
		o, _ := s.orders.GetOrdersByOrderID(ctx, ord.ID)
		orderitems = append(orderitems, o...)
	}
	userdata := UserProfileOutput{
		Name:      users.Name,
		Email:     users.Email,
		Role:      users.Role,
		CreatedAt: users.CreatedAt,
		Cart:      cart,
		Wishlist:  wish,
		Orders:    orderitems,
		Address:   addre,
	}
	return &userdata, nil
}

func (s *userUseCase) GetProduct(
	ctx context.Context,
	id uint,
) (*domain.Product, error) {
	product, err := s.productrepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *userUseCase) GetOrderDetail(ctx context.Context, orderID, userID uint) ([]domain.Order, error) {
	return s.orders.GetOrdersByUserIDAndOrderID(ctx, userID, orderID)
}
func (s *userUseCase) GetOrderItemDetails(ctx context.Context, orderID uint) ([]domain.OrderItem, error) {
	return s.orders.GetOrdersByOrderID(ctx, orderID)
}

func (s *userUseCase) CancelOrder(ctx context.Context, orderID, userID uint) error {
	return s.orders.CancelOrder(ctx, orderID, userID)
}

func (s *userUseCase) GetAllProducts(
	ctx context.Context,
) ([]domain.Product, error) {

	product, err := s.productrepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *userUseCase) SearchProducts(ctx context.Context, query string) ([]domain.Product, error) {
	cleanQuery := strings.TrimSpace(query)

	if cleanQuery == "" {
		return []domain.Product{}, nil
	}
	return s.productrepo.Search(ctx, cleanQuery)
}

func (s *userUseCase) FilterProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error) {
	return s.productrepo.GetProducts(ctx, filter)
}
