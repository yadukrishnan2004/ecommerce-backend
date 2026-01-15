package user

import "errors"

type User struct {
	ID       string
	Email    string
	Password string
}

func New(email, password string) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if len(password)< 8{
		return nil,errors.New("password must be at least 8 charactres")
	}

	return &User{
		Email: email,
		Password: password,
	},nil

}