package repository

import "localdev.me/authorizer/internal/domain/entity"

type PermissionRepository interface {
	Create(app *entity.Permission) error
	GetByID(id string) (*entity.Permission, error)
	GetBySlug(slug string) (*entity.Permission, error)
	GetByApplication(appID string) (*entity.Permission, error)
	Update(app *entity.Permission) error
	Delete(id string) error
	List(limit, offset int) ([]*entity.Permission, error)
}
