package middleware

import (
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"localdev.me/authorizer/config"

	"localdev.me/authorizer/pkg/response"
)

type contextKey string

const userContextKey = contextKey("user")

func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := config.GetConfig()
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return response.ErrorHandler(c, http.StatusUnauthorized, "Unauthorized", "token is missing")
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return response.ErrorHandler(c, http.StatusUnauthorized, "Unauthorized", "invalid token format")
		}

		token, err := jwt.ParseWithClaims(tokenParts[1], &Claims{}, func(token *jwt.Token) (any, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			return response.ErrorHandler(c, http.StatusUnauthorized, "Unauthorized", "invalid token")
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return response.ErrorHandler(c, http.StatusUnauthorized, "Unauthorized", "invalid token claims")
		}

		c.Set(string(userContextKey), claims)

		return next(c)
	}
}

func GetUserFromContext(c echo.Context) *Claims {
	if user, ok := c.Get(string(userContextKey)).(*Claims); ok {
		return user
	}
	return nil
}
