package dto

type AdminUpdateUserRequest struct {
	Name      *string `json:"name,omitempty"`
	Email     *string `json:"email,omitempty"`
	Role      *string `json:"role,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
	IsBlocked *bool   `json:"is_blocked,omitempty"`
}

type AddNewProduct struct {
	Name        string   `json:"name" validate:"required"`
	Price       int      `json:"price" validate:"required"`
	Description string   `json:"desc" validate:"required"`
	Category    string   `json:"category" validate:"required"`
	Offer       string   `json:"offer,omitempty"`
	OfferPrice  int      `json:"offerprice,omitempty"`
	Production  string   `json:"production,omitempty"`
	Images      []string `json:"images"`
}
