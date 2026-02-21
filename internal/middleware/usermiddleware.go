package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

func UserMiddleware(c *fiber.Ctx) error {
	token := c.Cookies("jwt")

	if token == "" {
		authHead := c.Get("Authorization")
		if strings.HasPrefix(authHead, "Bearer ") {
			token = strings.TrimPrefix(authHead, "Bearer ")
		}
	}
	if token == "" {
		return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "token is nil")
	}
	stoken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, err.Error())
	}
	if !stoken.Valid {
		return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "invalid token")
	}
	claims, ok := stoken.Claims.(jwt.MapClaims)
	if !ok {
		return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "invalid claims")
	}
	if claims["role"] !=  "user"{
		if claims["role"] !=  "admin"{
			return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "you cannot access this")
		}
	}
	c.Locals("userid", claims["sub"])
	fmt.Print(claims["sub"])
	return c.Next()
}
