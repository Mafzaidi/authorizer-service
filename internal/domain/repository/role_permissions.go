package repository

import "localdev.me/authorizer/internal/domain/entity"

type RolePermRepository interface {
	Grant(roleID, permID string) error
	Revoke(roleID, permID string) error
	Replace(roleID string, permIDs []string) error
	GetPermsByRole(roleID string) ([]*entity.Permission, error)
	GetPermsByRoles(roleIDs []string) ([]*entity.Permission, error)
}
