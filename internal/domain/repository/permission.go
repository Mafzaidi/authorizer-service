package repository

import "localdev.me/authorizer/internal/domain/entity"

type PermRepository interface {
	Create(app *entity.Permission) error
	Update(app *entity.Permission) error
	Upsert(perm *entity.Permission) error
	BulkUpsert(perms []*entity.Permission) error
	Delete(id string) error
	GetByID(id string) (*entity.Permission, error)
	GetByAppAndCode(appID, code string) (*entity.Permission, error)
	List(limit, offset int) ([]*entity.Permission, error)
	ListByApp(appID string) ([]*entity.Permission, error)
}
