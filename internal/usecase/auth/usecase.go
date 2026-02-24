package auth

import (
	"context"
	"errors"
	"time"

	"github.com/mafzaidi/authorizer/config"
	"github.com/mafzaidi/authorizer/internal/delivery/http/middleware"
	"github.com/mafzaidi/authorizer/internal/delivery/http/middleware/pwd"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	"github.com/mafzaidi/authorizer/internal/service"
)

type UserToken struct {
	User         *entity.User
	Token        string
	RefreshToken string
	Claims       *middleware.JWTClaims
}

type authUsecase struct {
	authRepo repository.AuthRepository
	userRepo repository.UserRepository
	jwtSvc   service.JWTService
}

func NewAuthUseCase(
	authRepo repository.AuthRepository,
	userRepo repository.UserRepository,
	jwtSvc service.JWTService,
) Usecase {
	return &authUsecase{
		authRepo: authRepo,
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
	}
}

func (uc *authUsecase) Login(
	ctx context.Context,
	appCode,
	email,
	password,
	validToken string,
	cfg *config.Config,
) (*UserToken, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil || !pwd.CheckHash(user.Password, password) {
		return nil, errors.New("email or password is invalid")
	}

	accessToken, claims, err := uc.jwtSvc.GenerateAccessToken(
		ctx,
		user.ID,
		appCode,
		validToken,
		cfg.JWT.PrivateKey,
		cfg.JWT.PublicKey,
		cfg.JWT.KeyID,
	)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := uc.jwtSvc.GenerateRefreshToken()
	if err != nil {
		return nil, errors.New("failed generating refresh token")
	}

	err = uc.authRepo.StoreRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		return nil, errors.New("failed saving refresh token")
	}

	token := &UserToken{
		User:         user,
		Token:        accessToken,
		RefreshToken: refreshToken,
		Claims:       claims,
	}

	return token, nil
}

func (uc *authUsecase) RefreshToken(
	ctx context.Context,
	refreshToken string,
	cfg *config.Config,
) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userID, err := uc.authRepo.GetUserIDByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	savedToken, err := uc.authRepo.GetRefreshToken(ctx, userID)
	if err != nil || savedToken != refreshToken {
		return "", "", errors.New("refresh token mismatch")
	}

	return "", "", nil
}
