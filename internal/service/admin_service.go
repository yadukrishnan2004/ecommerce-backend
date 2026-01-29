package service

import (
	"context"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler/dto"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type AdminService interface {
	UpdateUser(ctx context.Context, userId uint, req *dto.AdminUpdateUserRequest) (*domain.User, error)
	BlockUser(ctx context.Context, userId uint) (string, error)
}

type adminService struct {
	repo domain.UserRepositery
}

func NewAdminService(rep domain.UserRepositery) AdminService {
	return &adminService{
		repo: rep,
	}
}

func (s *adminService) UpdateUser(
	ctx context.Context,
	userId uint,
	req *dto.AdminUpdateUserRequest,
) (*domain.User, error) {

	user, err := s.repo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if req.IsBlocked != nil {
		user.IsBlocked = *req.IsBlocked
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *adminService) BlockUser(
	ctx context.Context,
	userId uint,
) (string, error) {

	user, err := s.repo.GetByID(ctx, userId)
	if err != nil {
		return "", err
	}

	if user.IsBlocked {
		user.IsBlocked = false
	} else {
		user.IsBlocked = true
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return "", err
	}

	if user.IsBlocked {
		return "user blocked", nil
	}

	return "user unblocked", nil
}
