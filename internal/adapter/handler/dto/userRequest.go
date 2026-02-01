package dto

type UserProfile struct {
    Name    	string `json:"name"`
	Email   	string `json:"email"`
	Role 		string `json:"role"`
	IsBlocked   bool   `json:"is_blocked"`
}
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type Otp struct {
	Otp string `json:"otp" validate:"required,min=6,max=6"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type  Getemail struct {
	Email string `json:"email" validate:"required,email"`
}

type Reset struct {
	Code        string `json:"code" validate:"required,min=6,max=6"`
	Newpassword string `json:"password" validate:"required,email"`
}

type UpdateUser struct {
	Name      *string `json:"name,omitempty"`
    Email     *string `json:"email,omitempty"`
}

type UpdateStatus struct{
	Status  string `json:"status" validate:"required"`
}

