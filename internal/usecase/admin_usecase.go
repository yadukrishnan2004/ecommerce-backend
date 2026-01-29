package usecase

import (
	"context"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type AdminUpdateUserInput struct {
	Name      *string
	Email     *string
	Role      *string
	IsActive  *bool
	IsBlocked *bool
}

type AdminUseCase interface {
	UpdateUser(ctx context.Context, userId uint, input AdminUpdateUserInput) (*domain.User, error)
	BlockUser(ctx context.Context, userId uint) (string, error)
}

type adminUseCase struct {
	repo domain.UserRepository
}

func NewAdminUseCase(rep domain.UserRepository) AdminUseCase {
	return &adminUseCase{
		repo: rep,
	}
}

func (s *adminUseCase) UpdateUser(
	ctx context.Context,
	userId uint,
	input AdminUpdateUserInput,
) (*domain.User, error) {

	user, err := s.repo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Email != nil {
		user.Email = *input.Email
	}

	if input.Role != nil {
		user.Role = *input.Role
	}

	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	if input.IsBlocked != nil {
		user.IsBlocked = *input.IsBlocked
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *adminUseCase) BlockUser(
	ctx context.Context,
	userId uint,
) (string, error) {

	user, err := s.repo.GetByID(ctx, userId)
	if err != nil {
		return "", err
	}

	// Toggle block status
	user.IsBlocked = !user.IsBlocked

	if err := s.repo.Update(ctx, user); err != nil {
		return "", err
	}

	if user.IsBlocked {
		return "user blocked", nil
	}

	return "user unblocked", nil
}
