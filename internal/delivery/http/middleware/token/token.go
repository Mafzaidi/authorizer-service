package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mafzaidi/authorizer/internal/delivery/http/middleware"
)

type JWTGen struct {
	Secret string
	Claims *middleware.JWTClaims
}

func GenerateAccessToken(t *JWTGen) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		t.Claims,
	)
	tokenString, err := token.SignedString([]byte(t.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateAccessToken(cookie, secret string) (*middleware.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(cookie, &middleware.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	claims, ok := token.Claims.(*middleware.JWTClaims)
	if !ok {
		return claims, errors.New("token invalid")
	}

	if err != nil || !token.Valid {
		return claims, err
	}
	return claims, nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32) // 256-bit
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
