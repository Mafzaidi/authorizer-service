package user

import (
	"context"

	"github.com/mafzaidi/authorizer/internal/domain/entity"
)

type Usecase interface {
	Register(ctx context.Context, req *RegisterInput) error
	GetDetail(ctx context.Context, userID string) (*entity.User, error)
	UpdateData(ctx context.Context, userID string, input *UpdateInput) error
	GetList(ctx context.Context, limit, offset int) ([]*entity.User, error)
	AssignRoles(ctx context.Context, userID, appID string, roles []string) error
}
