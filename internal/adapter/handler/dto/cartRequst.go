package dto

type CartRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}