package dto

type CartRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,min=1"`
}

type Cart struct {
	GrandTotal float32            `json:"grand_total"`
	Items      []CartItemResponse `json:"items"`
	Count      uint               `json:"count"`
}

type CartItemResponse struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	SubTotal    float64 `json:"sub_total"`
}

type UpdateReq struct {
	Quantity int `json:"quantity" validate:"required,min=0"`
}
