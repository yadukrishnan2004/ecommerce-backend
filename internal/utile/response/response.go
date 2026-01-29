package response

import (
	"github.com/gofiber/fiber/v2"
)



type ApiResponse struct {
	StatusCode int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}

func Response(
	c *fiber.Ctx,
	status int,
	message string,
	data interface{},
	err interface{},
) error {

	apiResponse := ApiResponse{
		StatusCode: status,
		Message:    message,
		Data:       data,
		Error:      err,
	}

	return c.Status(status).JSON(apiResponse)
}
