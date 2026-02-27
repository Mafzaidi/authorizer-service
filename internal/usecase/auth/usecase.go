package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mafzaidi/authorizer/internal/delivery/http/middleware"
	"github.com/mafzaidi/authorizer/internal/infrastructure/config"
	"github.com/mafzaidi/authorizer/internal/delivery/http/middleware/pwd"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	"github.com/mafzaidi/authorizer/internal/domain/service"
)

// JWTService defines the interface for JWT infrastructure service
// This allows the usecase to depend on the interface rather than concrete implementation
type JWTService interface {
	GenerateToken(ctx context.Context, claims *entity.Claims, privateKey *rsa.PrivateKey, keyID string) (string, error)
	ValidateToken(ctx context.Context, tokenString string, publicKey *rsa.PublicKey) (*entity.Claims, error)
}

type UserToken struct {
	User         *entity.User
	Token        string
	RefreshToken string
	Claims       *middleware.JWTClaims
}

type authUsecase struct {
	authRepo    repository.AuthRepository
	userRepo    repository.UserRepository
	authService service.AuthService
	jwtService  JWTService
	logger      service.Logger
}

func NewAuthUseCase(
	authRepo repository.AuthRepository,
	userRepo repository.UserRepository,
	authService service.AuthService,
	jwtService JWTService,
	logger service.Logger,
) Usecase {
	return &authUsecase{
		authRepo:    authRepo,
		userRepo:    userRepo,
		authService: authService,
		jwtService:  jwtService,
		logger:      logger,
	}
}

func (uc *authUsecase) Login(
	ctx context.Context,
	appCode,
	email,
	password,
	validToken string,
	cfg *config.Config,
) (*UserToken, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate input
	if email == "" {
		uc.logger.Warn("Login attempt with empty email", service.Fields{})
		return nil, errors.New("email cannot be empty")
	}

	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		uc.logger.Warn("Login failed: user not found", service.Fields{
			"email": email,
		})
		return nil, errors.New("email or password is invalid")
	}

	// Verify password
	if !pwd.CheckHash(user.Password, password) {
		uc.logger.Warn("Login failed: invalid password", service.Fields{
			"email":   email,
			"user_id": user.ID,
		})
		return nil, errors.New("email or password is invalid")
	}

	// Check if we can reuse existing valid token
	if validToken != "" {
		existingClaims, err := uc.jwtService.ValidateToken(ctx, validToken, cfg.JWT.PublicKey)
		if err == nil && existingClaims.Subject == user.ID {
			// Token is still valid and belongs to this user, reuse it
			uc.logger.Info("Reusing valid token", service.Fields{
				"user_id": user.ID,
				"email":   email,
			})

			// Convert entity.Claims to middleware.JWTClaims for backward compatibility
			middlewareClaims := convertToMiddlewareClaims(existingClaims)

			return &UserToken{
				User:         user,
				Token:        validToken,
				RefreshToken: "", // Not generating new refresh token
				Claims:       middlewareClaims,
			}, nil
		}
	}

	// Build claims using domain service
	claims, err := uc.authService.BuildClaims(ctx, user, appCode)
	if err != nil {
		uc.logger.Error("Failed to build claims", service.Fields{
			"user_id":  user.ID,
			"app_code": appCode,
			"error":    err.Error(),
		})
		return nil, errors.New("failed to build authorization claims")
	}

	// Generate access token using infrastructure service
	accessToken, err := uc.jwtService.GenerateToken(ctx, claims, cfg.JWT.PrivateKey, cfg.JWT.KeyID)
	if err != nil {
		uc.logger.Error("Failed to generate access token", service.Fields{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return nil, errors.New("failed to generate access token")
	}

	// Generate refresh token
	refreshToken, err := uc.generateRefreshToken()
	if err != nil {
		uc.logger.Error("Failed to generate refresh token", service.Fields{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return nil, errors.New("failed generating refresh token")
	}

	// Store refresh token
	err = uc.authRepo.StoreRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		uc.logger.Error("Failed to store refresh token", service.Fields{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return nil, errors.New("failed saving refresh token")
	}

	uc.logger.Info("User logged in successfully", service.Fields{
		"user_id":  user.ID,
		"email":    email,
		"app_code": appCode,
	})

	// Convert entity.Claims to middleware.JWTClaims for backward compatibility
	middlewareClaims := convertToMiddlewareClaims(claims)

	token := &UserToken{
		User:         user,
		Token:        accessToken,
		RefreshToken: refreshToken,
		Claims:       middlewareClaims,
	}

	return token, nil
}

func (uc *authUsecase) RefreshToken(
	ctx context.Context,
	refreshToken string,
	cfg *config.Config,
) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userID, err := uc.authRepo.GetUserIDByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	savedToken, err := uc.authRepo.GetRefreshToken(ctx, userID)
	if err != nil || savedToken != refreshToken {
		return "", "", errors.New("refresh token mismatch")
	}

	return "", "", nil
}

// generateRefreshToken generates a random refresh token
// This is a simple implementation that can be enhanced with more security
func (uc *authUsecase) generateRefreshToken() (string, error) {
	// For now, use a simple UUID-like token
	// In production, this should use crypto/rand for secure random generation
	return time.Now().Format("20060102150405") + "-refresh-token", nil
}

// convertToMiddlewareClaims converts entity.Claims to middleware.JWTClaims
// This maintains backward compatibility with existing middleware
func convertToMiddlewareClaims(claims *entity.Claims) *middleware.JWTClaims {
	if claims == nil {
		return nil
	}

	// Convert entity.Authorization to middleware.Authorization
	var middlewareAuth []middleware.Authorization
	for _, auth := range claims.Authorization {
		middlewareAuth = append(middlewareAuth, middleware.Authorization{
			App:         auth.App,
			Roles:       auth.Roles,
			Permissions: auth.Permissions,
		})
	}

	return &middleware.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    claims.Issuer,
			Subject:   claims.Subject,
			Audience:  claims.Audience,
			ExpiresAt: jwt.NewNumericDate(time.Unix(claims.ExpiresAt, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Unix(claims.IssuedAt, 0)),
		},
		UserID:        claims.Subject,
		Username:      claims.Username,
		Email:         claims.Email,
		Authorization: middlewareAuth,
	}
}
