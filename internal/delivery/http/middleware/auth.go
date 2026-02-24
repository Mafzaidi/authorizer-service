package middleware

import (
	"net/http"
	"slices"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/config"

	"github.com/mafzaidi/authorizer/pkg/response"
)

type contextKey string

const userContextKey = contextKey("user_claims")

func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := config.GetConfig()
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return response.ErrorHandler(c, http.StatusUnauthorized, "Unauthorized", "token is missing")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.ErrorHandler(c, http.StatusUnauthorized, "Unauthorized", "invalid token format")
		}

		rawToken := parts[1]

		token, err := jwt.ParseWithClaims(rawToken, &JWTClaims{}, func(token *jwt.Token) (any, error) {

			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, echo.ErrUnauthorized
			}

			return cfg.JWT.PublicKey, nil
		})

		if err != nil || !token.Valid {
			return response.ErrorHandler(c, http.StatusUnauthorized, "Unauthorized", "invalid token")
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return response.ErrorHandler(c, http.StatusUnauthorized, "Unauthorized", "invalid token claims")
		}

		c.Set(string(userContextKey), claims)

		return next(c)
	}
}

func RequirePermission(app, perm string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			claims := GetUserFromContext(c)
			if !HasPermission(claims, app, perm) {
				return response.ErrorHandler(c, http.StatusForbidden, "Forbidden", "missing required permission")
			}

			return next(c)
		}
	}
}

func HasPermission(claims *JWTClaims, app, perm string) bool {
	for _, a := range claims.Authorization {
		if a.App == "GLOBAL" {
			return true
		}
		if a.App == app {
			for _, p := range a.Permissions {
				if p == perm {
					return true
				}
			}
		}
	}
	return false
}

func GetUserFromContext(c echo.Context) *JWTClaims {
	if claims, ok := c.Get(string(userContextKey)).(*JWTClaims); ok {
		return claims
	}
	return nil
}

func HasRole(required string, userRoles []string) bool {
	return slices.Contains(userRoles, required)
}
