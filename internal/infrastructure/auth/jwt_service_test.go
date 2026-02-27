package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// generateTestKeys generates RSA key pair for testing
func generateTestKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "Failed to generate test RSA key")
	return privateKey, &privateKey.PublicKey
}

// createTestClaims creates sample claims for testing
func createTestClaims() *entity.Claims {
	now := time.Now()
	return &entity.Claims{
		Issuer:    "authorizer",
		Subject:   "user-123",
		Audience:  []string{"APP1", "APP2"},
		ExpiresAt: now.Add(time.Hour).Unix(),
		IssuedAt:  now.Unix(),
		Username:  "testuser",
		Email:     "test@example.com",
		Authorization: []entity.Authorization{
			{
				App:         "APP1",
				Roles:       []string{"admin", "user"},
				Permissions: []string{"read", "write"},
			},
			{
				App:         "APP2",
				Roles:       []string{"viewer"},
				Permissions: []string{"read"},
			},
		},
	}
}

func TestNewJWTService(t *testing.T) {
	log := logger.New()
	service := NewJWTService(log)
	
	assert.NotNil(t, service, "JWT service should not be nil")
}

func TestGenerateToken_Success(t *testing.T) {
	// Setup
	log := logger.New()
	service := NewJWTService(log)
	privateKey, _ := generateTestKeys(t)
	claims := createTestClaims()
	ctx := context.Background()
	
	// Execute
	token, err := service.GenerateToken(ctx, claims, privateKey, "test-key-id")
	
	// Verify
	assert.NoError(t, err, "GenerateToken should not return error")
	assert.NotEmpty(t, token, "Token should not be empty")
	assert.Contains(t, token, ".", "Token should be in JWT format (contains dots)")
}

func TestGenerateToken_NilClaims(t *testing.T) {
	// Setup
	log := logger.New()
	service := NewJWTService(log)
	privateKey, _ := generateTestKeys(t)
	ctx := context.Background()
	
	// Execute
	token, err := service.GenerateToken(ctx, nil, privateKey, "test-key-id")
	
	// Verify
	assert.Error(t, err, "GenerateToken should return error for nil claims")
	assert.Empty(t, token, "Token should be empty on error")
	assert.Contains(t, err.Error(), "claims cannot be nil")
}

func TestGenerateToken_NilPrivateKey(t *testing.T) {
	// Setup
	log := logger.New()
	service := NewJWTService(log)
	claims := createTestClaims()
	ctx := context.Background()
	
	// Execute
	token, err := service.GenerateToken(ctx, claims, nil, "test-key-id")
	
	// Verify
	assert.Error(t, err, "GenerateToken should return error for nil private key")
	assert.Empty(t, token, "Token should be empty on error")
	assert.Contains(t, err.Error(), "private key cannot be nil")
}

func TestValidateToken_Success(t *testing.T) {
	// Setup
	log := logger.New()
	service := NewJWTService(log)
	privateKey, publicKey := generateTestKeys(t)
	originalClaims := createTestClaims()
	ctx := context.Background()
	
	// Generate token
	token, err := service.GenerateToken(ctx, originalClaims, privateKey, "test-key-id")
	require.NoError(t, err, "Failed to generate token for test")
	
	// Execute
	validatedClaims, err := service.ValidateToken(ctx, token, publicKey)
	
	// Verify
	assert.NoError(t, err, "ValidateToken should not return error")
	assert.NotNil(t, validatedClaims, "Validated claims should not be nil")
	assert.Equal(t, originalClaims.Issuer, validatedClaims.Issuer)
	assert.Equal(t, originalClaims.Subject, validatedClaims.Subject)
	assert.Equal(t, originalClaims.Username, validatedClaims.Username)
	assert.Equal(t, originalClaims.Email, validatedClaims.Email)
	assert.Equal(t, originalClaims.Audience, validatedClaims.Audience)
	assert.Equal(t, len(originalClaims.Authorization), len(validatedClaims.Authorization))
}

func TestValidateToken_EmptyToken(t *testing.T) {
	// Setup
	log := logger.New()
	service := NewJWTService(log)
	_, publicKey := generateTestKeys(t)
	ctx := context.Background()
	
	// Execute
	claims, err := service.ValidateToken(ctx, "", publicKey)
	
	// Verify
	assert.Error(t, err, "ValidateToken should return error for empty token")
	assert.Nil(t, claims, "Claims should be nil on error")
	assert.Contains(t, err.Error(), "token string cannot be empty")
}

func TestValidateToken_NilPublicKey(t *testing.T) {
	// Setup
	log := logger.New()
	service := NewJWTService(log)
	ctx := context.Background()
	
	// Execute
	claims, err := service.ValidateToken(ctx, "some.token.string", nil)
	
	// Verify
	assert.Error(t, err, "ValidateToken should return error for nil public key")
	assert.Nil(t, claims, "Claims should be nil on error")
	assert.Contains(t, err.Error(), "public key cannot be nil")
}

func TestValidateToken_InvalidToken(t *testing.T) {
	// Setup
	log := logger.New()
	service := NewJWTService(log)
	_, publicKey := generateTestKeys(t)
	ctx := context.Background()
	
	// Execute
	claims, err := service.ValidateToken(ctx, "invalid.token.string", publicKey)
	
	// Verify
	assert.Error(t, err, "ValidateToken should return error for invalid token")
	assert.Nil(t, claims, "Claims should be nil on error")
}

func TestValidateToken_WrongPublicKey(t *testing.T) {
	// Setup
	log := logger.New()
	service := NewJWTService(log)
	privateKey1, _ := generateTestKeys(t)
	_, publicKey2 := generateTestKeys(t) // Different key pair
	originalClaims := createTestClaims()
	ctx := context.Background()
	
	// Generate token with first key
	token, err := service.GenerateToken(ctx, originalClaims, privateKey1, "test-key-id")
	require.NoError(t, err, "Failed to generate token for test")
	
	// Execute - try to validate with different public key
	claims, err := service.ValidateToken(ctx, token, publicKey2)
	
	// Verify
	assert.Error(t, err, "ValidateToken should return error for wrong public key")
	assert.Nil(t, claims, "Claims should be nil on error")
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	// Setup
	log := logger.New()
	service := NewJWTService(log)
	privateKey, publicKey := generateTestKeys(t)
	ctx := context.Background()
	
	// Create expired claims
	now := time.Now()
	expiredClaims := &entity.Claims{
		Issuer:    "authorizer",
		Subject:   "user-123",
		Audience:  []string{"APP1"},
		ExpiresAt: now.Add(-time.Hour).Unix(), // Expired 1 hour ago
		IssuedAt:  now.Add(-2 * time.Hour).Unix(),
		Username:  "testuser",
		Email:     "test@example.com",
		Authorization: []entity.Authorization{
			{
				App:         "APP1",
				Roles:       []string{"user"},
				Permissions: []string{"read"},
			},
		},
	}
	
	// Generate token with expired claims
	token, err := service.GenerateToken(ctx, expiredClaims, privateKey, "test-key-id")
	require.NoError(t, err, "Failed to generate token for test")
	
	// Execute
	claims, err := service.ValidateToken(ctx, token, publicKey)
	
	// Verify
	assert.Error(t, err, "ValidateToken should return error for expired token")
	assert.Nil(t, claims, "Claims should be nil on error")
}

func TestGenerateAndValidate_RoundTrip(t *testing.T) {
	// Setup
	log := logger.New()
	service := NewJWTService(log)
	privateKey, publicKey := generateTestKeys(t)
	ctx := context.Background()
	
	// Test with various claim configurations
	testCases := []struct {
		name   string
		claims *entity.Claims
	}{
		{
			name: "Single app with multiple roles",
			claims: &entity.Claims{
				Issuer:    "authorizer",
				Subject:   "user-1",
				Audience:  []string{"APP1"},
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
				IssuedAt:  time.Now().Unix(),
				Username:  "user1",
				Email:     "user1@example.com",
				Authorization: []entity.Authorization{
					{
						App:         "APP1",
						Roles:       []string{"admin", "user", "moderator"},
						Permissions: []string{"read", "write", "delete"},
					},
				},
			},
		},
		{
			name: "Multiple apps",
			claims: &entity.Claims{
				Issuer:    "authorizer",
				Subject:   "user-2",
				Audience:  []string{"APP1", "APP2", "APP3"},
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
				IssuedAt:  time.Now().Unix(),
				Username:  "user2",
				Email:     "user2@example.com",
				Authorization: []entity.Authorization{
					{
						App:         "APP1",
						Roles:       []string{"user"},
						Permissions: []string{"read"},
					},
					{
						App:         "APP2",
						Roles:       []string{"admin"},
						Permissions: []string{"read", "write"},
					},
					{
						App:         "APP3",
						Roles:       []string{"viewer"},
						Permissions: []string{"read"},
					},
				},
			},
		},
		{
			name: "Global roles",
			claims: &entity.Claims{
				Issuer:    "authorizer",
				Subject:   "user-3",
				Audience:  []string{"GLOBAL"},
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
				IssuedAt:  time.Now().Unix(),
				Username:  "superadmin",
				Email:     "superadmin@example.com",
				Authorization: []entity.Authorization{
					{
						App:         "GLOBAL",
						Roles:       []string{"superadmin"},
						Permissions: []string{"*"},
					},
				},
			},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate token
			token, err := service.GenerateToken(ctx, tc.claims, privateKey, "test-key-id")
			require.NoError(t, err, "Failed to generate token")
			require.NotEmpty(t, token, "Token should not be empty")
			
			// Validate token
			validatedClaims, err := service.ValidateToken(ctx, token, publicKey)
			require.NoError(t, err, "Failed to validate token")
			require.NotNil(t, validatedClaims, "Validated claims should not be nil")
			
			// Verify all fields match
			assert.Equal(t, tc.claims.Issuer, validatedClaims.Issuer)
			assert.Equal(t, tc.claims.Subject, validatedClaims.Subject)
			assert.Equal(t, tc.claims.Username, validatedClaims.Username)
			assert.Equal(t, tc.claims.Email, validatedClaims.Email)
			assert.Equal(t, tc.claims.Audience, validatedClaims.Audience)
			assert.Equal(t, tc.claims.ExpiresAt, validatedClaims.ExpiresAt)
			assert.Equal(t, tc.claims.IssuedAt, validatedClaims.IssuedAt)
			assert.Equal(t, len(tc.claims.Authorization), len(validatedClaims.Authorization))
			
			// Verify authorization details
			for i, auth := range tc.claims.Authorization {
				assert.Equal(t, auth.App, validatedClaims.Authorization[i].App)
				assert.Equal(t, auth.Roles, validatedClaims.Authorization[i].Roles)
				assert.Equal(t, auth.Permissions, validatedClaims.Authorization[i].Permissions)
			}
		})
	}
}
