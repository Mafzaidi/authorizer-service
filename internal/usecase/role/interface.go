package role

import "context"

type Usecase interface {
	Create(ctx context.Context, input *CreateInput) error
	GrantPerms(ctx context.Context, roleID string, perms []string) error
}
