package middleware

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/service"
	"github.com/mafzaidi/authorizer/internal/infrastructure/auth"
	"github.com/mafzaidi/authorizer/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockLogger implements service.Logger for testing
type mockLogger struct {
	lastMessage string
	lastFields  service.Fields
}

func (m *mockLogger) Info(message string, fields service.Fields) {
	m.lastMessage = message
	m.lastFields = fields
}

func (m *mockLogger) Warn(message string, fields service.Fields) {
	m.lastMessage = message
	m.lastFields = fields
}

func (m *mockLogger) Error(message string, fields service.Fields) {
	m.lastMessage = message
	m.lastFields = fields
}

func TestJWTAuthMiddleware_MissingToken(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Generate test keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	cfg := &config.Config{
		JWT: &config.JWT{
			PrivateKey: privateKey,
			PublicKey:  &privateKey.PublicKey,
			KeyID:      "test-key-id",
		},
	}

	logger := &mockLogger{}
	jwtService := auth.NewJWTService(logger)

	// Create middleware
	middleware := JWTAuthMiddleware(jwtService, cfg, logger)
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err = handler(c)

	// Assert - ErrorHandler returns nil but sets HTTP status
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, logger.lastMessage, "missing token")
}

func TestJWTAuthMiddleware_InvalidTokenFormat(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Generate test keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	cfg := &config.Config{
		JWT: &config.JWT{
			PrivateKey: privateKey,
			PublicKey:  &privateKey.PublicKey,
			KeyID:      "test-key-id",
		},
	}

	logger := &mockLogger{}
	jwtService := auth.NewJWTService(logger)

	// Create middleware
	middleware := JWTAuthMiddleware(jwtService, cfg, logger)
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err = handler(c)

	// Assert - ErrorHandler returns nil but sets HTTP status
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, logger.lastMessage, "invalid token format")
}

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
	// Setup
	e := echo.New()

	// Generate test keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	cfg := &config.Config{
		JWT: &config.JWT{
			PrivateKey: privateKey,
			PublicKey:  &privateKey.PublicKey,
			KeyID:      "test-key-id",
		},
	}

	logger := &mockLogger{}
	jwtService := auth.NewJWTService(logger)

	// Create valid token
	claims := &entity.Claims{
		Issuer:    "test-issuer",
		Subject:   "user-123",
		Audience:  []string{"test-app"},
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
		Username:  "testuser",
		Email:     "test@example.com",
		Authorization: []entity.Authorization{
			{
				App:         "test-app",
				Roles:       []string{"admin"},
				Permissions: []string{"read", "write"},
			},
		},
	}

	token, err := jwtService.GenerateToken(context.Background(), claims, privateKey, cfg.JWT.KeyID)
	require.NoError(t, err)

	// Create request with valid token
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware
	middleware := JWTAuthMiddleware(jwtService, cfg, logger)
	handler := middleware(func(c echo.Context) error {
		// Verify claims are set in context
		jwtClaims := GetUserFromContext(c)
		assert.NotNil(t, jwtClaims)
		assert.Equal(t, "user-123", jwtClaims.UserID)
		assert.Equal(t, "testuser", jwtClaims.Username)
		assert.Equal(t, "test@example.com", jwtClaims.Email)
		assert.Len(t, jwtClaims.Authorization, 1)
		assert.Equal(t, "test-app", jwtClaims.Authorization[0].App)
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err = handler(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestJWTAuthMiddleware_InvalidToken(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Generate test keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	cfg := &config.Config{
		JWT: &config.JWT{
			PrivateKey: privateKey,
			PublicKey:  &privateKey.PublicKey,
			KeyID:      "test-key-id",
		},
	}

	logger := &mockLogger{}
	jwtService := auth.NewJWTService(logger)

	// Create middleware
	middleware := JWTAuthMiddleware(jwtService, cfg, logger)
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err = handler(c)

	// Assert - ErrorHandler returns nil but sets HTTP status
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, logger.lastMessage, "token validation error")
}

func TestRequirePermission_WithPermission(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set claims in context (simulating JWT middleware)
	claims := &JWTClaims{
		UserID:   "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Authorization: []Authorization{
			{
				App:         "test-app",
				Roles:       []string{"admin"},
				Permissions: []string{"read", "write", "delete"},
			},
		},
	}
	c.Set(string(userContextKey), claims)

	// Create middleware
	middleware := RequirePermission("test-app", "write")
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequirePermission_WithoutPermission(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set claims in context with different permissions
	claims := &JWTClaims{
		UserID:   "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Authorization: []Authorization{
			{
				App:         "test-app",
				Roles:       []string{"viewer"},
				Permissions: []string{"read"},
			},
		},
	}
	c.Set(string(userContextKey), claims)

	// Create middleware requiring write permission
	middleware := RequirePermission("test-app", "write")
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert - ErrorHandler returns nil but sets HTTP status
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRequirePermission_GlobalPermission(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set claims with GLOBAL app (should have access to everything)
	claims := &JWTClaims{
		UserID:   "admin-123",
		Username: "admin",
		Email:    "admin@example.com",
		Authorization: []Authorization{
			{
				App:         "GLOBAL",
				Roles:       []string{"superadmin"},
				Permissions: []string{},
			},
		},
	}
	c.Set(string(userContextKey), claims)

	// Create middleware
	middleware := RequirePermission("test-app", "write")
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequirePermission_MissingClaims(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Don't set claims in context (simulating missing JWT middleware)

	// Create middleware
	middleware := RequirePermission("test-app", "write")
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert - ErrorHandler returns nil but sets HTTP status
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRequirePermission_WrongApp(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set claims for different app
	claims := &JWTClaims{
		UserID:   "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Authorization: []Authorization{
			{
				App:         "other-app",
				Roles:       []string{"admin"},
				Permissions: []string{"read", "write"},
			},
		},
	}
	c.Set(string(userContextKey), claims)

	// Create middleware for test-app
	middleware := RequirePermission("test-app", "write")
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert - ErrorHandler returns nil but sets HTTP status
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name     string
		claims   *JWTClaims
		app      string
		perm     string
		expected bool
	}{
		{
			name: "has permission",
			claims: &JWTClaims{
				Authorization: []Authorization{
					{
						App:         "test-app",
						Permissions: []string{"read", "write"},
					},
				},
			},
			app:      "test-app",
			perm:     "write",
			expected: true,
		},
		{
			name: "does not have permission",
			claims: &JWTClaims{
				Authorization: []Authorization{
					{
						App:         "test-app",
						Permissions: []string{"read"},
					},
				},
			},
			app:      "test-app",
			perm:     "write",
			expected: false,
		},
		{
			name: "global permission",
			claims: &JWTClaims{
				Authorization: []Authorization{
					{
						App:         "GLOBAL",
						Permissions: []string{},
					},
				},
			},
			app:      "test-app",
			perm:     "write",
			expected: true,
		},
		{
			name: "wrong app",
			claims: &JWTClaims{
				Authorization: []Authorization{
					{
						App:         "other-app",
						Permissions: []string{"read", "write"},
					},
				},
			},
			app:      "test-app",
			perm:     "write",
			expected: false,
		},
		{
			name:     "nil claims",
			claims:   nil,
			app:      "test-app",
			perm:     "write",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasPermission(tt.claims, tt.app, tt.perm)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasRole(t *testing.T) {
	tests := []struct {
		name      string
		required  string
		userRoles []string
		expected  bool
	}{
		{
			name:      "has role",
			required:  "admin",
			userRoles: []string{"admin", "user"},
			expected:  true,
		},
		{
			name:      "does not have role",
			required:  "admin",
			userRoles: []string{"user", "viewer"},
			expected:  false,
		},
		{
			name:      "empty roles",
			required:  "admin",
			userRoles: []string{},
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasRole(tt.required, tt.userRoles)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetUserFromContext(t *testing.T) {
	t.Run("claims exist", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expectedClaims := &JWTClaims{
			UserID:   "user-123",
			Username: "testuser",
			Email:    "test@example.com",
		}
		c.Set(string(userContextKey), expectedClaims)

		claims := GetUserFromContext(c)
		assert.NotNil(t, claims)
		assert.Equal(t, expectedClaims.UserID, claims.UserID)
		assert.Equal(t, expectedClaims.Username, claims.Username)
		assert.Equal(t, expectedClaims.Email, claims.Email)
	})

	t.Run("claims do not exist", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		claims := GetUserFromContext(c)
		assert.Nil(t, claims)
	})

	t.Run("wrong type in context", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.Set(string(userContextKey), "not a claims object")

		claims := GetUserFromContext(c)
		assert.Nil(t, claims)
	})
}
