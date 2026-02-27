package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/service"
)

// JWTService handles JWT token generation and validation (infrastructure concern)
// This service is responsible for the technical aspects of JWT tokens:
// - Signing tokens with RSA private keys
// - Validating token signatures with RSA public keys
// - Parsing and verifying JWT structure
//
// Business logic for building claims is handled by the domain service.
type JWTService interface {
	// GenerateToken creates a signed JWT token from the provided claims
	// Parameters:
	//   - ctx: context for cancellation and timeout
	//   - claims: the claims to encode in the token (built by domain service)
	//   - privateKey: RSA private key for signing
	//   - keyID: key identifier for JWKS
	// Returns:
	//   - string: the signed JWT token
	//   - error: if signing fails
	GenerateToken(ctx context.Context, claims *entity.Claims, privateKey *rsa.PrivateKey, keyID string) (string, error)

	// ValidateToken validates a JWT token and returns the claims
	// Parameters:
	//   - ctx: context for cancellation and timeout
	//   - tokenString: the JWT token to validate
	//   - publicKey: RSA public key for verification
	// Returns:
	//   - *entity.Claims: the validated claims
	//   - error: if validation fails
	ValidateToken(ctx context.Context, tokenString string, publicKey *rsa.PublicKey) (*entity.Claims, error)
}

type jwtService struct {
	logger service.Logger
}

// NewJWTService creates a new JWT service instance
func NewJWTService(logger service.Logger) JWTService {
	return &jwtService{
		logger: logger,
	}
}

// jwtClaims is an internal struct that bridges between entity.Claims and jwt.Claims
// It embeds jwt.RegisteredClaims for standard JWT fields and adds custom fields
type jwtClaims struct {
	jwt.RegisteredClaims
	Username      string                 `json:"username"`
	Email         string                 `json:"email"`
	Authorization []entity.Authorization `json:"authorization"`
}

// GenerateToken creates a signed JWT token from the provided claims
func (s *jwtService) GenerateToken(ctx context.Context, claims *entity.Claims, privateKey *rsa.PrivateKey, keyID string) (string, error) {
	if claims == nil {
		s.logger.Error("GenerateToken called with nil claims", service.Fields{})
		return "", errors.New("claims cannot be nil")
	}

	if privateKey == nil {
		s.logger.Error("GenerateToken called with nil private key", service.Fields{})
		return "", errors.New("private key cannot be nil")
	}

	// Convert entity.Claims to jwt.Claims
	jwtClaims := &jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   claims.Issuer,
			Subject:  claims.Subject,
			Audience: claims.Audience,
		},
		Username:      claims.Username,
		Email:         claims.Email,
		Authorization: claims.Authorization,
	}

	// Set timestamps from Unix timestamps
	jwtClaims.ExpiresAt = jwt.NewNumericDate(time.Unix(claims.ExpiresAt, 0))
	jwtClaims.IssuedAt = jwt.NewNumericDate(time.Unix(claims.IssuedAt, 0))

	// Create token with RS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaims)

	// Set kid in token header for JWKS key identification
	token.Header["kid"] = keyID

	// Sign the token
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		s.logger.Error("Failed to sign JWT token", service.Fields{
			"error":   err.Error(),
			"subject": claims.Subject,
		})
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *jwtService) ValidateToken(ctx context.Context, tokenString string, publicKey *rsa.PublicKey) (*entity.Claims, error) {
	if tokenString == "" {
		s.logger.Warn("ValidateToken called with empty token string", service.Fields{})
		return nil, errors.New("token string cannot be empty")
	}

	if publicKey == nil {
		s.logger.Error("ValidateToken called with nil public key", service.Fields{})
		return nil, errors.New("public key cannot be nil")
	}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method is RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			s.logger.Warn("Unexpected signing method", service.Fields{
				"method": token.Method.Alg(),
			})
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		s.logger.Warn("Failed to parse JWT token", service.Fields{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		s.logger.Error("Failed to extract claims from token", service.Fields{})
		return nil, errors.New("invalid token claims")
	}

	// Verify token is valid
	if !token.Valid {
		s.logger.Warn("Token is not valid", service.Fields{
			"subject": claims.Subject,
		})
		return nil, errors.New("token is not valid")
	}

	// Convert jwt.Claims back to entity.Claims
	entityClaims := &entity.Claims{
		Issuer:        claims.Issuer,
		Subject:       claims.Subject,
		Audience:      claims.Audience,
		ExpiresAt:     claims.ExpiresAt.Unix(),
		IssuedAt:      claims.IssuedAt.Unix(),
		Username:      claims.Username,
		Email:         claims.Email,
		Authorization: claims.Authorization,
	}

	return entityClaims, nil
}
