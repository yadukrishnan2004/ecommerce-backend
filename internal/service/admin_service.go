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
	err:=s.repo.Update(ctx,&req)
	if err != nil{
		return nil,err
	}
	user,erro:=s.repo.GetByID(ctx,userId)
	if erro !=nil{
		return nil,erro
	}
	return user,nil
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