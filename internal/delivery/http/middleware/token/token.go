package token

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"localdev.me/authorizer/internal/delivery/http/middleware"
)

type JWTGen struct {
	Secret string
	Claims *middleware.Claims
}

func Generate(t *JWTGen) (string, error) {

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

func Validate(cookie, secret string) (*middleware.Claims, error) {
	token, err := jwt.ParseWithClaims(cookie, &middleware.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	claims, ok := token.Claims.(*middleware.Claims)
	if !ok {
		return claims, errors.New("token invalid")
	}

	if err != nil || !token.Valid {
		return claims, err
	}
	return claims, nil
}
