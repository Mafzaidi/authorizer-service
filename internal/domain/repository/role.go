package repository

import (
	"context"

	"github.com/mafzaidi/authorizer/internal/domain/entity"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entity.Role, error)
	GetByAppAndCode(ctx context.Context, appID, code string) (*entity.Role, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Role, error)
	ListByApp(ctx context.Context, appID string) ([]*entity.Role, error)
}
