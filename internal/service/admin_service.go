package service

import (
	"context"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type adminService struct {
	repo domain.UserRepositery
}

func NewAdminService(rep domain.UserRepositery) domain.AdminService{
	return &adminService{
		repo:rep,
	}
}	

func(s *adminService) UpdateUser(ctx context.Context,userId uint,req domain.User)(*domain.User,error){
  user, err := s.repo.GetByID(ctx, userId)
    if err != nil {
        return nil, err
    }

    if req.Name != "" {
        user.Name = req.Name
    }
    if req.Email != "" {
        user.Email = req.Email
    }
    if req.Role != "" {
        user.Role = req.Role
    }

    user.IsActive = req.IsActive
    user.IsBlocked = req.IsBlocked

    if err := s.repo.Update(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

func(s *adminService) BlockUser(ctx context.Context,userId uint)error{
	user,err:=s.repo.GetByID(ctx,userId)
	if err != nil{
		return err
	}
	if !user.IsBlocked{
		user.IsBlocked=true
	}else{
		user.IsBlocked=false
	}
	return s.repo.Update(ctx,user)
}