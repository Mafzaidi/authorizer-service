package repository

import "localdev.me/authorizer/internal/domain/entity"

type RoleRepository interface {
	Create(role *entity.Role) error
	Update(role *entity.Role) error
	Delete(id string) error
	GetByID(id string) (*entity.Role, error)
	GetByName(name string) (*entity.Role, error)
	List(limit, offset int) ([]*entity.Role, error)
}
