package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ResetMiddleware(c *fiber.Ctx) error{
	token:=c.Cookies("forgetpassword")

	if token == ""{
		authHead:=c.Get("Authorization")
		if strings.HasPrefix(authHead,"Bearer "){
			token=strings.TrimPrefix(authHead,"Bearer ")
		}
	}

	if token==""{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	stoken,err:=jwt.Parse(token,func(t *jwt.Token)(any,error){
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err != nil || !stoken.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Invalid token"})
	}
	claims, ok := stoken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Invalid claims"})
	}
	c.Locals("email", claims["sub"])
	return c.Next()
}