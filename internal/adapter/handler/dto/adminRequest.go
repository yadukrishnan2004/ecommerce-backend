package dto

type AdminUpdateUserRequest struct {
    Name      *string `json:"name,omitempty"`
    Email     *string `json:"email,omitempty"`
    Role      *string `json:"role,omitempty"`
    IsActive  *bool   `json:"is_active,omitempty"`
    IsBlocked *bool   `json:"is_blocked,omitempty"`
}