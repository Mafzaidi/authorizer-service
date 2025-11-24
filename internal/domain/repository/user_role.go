package repository

import "localdev.me/authorizer/internal/domain/entity"

type UserRoleRepository interface {
	Assign(userID, roleID string) error
	Unassign(userID, roleID string) error
	GetRolesByUser(userID string) ([]*entity.Role, error)
	GetUsersByRole(roleID string) ([]*entity.User, error)
}
