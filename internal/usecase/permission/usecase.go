package permission

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	"github.com/mafzaidi/authorizer/pkg/idgen"
)

type permUsecase struct {
	permRepo repository.PermRepository
	appRepo  repository.AppRepository
}

func NewPermUsecase(
	permRepo repository.PermRepository,
	appRepo repository.AppRepository,
) Usecase {
	return &permUsecase{
		permRepo: permRepo,
		appRepo:  appRepo,
	}
}

func (uc *permUsecase) Create(ctx context.Context, in *CreateInput) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if in.AppID == "" {
		return errors.New("application ID is required")
	}

	if in.Code == "" {
		return errors.New("resource is required")
	}

	app, err := uc.appRepo.GetByID(ctx, in.AppID)
	if err != nil {
		return fmt.Errorf("failed: %w", err)
	}

	existingPerm, _ := uc.permRepo.GetByAppAndCode(ctx, app.ID, in.Code)
	if existingPerm != nil {
		return errors.New("permission already exists")
	}

	perm := &entity.Permission{
		ID:            idgen.NewUUIDv7(),
		ApplicationID: &app.ID,
		Code:          in.Code,
		Description:   &in.Description,
		Version:       1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return uc.permRepo.Create(ctx, perm)
}

func (uc *permUsecase) SyncPermissions(ctx context.Context, in *SyncInput) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if in.AppCode == "" {
		return errors.New("application code is required")
	}

	app, err := uc.appRepo.GetByCode(ctx, in.AppCode)
	if err != nil {
		return fmt.Errorf("failed: %w", err)
	}

	var perms []*entity.Permission
	for _, v := range in.Permissions {
		perm := &entity.Permission{
			ID:            idgen.NewUUIDv7(),
			ApplicationID: &app.ID,
			Code:          v.Code,
			Description:   &v.Description,
			Version:       1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		perms = append(perms, perm)
	}

	return uc.permRepo.BulkUpsert(ctx, perms)
}
