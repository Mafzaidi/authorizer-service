package permission

import (
	"errors"
	"fmt"
	"time"

	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
	"localdev.me/authorizer/pkg/idgen"
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

func (u *permUsecase) Create(in *CreateInput) error {

	if in.AppID == "" {
		return errors.New("application ID is required")
	}

	if in.Code == "" {
		return errors.New("resource is required")
	}

	app, err := u.appRepo.GetByID(in.AppID)
	if err != nil {
		return fmt.Errorf("failed: %w", err)
	}

	existingPerm, _ := u.permRepo.GetByAppAndCode(app.ID, in.Code)
	if existingPerm != nil {
		return errors.New("permission already exists")
	}

	perm := &entity.Permission{
		ID:            idgen.NewUUIDv7(),
		ApplicationID: app.ID,
		Code:          in.Code,
		Description:   in.Description,
		Version:       1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return u.permRepo.Create(perm)
}

func (u *permUsecase) SyncPermissions(in *SyncInput) error {
	if in.AppCode == "" {
		return errors.New("application code is required")
	}

	app, err := u.appRepo.GetByCode(in.AppCode)
	if err != nil {
		return fmt.Errorf("failed: %w", err)
	}

	var perms []*entity.Permission
	for _, v := range in.Permissions {
		perm := &entity.Permission{
			ID:            idgen.NewUUIDv7(),
			ApplicationID: app.ID,
			Code:          v.Code,
			Description:   v.Description,
			Version:       1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		perms = append(perms, perm)
	}

	return u.permRepo.BulkUpsert(perms)
}
