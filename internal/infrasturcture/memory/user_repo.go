package memory

import (
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain/user"
)


//Temporary in memory repository 

type UserRepo struct {
	users map[string]*user.User
}

func NewUserRepo() *UserRepo{
	return &UserRepo{
		users:make(map[string]*user.User),
	}
}

func (r *UserRepo) FindByEmail(email string) (*user.User,error){
	return r.users[email],nil
}

func (r *UserRepo) Create(u *user.User) error{
	r.users[u.Email]=u
	return nil
}