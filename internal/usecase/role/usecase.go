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

func (u *roleUsecase) Create(name, description, application string) error {

	if name == "" || description == "" {
		return errors.New("role name and description is required")
	}

	existing, _ := u.repo.GetByName(name)
	if existing != nil {
		return errors.New("role already exists")
	}

	newRole := &entity.Role{
		ID:          idgen.NewUUIDv7(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return u.repo.Create(newRole)
}
