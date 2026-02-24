package repository

import (
	"context"

	"github.com/mafzaidi/authorizer/internal/domain/entity"
)

type AppRepository interface {
	Create(ctx context.Context, app *entity.Application) error
	GetByID(ctx context.Context, id string) (*entity.Application, error)
	GetByCode(ctx context.Context, code string) (*entity.Application, error)
	GetAll(ctx context.Context) ([]*entity.Application, error)
	Update(ctx context.Context, app *entity.Application) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Application, error)
}
