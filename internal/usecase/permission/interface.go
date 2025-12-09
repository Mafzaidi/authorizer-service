package permission

import "context"

type Usecase interface {
	Create(ctx context.Context, input *CreateInput) error
	SyncPermissions(ctx context.Context, i *SyncInput) error
}
