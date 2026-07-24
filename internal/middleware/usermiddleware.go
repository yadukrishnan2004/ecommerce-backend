package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/repository"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/config"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
	"gorm.io/gorm"
)

// UserMiddleware validates the JWT access token (or auto-refreshes via refresh token)
// and injects the authenticated user's ID into c.Locals("userid").
// The JWT secret is sourced from cfg.JWT, never from os.Getenv at request time.
func UserMiddleware(db *gorm.DB, cfg *config.Config) fiber.Handler {
	secret := []byte(cfg.JWT.Secret)

	return func(c *fiber.Ctx) error {
		token := c.Cookies("jwt")

		if token == "" {
			authHead := c.Get("Authorization")
			if strings.HasPrefix(authHead, "Bearer ") {
				token = strings.TrimPrefix(authHead, "Bearer ")
			}
		}

		var stoken *jwt.Token
		var err error
		if token != "" {
			stoken, err = jwt.Parse(token, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fiber.ErrUnauthorized
				}
				return secret, nil
			})
		}

		// Try to refresh token if access token is missing or expired
		if token == "" || (err != nil && errors.Is(err, jwt.ErrTokenExpired)) {
			refreshToken := c.Cookies("refresh")
			if refreshToken == "" {
				refreshToken = c.Get("X-Refresh-Token")
			}
			if refreshToken == "" {
				return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "session expired")
			}

			rtoken, rerr := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fiber.ErrUnauthorized
				}
				return secret, nil
			})

			if rerr != nil {
				return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, rerr.Error())
			}
			if !rtoken.Valid {
				return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "invalid refresh token")
			}

			cl, ok := rtoken.Claims.(jwt.MapClaims)
			if !ok {
				return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "invalid refresh claims")
			}

			// Validate token type — must be "refresh"
			if cl["type"] != "refresh" {
				return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "invalid token type")
			}

			newClaims := jwt.MapClaims{
				"sub":  cl["sub"],
				"role": cl["role"],
				"type": "access",
				"exp":  time.Now().Add(15 * time.Minute).Unix(),
				"iat":  time.Now().Unix(),
			}

			t := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
			signedToken, signErr := t.SignedString(secret)
			if signErr != nil {
				return response.Response(c, http.StatusInternalServerError, "error", nil, "failed to generate new access token")
			}

			c.Cookie(&fiber.Cookie{
				Name:     "jwt",
				Value:    signedToken,
				Expires:  time.Now().Add(15 * time.Minute),
				HTTPOnly: true,
				Path:     "/",
			})
			c.Set("Authorization", "Bearer "+signedToken)

			role, _ := cl["role"].(string)
			if role != "user" && role != "admin" {
				return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "you cannot access this")
			}

			userID := getIDFromClaim(cl["sub"])
			if userID == 0 {
				return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "invalid user claims")
			}

			var user repository.User
			if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
				return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "user not found")
			}
			if user.IsBlocked {
				return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "your account has been blocked")
			}
			if !user.IsActive {
				return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "your account is inactive")
			}

			c.Locals("userid", float64(userID))
			return c.Next()
		}

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
		if claims["role"] != "user" && claims["role"] != "admin" {
			return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "you cannot access this")
		}

		userID := getIDFromClaim(claims["sub"])
		if userID == 0 {
			return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "invalid user claims")
		}

		var user repository.User
		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "user not found")
		}
		if user.IsBlocked {
			return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "your account has been blocked")
		}
		if !user.IsActive {
			return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, "your account is inactive")
		}

		c.Locals("userid", float64(userID))
		return c.Next()
	}
}