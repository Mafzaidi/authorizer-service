package repository

import (
	"context"

	"github.com/mafzaidi/authorizer/internal/domain/entity"
)

type UserRoleRepository interface {
	Assign(ctx context.Context, userID string, roleIDs []string) error
	Unassign(ctx context.Context, userID, roleID string) error
	Replace(ctx context.Context, userID string, roleIDs []string) error
	GetRolesByUser(ctx context.Context, userID string) ([]*entity.Role, error)
	GetRolesByUserAndApp(ctx context.Context, userID, appID string) ([]*entity.Role, error)
	GetGlobalRolesByUser(ctx context.Context, userID string) ([]*entity.Role, error)
	GetUsersByRole(ctx context.Context, roleID string) ([]*entity.User, error)
}
