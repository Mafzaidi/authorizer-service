package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"localdev.me/authorizer/config"
	"localdev.me/authorizer/internal/delivery/http/middleware"
	"localdev.me/authorizer/internal/delivery/http/middleware/pwd"
	"localdev.me/authorizer/internal/delivery/http/middleware/token"
	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
)

type UserToken struct {
	User   *entity.User
	Token  string
	Claims *middleware.Claims
}

type AuthUseCase interface {
	Login(email, password, validToken string, conf *config.Config) (*UserToken, error)
}

type authUC struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(userRepo repository.UserRepository) AuthUseCase {
	return &authUC{
		userRepo: userRepo,
	}
}

func (u *authUC) Login(email, password, validToken string, cfg *config.Config) (*UserToken, error) {

	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	user, err := u.userRepo.GetByEmail(email)
	if err != nil || !pwd.CheckHash(user.PasswordHash, password) {
		return nil, errors.New("email or password is invalid")
	}

	var claims *middleware.Claims
	if validToken != "" {
		claims, _ = token.Validate(validToken, cfg.JWT.Secret)
	}

	if claims != nil && claims.Email == email {
		ut := &UserToken{
			User:   user,
			Token:  validToken,
			Claims: claims,
		}
		return ut, nil
	}

	claims = &middleware.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.App.Name,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(1) * time.Hour)),
		},
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	jwtPayload := &token.JWTGen{
		Secret: cfg.JWT.Secret,
		Claims: claims,
	}

	token, err := token.Generate(jwtPayload)

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("failed to generate jwt")
	}

	ut := &UserToken{
		User:   user,
		Token:  token,
		Claims: claims,
	}

	return ut, nil
}
