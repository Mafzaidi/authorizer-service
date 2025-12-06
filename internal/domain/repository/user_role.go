package repository

import "localdev.me/authorizer/internal/domain/entity"

type UserRoleRepository interface {
	Assign(userID string, roleIDs []string) error
	Unassign(userID, roleID string) error
	Replace(userID string, roleIDs []string) error
	GetRolesByUserAndApp(userID, appID string) ([]*entity.Role, error)
	GetUsersByRole(roleID string) ([]*entity.User, error)
}
