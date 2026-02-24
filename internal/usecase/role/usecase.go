package role

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	"github.com/mafzaidi/authorizer/pkg/idgen"
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

func (uc *roleUsecase) Create(ctx context.Context, in *CreateInput) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if in.AppID == "" {
		return errors.New("app ID is required")
	}

	if in.Code == "" || in.Name == "" {
		return errors.New("code and name is required")
	}

	existingRole, _ := uc.roleRepo.GetByAppAndCode(ctx, in.AppID, in.Code)
	if existingRole != nil {
		return errors.New("role already exists")
	}

	role := &entity.Role{
		ID:            idgen.NewUUIDv7(),
		ApplicationID: &in.AppID,
		Code:          in.Code,
		Name:          in.Name,
		Description:   &in.Description,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return uc.roleRepo.Create(ctx, role)
}

func (uc *roleUsecase) GrantPerms(ctx context.Context, roleID string, perms []string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if roleID == "" {
		return errors.New("roleID is required")
	}

	if len(perms) == 0 {
		return errors.New("permissions is required")
	}

	role, err := uc.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed: %w", err)
	}

	app, err := uc.appRepo.GetByID(ctx, *role.ApplicationID)
	if err != nil {
		return fmt.Errorf("failed: %w", err)
	}

	var permIDs []string
	for _, v := range perms {
		perm, err := uc.permRepo.GetByAppAndCode(ctx, app.ID, v)
		if err != nil {
			return fmt.Errorf("failed: %w", err)
		}
		permIDs = append(permIDs, perm.ID)
	}

	return uc.rolePermRepo.Replace(ctx, role.ID, permIDs)
}
