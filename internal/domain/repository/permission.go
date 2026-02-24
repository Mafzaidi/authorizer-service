package repository

import (
	"context"

	"github.com/mafzaidi/authorizer/internal/domain/entity"
)

type PermRepository interface {
	Create(ctx context.Context, app *entity.Permission) error
	Update(ctx context.Context, app *entity.Permission) error
	Upsert(ctx context.Context, perm *entity.Permission) error
	BulkUpsert(ctx context.Context, perms []*entity.Permission) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entity.Permission, error)
	GetByAppAndCode(ctx context.Context, appID, code string) (*entity.Permission, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Permission, error)
	ListByApp(ctx context.Context, appID string) ([]*entity.Permission, error)
}
