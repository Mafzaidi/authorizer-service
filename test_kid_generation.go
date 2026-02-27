package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mafzaidi/authorizer/internal/delivery/http/middleware"
	"github.com/mafzaidi/authorizer/internal/infrastructure/config"
)

func main() {
	// Load config to get private key and kid
	cfg := config.GetConfig()

	// Create test claims
	claims := &middleware.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "authorizer",
			Subject:   "test-user-123",
			Audience:  []string{"STACKFORGE"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Username: "testuser",
		Email:    "test@example.com",
		Authorization: []middleware.Authorization{
			{
				App:         "STACKFORGE",
				Roles:       []string{"user"},
				Permissions: []string{"todo.read", "todo.write"},
			},
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Set kid in header
	token.Header["kid"] = cfg.JWT.KeyID

	// Sign token
	tokenString, err := token.SignedString(cfg.JWT.PrivateKey)
	if err != nil {
		fmt.Printf("Error signing token: %v\n", err)
		return
	}

	fmt.Println("=== Token Generation Test ===")
	fmt.Println()
	fmt.Printf("Generated Token:\n%s\n", tokenString)
	fmt.Println()

	// Decode header to verify kid
	parts := strings.Split(tokenString, ".")
	if len(parts) < 2 {
		fmt.Println("Invalid token format")
		return
	}

	// Decode header (add padding if needed)
	header := parts[0]
	for len(header)%4 != 0 {
		header += "="
	}

	// Parse header
	var headerMap map[string]interface{}
	headerBytes, err := jwt.DecodeSegment(parts[0])
	if err != nil {
		fmt.Printf("Error decoding header: %v\n", err)
		return
	}
	err = json.Unmarshal(headerBytes, &headerMap)
	if err != nil {
		fmt.Printf("Error unmarshaling header: %v\n", err)
		return
	}

	fmt.Println("Token Header:")
	headerJSON, _ := json.MarshalIndent(headerMap, "", "  ")
	fmt.Println(string(headerJSON))
	fmt.Println()

	// Check kid
	if kid, ok := headerMap["kid"].(string); ok {
		fmt.Printf("✓ KID found in header: %s\n", kid)
		fmt.Printf("✓ Config KID: %s\n", cfg.JWT.KeyID)
		if kid == cfg.JWT.KeyID {
			fmt.Println("✓ KID matches config!")
		} else {
			fmt.Println("✗ KID does not match config")
		}
	} else {
		fmt.Println("✗ KID not found in header")
	}
}
