package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/repository"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
	"gorm.io/gorm"
)

func Adminmiddleware(db *gorm.DB) fiber.Handler {
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
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
		}

		if token == "" || (err != nil && errors.Is(err, jwt.ErrTokenExpired)) {
			refreshToken := c.Cookies("refresh")
			if refreshToken == "" {
				refreshToken = c.Get("X-Refresh-Token")
			}
			if refreshToken == "" {
				return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "refresh token is missing")
			}

			rtoken, rerr := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fiber.ErrUnauthorized
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if rerr != nil {
				return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, rerr.Error())
			}
			if !rtoken.Valid {
				return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "invalid refresh token")
			}

			cl, ok := rtoken.Claims.(jwt.MapClaims)
			if !ok {
				return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "invalid refresh claims")
			}

			secret := os.Getenv("JWT_SECRET")

			newClaims := jwt.MapClaims{
				"sub":  cl["sub"],
				"role": cl["role"],
				"exp":  time.Now().Add(15 * time.Minute).Unix(),
				"iat":  time.Now().Unix(),
			}

			t := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
			signedToken, signErr := t.SignedString([]byte(secret))
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
			if role != "admin" {
				return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "you cannot access this")
			}

			// Check user status in database
			var user repository.User
			if err := db.Where("id = ?", cl["sub"]).First(&user).Error; err != nil {
				return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "user not found")
			}
			if user.IsBlocked {
				return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "your account has been blocked")
			}
			if !user.IsActive {
				return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "your account is inactive")
			}

			c.Locals("email", cl["sub"])
			return c.Next()
		}

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
		if claims["role"] != "admin" {
			return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "you cannot access this")
		}

		// Check user status in database
		var user repository.User
		if err := db.Where("id = ?", claims["sub"]).First(&user).Error; err != nil {
			return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "user not found")
		}
		if user.IsBlocked {
			return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "your account has been blocked")
		}
		if !user.IsActive {
			return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "your account is inactive")
		}

		c.Locals("email", claims["sub"])
		return c.Next()
	}
}
