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
	Claims *middleware.JWTClaims
}

type RoleInfo struct {
	ID   string
	Code string
}

type authUsecase struct {
	userRepo     repository.UserRepository
	appRepo      repository.AppRepository
	userRoleRepo repository.UserRoleRepository
	roleRepo     repository.RoleRepository
	rolePermRepo repository.RolePermRepository
	permRepo     repository.PermRepository
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	appRepo repository.AppRepository,
	userRoleRepo repository.UserRoleRepository,
	roleRepo repository.RoleRepository,
	rolePermRepo repository.RolePermRepository,
	permRepo repository.PermRepository,
) Usecase {
	return &authUsecase{
		userRepo:     userRepo,
		appRepo:      appRepo,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
		rolePermRepo: rolePermRepo,
		permRepo:     permRepo,
	}
}

func (u *authUsecase) Login(application, email, password, validToken string, cfg *config.Config) (*UserToken, error) {
	if application == "" {
		return nil, errors.New("application cannot be empty")
	}

	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	user, err := u.userRepo.GetByEmail(email)
	if err != nil || !pwd.CheckHash(user.Password, password) {
		return nil, errors.New("email or password is invalid")
	}
	var roles []RoleInfo

	app, _ := u.appRepo.GetByCode(application)
	if app != nil {
		userRoles, _ := u.userRoleRepo.GetRolesByUserAndApp(user.ID, app.ID)

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
		rolePerms, _ := u.rolePermRepo.GetPermsByRoles(roleIDs)
		for _, p := range rolePerms {
			permissions = append(permissions, p.Code)
		}
	}

	var claims *middleware.JWTClaims
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

	claims = &middleware.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "authorizer-service",
			Subject:   user.ID,
			Audience:  []string{app.Name},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(1) * time.Hour)),
		},
		Username:    user.Username,
		Email:       user.Email,
		Roles:       roleNames,
		Permissions: permissions,
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
