package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mafzaidi/authorizer/internal/delivery/http/middleware"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
)

type JWTService interface {
	GenerateAccessToken(ctx context.Context, userID, appCode, validToken string, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, keyID string) (string, *middleware.JWTClaims, error)
	GenerateRefreshToken() (string, error)
}

type JWTpayload struct {
	PrivateKey *rsa.PrivateKey
	Claims     *middleware.JWTClaims
	KeyID      string
}

type UserToken struct {
	User         *entity.User
	Token        string
	RefreshToken string
	Claims       *middleware.JWTClaims
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

func (s *jwtService) GenerateAccessToken(
	ctx context.Context,
	userID,
	appCode,
	validToken string,
	privateKey *rsa.PrivateKey,
	publicKey *rsa.PublicKey,
	keyID string,
) (string, *middleware.JWTClaims, error) {

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", nil, err
	}

	var existingClaims *middleware.JWTClaims
	if validToken != "" {
		existingClaims, _ = s.validateAccessToken(validToken, publicKey)
	}

	if existingClaims != nil && existingClaims.Email == user.Email {
		return validToken, existingClaims, nil
	}

	var authorizations []middleware.Authorization
	var audiences []string

	globalRoles, _ := s.userRoleRepo.GetGlobalRolesByUser(ctx, user.ID)
	if len(globalRoles) > 0 {
		var roles []string
		for _, r := range globalRoles {
			roles = append(roles, r.Code)
		}

		authorizations = append(authorizations, middleware.Authorization{
			App:         "GLOBAL",
			Roles:       roles,
			Permissions: []string{"*"},
		})

		audiences = append(audiences, "GLOBAL")
	}

	apps, err := s.resolveApps(ctx, appCode)
	if err != nil {
		return "", nil, err
	}

	for _, app := range apps {
		appRoles, _ := s.userRoleRepo.GetRolesByUserAndApp(ctx, userID, app.ID)
		if len(appRoles) == 0 {
			continue
		}

		roleSet := make(map[string]struct{})
		permSet := make(map[string]struct{})

		for _, r := range appRoles {
			roleSet[r.Code] = struct{}{}

			perms, _ := s.rolePermRepo.GetPermsByRole(ctx, r.ID)
			for _, p := range perms {
				permSet[p.Code] = struct{}{}
			}
		}

		authorizations = append(authorizations, middleware.Authorization{
			App:         app.Code,
			Roles:       mapKeys(roleSet),
			Permissions: mapKeys(permSet),
		})

		audiences = append(audiences, app.Code)
	}

	claims := &middleware.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "authorizer",
			Subject:   userID,
			Audience:  audiences,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Username:      user.Username,
		Email:         user.Email,
		Authorization: authorizations,
	}

	jwtPayload := &JWTpayload{
		PrivateKey: privateKey,
		Claims:     claims,
		KeyID:      keyID,
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
		jwt.SigningMethodRS256,
		pl.Claims,
	)

	// Set kid in token header for JWKS key identification
	token.Header["kid"] = pl.KeyID

	tokenString, err := token.SignedString(pl.PrivateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *jwtService) validateAccessToken(cookie string, publicKey *rsa.PublicKey) (*middleware.JWTClaims, error) {

	token, err := jwt.ParseWithClaims(cookie, &middleware.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*middleware.JWTClaims)
	if !ok {
		return claims, errors.New("token invalid")
	}

	if !token.Valid {
		return claims, errors.New("token is not valid")
	}
	return claims, nil
}

func (s *jwtService) resolveApps(ctx context.Context, appCode string) ([]*entity.Application, error) {
	if appCode == "" {
		return s.appRepo.GetAll(ctx)
	}

	app, err := s.appRepo.GetByCode(ctx, appCode)
	if err != nil {
		return nil, err
	}
	return []*entity.Application{app}, nil
}

func mapKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
