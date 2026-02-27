package permission

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	"github.com/mafzaidi/authorizer/internal/domain/service"
	"github.com/mafzaidi/authorizer/pkg/idgen"
)

type permUsecase struct {
	permRepo repository.PermRepository
	appRepo  repository.AppRepository
	logger   service.Logger
}

func NewPermUsecase(
	permRepo repository.PermRepository,
	appRepo repository.AppRepository,
	logger service.Logger,
) Usecase {
	return &permUsecase{
		permRepo: permRepo,
		appRepo:  appRepo,
		logger:   logger,
	}
}

func (uc *permUsecase) Create(ctx context.Context, in *CreateInput) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if in.AppID == "" {
		uc.logger.Warn("Create permission failed: application ID is required", service.Fields{})
		return errors.New("application ID is required")
	}

	if in.Code == "" {
		uc.logger.Warn("Create permission failed: resource is required", service.Fields{
			"app_id": in.AppID,
		})
		return errors.New("resource is required")
	}

	app, err := uc.appRepo.GetByID(ctx, in.AppID)
	if err != nil {
		uc.logger.Error("Failed to get application by ID", service.Fields{
			"app_id": in.AppID,
			"error":  err.Error(),
		})
		return fmt.Errorf("failed: %w", err)
	}

	existingPerm, _ := uc.permRepo.GetByAppAndCode(ctx, app.ID, in.Code)
	if existingPerm != nil {
		uc.logger.Warn("Create permission failed: permission already exists", service.Fields{
			"app_id": app.ID,
			"code":   in.Code,
		})
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

	err = uc.permRepo.Create(ctx, perm)
	if err != nil {
		uc.logger.Error("Failed to create permission", service.Fields{
			"app_id": app.ID,
			"code":   in.Code,
			"error":  err.Error(),
		})
		return err
	}

	uc.logger.Info("Permission created successfully", service.Fields{
		"permission_id": perm.ID,
		"app_id":        app.ID,
		"code":          in.Code,
	})

	return nil
}

func (uc *permUsecase) SyncPermissions(ctx context.Context, in *SyncInput) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if in.AppCode == "" {
		uc.logger.Warn("Sync permissions failed: application code is required", service.Fields{})
		return errors.New("application code is required")
	}

	app, err := uc.appRepo.GetByCode(ctx, in.AppCode)
	if err != nil {
		uc.logger.Error("Failed to get application by code", service.Fields{
			"app_code": in.AppCode,
			"error":    err.Error(),
		})
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

	err = uc.permRepo.BulkUpsert(ctx, perms)
	if err != nil {
		uc.logger.Error("Failed to sync permissions", service.Fields{
			"app_code":         in.AppCode,
			"app_id":           app.ID,
			"permissions_count": len(perms),
			"error":            err.Error(),
		})
		return err
	}

	uc.logger.Info("Permissions synced successfully", service.Fields{
		"app_code":         in.AppCode,
		"app_id":           app.ID,
		"permissions_count": len(perms),
	})

	return nil
}
