package service

import (
	"context"
	"time"

	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
)

// AuthService defines the interface for authentication domain service
// This service handles pure business logic for building JWT claims
// from user data and authorization rules
type AuthService interface {
	// BuildClaims constructs JWT claims from user data and authorization rules
	// It queries user roles and permissions for the specified application
	// and builds the authorization array for JWT claims
	//
	// Parameters:
	//   - ctx: context for cancellation and timeout
	//   - user: the authenticated user
	//   - appCode: the application code for which to build claims
	//
	// Returns:
	//   - *entity.Claims: the constructed claims with authorization data
	//   - error: if there's an error querying roles/permissions or building claims
	BuildClaims(ctx context.Context, user *entity.User, appCode string) (*entity.Claims, error)
}

// authService implements the AuthService interface
type authService struct {
	userRoleRepo repository.UserRoleRepository
	roleRepo     repository.RoleRepository
	rolePermRepo repository.RolePermRepository
	appRepo      repository.AppRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(
	userRoleRepo repository.UserRoleRepository,
	roleRepo repository.RoleRepository,
	rolePermRepo repository.RolePermRepository,
	appRepo repository.AppRepository,
) AuthService {
	return &authService{
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
		rolePermRepo: rolePermRepo,
		appRepo:      appRepo,
	}
}

// BuildClaims constructs JWT claims from user data and authorization rules
func (s *authService) BuildClaims(ctx context.Context, user *entity.User, appCode string) (*entity.Claims, error) {
	var authorizations []entity.Authorization
	var audiences []string

	// Handle global roles
	globalRoles, _ := s.userRoleRepo.GetGlobalRolesByUser(ctx, user.ID)
	if len(globalRoles) > 0 {
		var roles []string
		for _, r := range globalRoles {
			roles = append(roles, r.Code)
		}

		authorizations = append(authorizations, entity.Authorization{
			App:         "GLOBAL",
			Roles:       roles,
			Permissions: []string{"*"},
		})

		audiences = append(audiences, "GLOBAL")
	}

	// Resolve applications
	apps, err := s.resolveApps(ctx, appCode)
	if err != nil {
		return nil, err
	}

	// Build authorizations for each app
	for _, app := range apps {
		appRoles, _ := s.userRoleRepo.GetRolesByUserAndApp(ctx, user.ID, app.ID)
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

		authorizations = append(authorizations, entity.Authorization{
			App:         app.Code,
			Roles:       mapKeys(roleSet),
			Permissions: mapKeys(permSet),
		})

		audiences = append(audiences, app.Code)
	}

	// Build claims
	now := time.Now()
	claims := &entity.Claims{
		Issuer:        "authorizer",
		Subject:       user.ID,
		Audience:      audiences,
		IssuedAt:      now.Unix(),
		ExpiresAt:     now.Add(time.Hour).Unix(),
		Username:      user.Username,
		Email:         user.Email,
		Authorization: authorizations,
	}

	return claims, nil
}

// resolveApps resolves the applications based on the appCode
// If appCode is empty, returns all applications
// Otherwise, returns the specific application
func (s *authService) resolveApps(ctx context.Context, appCode string) ([]*entity.Application, error) {
	if appCode == "" {
		return s.appRepo.GetAll(ctx)
	}

	app, err := s.appRepo.GetByCode(ctx, appCode)
	if err != nil {
		return nil, err
	}
	return []*entity.Application{app}, nil
}

// mapKeys extracts keys from a map into a slice
func mapKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
