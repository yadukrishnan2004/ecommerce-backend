package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler/dto"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/usecase"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

type OrderHandler struct {
    svc usecase.OrderService
}

func NewOrderHandler(svc usecase.OrderService) *OrderHandler {
    return &OrderHandler{svc: svc}
}

func (h *OrderHandler) PlaceOrder(c *fiber.Ctx) error {
    userIDFloat, ok := c.Locals("userid").(float64)
    if !ok {
        return response.Response(c,fiber.StatusUnauthorized,"no user found",nil,nil)
    }

    err := h.svc.PlaceOrder(c.Context(), uint(userIDFloat))
    if err != nil {
        return response.Response(c,fiber.StatusBadRequest,"server",nil,err.Error())
    }
	return response.Response(c,fiber.StatusOK,"Order placed successfully",nil,nil)
}

func (h *OrderHandler) GetOrderHistory(c *fiber.Ctx) error {
    
    userIDFloat, ok := c.Locals("userid").(float64)
    if !ok {
        return response.Response(c,fiber.StatusUnauthorized,"no user found",nil,nil)
    }

    
    orders, err := h.svc.GetOrderHistory(c.Context(), uint(userIDFloat))
    if err != nil {
       return response.Response(c,fiber.StatusBadRequest,"server",nil,err.Error())
    }
    history:=dto.Orders{
        Count: len(orders),
        Items: orders,
    }

    return response.Response(c,fiber.StatusOK,"get Order successfully",history,nil)

}

func (h *OrderHandler) BuyNow(c *fiber.Ctx) error {
    userIDFloat, ok := c.Locals("userid").(float64)
    if !ok {
        return response.Response(c,fiber.StatusUnauthorized,"no user found",nil,nil)
    }

    type BuyNowReq struct {
        ProductID uint `json:"product_id"`
        Quantity  int  `json:"quantity"`
    }
    req := new(BuyNowReq)
    if err := c.BodyParser(req); err != nil {
        return response.Response(c,http.StatusBadRequest,"invalid input value",nil,nil)
    }

    if req.Quantity <= 0 {
        return response.Response(c,http.StatusBadRequest, "Quantity must be greater than 0",nil,nil)
    }

    err := h.svc.BuyNow(c.Context(), uint(userIDFloat), req.ProductID, req.Quantity)
    if err != nil {
        return response.Response(c,http.StatusInternalServerError,"server error",nil,err.Error())
    }

    return response.Response(c,http.StatusOK, "Order placed successfully",nil,nil)
}

