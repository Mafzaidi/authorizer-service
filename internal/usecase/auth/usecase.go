package auth

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

type authUsecase struct {
	userRepo     repository.UserRepository
	userRoleRepo repository.UserRoleRepository
	roleRepo     repository.RoleRepository
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	userRoleRepo repository.UserRoleRepository,
	roleRepo repository.RoleRepository,
) Usecase {
	return &authUsecase{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
	}
}

func (u *authUsecase) Login(email, password, validToken string, cfg *config.Config) (*UserToken, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	user, err := u.userRepo.GetByEmail(email)
	if err != nil || !pwd.CheckHash(user.Password, password) {
		return nil, errors.New("email or password is invalid")
	}

	r, err := u.userRoleRepo.GetRolesByUser(user.ID)
	if err != nil {
		return nil, errors.New("something went wrong in server")
	}

	roles := []string{"user"}

	for _, role := range r {
		roles = append(roles, role.Name)
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
		Roles:    roles,
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
