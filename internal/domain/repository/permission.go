package repository

import "localdev.me/authorizer/internal/domain/entity"

type PermRepository interface {
	Create(app *entity.Permission) error
	GetByID(id string) (*entity.Permission, error)
	GetByCode(code string) (*entity.Permission, error)
	GetByApp(appID string) (*entity.Permission, error)
	Update(app *entity.Permission) error
	Upsert(perm *entity.Permission) error
	Delete(id string) error
	List(limit, offset int) ([]*entity.Permission, error)
}
