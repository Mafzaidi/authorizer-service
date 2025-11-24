package repository

import "localdev.me/authorizer/internal/domain/entity"

type ApplicationRepository interface {
	Create(app *entity.Application) error
	GetByID(id string) (*entity.Application, error)
	GetBySlug(slug string) (*entity.Application, error)
	Update(app *entity.Application) error
	Delete(id string) error
	List(limit, offset int) ([]*entity.Application, error)
}
