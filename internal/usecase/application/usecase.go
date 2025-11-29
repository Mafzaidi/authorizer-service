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

func (u *appUsecase) Create(input *CreateInput) error {

	if input.Code == "" || input.Name == "" {
		return errors.New("code and name is required")
	}

	existing, _ := u.repo.GetByCode(input.Code)
	if existing != nil {
		return errors.New("application already exists")
	}

	newApp := &entity.Application{
		ID:          idgen.NewUUIDv7(),
		Code:        input.Code,
		Description: input.Description,
		Name:        input.Name,
		Metadata:    input.Metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return u.repo.Create(newApp)
}
