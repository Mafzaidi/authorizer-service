package role

import (
	"errors"
	"time"

	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
	"localdev.me/authorizer/pkg/idgen"
)

type roleUsecase struct {
	repo repository.RoleRepository
}

func NewRoleUsecase(repo repository.RoleRepository) Usecase {
	return &roleUsecase{
		repo: repo,
	}
}

func (u *roleUsecase) Create(input *CreateInput) error {

	if input.AppID == "" {
		return errors.New("app ID is required")
	}

	if input.Code == "" || input.Name == "" {
		return errors.New("code and name is required")
	}

	existing, _ := u.repo.GetByAppAndCode(input.AppID, input.Code)
	if existing != nil {
		return errors.New("role already exists")
	}

	newRole := &entity.Role{
		ID:            idgen.NewUUIDv7(),
		ApplicationID: input.AppID,
		Code:          input.Code,
		Name:          input.Name,
		Description:   input.Description,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return u.repo.Create(newRole)
}
