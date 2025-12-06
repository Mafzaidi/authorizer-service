package application

import (
	"errors"
	"time"

	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
	"localdev.me/authorizer/pkg/idgen"
)

type appUsecase struct {
	repo repository.AppRepository
}

func NewAppUsecase(repo repository.AppRepository) Usecase {
	return &appUsecase{
		repo: repo,
	}
}

func (u *appUsecase) Create(in *CreateInput) error {

	if in.Code == "" || in.Name == "" {
		return errors.New("code and name is required")
	}

	existingApp, _ := u.repo.GetByCode(in.Code)
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

	return u.repo.Create(app)
}
