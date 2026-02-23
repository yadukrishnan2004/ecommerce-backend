package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

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
		if errors.Is(err, jwt.ErrTokenExpired) {
			refreshToken := c.Cookies("refresh")
			if refreshToken == "" {
				authHead := c.Get("Authorization")
				if strings.HasPrefix(authHead, "Bearer ") {
					refreshToken = strings.TrimPrefix(authHead, "Bearer ")
				}
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
			})
			c.Set("Authorization", "Bearer "+signedToken)

			role, _ := cl["role"].(string)
			if role != "user" && role != "admin" {
				return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "you cannot access this")
			}

			c.Locals("userid", cl["sub"])
			fmt.Print(cl["sub"])
			return c.Next()
		}
		return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, err.Error())
	}
	if !stoken.Valid {
		return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "invalid token")
	}
	claims, ok := stoken.Claims.(jwt.MapClaims)
	if !ok {
		return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "invalid claims")
	}
	if claims["role"] != "user" && claims["role"] != "admin" {
		return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, "you cannot access this")
	}
	c.Locals("userid", claims["sub"])
	fmt.Print(claims["sub"])
	return c.Next()
}
