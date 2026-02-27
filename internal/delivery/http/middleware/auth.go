package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/service"
	"github.com/mafzaidi/authorizer/internal/infrastructure/auth"
	"github.com/mafzaidi/authorizer/internal/infrastructure/config"
	"github.com/mafzaidi/authorizer/internal/infrastructure/logger"
	"github.com/mafzaidi/authorizer/pkg/response"
)

type contextKey string

const userContextKey = contextKey("user_claims")

// JWTAuthMiddleware creates a JWT authentication middleware with explicit dependencies
// Parameters:
//   - jwtService: JWT service for token validation
//   - cfg: configuration containing JWT public key
//   - log: logger for structured logging of auth failures
//
// Returns:
//   - echo.MiddlewareFunc: middleware function that validates JWT tokens
func JWTAuthMiddleware(jwtService auth.JWTService, cfg *config.Config, log service.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				log.Warn("Authentication failed: missing token", service.Fields{
					"path":   c.Request().URL.Path,
					"method": c.Request().Method,
				})
				return response.ErrorHandler(c, http.StatusUnauthorized, "Unauthorized", "token is missing")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Warn("Authentication failed: invalid token format", service.Fields{
					"path":   c.Request().URL.Path,
					"method": c.Request().Method,
				})
				return response.ErrorHandler(c, http.StatusUnauthorized, "Unauthorized", "invalid token format")
			}

			rawToken := parts[1]

			// Use JWT service to validate token
			ctx := context.Background()
			claims, err := jwtService.ValidateToken(ctx, rawToken, cfg.JWT.PublicKey)
			if err != nil {
				log.Warn("Authentication failed: token validation error", service.Fields{
					"path":  c.Request().URL.Path,
					"method": c.Request().Method,
					"error": err.Error(),
					"token": logger.TruncateToken(rawToken),
				})
				return response.ErrorHandler(c, http.StatusUnauthorized, "Unauthorized", "invalid token")
			}

			// Convert entity.Claims to JWTClaims for backward compatibility with existing code
			jwtClaims := &JWTClaims{
				UserID:        claims.Subject,
				Username:      claims.Username,
				Email:         claims.Email,
				Authorization: convertAuthorization(claims.Authorization),
			}

			c.Set(string(userContextKey), jwtClaims)

			return next(c)
		}
	}
}

// convertAuthorization converts entity.Authorization to middleware.Authorization
func convertAuthorization(entityAuth []entity.Authorization) []Authorization {
	result := make([]Authorization, len(entityAuth))
	for i, auth := range entityAuth {
		result[i] = Authorization{
			App:         auth.App,
			Roles:       auth.Roles,
			Permissions: auth.Permissions,
		}
	}
	return result
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
	if claims == nil {
		return false
	}
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
