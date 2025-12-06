package role

import (
	"errors"
	"fmt"
	"time"

	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
	"localdev.me/authorizer/pkg/idgen"
)

type roleUsecase struct {
	roleRepo     repository.RoleRepository
	appRepo      repository.AppRepository
	permRepo     repository.PermRepository
	rolePermRepo repository.RolePermRepository
}

func NewRoleUsecase(
	roleRepo repository.RoleRepository,
	appRepo repository.AppRepository,
	permRepo repository.PermRepository,
	rolePermRepo repository.RolePermRepository,
) Usecase {
	return &roleUsecase{
		roleRepo:     roleRepo,
		appRepo:      appRepo,
		permRepo:     permRepo,
		rolePermRepo: rolePermRepo,
	}
}

func (u *roleUsecase) Create(in *CreateInput) error {

	if in.AppID == "" {
		return errors.New("app ID is required")
	}

	if in.Code == "" || in.Name == "" {
		return errors.New("code and name is required")
	}

	existingRole, _ := u.roleRepo.GetByAppAndCode(in.AppID, in.Code)
	if existingRole != nil {
		return errors.New("role already exists")
	}

	role := &entity.Role{
		ID:            idgen.NewUUIDv7(),
		ApplicationID: in.AppID,
		Code:          in.Code,
		Name:          in.Name,
		Description:   in.Description,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return u.roleRepo.Create(role)
}

func (u *roleUsecase) GrantPerms(roleID string, perms []string) error {

	if roleID == "" {
		return errors.New("roleID is required")
	}

	if len(perms) == 0 {
		return errors.New("permissions is required")
	}

	role, err := u.roleRepo.GetByID(roleID)
	if err != nil {
		return fmt.Errorf("failed: %w", err)
	}

	app, err := u.appRepo.GetByID(role.ApplicationID)
	if err != nil {
		return fmt.Errorf("failed: %w", err)
	}

	var permIDs []string
	for _, v := range perms {
		perm, err := u.permRepo.GetByAppAndCode(app.ID, v)
		if err != nil {
			return fmt.Errorf("failed: %w", err)
		}
		permIDs = append(permIDs, perm.ID)
	}

	return u.rolePermRepo.Replace(role.ID, permIDs)
}
