package repository

import "localdev.me/authorizer/internal/domain/entity"

type AppRepository interface {
	Create(app *entity.Application) error
	GetByID(id string) (*entity.Application, error)
	GetByCode(code string) (*entity.Application, error)
	Update(app *entity.Application) error
	Delete(id string) error
	List(limit, offset int) ([]*entity.Application, error)
}
