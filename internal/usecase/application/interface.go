package application

import "context"

type Usecase interface {
	Create(ctx context.Context, input *CreateInput) error
}
