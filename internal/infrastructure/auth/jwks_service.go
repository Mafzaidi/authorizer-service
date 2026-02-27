package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
)

// JWKSService handles JWKS (JSON Web Key Set) operations
// This service is responsible for converting RSA public keys to JWKS format
// for public key distribution to JWT consumers.
type JWKSService interface {
	// GetJWKS converts an RSA public key to JWKS format
	// Parameters:
	//   - publicKey: RSA public key to convert
	//   - keyID: key identifier for the JWK
	// Returns:
	//   - *JWKSResponse: the JWKS response containing the public key
	//   - error: if conversion fails
	GetJWKS(publicKey *rsa.PublicKey, keyID string) (*JWKSResponse, error)
}

// JWKSResponse represents a JSON Web Key Set response
// It contains an array of JWK (JSON Web Key) objects
type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key
// It contains the public key information in JWK format
type JWK struct {
	Kty string `json:"kty"` // Key type (e.g., "RSA")
	Use string `json:"use"` // Public key use (e.g., "sig" for signature)
	Alg string `json:"alg"` // Algorithm (e.g., "RS256")
	Kid string `json:"kid"` // Key ID
	N   string `json:"n"`   // RSA modulus (base64url encoded)
	E   string `json:"e"`   // RSA exponent (base64url encoded)
}

type jwksService struct{}

// NewJWKSService creates a new JWKS service instance
func NewJWKSService() JWKSService {
	return &jwksService{}
}

// GetJWKS converts an RSA public key to JWKS format
func (s *jwksService) GetJWKS(publicKey *rsa.PublicKey, keyID string) (*JWKSResponse, error) {
	if publicKey == nil {
		return nil, ErrNilPublicKey
	}

	// Convert RSA modulus (N) to base64url encoding
	nBytes := publicKey.N.Bytes()
	n := base64.RawURLEncoding.EncodeToString(nBytes)

	// Convert RSA exponent (E) to base64url encoding
	eBytes := big.NewInt(int64(publicKey.E)).Bytes()
	e := base64.RawURLEncoding.EncodeToString(eBytes)

	// Create JWK
	jwk := JWK{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		Kid: keyID,
		N:   n,
		E:   e,
	}

	// Create JWKS response with single key
	response := &JWKSResponse{
		Keys: []JWK{jwk},
	}

	return response, nil
}

// ErrNilPublicKey is returned when a nil public key is provided
var ErrNilPublicKey = &jwksError{message: "public key cannot be nil"}

type jwksError struct {
	message string
}

func (e *jwksError) Error() string {
	return e.message
}
