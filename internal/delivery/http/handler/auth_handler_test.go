package handler

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/internal/delivery/http/middleware"
	"github.com/mafzaidi/authorizer/internal/infrastructure/config"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/infrastructure/auth"
	"github.com/mafzaidi/authorizer/internal/infrastructure/logger"
	authUsecase "github.com/mafzaidi/authorizer/internal/usecase/auth"
)

// MockAuthUseCase is a mock implementation of auth.Usecase
type MockAuthUseCase struct {
	LoginFunc        func(ctx context.Context, appCode, email, password, validToken string, cfg *config.Config) (*authUsecase.UserToken, error)
	RefreshTokenFunc func(ctx context.Context, refreshToken string, cfg *config.Config) (string, string, error)
}

func (m *MockAuthUseCase) Login(ctx context.Context, appCode, email, password, validToken string, cfg *config.Config) (*authUsecase.UserToken, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, appCode, email, password, validToken, cfg)
	}
	return nil, errors.New("not implemented")
}

func (m *MockAuthUseCase) RefreshToken(ctx context.Context, refreshToken string, cfg *config.Config) (string, string, error) {
	if m.RefreshTokenFunc != nil {
		return m.RefreshTokenFunc(ctx, refreshToken, cfg)
	}
	return "", "", errors.New("not implemented")
}

// MockJWKSService is a mock implementation of auth.JWKSService
type MockJWKSService struct {
	GetJWKSFunc func(publicKey *rsa.PublicKey, keyID string) (*auth.JWKSResponse, error)
}

func (m *MockJWKSService) GetJWKS(publicKey *rsa.PublicKey, keyID string) (*auth.JWKSResponse, error) {
	if m.GetJWKSFunc != nil {
		return m.GetJWKSFunc(publicKey, keyID)
	}
	return nil, errors.New("not implemented")
}

func TestNewAuthHandler(t *testing.T) {
	mockAuthUC := &MockAuthUseCase{}
	mockJWKSService := &MockJWKSService{}
	cfg := &config.Config{}
	log := logger.New()

	handler := NewAuthHandler(mockAuthUC, mockJWKSService, cfg, log)

	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}

	if handler.cfg != cfg {
		t.Error("Expected cfg to be set correctly")
	}

	if handler.logger != log {
		t.Error("Expected logger to be set correctly")
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	// Setup
	mockAuthUC := &MockAuthUseCase{
		LoginFunc: func(ctx context.Context, appCode, email, password, validToken string, cfg *config.Config) (*authUsecase.UserToken, error) {
			return &authUsecase.UserToken{
				User: &entity.User{
					ID:       "user-123",
					Email:    "test@example.com",
					FullName: "Test User",
				},
				Token:        "test-token",
				RefreshToken: "refresh-token",
				Claims: &middleware.JWTClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						Subject:   "user-123",
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
					},
					Username: "testuser",
					Email:    "test@example.com",
					Authorization: []middleware.Authorization{
						{
							App:         "APP1",
							Roles:       []string{"admin"},
							Permissions: []string{"read", "write"},
						},
					},
				},
			}, nil
		},
	}
	mockJWKSService := &MockJWKSService{}
	cfg := &config.Config{}
	log := logger.New()

	handler := NewAuthHandler(mockAuthUC, mockJWKSService, cfg, log)

	// Create request
	loginReq := LoginRequest{
		Application: "APP1",
		Email:       "test@example.com",
		Password:    "password123",
	}
	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.Login()(c)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Verify response body is not empty
	if rec.Body.Len() == 0 {
		t.Error("Expected response body, got empty")
	}
}

func TestAuthHandler_Login_InvalidRequest(t *testing.T) {
	// Setup
	mockAuthUC := &MockAuthUseCase{}
	mockJWKSService := &MockJWKSService{}
	cfg := &config.Config{}
	log := logger.New()

	handler := NewAuthHandler(mockAuthUC, mockJWKSService, cfg, log)

	// Create invalid request (malformed JSON)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.Login()(c)

	// Assert - response.ErrorHandler returns nil (it writes to response recorder)
	// Check the status code instead
	if err != nil {
		t.Fatalf("Expected no error from handler, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for invalid JSON, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	// Setup
	mockAuthUC := &MockAuthUseCase{}
	mockJWKSService := &MockJWKSService{}
	cfg := &config.Config{}
	log := logger.New()

	handler := NewAuthHandler(mockAuthUC, mockJWKSService, cfg, log)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.Logout()(c)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Verify cookie is set to expire
	cookies := rec.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "jwt_user_token" {
			found = true
			if cookie.MaxAge != -1 {
				t.Error("Expected cookie MaxAge to be -1 (expired)")
			}
			if cookie.Value != "" {
				t.Error("Expected cookie value to be empty")
			}
		}
	}

	if !found {
		t.Error("Expected jwt_user_token cookie to be set")
	}
}

func TestAuthHandler_GetJWKS_Success(t *testing.T) {
	// Setup
	mockAuthUC := &MockAuthUseCase{}
	mockJWKSService := &MockJWKSService{
		GetJWKSFunc: func(publicKey *rsa.PublicKey, keyID string) (*auth.JWKSResponse, error) {
			return &auth.JWKSResponse{
				Keys: []auth.JWK{
					{
						Kty: "RSA",
						Use: "sig",
						Alg: "RS256",
						Kid: "test-key-id",
						N:   "test-n",
						E:   "test-e",
					},
				},
			}, nil
		},
	}
	cfg := &config.Config{
		JWT: &config.JWT{
			PublicKey: &rsa.PublicKey{},
			KeyID:     "test-key",
		},
	}
	log := logger.New()

	handler := NewAuthHandler(mockAuthUC, mockJWKSService, cfg, log)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.GetJWKS()(c)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Verify response structure
	var response JWKSResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.Keys) != 1 {
		t.Errorf("Expected 1 key, got %d", len(response.Keys))
	}

	if response.Keys[0].Kty != "RSA" {
		t.Errorf("Expected Kty 'RSA', got %s", response.Keys[0].Kty)
	}
}

func TestAuthHandler_GetJWKS_ServiceError(t *testing.T) {
	// Setup
	mockAuthUC := &MockAuthUseCase{}
	mockJWKSService := &MockJWKSService{
		GetJWKSFunc: func(publicKey *rsa.PublicKey, keyID string) (*auth.JWKSResponse, error) {
			return nil, errors.New("service error")
		},
	}
	cfg := &config.Config{
		JWT: &config.JWT{
			PublicKey: &rsa.PublicKey{},
			KeyID:     "test-key",
		},
	}
	log := logger.New()

	handler := NewAuthHandler(mockAuthUC, mockJWKSService, cfg, log)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.GetJWKS()(c)

	// Assert - response.ErrorHandler returns nil (it writes to response recorder)
	if err != nil {
		t.Fatalf("Expected no error from handler, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d when JWKS service fails, got %d", http.StatusInternalServerError, rec.Code)
	}
}
