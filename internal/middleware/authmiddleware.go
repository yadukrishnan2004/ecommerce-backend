package middleware

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/config"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

// ResetMiddleware validates the short-lived auth token issued during ForgotPassword.
// It expects a "type": "auth" claim and injects the user's email into c.Locals("email").
// The JWT secret is sourced from cfg.JWT, never from os.Getenv at request time.
func ResetMiddleware(cfg *config.Config) fiber.Handler {
	secret := []byte(cfg.JWT.Secret)

	return func(c *fiber.Ctx) error {
		token := c.Cookies("jwtverify")

		if token == "" {
			authHead := c.Get("Authorization")
			if strings.HasPrefix(authHead, "Bearer ") {
				token = strings.TrimPrefix(authHead, "Bearer ")
			}
		}
		if token == "" {
			return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "token is missing")
		}

		stoken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return secret, nil
		})
		if err != nil {
			return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, err.Error())
		}

		if !stoken.Valid {
			return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "invalid token")
		}
		claims, ok := stoken.Claims.(jwt.MapClaims)
		if !ok {
			return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "invalid claims")
		}

		// Validate token type — must be "auth" (not a regular access or refresh token)
		if claims["type"] != "auth" {
			return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "invalid token type")
		}

		email, ok := claims["sub"].(string)
		if !ok || email == "" {
			return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "invalid email in token")
		}

		c.Locals("email", email)
		return c.Next()
	}
}
