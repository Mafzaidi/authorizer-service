package application

import (
	"context"
	"errors"
	"time"

	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	"github.com/mafzaidi/authorizer/pkg/idgen"
)

type appUsecase struct {
	repo repository.AppRepository
}

func NewAppUsecase(repo repository.AppRepository) Usecase {
	return &appUsecase{
		repo: repo,
	}
}

func (uc *appUsecase) Create(ctx context.Context, in *CreateInput) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if in.Code == "" || in.Name == "" {
		return errors.New("code and name is required")
	}

	existingApp, _ := uc.repo.GetByCode(ctx, in.Code)
	if existingApp != nil {
		return errors.New("application already exists")
	}

	app := &entity.Application{
		ID:          idgen.NewUUIDv7(),
		Code:        in.Code,
		Description: in.Description,
		Name:        in.Name,
		Metadata:    in.Metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return uc.repo.Create(ctx, app)
}
