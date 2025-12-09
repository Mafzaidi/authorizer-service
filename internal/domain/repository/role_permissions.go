package repository

import (
	"context"

	"localdev.me/authorizer/internal/domain/entity"
)

type RolePermRepository interface {
	Grant(ctx context.Context, roleID, permID string) error
	Revoke(ctx context.Context, roleID, permID string) error
	Replace(ctx context.Context, roleID string, permIDs []string) error
	GetPermsByRole(ctx context.Context, roleID string) ([]*entity.Permission, error)
	GetPermsByRoles(ctx context.Context, roleIDs []string) ([]*entity.Permission, error)
}
