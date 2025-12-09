package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"localdev.me/authorizer/internal/delivery/http/middleware"
	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
)

type JWTService interface {
	GenerateAccessToken(ctx context.Context, userID, appCode, validToken, secret string) (string, *middleware.JWTClaims, error)
	GenerateRefreshToken() (string, error)
}

type JWTpayload struct {
	Secret string
	Claims *middleware.JWTClaims
}

type UserToken struct {
	User         *entity.User
	Token        string
	RefreshToken string
	Claims       *middleware.JWTClaims
}

type RoleInfo struct {
	ID   string
	Code string
}

type jwtService struct {
	authRepo     repository.AuthRepository
	userRepo     repository.UserRepository
	appRepo      repository.AppRepository
	userRoleRepo repository.UserRoleRepository
	roleRepo     repository.RoleRepository
	rolePermRepo repository.RolePermRepository
	permRepo     repository.PermRepository
}

func NewJWTService(
	authRepo repository.AuthRepository,
	userRepo repository.UserRepository,
	appRepo repository.AppRepository,
	userRoleRepo repository.UserRoleRepository,
	roleRepo repository.RoleRepository,
	rolePermRepo repository.RolePermRepository,
	permRepo repository.PermRepository,
) JWTService {
	return &jwtService{
		authRepo:     authRepo,
		userRepo:     userRepo,
		appRepo:      appRepo,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
		rolePermRepo: rolePermRepo,
		permRepo:     permRepo,
	}
}

func (s *jwtService) GenerateAccessToken(ctx context.Context, userID, appCode, validToken, secret string) (string, *middleware.JWTClaims, error) {

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", nil, err
	}

	var existingClaims *middleware.JWTClaims
	if validToken != "" {
		existingClaims, _ = s.validateAccessToken(validToken, secret)
	}

	if existingClaims != nil && existingClaims.Email == user.Email {
		return validToken, existingClaims, nil
	}

	var roles []RoleInfo

	app, _ := s.appRepo.GetByCode(ctx, appCode)
	if app != nil {
		userRoles, _ := s.userRoleRepo.GetRolesByUserAndApp(ctx, user.ID, app.ID)

		roles = make([]RoleInfo, len(userRoles))
		for i, r := range userRoles {
			roles[i] = RoleInfo{ID: r.ID, Code: r.Code}
		}
	}

	roleIDs := make([]string, len(roles))
	roleNames := make([]string, len(roles))

	for i, r := range roles {
		roleIDs[i] = r.ID
		roleNames[i] = r.Code
	}

	permissions := make([]string, 0)

	if len(roleIDs) > 0 {
		rolePerms, _ := s.rolePermRepo.GetPermsByRoles(ctx, roleIDs)
		for _, p := range rolePerms {
			permissions = append(permissions, p.Code)
		}
	}

	claims := &middleware.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "authorizer-service",
			Subject:  userID,
			Audience: []string{app.Code},
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(1 * time.Hour),
			),
		},
		Username:    user.Username,
		Email:       user.Email,
		Roles:       roleNames,
		Permissions: permissions,
	}

	jwtPayload := &JWTpayload{
		Secret: secret,
		Claims: claims,
	}

	token, err := s.generateJWT(jwtPayload)
	return token, claims, err
}

func (s *jwtService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *jwtService) generateJWT(pl *JWTpayload) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		pl.Claims,
	)
	tokenString, err := token.SignedString([]byte(pl.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *jwtService) validateAccessToken(cookie, secret string) (*middleware.JWTClaims, error) {

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
